package examples

import (
	"fmt"
	tgCommands "fun-coice/internal/domain/commands/tg"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (d *data) help(msg *tgbotapi.Message, commandName string, param string, params []string) tgCommands.HandlerResult {
	commandsList := "Examples\n"
	commandsList += "Commands:\n"
	for _, commandsItem := range d.Commands() {
		commandsList += commandsItem.Command + " - " + commandsItem.Description + "\n"
	}
	return tgCommands.Simple(msg.Chat.ID, commandsList)
}

func (d *data) exampleText(msg *tgbotapi.Message, commandName string, param string, params []string) tgCommands.HandlerResult {
	return tgCommands.Simple(msg.Chat.ID, "just text")
}

func (d *data) exampleReply(msg *tgbotapi.Message, commandName string, param string, params []string) tgCommands.HandlerResult {
	return tgCommands.SimpleReply(msg.Chat.ID, "text with reply", msg.MessageID)
}

func (d *data) exampleShowInlineButtons(msg *tgbotapi.Message, commandName string, param string, params []string) tgCommands.HandlerResult {
	newMsg := tgbotapi.NewMessage(msg.Chat.ID, "Text under buttons")
	newMsg.ReplyMarkup = tgCommands.GetTGButtons(tgCommands.KBRows(tgCommands.KBButs(
		tgCommands.KeyBoardButtonTG{Text: "Just button", Data: "example_button"},
		tgCommands.KeyBoardButtonTG{Text: "Edit buttons", Data: "example_buttons_edit"},
		tgCommands.KeyBoardButtonTG{Text: "Remove this message", Data: "example_remove_buttons_trigger"},
	)))
	return tgCommands.PreparedCommand(newMsg)
}

func (d *data) exampleEditInlineButtons(msg *tgbotapi.Message, commandName string, param string, params []string) tgCommands.HandlerResult {
	newMsg := tgbotapi.NewEditMessageText(msg.Chat.ID, msg.MessageID, "Text under buttons edited")
	newMsg.ReplyMarkup = tgbotapi.NewEditMessageReplyMarkup(
		msg.Chat.ID,
		msg.MessageID,
		tgCommands.GetTGButtons(tgCommands.KBRows(tgCommands.KBButs(
			tgCommands.KeyBoardButtonTG{Text: fmt.Sprintf("Counter (%v)", d.Counter()), Data: "example_button_counter"},
			tgCommands.KeyBoardButtonTG{Text: "Remove this message", Data: "example_remove_buttons_trigger"},
		)))).ReplyMarkup
	return tgCommands.PreparedCommand(newMsg)
}

func (d *data) exampleRemoveButtons(msg *tgbotapi.Message, commandName string, param string, params []string) tgCommands.HandlerResult {
	fmt.Println("example_remove_buttons_trigger", msg.Chat.ID, msg.MessageID)
	newMsg := tgbotapi.NewDeleteMessage(msg.Chat.ID, msg.MessageID)
	return tgCommands.PreparedCommand(newMsg)
}

func (d *data) exampleNotify(msg *tgbotapi.Message, commandName string, param string, params []string) tgCommands.HandlerResult {
	newMsg := tgbotapi.NewChatAction(msg.Chat.ID, "typing")
	return tgCommands.PreparedCommand(newMsg)
}

func (d *data) exampleCounterIncrement(msg *tgbotapi.Message, commandName string, param string, params []string) tgCommands.HandlerResult {
	fmt.Println("exampleCounterIncrement", msg.Chat.ID, msg.MessageID)
	newMsg := tgbotapi.NewEditMessageText(msg.Chat.ID, msg.MessageID, "Text under buttons, click counter")
	newMsg.ReplyMarkup = tgbotapi.NewEditMessageReplyMarkup(
		msg.Chat.ID,
		msg.MessageID,
		tgCommands.GetTGButtons(tgCommands.KBRows(tgCommands.KBButs(
			tgCommands.KeyBoardButtonTG{Text: fmt.Sprintf("Counter (%v)", d.Counter()), Data: "example_button_counter"},
			tgCommands.KeyBoardButtonTG{Text: "Action test", Data: "example_notify"},
			tgCommands.KeyBoardButtonTG{Text: "Remove this message", Data: "example_remove_buttons_trigger"},
		)))).ReplyMarkup
	return tgCommands.PreparedCommand(newMsg)
}
