package b64

import (
	"encoding/base64"

	tgModel "fun-coice/internal/domain/commands/tg"
	tgbotapi "fun-coice/pkg/telegram-bot-api"
)

type data struct {
	list tgModel.Commands
}

func New() tgModel.Service {
	serv := data{
		list: tgModel.NewCommands(),
	}

	serv.list.AddSimple("b64",
		"Encode string to base64",
		serv.encode,
		"base64", "base64encode", "base64_encode",
	).AddSimple("b64d",
		"Decode string from base64",
		serv.decode,
		"base64d", "base64decode", "base64_decode",
	)

	return &serv
}

func (d *data) Commands() tgModel.Commands {
	return d.list
}

func (d *data) Name() string {
	return "b64"
}

func (d *data) Destroy() {}

func (d *data) Dependency() *tgModel.ServiceDepends {
	return nil
}

func (d *data) Configure(_ tgModel.ServiceConfig) error {
	return nil
}

func (d *data) encode(msg *tgbotapi.Message, command *tgModel.Command) *tgModel.HandlerResult {
	b64result := base64.StdEncoding.EncodeToString([]byte(command.Arguments.Raw))
	return tgModel.SimpleReply(msg.Chat.ID, b64result, msg.MessageID)
}

func (d *data) decode(msg *tgbotapi.Message, command *tgModel.Command) *tgModel.HandlerResult {
	b64result, _ := base64.StdEncoding.DecodeString(command.Arguments.Raw)
	return tgModel.SimpleReply(msg.Chat.ID, string(b64result), msg.MessageID)
}
