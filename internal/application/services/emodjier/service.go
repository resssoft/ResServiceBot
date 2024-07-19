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
}

func New() tgModel.Service {
	serv := data{
		list:            tgModel.NewCommands(),
		repeatReactions: make(map[string]string),
	}
	serv.repeatReactions["💩"] = "💩"
	tgModel.NewCommand().
		Simple("repeatEmoji", "Add repeat emoji", serv.addRepeatReaction).
		Push(serv.list)
	tgModel.NewCommand().
		Simple("delRepeatEmoji", "Del repeat emoji", serv.delRepeatReaction).
		Push(serv.list)

	serv.list.AddEvent(tgModel.MessageReactionEvent, serv.reactionEvent)
	serv.list["event:"+tgModel.MessageReactionEvent] = tgModel.Command{
		Command: "/event:" + tgModel.MessageReactionEvent,
		IsEvent: true,
		Handler: serv.reactionEvent,
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

func (d *data) encode(msg *tgbotapi.Message, command *tgModel.Command) *tgModel.HandlerResult {
	return tgModel.SimpleReply(msg.Chat.ID, "b64result", msg.MessageID)
}

func (d *data) reactionEvent(msg *tgbotapi.Message, command *tgModel.Command) *tgModel.HandlerResult {
	//👌😱💯🔥👎❤️👍💩
	for oldReaction, newReaction := range d.repeatReactions {
		switch {
		case d.IsNewReaction(msg, oldReaction):
			return tgModel.Reaction(msg.Chat.ID, msg.MessageID, newReaction)
		}
	}
	return tgModel.EmptyCommand()
}

func (d *data) addRepeatReaction(msg *tgbotapi.Message, command *tgModel.Command) *tgModel.HandlerResult {
	items := strings.Split(command.Arguments.Raw, ":")
	if len(items) < 2 {
		return tgModel.SimpleReply(msg.Chat.ID, "Incorrect! Format: Samara:Europe/Samara", msg.MessageID)
	}
	d.repeatReactions[items[0]] = items[1]
	return tgModel.Reaction(msg.Chat.ID, msg.MessageID, "👌")
}

func (d *data) delRepeatReaction(msg *tgbotapi.Message, command *tgModel.Command) *tgModel.HandlerResult {
	_, ok := d.repeatReactions[command.Arguments.Raw]
	if ok {
		delete(d.repeatReactions, command.Arguments.Raw)
		return tgModel.Reaction(msg.Chat.ID, msg.MessageID, "👌")
	}
	return tgModel.SimpleReply(msg.Chat.ID, "Not found!!", msg.MessageID)
}

func (d *data) IsNewReaction(msg *tgbotapi.Message, reaction string) bool {
	fmt.Println(msg.Text, reaction)
	return strings.Contains(msg.Text, reaction) && !strings.Contains(msg.Caption, reaction)
}
