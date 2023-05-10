package admins

import (
	"encoding/json"
	"fmt"
	"fun-coice/config"
	tgCommands "fun-coice/internal/domain/commands/tg"
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

type data struct {
	admin tgCommands.Commands
	user  tgCommands.Commands
	bot   *tgbotapi.BotAPI
	DB    *scribble.Driver
}

var _ = (tgCommands.Service)(&data{})

func New(bot *tgbotapi.BotAPI, DB *scribble.Driver, userCommands tgCommands.Commands) tgCommands.Service {
	result := data{
		bot:  bot,
		DB:   DB,
		user: userCommands,
	}
	commandsList := make(tgCommands.Commands)
	commandsList["set"] = tgCommands.Command{
		Command:     "/set",
		Description: "Set var",
		CommandType: "text",
		Permissions: tgCommands.AdminPerms,
		Handler:     result.set,
	}
	commandsList["get"] = tgCommands.Command{
		Command:     "/get",
		Description: "Set var",
		CommandType: "text",
		Permissions: tgCommands.AdminPerms,
		Handler:     result.get,
	}
	commandsList["member"] = tgCommands.Command{
		Command:     "/member",
		Description: "member info",
		CommandType: "text",
		Permissions: tgCommands.AdminPerms,
		Handler:     result.member,
	}
	commandsList["rebuild"] = tgCommands.Command{
		Command:     "/rebuild",
		Description: "rebuild app",
		CommandType: "text",
		Permissions: tgCommands.AdminPerms,
		Handler:     result.rebuild,
	}
	commandsList["users"] = tgCommands.Command{
		Command:     "/users",
		Description: "users list",
		CommandType: "text",
		Permissions: tgCommands.AdminPerms,
		Handler:     result.users,
	}
	commandsList["addFeature"] = tgCommands.Command{
		Command:     "/addFeature",
		Synonyms:    []string{"фича"},
		Description: "Создание описание фичи",
		CommandType: "text",
		Permissions: tgCommands.AdminPerms,
		Handler:     result.addFeature,
	}
	commandsList["features"] = tgCommands.Command{
		Command:     "/features",
		Synonyms:    []string{"фичи"},
		Description: "Список фич приложения",
		CommandType: "text",
		Permissions: tgCommands.AdminPerms,
		Handler:     result.features,
	}
	commandsList["features"] = tgCommands.Command{
		Command:     "/features",
		Synonyms:    []string{"фичи"},
		Description: "Список фич приложения",
		CommandType: "text",
		Permissions: tgCommands.AdminPerms,
		Handler:     result.features,
	}

	commandsList["admin"] = tgCommands.Command{
		Command:     "/admin",
		Synonyms:    []string{"admins"},
		Description: "Admin info",
		CommandType: "text",
		Permissions: tgCommands.FreePerms,
		Handler:     result.info,
	}
	commandsList["command"] = tgCommands.Command{
		Command:     "/command",
		Description: "Command info",
		CommandType: "text",
		Permissions: tgCommands.FreePerms,
		Handler:     result.commandInfo,
	}
	commandsList["commands"] = tgCommands.Command{
		Command:     "/commands",
		Synonyms:    []string{"help"},
		Description: "Список комманд",
		CommandType: "text",
		Permissions: tgCommands.FreePerms,
		Handler:     result.commandsList,
	}
	commandsList["scanChat"] = tgCommands.Command{
		Command:     "/scanChat",
		Description: "scan Chat",
		CommandType: "text",
		Permissions: tgCommands.AdminPerms,
		Handler:     result.scanChat,
	}
	/*
		commandsList["scanChat"] = tgCommands.Command{
			Command:     "/scanChat",
			Description: "scan Chat",
			CommandType: "text",
			Permissions: tgCommands.AdminPerms,
			Handler:     result.scanChat,
		}
		/*
	*/
	// "/vars"

	result.admin = commandsList
	return &result
}

func (d data) Commands() tgCommands.Commands {
	return d.admin
}

func (d data) commandInfo(msg *tgbotapi.Message, commandName string, param string, params []string) (tgbotapi.Chattable, bool) {
	if len(params) < 2 {
		return tgbotapi.NewMessage(msg.Chat.ID, "Not found"), true
	}
	currentCommand, founded := d.user[params[1]]
	if !founded {
		return tgbotapi.NewMessage(msg.Chat.ID, "Not found"), true
	}
	info := fmt.Sprintf("Command: %s\nSynonyms: %s\nTriggers: %s\n\n%s",
		currentCommand.Command,
		strings.Join(currentCommand.Synonyms, ", "),
		currentCommand.Triggers,
		currentCommand.Description,
	)
	return tgbotapi.NewMessage(msg.Chat.ID, info), true
}

func (d data) commandsList(msg *tgbotapi.Message, commandName string, param string, params []string) (tgbotapi.Chattable, bool) {
	commandsList := "Commands:\n"
	for _, commandsItem := range d.user {
		commandsList += commandsItem.Command + " - " + commandsItem.Description + "\n"
	}
	return tgbotapi.NewMessage(msg.Chat.ID, commandsList), true
}

func (d data) info(msg *tgbotapi.Message, commandName string, param string, params []string) (tgbotapi.Chattable, bool) {
	return tgbotapi.NewMessage(msg.Chat.ID, "Admin is @"+config.TelegramAdminLogin()), true
}

func (d data) vars(msg *tgbotapi.Message, commandName string, param string, params []string) (tgbotapi.Chattable, bool) {
	if len(params) >= 3 {
		config.Set(params[1], params[2])
		return tgbotapi.NewMessage(msg.Chat.ID, "set "+params[1]+""+params[2]), true
	}
	return nil, false
}

func (d data) set(msg *tgbotapi.Message, commandName string, param string, params []string) (tgbotapi.Chattable, bool) {
	if len(params) >= 3 {
		config.Set(params[1], params[2])
		return tgbotapi.NewMessage(msg.Chat.ID, "set "+params[1]+""+params[2]), true
	}
	return nil, false
}

func (d data) get(msg *tgbotapi.Message, commandName string, param string, params []string) (tgbotapi.Chattable, bool) {
	if len(params) >= 3 {
		config.Set(params[1], params[2])
		return tgbotapi.NewMessage(msg.Chat.ID, "set "+params[1]+""+params[2]), true
	}
	return nil, false
}

func (d data) member(msg *tgbotapi.Message, commandName string, param string, params []string) (tgbotapi.Chattable, bool) {
	from := msg.From
	chat := msg.Chat
	chatConfigWithUser := tgbotapi.ChatConfigWithUser{
		ChatID:             chat.ID,
		SuperGroupUsername: "",
		UserID:             from.ID,
	}
	chatMember, _ := d.bot.GetChatMember(tgbotapi.GetChatMemberConfig{chatConfigWithUser})

	userInfo := fmt.Sprintf(
		"--== UserInfo==--\n"+
			"ID: %v\nUserName: %s\nFirstName: %s\nLastName: %s\nLanguageCode: %s"+
			"\n--==ChatInfo==--\n"+
			"ID: %v\nTitle: %s\nType: %s"+
			"\n--== MemberInfo==--\n"+
			"Status: %s",
		from.ID,
		from.UserName,
		from.FirstName,
		from.LastName,
		from.LanguageCode,
		chat.ID,
		chat.Title,
		chat.Type,
		chatMember.Status,
	)
	return tgbotapi.NewMessage(chat.ID, userInfo), true
}

func (d data) rebuild(msg *tgbotapi.Message, commandName string, param string, params []string) (tgbotapi.Chattable, bool) {
	dir, err := os.Getwd()
	if err != nil {
		log.Printf("Failed to get dir: %v", err)
		return tgbotapi.NewMessage(msg.Chat.ID, "Failed to get dir: "+err.Error()), true
	}
	cmd := exec.Command("/bin/sh", dir+"/rebuild.sh")
	if err := cmd.Start(); err != nil {
		log.Printf("Failed to start cmd: %v", err)
		return tgbotapi.NewMessage(msg.Chat.ID, "Failed to start cmd: "+err.Error()), true
	}
	log.Println("Exit by command rebuild...")

	os.Exit(3)
	return nil, false
}

func (d data) users(msg *tgbotapi.Message, commandName string, param string, params []string) (tgbotapi.Chattable, bool) {
	records, err := d.DB.ReadAll("user")
	if err != nil {
		fmt.Println("Error", err)
	}

	userList := []string{}
	for _, f := range records {
		userFound := tgCommands.User{}
		if err := json.Unmarshal([]byte(f), &userFound); err != nil {
			fmt.Println("Error", err)
		}
		userList = append(userList, "["+strconv.FormatInt(config.TelegramAdminId(), 10)+"] "+userFound.Name)
	}
	return tgbotapi.NewMessage(msg.Chat.ID, strings.Join(userList, "\n")), true
}

func (d data) addFeature(msg *tgbotapi.Message, commandName string, param string, params []string) (tgbotapi.Chattable, bool) {
	formattedMessage := ""
	d.DB.Read("features", "features", &formattedMessage)
	currentTime := time.Now().Format(time.RFC3339)
	formattedMessage += currentTime + " [" + appStat.Version + "]: " + param

	if err := d.DB.Write("features", "features", formattedMessage); err != nil {
		fmt.Println("add command error", err)
		return tgbotapi.NewMessage(msg.Chat.ID, "Err: "+err.Error()), true
	}
	return tgbotapi.NewMessage(msg.Chat.ID, "saved"), true
}

func (d data) features(msg *tgbotapi.Message, commandName string, param string, params []string) (tgbotapi.Chattable, bool) {
	formattedMessage := "-"
	d.DB.Read("features", "features", &formattedMessage)
	return tgbotapi.NewMessage(msg.Chat.ID, formattedMessage), true
}

func (d data) scanChat(msg *tgbotapi.Message, commandName string, param string, params []string) (tgbotapi.Chattable, bool) {
	fmt.Println("commandName", commandName)
	fmt.Println("param", param)
	fmt.Println("params", params)
	result := ""
	if len(params) < 2 {
		result = "Incorrect params"
		return tgbotapi.NewMessage(msg.Chat.ID, result), true
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
		return tgbotapi.NewMessage(msg.Chat.ID, result), true
	}

	chatMembersCount, err := d.bot.GetChatMembersCount(tgbotapi.ChatMemberCountConfig{
		tgbotapi.ChatConfig{
			ChatID: chatId,
			//SuperGroupUsername: "",
		},
	})
	if err != nil {
		fmt.Println(err.Error())
		return tgbotapi.NewMessage(msg.Chat.ID, result), true
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
	return tgbotapi.NewMessage(msg.Chat.ID, result), true
}

//wait command
func (d data) fillChatUsersInfo(msg *tgbotapi.Message, commandName string, param string, params []string) (tgbotapi.Chattable, bool) {
	var from int64
	var fromChat int64
	result := ""
	if msg.ForwardFrom == nil {
		return tgbotapi.NewMessage(msg.Chat.ID, "u need forward message"), true
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
		return tgbotapi.NewMessage(msg.Chat.ID, "Get user err"+err.Error()), true
	}
	err = d.DB.Write("chat"+strconv.FormatInt(fromChat, 10), strconv.FormatInt(from, 10), chatMemberInfo)
	if err != nil {
		fmt.Println(err.Error())
		result += err.Error()
	}
	return tgbotapi.NewMessage(msg.Chat.ID, result), true
}
