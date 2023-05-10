package text

import (
	tgCommands "fun-coice/internal/domain/commands/tg"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

type data struct {
	list tgCommands.Commands
}

func New() tgCommands.Service {
	result := data{}
	commandsList := make(tgCommands.Commands)
	commandsList["toLower"] = tgCommands.Command{
		Command:     "/toLower",
		Description: "String to low case",
		CommandType: "text",
		Permissions: tgCommands.FreePerms,
		Handler:     result.ToLower,
	}
	commandsList["toUpper"] = tgCommands.Command{
		Command:     "/toUpper",
		Description: "String to upper case",
		CommandType: "text",
		Permissions: tgCommands.FreePerms,
		Handler:     result.toUpper,
	}

	result.list = commandsList
	return &result
}

func (d data) Commands() tgCommands.Commands {
	return d.list
}

func (d data) ToLower(msg *tgbotapi.Message, commandName string, param string, params []string) (tgbotapi.Chattable, bool) {
	message := tgbotapi.NewMessage(msg.Chat.ID, strings.ToLower(param))
	return message, true
}

func (d data) toUpper(msg *tgbotapi.Message, commandName string, param string, params []string) (tgbotapi.Chattable, bool) {
	message := tgbotapi.NewMessage(msg.Chat.ID, strings.ToUpper(param))
	return message, true
}
