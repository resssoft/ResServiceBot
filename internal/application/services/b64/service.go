package b64

import (
	"encoding/base64"
	tgModel "fun-coice/internal/domain/commands/tg"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type data struct {
	list tgModel.Commands
}

func New() tgModel.Service {
	result := data{}
	commandsList := tgModel.NewCommands()
	commandsList["b64"] = tgModel.Command{
		Command:     "/b64",
		Synonyms:    []string{"base64", "base64encode", "base64_encode"},
		Description: "Encode string to base64",
		CommandType: "text",
		Permissions: tgModel.FreePerms,
		Handler:     result.encode,
	}
	commandsList["b64d"] = tgModel.Command{
		Command:     "/b64d",
		Synonyms:    []string{"base64d", "base64decode", "base64_decode"},
		Description: "Decode string from base64",
		CommandType: "text",
		Permissions: tgModel.FreePerms,
		Handler:     result.decode,
	}

	result.list = commandsList
	return &result
}

func (d data) Commands() tgModel.Commands {
	return d.list
}

func (d data) Name() string {
	return "b64"
}

func (d data) decode(msg *tgbotapi.Message, command *tgModel.Command) *tgModel.HandlerResult {
	b64result, _ := base64.StdEncoding.DecodeString(command.Arguments.Raw)
	return tgModel.SimpleReply(msg.Chat.ID, string(b64result), msg.MessageID)
}

func (d data) encode(msg *tgbotapi.Message, command *tgModel.Command) *tgModel.HandlerResult {
	b64result := base64.StdEncoding.EncodeToString([]byte(command.Arguments.Raw))
	return tgModel.SimpleReply(msg.Chat.ID, b64result, msg.MessageID)
}
