package admins

import (
	"encoding/json"
	"fmt"
	"fun-coice/config"
	tgModel "fun-coice/internal/domain/commands/tg"
	"fun-coice/pkg/appStat"
	"fun-coice/pkg/scribble"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// TODO: rename admins to administration
type data struct {
	admin   tgModel.Commands
	user    tgModel.Commands
	bot     *tgbotapi.BotAPI
	DB      *scribble.Driver
	botName string
}

var _ = (tgModel.Service)(&data{})

func New(bot *tgbotapi.BotAPI, DB *scribble.Driver, userCommands tgModel.Commands, botName string) tgModel.Service {
	result := data{
		bot:     bot,
		DB:      DB,
		user:    userCommands,
		botName: botName,
	}
	commandsList := tgModel.NewCommands()

	commandsList["set"] = tgModel.Command{
		Command:     "/set",
		Description: "Set var",
		CommandType: "text",
		Permissions: tgModel.AdminPerms,
		Handler:     result.set,
	}
	commandsList["get"] = tgModel.Command{
		Command:     "/get",
		Description: "Set var",
		CommandType: "text",
		Permissions: tgModel.AdminPerms,
		Handler:     result.get,
	}
	commandsList["member"] = tgModel.Command{
		Command:     "/member",
		Description: "member info",
		CommandType: "text",
		Permissions: tgModel.AdminPerms,
		Handler:     result.member,
	}
	commandsList["rebuild"] = tgModel.Command{
		Command:     "/rebuild",
		Description: "rebuild app",
		CommandType: "text",
		Permissions: tgModel.AdminPerms,
		Handler:     result.rebuild,
	}
	commandsList["users"] = tgModel.Command{
		Command:     "/users",
		Description: "users list",
		CommandType: "text",
		Permissions: tgModel.AdminPerms,
		Handler:     result.users,
	}
	commandsList["addFeature"] = tgModel.Command{
		Command:     "/addFeature",
		Synonyms:    []string{"фича"},
		Description: "Создание описание фичи",
		CommandType: "text",
		Permissions: tgModel.AdminPerms,
		Handler:     result.addFeature,
	}
	commandsList["features"] = tgModel.Command{
		Command:     "/features",
		Synonyms:    []string{"фичи"},
		Description: "Список фич приложения",
		CommandType: "text",
		Permissions: tgModel.AdminPerms,
		Handler:     result.features,
	}
	commandsList["features"] = tgModel.Command{
		Command:     "/features",
		Synonyms:    []string{"фичи"},
		Description: "Список фич приложения",
		CommandType: "text",
		Permissions: tgModel.AdminPerms,
		Handler:     result.features,
	}

	commandsList["admin"] = tgModel.Command{
		Command:     "/admin",
		Synonyms:    []string{"admins"},
		Description: "Admin info",
		CommandType: "text",
		Permissions: tgModel.FreePerms,
		Handler:     result.info,
	}
	commandsList["command"] = tgModel.Command{
		Command:     "/command",
		Description: "Command info",
		CommandType: "text",
		Permissions: tgModel.FreePerms,
		Handler:     result.commandInfo,
	}
	commandsList["commands"] = tgModel.Command{
		Command:     "/commands",
		Synonyms:    []string{"help"},
		Description: "Список комманд",
		CommandType: "text",
		Permissions: tgModel.FreePerms,
		Handler:     result.commandsList,
	}
	commandsList["scanChat"] = tgModel.Command{
		Command:     "/scanChat",
		Description: "scan Chat",
		CommandType: "text",
		Permissions: tgModel.AdminPerms,
		Handler:     result.scanChat,
	}
	/*
		commandsList["scanChat"] = tgModel.Command{
			Command:     "/scanChat",
			Description: "scan Chat",
			CommandType: "text",
			Permissions: tgModel.AdminPerms,
			Handler:     result.scanChat,
		}
		/*
	*/
	// "/vars"

	result.admin = commandsList
	return &result
}

func (d data) Commands() tgModel.Commands {
	return d.admin
}

func (d data) Name() string {
	return "admins"
}

func (d data) commandInfo(msg *tgbotapi.Message, command *tgModel.Command) tgModel.HandlerResult {
	params := strings.Split(command.Arguments.Raw, " ")
	if len(params) < 2 {
		return tgModel.Simple(msg.Chat.ID, "Not found")
	}
	currentCommand, founded := d.user[params[1]]
	if !founded {
		return tgModel.Simple(msg.Chat.ID, "Not found")
	}
	info := fmt.Sprintf("Command: %s\nSynonyms: %s\nTriggers: %s\n\n%s",
		currentCommand.Command,
		strings.Join(currentCommand.Synonyms, ", "),
		currentCommand.Triggers,
		currentCommand.Description,
	)
	return tgModel.Simple(msg.Chat.ID, info)
}

func (d data) commandsList(msg *tgbotapi.Message, command *tgModel.Command) tgModel.HandlerResult {
	commandsList := "Commands:\n"
	for _, commandsItem := range d.user {
		if commandsItem.ListExclude {
			continue
		}
		commandsList += commandsItem.Command + " - " + commandsItem.Description + "\n"
	}
	return tgModel.Simple(msg.Chat.ID, commandsList)
}

func (d data) info(msg *tgbotapi.Message, command *tgModel.Command) tgModel.HandlerResult {
	return tgModel.Simple(msg.Chat.ID, "Admin is @"+config.TelegramAdminLogin(d.botName))
}

func (d data) vars(msg *tgbotapi.Message, command *tgModel.Command) tgModel.HandlerResult {
	params := strings.Split(command.Arguments.Raw, " ")
	if len(params) >= 3 {
		config.Set(params[1], params[2])
		return tgModel.Simple(msg.Chat.ID, "set "+params[1]+""+params[2])
	}
	return tgModel.EmptyCommand()
}

func (d data) set(msg *tgbotapi.Message, command *tgModel.Command) tgModel.HandlerResult {
	params := strings.Split(command.Arguments.Raw, " ")
	if len(params) >= 3 {
		config.Set(params[1], params[2])
		return tgModel.Simple(msg.Chat.ID, "set "+params[1]+""+params[2])
	}
	return tgModel.EmptyCommand()
}

func (d data) get(msg *tgbotapi.Message, command *tgModel.Command) tgModel.HandlerResult {
	params := strings.Split(command.Arguments.Raw, " ")
	if len(params) >= 3 {
		config.Set(params[1], params[2])
		return tgModel.Simple(msg.Chat.ID, "set "+params[1]+""+params[2])
	}
	return tgModel.EmptyCommand()
}

func (d data) member(msg *tgbotapi.Message, command *tgModel.Command) tgModel.HandlerResult {
	chatConfigWithUser := tgbotapi.ChatConfigWithUser{
		ChatID: msg.Chat.ID,
		UserID: msg.From.ID,
	}
	chatMember, _ := d.bot.GetChatMember(tgbotapi.GetChatMemberConfig{chatConfigWithUser})
	userInfo := tgModel.UserAndChatInfo(msg.From, msg.Chat) + fmt.Sprintf("\nMemberStatus: %s", chatMember.Status)
	return tgModel.Simple(msg.Chat.ID, userInfo)
}

func (d data) rebuild(msg *tgbotapi.Message, command *tgModel.Command) tgModel.HandlerResult {
	dir, err := os.Getwd()
	if err != nil {
		log.Printf("Failed to get dir: %v", err)
		return tgModel.Simple(msg.Chat.ID, "Failed to get dir: "+err.Error())
	}
	cmd := exec.Command("/bin/sh", dir+"/rebuild.sh")
	if err := cmd.Start(); err != nil {
		log.Printf("Failed to start cmd: %v", err)
		return tgModel.Simple(msg.Chat.ID, "Failed to start cmd: "+err.Error())
	}
	log.Println("Exit by command rebuild...")

	os.Exit(3)
	return tgModel.EmptyCommand()
}

func (d data) users(msg *tgbotapi.Message, command *tgModel.Command) tgModel.HandlerResult {
	records, err := d.DB.ReadAll("user")
	if err != nil {
		fmt.Println("Error", err)
	}

	userList := []string{}
	for _, f := range records {
		userFound := tgModel.User{}
		if err := json.Unmarshal([]byte(f), &userFound); err != nil {
			fmt.Println("Error", err)
		}
		userList = append(userList, "["+strconv.FormatInt(config.TelegramAdminId(d.botName), 10)+"] "+userFound.Name)
	}
	return tgModel.Simple(msg.Chat.ID, strings.Join(userList, "\n"))
}

func (d data) addFeature(msg *tgbotapi.Message, command *tgModel.Command) tgModel.HandlerResult {
	formattedMessage := ""
	d.DB.Read("features", "features", &formattedMessage)
	currentTime := time.Now().Format(time.RFC3339)
	formattedMessage += currentTime + " [" + appStat.Version + "]: " + command.Arguments.Raw

	if err := d.DB.Write("features", "features", formattedMessage); err != nil {
		fmt.Println("add command error", err)
		return tgModel.Simple(msg.Chat.ID, "Err: "+err.Error())
	}
	return tgModel.Simple(msg.Chat.ID, "saved")
}

func (d data) features(msg *tgbotapi.Message, command *tgModel.Command) tgModel.HandlerResult {
	formattedMessage := "-"
	d.DB.Read("features", "features", &formattedMessage)
	return tgModel.Simple(msg.Chat.ID, formattedMessage)
}

func (d data) scanChat(msg *tgbotapi.Message, command *tgModel.Command) tgModel.HandlerResult {
	params := strings.Split(command.Arguments.Raw, " ")
	fmt.Println("commandName", command.Command)
	fmt.Println("param", command.Arguments.Raw)
	fmt.Println("params", params)
	result := ""
	if len(params) < 2 {
		result = "Incorrect params"
		return tgModel.Simple(msg.Chat.ID, result)
	}
	chatId, _ := strconv.ParseInt(params[1], 10, 64)
	chat, err := d.bot.GetChat(tgbotapi.ChatInfoConfig{
		tgbotapi.ChatConfig{
			ChatID: chatId,
			//SuperGroupUsername: "",
		},
	})
	if err != nil {
		fmt.Println(err.Error())
		return tgModel.Simple(msg.Chat.ID, result)
	}

	chatMembersCount, err := d.bot.GetChatMembersCount(tgbotapi.ChatMemberCountConfig{
		tgbotapi.ChatConfig{
			ChatID: chatId,
			//SuperGroupUsername: "",
		},
	})
	if err != nil {
		fmt.Println(err.Error())
		return tgModel.Simple(msg.Chat.ID, result)
	}
	dop := chat.Type
	if chat.HasProtectedContent {
		dop += " ProtectedContent "
	}
	if chat.InviteLink != "" {
		dop += " Link: " + chat.InviteLink
	}
	result = fmt.Sprintf("Chat[%v] Users [%v] \n%v\nTitle: %s \n %s",
		chat.ID,
		dop,
		chatMembersCount,
		chat.Title,
		chat.Description,
	)
	err = d.DB.Write("chats", strconv.FormatInt(chat.ID, 10), chat)
	if err != nil {
		fmt.Println(err.Error())
		result += err.Error()
	}
	return tgModel.Simple(msg.Chat.ID, result)
}

// wait command
func (d data) fillChatUsersInfo(msg *tgbotapi.Message, command *tgModel.Command) tgModel.HandlerResult {
	var from int64
	var fromChat int64
	result := ""
	if msg.ForwardFrom == nil {
		return tgModel.Simple(msg.Chat.ID, "u need forward message")
	} else {
		from = msg.ForwardFrom.ID
		fromChat = msg.ForwardFromChat.ID
	}

	chatMemberInfo, err := d.bot.GetChatMember(tgbotapi.GetChatMemberConfig{
		ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
			ChatID: fromChat,
			UserID: from,
		},
	})
	if err != nil {
		fmt.Println(err.Error())
		return tgModel.Simple(msg.Chat.ID, "Get user err"+err.Error())
	}
	err = d.DB.Write("chat"+strconv.FormatInt(fromChat, 10), strconv.FormatInt(from, 10), chatMemberInfo)
	if err != nil {
		fmt.Println(err.Error())
		result += err.Error()
	}
	return tgModel.Simple(msg.Chat.ID, result)
}
