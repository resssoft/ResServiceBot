package emojier

import (
	"fmt"
	"log"
	"strings"

	tgModel "fun-coice/internal/domain/commands/tg"
	tgbotapi "fun-coice/pkg/telegram-bot-api"
)

type data struct {
	list      tgModel.Commands
	reactions map[string]botsData
}

type botsData struct {
	repeatReactions map[string]string
	textReactions   map[string]string
}

// 2024.07.19
var supportedEmoji = "👍👎❤️🔥🥰👏😁🤔🤯😱🤬😢🎉🤩🤮💩🙏👌🕊🤡🥱🥴😍🐳❤‍🔥🌚🌭💯🤣⚡️🍌🏆💔🤨😐🍓🍾💋🖕😈😴😭🤓👻👨‍💻👀🎃🙈😇😨🤝✍️🤗\U0001FAE1🎅🎄☃️💅🤪🗿🆒💘🙉🦄😘💊🙊😎👾🤷‍♂🤷🤷‍♀😡"

func New() tgModel.Service {
	serv := data{
		list:      tgModel.NewCommands(),
		reactions: make(map[string]botsData),
	}
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
		Simple("start", "Start bot", serv.description).
		Push(serv.list)
	tgModel.NewCommand().
		Simple("about", "About", serv.description, "help").
		Push(serv.list)
	tgModel.NewCommand().
		Simple("testReaction", "Try to set emoji \nExample: \n/testReaction 👌", serv.testReaction).
		Push(serv.list)
	tgModel.NewCommand().
		Simple("available", "Show available emoji", serv.available).
		Push(serv.list)

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

func (d *data) Configure(sc tgModel.ServiceConfig) error {
	d.reactions[sc.MessageSender.BotName()] = botsData{
		repeatReactions: make(map[string]string),
		textReactions:   make(map[string]string),
	}
	d.reactions[sc.MessageSender.BotName()].repeatReactions["💩"] = "💩"
	d.reactions[sc.MessageSender.BotName()].textReactions["говно"] = "💩"
	d.reactions[sc.MessageSender.BotName()].textReactions["гавно"] = "💩"
	d.reactions[sc.MessageSender.BotName()].textReactions["shit"] = "💩"
	d.reactions[sc.MessageSender.BotName()].textReactions["shit"] = "💩"
	d.reactions[sc.MessageSender.BotName()].textReactions["моряша"] = "❤️"
	d.reactions[sc.MessageSender.BotName()].textReactions["кот"] = "❤️"
	d.reactions[sc.MessageSender.BotName()].textReactions["котик"] = "❤️"
	return nil
}

func (d *data) description(msg *tgbotapi.Message, command *tgModel.Command) *tgModel.HandlerResult {
	log.Println("============ description", command.Command)
	return tgModel.SimpleReply(msg.Chat.ID,
		"Bot can repeat users emoji and set reactions by words triggers, commands list: /emojiCommands",
		msg.MessageID)
}

func (d *data) available(msg *tgbotapi.Message, command *tgModel.Command) *tgModel.HandlerResult {
	return tgModel.SimpleReply(msg.Chat.ID,
		"👍👎❤🔥🥰👏😁🤔🤯😱🤬😢🎉🤩🤮💩🙏👌🕊🤡🥱🥴😍🐳❤‍🔥🌚🌭💯🤣⚡🍌🏆💔🤨😐🍓🍾💋🖕😈😴😭🤓👻👨‍💻👀🎃🙈😇😨🤝✍🤗\U0001FAE1🎅🎄☃💅🤪🗿🆒💘🙉🦄😘💊🙊😎👾🤷‍♂🤷🤷‍♀😡",
		msg.MessageID)
}

func (d *data) reactionEvent(msg *tgbotapi.Message, command *tgModel.Command) *tgModel.HandlerResult {
	//👌😱💯🔥👎❤️👍💩
	for oldReaction, newReaction := range d.reactions[command.BotName].repeatReactions {
		if d.IsNewReaction(msg, oldReaction) {
			return tgModel.Reaction(msg.Chat.ID, msg.MessageID, newReaction)
		}
	}
	return tgModel.EmptyCommand()
}

func (d *data) textReactionEvent(msg *tgbotapi.Message, command *tgModel.Command) *tgModel.HandlerResult {
	for trigger, newReaction := range d.reactions[command.BotName].textReactions {
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
	if !d.IsSupport(strings.TrimSpace(items[0])) {
		return tgModel.SimpleReply(msg.Chat.ID, "Incorrect! emoji1 is not support, see support list: /available", msg.MessageID)
	}
	if !d.IsSupport(strings.TrimSpace(items[1])) {
		return tgModel.SimpleReply(msg.Chat.ID, "Incorrect! emoji2 is not support, see support list: /available", msg.MessageID)
	}
	d.reactions[command.BotName].repeatReactions[items[0]] = strings.TrimSpace(items[1])
	return tgModel.Reaction(msg.Chat.ID, msg.MessageID, strings.TrimSpace(items[1]))
}

func (d *data) delRepeatReaction(msg *tgbotapi.Message, command *tgModel.Command) *tgModel.HandlerResult {
	reaction := strings.TrimSpace(command.Arguments.Raw)
	_, ok := d.reactions[command.BotName].repeatReactions[reaction]
	if ok {
		delete(d.reactions[command.BotName].repeatReactions, reaction)
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
	if !d.IsSupport(strings.TrimSpace(items[1])) {
		return tgModel.SimpleReply(msg.Chat.ID, "Incorrect! emoji is not support, see support list: /available", msg.MessageID)
	}
	d.reactions[command.BotName].textReactions[items[0]] = strings.TrimSpace(items[1])
	return tgModel.Reaction(msg.Chat.ID, msg.MessageID, strings.TrimSpace(items[1]))
}

func (d *data) delTextReaction(msg *tgbotapi.Message, command *tgModel.Command) *tgModel.HandlerResult {
	reaction := strings.TrimSpace(command.Arguments.Raw)
	_, ok := d.reactions[command.BotName].textReactions[reaction]
	if ok {
		delete(d.reactions[command.BotName].textReactions, reaction)
		return tgModel.Reaction(msg.Chat.ID, msg.MessageID, "👌")
	}
	return tgModel.SimpleReply(msg.Chat.ID, "Not found!!", msg.MessageID)
}

func (d *data) TextReactions(msg *tgbotapi.Message, command *tgModel.Command) *tgModel.HandlerResult {
	result := ""
	for trigger, newReaction := range d.reactions[command.BotName].textReactions {
		result += fmt.Sprintf("[%s|=>|%s]", trigger, newReaction)
	}
	return tgModel.SimpleReply(msg.Chat.ID, result, msg.MessageID)
}

func (d *data) RepeatReactions(msg *tgbotapi.Message, command *tgModel.Command) *tgModel.HandlerResult {
	result := ""
	for trigger, newReaction := range d.reactions[command.BotName].repeatReactions {
		result += fmt.Sprintf("[%s|=>|%s]", trigger, newReaction)
	}
	return tgModel.SimpleReply(msg.Chat.ID, result, msg.MessageID)
}

func (d *data) testReaction(msg *tgbotapi.Message, command *tgModel.Command) *tgModel.HandlerResult {
	reaction := strings.TrimSpace(command.Arguments.Raw)
	return tgModel.Reaction(msg.Chat.ID, msg.MessageID, reaction)
}

func (d *data) IsNewReaction(msg *tgbotapi.Message, reaction string) bool {
	fmt.Println(msg.Text, reaction)
	return strings.Contains(msg.Text, reaction) && !strings.Contains(msg.Caption, reaction)
}

func (d *data) IsSupport(reaction string) bool {
	fmt.Println(fmt.Sprintf("IsSupport[%s]", reaction))
	return strings.Contains(supportedEmoji, reaction)
}
