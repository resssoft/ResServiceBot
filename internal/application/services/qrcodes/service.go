package qrcodes

import (
	tgModel "fun-coice/internal/domain/commands/tg"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/skip2/go-qrcode"
)

type data struct {
	list tgModel.Commands
}

func New() tgModel.Service {
	result := data{}
	commandsList := tgModel.NewCommands()
	commandsList["qr"] = tgModel.Command{
		Command:     "/qr",
		Description: "String to QR image",
		CommandType: "text",
		Permissions: tgModel.FreePerms,
		Handler:     result.qr,
	}
	commandsList["qr256"] = tgModel.Command{
		Command:     "/qr256",
		Description: "String to QR image - 256px",
		CommandType: "text",
		Permissions: tgModel.FreePerms,
		Handler:     result.qr256,
	}

	result.list = commandsList
	return &result
}

func (d data) Commands() tgModel.Commands {
	return d.list
}

func (d data) Name() string {
	return "qr"
}

func (d data) qr(msg *tgbotapi.Message, command *tgModel.Command) *tgModel.HandlerResult {
	var qr []byte
	qr, err := qrcode.Encode(command.Arguments.Raw, qrcode.Medium, 1024)
	if err != nil {
		return tgModel.Simple(msg.Chat.ID, "failed to encode qrcode")
	}
	file := tgbotapi.FileBytes{Name: "qr.png", Bytes: qr}
	message := tgbotapi.NewPhoto(msg.Chat.ID, file)
	message.Caption = command.Arguments.Raw
	return tgModel.PreparedCommand(message)
}

func (d data) qr256(msg *tgbotapi.Message, command *tgModel.Command) *tgModel.HandlerResult {
	var qr []byte
	qr, err := qrcode.Encode(command.Arguments.Raw, qrcode.Medium, 256)
	if err != nil {
		return tgModel.Simple(msg.Chat.ID, "failed to encode qrcode")
	}
	file := tgbotapi.FileBytes{Name: "qr256.png", Bytes: qr}
	message := tgbotapi.NewPhoto(msg.Chat.ID, file)
	message.Caption = command.Arguments.Raw
	return tgModel.PreparedCommand(message)
}
