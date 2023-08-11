package tgCommands

import (
	"bufio"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"os"
	"regexp"
	"strings"
)

type User struct {
	UserID  int64
	ChatId  int64
	Login   string
	Name    string
	IsAdmin bool
}

type Command struct {
	Command     string
	Synonyms    []string
	Triggers    []string
	Templates   []string
	Description string
	CommandType string
	ListExclude bool
	Permissions CommandPermissions
	Handler     HandlerFunc //TODO: REMOVE params []string and fix in the implements
	//Service string  // set in the bot only
}

//OLD func(*tgbotapi.Message, string, string, []string) (tgbotapi.Chattable, bool) tgCommands.HandlerResult tgCommands.PreparedCommand( tgCommands.PreparedCommand
// tgCommands.PreparedCommand(tgbotapi.NewMessage ->  tgCommands.Simple
//TODO: Handler     func(*tgbotapi.Message, string, string, []string) (tgbotapi.Chattable, HandlerResult)

type HandlerResult struct {
	Prepared bool   // command is prepared for sending
	Wait     bool   // wait next command
	Next     string // next command
	Messages []tgbotapi.Chattable
	Events   []Event
}

// TODO: REMOVE params []string and fix in the implements
type HandlerFunc func(*tgbotapi.Message, string, string, []string) HandlerResult

func EmptyCommand() HandlerResult {
	return HandlerResult{
		Messages: nil,
	}
}

func PreparedCommand(chatEvents ...tgbotapi.Chattable) HandlerResult {
	return HandlerResult{
		Prepared: true,
		Messages: chatEvents,
	}
}

func Simple(chatId int64, text string) HandlerResult {
	return PreparedCommand(tgbotapi.NewMessage(chatId, text))
}

func SimpleReply(chatId int64, text string, replyTo int) HandlerResult {
	newMsg := tgbotapi.NewMessage(chatId, text)
	newMsg.ReplyToMessageID = replyTo
	return PreparedCommand(newMsg)
}

func UnPreparedCommand(chatEvent tgbotapi.Chattable) HandlerResult {
	return HandlerResult{
		Messages: []tgbotapi.Chattable{chatEvent},
	}
}

func WaitingCommand(command string) HandlerResult {
	return HandlerResult{
		Wait: true,
		Next: command,
	}
}

func WaitingWithText(chatId int64, text, command string) HandlerResult {
	return HandlerResult{
		Wait:     true,
		Prepared: true,
		Messages: []tgbotapi.Chattable{tgbotapi.NewMessage(chatId, text)},
		Next:     command,
	}
}

func WaitingPreparedCommand(chatEvent tgbotapi.Chattable) HandlerResult {
	return HandlerResult{
		Wait:     true,
		Prepared: true,
		Messages: []tgbotapi.Chattable{chatEvent},
	}
}

func (hr HandlerResult) WithEvent(newEvent Event) HandlerResult {
	hr.Events = append(hr.Events, newEvent)
	return hr
}

func (t *Command) IsImplemented(msg, botName string) bool {
	if IsCommand(t.Command, msg, botName) {
		return true
	}
	for _, synonym := range t.Synonyms {
		if IsCommand(synonym, msg, botName) {
			return true
		}
	}
	return false
}

func (t *Command) IsMatched(msg, botName string) bool {
	if len(t.Templates) > 0 {
		for _, template := range t.Templates {
			templateMatched, _ := regexp.MatchString(template, msg)
			if templateMatched {
				return true
			}
		}
	}
	return false
}

func (t *Command) Permission(messageItem *tgbotapi.Message, adminId int64) bool {
	if messageItem != nil {
		if messageItem == nil {
			return false
		}
		switch messageItem.Chat.Type {
		case "private":
			if t.Permissions.Check(messageItem.From, adminId) {
				return true
			}
		case "chat":
			if t.Permissions.Check(messageItem.From, adminId) {
				return true
			}
		}
	}
	return false
}

func IsCommand(command, msg, botName string) bool {
	return msg == command || msg == fmt.Sprintf("%s@%s", command, botName)
}

func (t *Command) ParsedArgs(args string) []string {
	return strings.Split(args, " ")
}

func (tgp *CommandPermissions) Check(user *tgbotapi.User, adminId int64) bool {
	if tgp.UserPermissions == "all" {
		return true
	}
	if tgp.UserPermissions == "admin" && user.ID == adminId {
		return true
	}
	return false
}

type Commands map[string]Command

type CommandPermissions struct {
	UserPermissions string
	ChatPermissions string
}

var FreePerms = CommandPermissions{
	ChatPermissions: "all",
	UserPermissions: "all",
}

var AdminPerms = CommandPermissions{
	ChatPermissions: "admin",
	UserPermissions: "admin",
}

var ModerPerms = CommandPermissions{
	ChatPermissions: "moder",
	UserPermissions: "moder",
}

func NewCommands() Commands {
	return make(Commands)
}

func (cs Commands) Merge(list Commands) Commands {
	merged := make(Commands)
	for key, value := range cs {
		merged[key] = value
	}
	for key, value := range list {
		merged[key] = value
	}
	return merged
}

func (cs Commands) Add(
	name, description, commandType string,
	synonyms, triggers, templates []string,
	listExclude bool,
	permissions CommandPermissions,
	handler HandlerFunc,
) Commands {
	if cs == nil {
		cs = make(Commands)
	}
	cs[name] = Command{
		Command:     "/" + name,
		Synonyms:    synonyms,
		Triggers:    triggers,
		Templates:   templates,
		Description: description,
		CommandType: commandType,
		ListExclude: listExclude,
		Permissions: permissions,
		Handler:     handler,
	}
	return cs
}

func (cs Commands) AddSimple(
	name, description string,
	handler HandlerFunc,
) Commands {
	if cs == nil {
		cs = make(Commands)
	}
	cs[name] = Command{
		Command:     "/" + name,
		Description: description,
		CommandType: "text",
		Permissions: FreePerms,
		Handler:     handler,
	}
	return cs
}

func (cs Commands) Exclude() Commands {
	for index, item := range cs {
		item.ListExclude = true
		cs[index] = item
	}
	return cs
}

func writeLines(lines []string, path string) error {

	// overwrite file if it exists
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}

	defer file.Close()

	// new writer w/ default 4096 buffer size
	w := bufio.NewWriter(file)

	for _, line := range lines {
		_, err := w.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}

	// flush outstanding data
	return w.Flush()
}

type TgFileInfo struct {
	Ok     bool `json:"ok,omitempty"`
	Result struct {
		FileId       string `json:"file_id,omitempty"`
		FileUniqueId string `json:"file_unique_id,omitempty"`
		FileSize     int    `json:"file_size,omitempty"`
		FilePath     string `json:"file_path,omitempty"`
	} `json:"result,omitempty"`
}

type KeyBoardTG struct {
	Rows []KeyBoardRowTG
}

type KeyBoardRowTG struct {
	Buttons []KeyBoardButtonTG
}

type KeyBoardButtonTG struct {
	Text string
	Data string
}

type TGUser struct {
	UserID  int64
	ChatId  int64
	Login   string
	Name    string
	IsAdmin bool
}

type Event struct {
	Name ChatEvent
	Msg  *tgbotapi.Message
}

func NewEvent(name ChatEvent, msg *tgbotapi.Message) Event {
	return Event{
		Name: name,
		Msg:  msg,
	}
}

type ChatEvent string

const (
	StartBotEvent        ChatEvent = "start" //triggered by /start command from the bot
	UserLeaveChantEvent  ChatEvent = "user_leave_chat"
	UserJoinedChantEvent ChatEvent = "user_joined_chat"
	TextMsgBotEvent      ChatEvent = "text_msg" //triggered by /start command from the bot
)

func (ce ChatEvent) String() string {
	return string(ce)
}

func UserAndChatInfo(user *tgbotapi.User, chat *tgbotapi.Chat) string {
	return UserInfo(user) + ChatInfo(chat)
}

func UserInfo(user *tgbotapi.User) string {
	if user == nil {
		return ""
	}
	userLogin := ""
	if user.UserName != "" {
		userLogin += fmt.Sprintf("(@%s)", user.UserName)
	}
	if user.LanguageCode != "" {
		userLogin += fmt.Sprintf("(%s)", user.LanguageCode)
	}
	userInfo := fmt.Sprintf("User: [%v] %s %s %s",
		user.ID,
		userLogin,
		user.FirstName,
		user.LastName,
	)
	return userInfo
}

func ChatInfo(chat *tgbotapi.Chat) string {
	if chat != nil {
		return fmt.Sprintf("\nChat: [%v] (%s): %s",
			chat.ID,
			chat.Type,
			chat.Title,
		)
	}
	return ""
}

type SentMessages chan<- tgbotapi.Chattable
