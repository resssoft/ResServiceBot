package text

import (
	tgModel "fun-coice/internal/domain/commands/tg"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

type data struct {
	list tgModel.Commands
}

func New() tgModel.Service {
	result := data{}
	commandsList := tgModel.NewCommands()
	commandsList["toLower"] = tgModel.Command{
		Command:     "/toLower",
		Description: "String to low case",
		CommandType: "text",
		Permissions: tgModel.FreePerms,
		Handler:     result.ToLower,
	}
	commandsList["toUpper"] = tgModel.Command{
		Command:     "/toUpper",
		Description: "String to upper case",
		CommandType: "text",
		Permissions: tgModel.FreePerms,
		Handler:     result.toUpper,
	}

	result.list = commandsList
	return &result
}

func (d data) Commands() tgModel.Commands {
	return d.list
}

func (d data) ToLower(msg *tgbotapi.Message, command *tgModel.Command) tgModel.HandlerResult {
	return tgModel.Simple(msg.Chat.ID, strings.ToLower(command.Arguments.Raw))
}

func (d data) toUpper(msg *tgbotapi.Message, command *tgModel.Command) tgModel.HandlerResult {
	return tgModel.Simple(msg.Chat.ID, strings.ToUpper(command.Arguments.Raw))
}
