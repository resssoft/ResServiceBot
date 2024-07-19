package emojier

import (
	"fmt"
	"strings"

	tgModel "fun-coice/internal/domain/commands/tg"
	tgbotapi "fun-coice/pkg/telegram-bot-api"
)

type data struct {
	list            tgModel.Commands
	repeatReactions map[string]string
	textReactions   map[string]string
}

func New() tgModel.Service {
	serv := data{
		list:            tgModel.NewCommands(),
		repeatReactions: make(map[string]string),
		textReactions:   make(map[string]string),
	}
	serv.repeatReactions["💩"] = "💩"
	serv.textReactions["говно"] = "💩"
	serv.textReactions["гавно"] = "💩"
	serv.textReactions["shit"] = "💩"
	serv.textReactions["shit"] = "💩"
	serv.textReactions["моряша"] = "❤️"
	serv.textReactions["кот"] = "❤️"
	serv.textReactions["котик"] = "❤️"
	tgModel.NewCommand().
		Simple("repeatEmoji", "Add repeat emoji \nExample: \n/repeatEmoji 👌:👌", serv.addRepeatReaction).
		Push(serv.list)
	tgModel.NewCommand().
		Simple("delRepeatEmoji", "Del repeat emoji \nExample: \n/delRepeatEmoji 👌", serv.delRepeatReaction).
		Push(serv.list)
	tgModel.NewCommand().
		Simple("addTextReaction", "Add text trigger \nExample: \n/addTextReaction love:❤️", serv.addTextReaction).
		Push(serv.list)
	tgModel.NewCommand().
		Simple("delTextReaction", "Del text trigger \nExample: \n/delTextReaction love", serv.delTextReaction).
		Push(serv.list)
	tgModel.NewCommand().
		Simple("textReactions", "Show emoji triggers", serv.TextReactions).
		Push(serv.list)
	tgModel.NewCommand().
		Simple("repeatReactions", "Show emoji repeats", serv.RepeatReactions).
		Push(serv.list)
	tgModel.NewCommand().
		Simple("emojiCommands", "Show emoji repeats", serv.RepeatReactions).
		Push(serv.list)
	tgModel.NewCommand().
		Simple("start", "Start bot", serv.description).
		PushSafety(serv.list)
	tgModel.NewCommand().
		Simple("about", "Show emoji repeats", serv.description, "help").
		PushSafety(serv.list)

	serv.list.AddEvent(tgModel.MessageReactionEvent, serv.reactionEvent)
	serv.list["event:"+tgModel.MessageReactionEvent] = tgModel.Command{
		Command: "/event:" + tgModel.MessageReactionEvent,
		IsEvent: true,
		Handler: serv.reactionEvent,
	}

	serv.list.AddEvent(tgModel.TextMsgBotEvent, serv.textReactionEvent)
	serv.list["event:"+tgModel.TextMsgBotEvent] = tgModel.Command{
		Command: "/event:" + tgModel.TextMsgBotEvent,
		IsEvent: true,
		Handler: serv.textReactionEvent,
	}

	return &serv
}

func (d *data) Commands() tgModel.Commands {
	return d.list
}

func (d *data) Name() string {
	return "emojier"
}

func (d *data) Destroy() {}

func (d *data) Dependency() *tgModel.ServiceDepends {
	return nil
}

func (d *data) Configure(_ tgModel.ServiceConfig) error {
	return nil
}

func (d *data) description(msg *tgbotapi.Message, command *tgModel.Command) *tgModel.HandlerResult {
	return tgModel.SimpleReply(msg.Chat.ID,
		"Bot can repeat users emoji and set reactions by words triggers, commands list: /emojiCommands",
		msg.MessageID)
}

func (d *data) reactionEvent(msg *tgbotapi.Message, command *tgModel.Command) *tgModel.HandlerResult {
	//👌😱💯🔥👎❤️👍💩
	for oldReaction, newReaction := range d.repeatReactions {
		if d.IsNewReaction(msg, oldReaction) {
			return tgModel.Reaction(msg.Chat.ID, msg.MessageID, newReaction)
		}
	}
	return tgModel.EmptyCommand()
}

func (d *data) textReactionEvent(msg *tgbotapi.Message, command *tgModel.Command) *tgModel.HandlerResult {
	for trigger, newReaction := range d.textReactions {
		if strings.Contains(strings.ToLower(msg.Text), trigger) {
			return tgModel.Reaction(msg.Chat.ID, msg.MessageID, newReaction)
		}
	}
	return tgModel.EmptyCommand()
}

func (d *data) addRepeatReaction(msg *tgbotapi.Message, command *tgModel.Command) *tgModel.HandlerResult {
	items := strings.Split(command.Arguments.Raw, ":")
	if len(items) < 2 {
		return tgModel.SimpleReply(msg.Chat.ID, "Incorrect! Format: 👌:👌", msg.MessageID)
	}
	if strings.TrimSpace(items[0]) == "" {
		return tgModel.SimpleReply(msg.Chat.ID, "Incorrect! emoji1 empty or premium", msg.MessageID)
	}
	if strings.TrimSpace(items[1]) == "" {
		return tgModel.SimpleReply(msg.Chat.ID, "Incorrect! emoji2 empty or premium", msg.MessageID)
	}
	d.repeatReactions[items[0]] = strings.TrimSpace(items[1])
	return tgModel.Reaction(msg.Chat.ID, msg.MessageID, strings.TrimSpace(items[1]))
}

func (d *data) delRepeatReaction(msg *tgbotapi.Message, command *tgModel.Command) *tgModel.HandlerResult {
	reaction := strings.TrimSpace(command.Arguments.Raw)
	_, ok := d.repeatReactions[reaction]
	if ok {
		delete(d.repeatReactions, reaction)
		return tgModel.Reaction(msg.Chat.ID, msg.MessageID, "👌")
	}
	return tgModel.SimpleReply(msg.Chat.ID, "Not found!!", msg.MessageID)
}

func (d *data) addTextReaction(msg *tgbotapi.Message, command *tgModel.Command) *tgModel.HandlerResult {
	items := strings.Split(command.Arguments.Raw, ":")
	if len(items) < 2 {
		return tgModel.SimpleReply(msg.Chat.ID, "Incorrect! Format: love:❤️", msg.MessageID)
	}
	if strings.TrimSpace(items[1]) == "" {
		return tgModel.SimpleReply(msg.Chat.ID, "Incorrect! emoji empty or premium", msg.MessageID)
	}
	d.textReactions[items[0]] = strings.TrimSpace(items[1])
	return tgModel.Reaction(msg.Chat.ID, msg.MessageID, strings.TrimSpace(items[1]))
}

func (d *data) delTextReaction(msg *tgbotapi.Message, command *tgModel.Command) *tgModel.HandlerResult {
	reaction := strings.TrimSpace(command.Arguments.Raw)
	_, ok := d.textReactions[reaction]
	if ok {
		delete(d.textReactions, reaction)
		return tgModel.Reaction(msg.Chat.ID, msg.MessageID, "👌")
	}
	return tgModel.SimpleReply(msg.Chat.ID, "Not found!!", msg.MessageID)
}

func (d *data) TextReactions(msg *tgbotapi.Message, command *tgModel.Command) *tgModel.HandlerResult {
	result := ""
	for trigger, newReaction := range d.textReactions {
		result += fmt.Sprintf("[%s|=>|%s]", trigger, newReaction)
	}
	return tgModel.SimpleReply(msg.Chat.ID, result, msg.MessageID)
}

func (d *data) RepeatReactions(msg *tgbotapi.Message, command *tgModel.Command) *tgModel.HandlerResult {
	result := ""
	for trigger, newReaction := range d.repeatReactions {
		result += fmt.Sprintf("[%s|=>|%s]", trigger, newReaction)
	}
	return tgModel.SimpleReply(msg.Chat.ID, result, msg.MessageID)
}

func (d *data) IsNewReaction(msg *tgbotapi.Message, reaction string) bool {
	fmt.Println(msg.Text, reaction)
	return strings.Contains(msg.Text, reaction) && !strings.Contains(msg.Caption, reaction)
}
