package qrcodes

import (
	tgCommands "fun-coice/internal/domain/commands/tg"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/skip2/go-qrcode"
)

type data struct {
	list tgCommands.Commands
}

func New() tgCommands.Service {
	result := data{}
	commandsList := make(tgCommands.Commands)
	commandsList["qr"] = tgCommands.Command{
		Command:     "/qr",
		Description: "String to QR image",
		CommandType: "text",
		Permissions: tgCommands.FreePerms,
		Handler:     result.qr,
	}
	commandsList["qr256"] = tgCommands.Command{
		Command:     "/qr256",
		Description: "String to QR image - 256px",
		CommandType: "text",
		Permissions: tgCommands.FreePerms,
		Handler:     result.qr256,
	}

	result.list = commandsList
	return &result
}

func (d data) Commands() tgCommands.Commands {
	return d.list
}

func (d data) qr(msg *tgbotapi.Message, commandName string, param string, params []string) (tgbotapi.Chattable, bool) {
	var qr []byte
	qr, err := qrcode.Encode(param, qrcode.Medium, 1024)
	if err != nil {
		return tgbotapi.NewMessage(msg.Chat.ID, "failed to encode qrcode"), true
	}
	file := tgbotapi.FileBytes{Name: "qr.png", Bytes: qr}
	message := tgbotapi.NewPhoto(msg.Chat.ID, file)
	message.Caption = param
	return message, true
}

func (d data) qr256(msg *tgbotapi.Message, commandName string, param string, params []string) (tgbotapi.Chattable, bool) {
	var qr []byte
	qr, err := qrcode.Encode(param, qrcode.Medium, 256)
	if err != nil {
		return tgbotapi.NewMessage(msg.Chat.ID, "failed to encode qrcode"), true
	}
	file := tgbotapi.FileBytes{Name: "qr256.png", Bytes: qr}
	message := tgbotapi.NewPhoto(msg.Chat.ID, file)
	message.Caption = param
	return message, true
}
