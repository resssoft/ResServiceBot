package admins

import (
	"fmt"
	"fun-coice/config"
	tgCommands "fun-coice/internal/domain/commands/tg"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type data struct {
	list tgCommands.Commands
	bot  *tgbotapi.BotAPI
}

var perm = tgCommands.AdminPerms

func New(bot *tgbotapi.BotAPI) tgCommands.Service {
	result := data{
		bot: bot,
	}
	commandsList := make(tgCommands.Commands)
	commandsList["admin"] = tgCommands.Command{
		Command:     "/admin",
		Synonyms:    []string{"admins"},
		Description: "Admin info",
		CommandType: "text",
		Permissions: tgCommands.FreePerms,
		Handler:     result.info,
	}
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
	// "/vars"

	result.list = commandsList
	return &result
}

func (d data) Commands() tgCommands.Commands {
	return d.list
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
