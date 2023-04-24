package translate

import (
	"fmt"
	tgCommands "fun-coice/internal/domain/commands/tg"
	gt "github.com/bas24/googletranslatefree"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type data struct {
	list tgCommands.Commands
}

func New() tgCommands.Service {
	result := data{}
	commandsList := make(tgCommands.Commands)
	commandsList["tr"] = tgCommands.Command{
		Command:     "/tr",
		Synonyms:    []string{"tran", "translate", "gtr"},
		Description: "Translate",
		CommandType: "text",
		Permissions: tgCommands.FreePerms,
		Handler:     result.tr,
	}
	commandsList["tr_hy_ru"] = tgCommands.Command{
		Command:     "/tr_hy_ru",
		Description: "Translate from armenian to russian",
		CommandType: "text",
		Permissions: tgCommands.FreePerms,
		Handler:     result.tr_hy_ru,
	}
	commandsList["tr_ru_hy"] = tgCommands.Command{
		Command:     "/tr_ru_hy",
		Description: "Translate from russian to armenian",
		CommandType: "text",
		Permissions: tgCommands.FreePerms,
		Handler:     result.tr_ru_hy,
	}

	result.list = commandsList
	return &result
}

func (d data) Commands() tgCommands.Commands {
	return d.list
}

func (d data) tr(msg *tgbotapi.Message, commandName string, param string, params []string) (tgbotapi.Chattable, bool) {
	result, err := gt.Translate(param, params[1], params[2])
	if err != nil {
		result = err.Error()
		fmt.Println(err)
	}
	msgNew := tgbotapi.NewMessage(msg.Chat.ID, result)
	msgNew.ReplyToMessageID = msg.MessageID
	return msgNew, true
}

func (d data) tr_hy_ru(msg *tgbotapi.Message, commandName string, param string, params []string) (tgbotapi.Chattable, bool) {
	result, err := gt.Translate(param, "hy", "ru")
	if err != nil {
		result = err.Error()
		fmt.Println(err)
	}
	msgNew := tgbotapi.NewMessage(msg.Chat.ID, result)
	msgNew.ReplyToMessageID = msg.MessageID
	return msgNew, true
}

func (d data) tr_ru_hy(msg *tgbotapi.Message, commandName string, param string, params []string) (tgbotapi.Chattable, bool) {
	result, err := gt.Translate(param, "ru", "hy")
	if err != nil {
		result = err.Error()
		fmt.Println(err)
	}
	msgNew := tgbotapi.NewMessage(msg.Chat.ID, result)
	msgNew.ReplyToMessageID = msg.MessageID
	return msgNew, true
}
