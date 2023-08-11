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
	commandsList := tgCommands.NewCommands()
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
	commandsList["tr_ru_en"] = tgCommands.Command{
		Command:     "/tr_ru_en",
		Description: "Translate from russian to english",
		CommandType: "text",
		Permissions: tgCommands.FreePerms,
		Handler:     result.tr_ru_en,
	}
	commandsList["tr_en_ru"] = tgCommands.Command{
		Command:     "/tr_en_ru",
		Description: "Translate from english to russian",
		CommandType: "text",
		Permissions: tgCommands.FreePerms,
		Handler:     result.tr_en_ru,
	}
	commandsList["перевод"] = tgCommands.Command{
		Command:     "/перевод",
		Synonyms:    []string{"перевод", "переведи"},
		Description: "Translate from english to russian",
		CommandType: "text",
		Permissions: tgCommands.FreePerms,
		Handler:     result.trNext,
	}

	result.list = commandsList
	return &result
}

func (d data) Commands() tgCommands.Commands {
	return d.list
}

func (d data) tr(msg *tgbotapi.Message, commandName string, param string, params []string) tgCommands.HandlerResult {
	result, err := gt.Translate(param, params[1], params[2])
	if err != nil {
		result = err.Error()
		fmt.Println(err)
	}
	return tgCommands.SimpleReply(msg.Chat.ID, result, msg.MessageID)
}

func (d data) tr_hy_ru(msg *tgbotapi.Message, commandName string, param string, params []string) tgCommands.HandlerResult {
	result, err := gt.Translate(param, "hy", "ru")
	if err != nil {
		result = err.Error()
		fmt.Println(err)
	}
	return tgCommands.SimpleReply(msg.Chat.ID, result, msg.MessageID)
}

func (d data) tr_ru_hy(msg *tgbotapi.Message, commandName string, param string, params []string) tgCommands.HandlerResult {
	result, err := gt.Translate(param, "ru", "hy")
	if err != nil {
		result = err.Error()
		fmt.Println(err)
	}
	return tgCommands.SimpleReply(msg.Chat.ID, result, msg.MessageID)
}

func (d data) tr_ru_en(msg *tgbotapi.Message, commandName string, param string, params []string) tgCommands.HandlerResult {
	result, err := gt.Translate(param, "ru", "en")
	if err != nil {
		result = err.Error()
		fmt.Println(err)
	}
	return tgCommands.SimpleReply(msg.Chat.ID, result, msg.MessageID)
}

func (d data) tr_en_ru(msg *tgbotapi.Message, commandName string, param string, params []string) tgCommands.HandlerResult {
	result, err := gt.Translate(param, "en", "ru")
	if err != nil {
		result = err.Error()
		fmt.Println(err)
	}
	return tgCommands.SimpleReply(msg.Chat.ID, result, msg.MessageID)
}

func (d data) trNext(msg *tgbotapi.Message, commandName string, param string, params []string) tgCommands.HandlerResult {
	return tgCommands.WaitingWithText(msg.Chat.ID, "Enter for translate from EN to RU", "tr_en_ru")
}
