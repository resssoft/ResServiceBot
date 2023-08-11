package b64

import (
	"encoding/base64"
	tgCommands "fun-coice/internal/domain/commands/tg"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type data struct {
	list tgCommands.Commands
}

func New() tgCommands.Service {
	result := data{}
	commandsList := tgCommands.NewCommands()
	commandsList["b64"] = tgCommands.Command{
		Command:     "/b64",
		Synonyms:    []string{"base64", "base64encode", "base64_encode"},
		Description: "Encode string to base64",
		CommandType: "text",
		Permissions: tgCommands.FreePerms,
		Handler:     result.encode,
	}
	commandsList["b64d"] = tgCommands.Command{
		Command:     "/b64d",
		Synonyms:    []string{"base64d", "base64decode", "base64_decode"},
		Description: "Decode string from base64",
		CommandType: "text",
		Permissions: tgCommands.FreePerms,
		Handler:     result.decode,
	}

	result.list = commandsList
	return &result
}

func (d data) Commands() tgCommands.Commands {
	return d.list
}

func (d data) decode(msg *tgbotapi.Message, commandName string, param string, params []string) tgCommands.HandlerResult {
	b64result, _ := base64.StdEncoding.DecodeString(param)
	return tgCommands.SimpleReply(msg.Chat.ID, string(b64result), msg.MessageID)
}

func (d data) encode(msg *tgbotapi.Message, commandName string, param string, params []string) tgCommands.HandlerResult {
	b64result := base64.StdEncoding.EncodeToString([]byte(param))
	return tgCommands.SimpleReply(msg.Chat.ID, b64result, msg.MessageID)
}
