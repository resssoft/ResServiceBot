package translate

import (
	"fmt"
	tgModel "fun-coice/internal/domain/commands/tg"
	gt "github.com/bas24/googletranslatefree"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type data struct {
	list tgModel.Commands
}

func New() tgModel.Service {
	result := data{}
	commandsList := tgModel.NewCommands()
	commandsList["tr"] = tgModel.Command{
		Command:     "/tr",
		Synonyms:    []string{"tran", "translate", "gtr"},
		Description: "Translate",
		CommandType: "text",
		Permissions: tgModel.FreePerms,
		Handler:     result.tr,
	}
	commandsList["tr_hy_ru"] = tgModel.Command{
		Command:     "/tr_hy_ru",
		Description: "Translate from armenian to russian",
		CommandType: "text",
		Permissions: tgModel.FreePerms,
		Handler:     result.tr_hy_ru,
	}
	commandsList["tr_ru_hy"] = tgModel.Command{
		Command:     "/tr_ru_hy",
		Description: "Translate from russian to armenian",
		CommandType: "text",
		Permissions: tgModel.FreePerms,
		Handler:     result.tr_ru_hy,
	}
	commandsList["tr_ru_en"] = tgModel.Command{
		Command:     "/tr_ru_en",
		Description: "Translate from russian to english",
		CommandType: "text",
		Permissions: tgModel.FreePerms,
		Handler:     result.tr_ru_en,
	}
	commandsList["tr_en_ru"] = tgModel.Command{
		Command:     "/tr_en_ru",
		Description: "Translate from english to russian",
		CommandType: "text",
		Permissions: tgModel.FreePerms,
		Handler:     result.tr_en_ru,
	}
	commandsList["перевод"] = tgModel.Command{
		Command:     "/перевод",
		Synonyms:    []string{"перевод", "переведи"},
		Description: "Translate from english to russian",
		CommandType: "text",
		Permissions: tgModel.FreePerms,
		Handler:     result.trNext,
	}
	commandsList["tr_settings"] = tgModel.Command{
		Command:     "/tr_settings",
		Description: "Translate settings",
		CommandType: "text",
		Permissions: tgModel.FreePerms,
		Handler:     result.trSettings,
	}
	//TODO: /tr_revert - translate revert direction for last message
	//TODO: add for messages commands info and bot ADDV
	//TODO: use buttongs for choice direction
	result.list = commandsList
	return &result
}

func (d data) Commands() tgModel.Commands {
	return d.list
}

func (d data) Name() string {
	return "translate"
}

func (d data) tr(msg *tgbotapi.Message, command *tgModel.Command) *tgModel.HandlerResult {
	params := command.Arguments.Parse()
	if len(params) == 0 {
		//TODO: add settings by user for default translate (choice 2 langs OR choice old translates options) OR BUTTONS
		return tgModel.DeferredWithText(msg.Chat.ID, "Enter for translate from EN to RU (default for /tr_settings)", "tr_en_ru", "", nil)
	}
	result, err := gt.Translate(command.Arguments.Raw, params[1], params[2])
	if err != nil {
		result = err.Error()
		fmt.Println(err)
	}
	return tgModel.SimpleReply(msg.Chat.ID, result, msg.MessageID)
}

func (d data) trSettings(msg *tgbotapi.Message, command *tgModel.Command) *tgModel.HandlerResult {
	return tgModel.SimpleReply(msg.Chat.ID, "Not implemented, it will be later", msg.MessageID)
}

func (d data) tr_hy_ru(msg *tgbotapi.Message, command *tgModel.Command) *tgModel.HandlerResult {
	if command.Arguments.Raw == "" {
		return tgModel.SimpleReply(msg.Chat.ID, "Error: empty message!", msg.MessageID)
	}
	result, err := gt.Translate(command.Arguments.Raw, "hy", "ru")
	if err != nil {
		result = err.Error()
		fmt.Println(err)
	}
	return tgModel.SimpleReply(msg.Chat.ID, result, msg.MessageID)
}

func (d data) tr_ru_hy(msg *tgbotapi.Message, command *tgModel.Command) *tgModel.HandlerResult {
	if command.Arguments.Raw == "" {
		return tgModel.SimpleReply(msg.Chat.ID, "Error: empty message!", msg.MessageID)
	}
	result, err := gt.Translate(command.Arguments.Raw, "ru", "hy")
	if err != nil {
		result = err.Error()
		fmt.Println(err)
	}
	return tgModel.SimpleReply(msg.Chat.ID, result, msg.MessageID)
}

func (d data) tr_ru_en(msg *tgbotapi.Message, command *tgModel.Command) *tgModel.HandlerResult {
	if command.Arguments.Raw == "" {
		return tgModel.SimpleReply(msg.Chat.ID, "Error: empty message!", msg.MessageID)
	}
	result, err := gt.Translate(command.Arguments.Raw, "ru", "en")
	if err != nil {
		result = err.Error()
		fmt.Println(err)
	}
	return tgModel.SimpleReply(msg.Chat.ID, result, msg.MessageID)
}

func (d data) tr_en_ru(msg *tgbotapi.Message, command *tgModel.Command) *tgModel.HandlerResult {
	if command.Arguments.Raw == "" {
		return tgModel.SimpleReply(msg.Chat.ID, "Error: empty message!", msg.MessageID)
	}
	result, err := gt.Translate(command.Arguments.Raw, "en", "ru")
	if err != nil {
		result = err.Error()
		fmt.Println(err)
	}
	return tgModel.SimpleReply(msg.Chat.ID, result, msg.MessageID)
}

func (d data) trNext(msg *tgbotapi.Message, command *tgModel.Command) *tgModel.HandlerResult {
	return tgModel.DeferredWithText(msg.Chat.ID, "Enter for translate from EN to RU", "tr_en_ru", "", nil)
}
