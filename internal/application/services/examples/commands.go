package examples

import (
	"fmt"
	tgModel "fun-coice/internal/domain/commands/tg"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (d *data) help(msg *tgbotapi.Message, _ *tgModel.Command) *tgModel.HandlerResult {
	commandsList := "Examples\n"
	commandsList += "Commands:\n"
	for _, commandsItem := range d.Commands() {
		commandsList += "/" + commandsItem.Command + " - " + commandsItem.Description + "\n"
	}
	return tgModel.Simple(msg.Chat.ID, commandsList)
}

func (d *data) exampleText(msg *tgbotapi.Message, _ *tgModel.Command) *tgModel.HandlerResult {
	return tgModel.Simple(msg.Chat.ID, "just text")
}

func (d *data) exampleReply(msg *tgbotapi.Message, _ *tgModel.Command) *tgModel.HandlerResult {
	return tgModel.SimpleReply(msg.Chat.ID, "text with reply", msg.MessageID)
}

func (d *data) exampleShowInlineButtons(msg *tgbotapi.Message, _ *tgModel.Command) *tgModel.HandlerResult {
	newMsg := tgbotapi.NewMessage(msg.Chat.ID, "Text under buttons")
	newMsg.ReplyMarkup = tgModel.GetTGButtons(tgModel.KBRows(tgModel.KBButs(
		tgModel.KeyBoardButtonTG{Text: "Just button", Data: "example_button"},
		tgModel.KeyBoardButtonTG{Text: "Edit buttons", Data: "example_buttons_edit"},
		tgModel.KeyBoardButtonTG{Text: "Remove this message", Data: "example_remove_buttons_trigger"},
	)))
	return tgModel.PreparedCommand(newMsg)
}

func (d *data) exampleEditInlineButtons(msg *tgbotapi.Message, _ *tgModel.Command) *tgModel.HandlerResult {
	newMsg := tgbotapi.NewEditMessageText(msg.Chat.ID, msg.MessageID, "Text under buttons edited")
	newMsg.ReplyMarkup = tgbotapi.NewEditMessageReplyMarkup(
		msg.Chat.ID,
		msg.MessageID,
		*tgModel.GetTGButtons(tgModel.KBRows(tgModel.KBButs(
			tgModel.KeyBoardButtonTG{Text: fmt.Sprintf("Counter (%v)", d.Counter()), Data: "example_button_counter"},
			tgModel.KeyBoardButtonTG{Text: "Remove this message", Data: "example_remove_buttons_trigger"},
		)))).ReplyMarkup
	return tgModel.PreparedCommand(newMsg)
}

func (d *data) exampleRemoveButtons(msg *tgbotapi.Message, _ *tgModel.Command) *tgModel.HandlerResult {
	fmt.Println("example_remove_buttons_trigger", msg.Chat.ID, msg.MessageID)
	newMsg := tgbotapi.NewDeleteMessage(msg.Chat.ID, msg.MessageID)
	return tgModel.PreparedCommand(newMsg)
}

func (d *data) exampleNotify(msg *tgbotapi.Message, _ *tgModel.Command) *tgModel.HandlerResult {
	newMsg := tgbotapi.NewChatAction(msg.Chat.ID, "typing")
	return tgModel.PreparedCommand(newMsg)
}

func (d *data) exampleCounterIncrement(msg *tgbotapi.Message, _ *tgModel.Command) *tgModel.HandlerResult {
	fmt.Println("exampleCounterIncrement", msg.Chat.ID, msg.MessageID)
	newMsg := tgbotapi.NewEditMessageText(msg.Chat.ID, msg.MessageID, "Text under buttons, click counter")
	newMsg.ReplyMarkup = tgbotapi.NewEditMessageReplyMarkup(
		msg.Chat.ID,
		msg.MessageID,
		*tgModel.GetTGButtons(tgModel.KBRows(tgModel.KBButs(
			tgModel.KeyBoardButtonTG{Text: fmt.Sprintf("Counter (%v)", d.Counter()), Data: "example_button_counter"},
			tgModel.KeyBoardButtonTG{Text: "Action test", Data: "example_notify"},
			tgModel.KeyBoardButtonTG{Text: "Remove this message", Data: "example_remove_buttons_trigger"},
		)))).ReplyMarkup
	return tgModel.PreparedCommand(newMsg)
}
