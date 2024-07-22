package emojier

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"fun-coice/internal/application/services/emodjier/repository/model"
	repository "fun-coice/internal/application/services/emodjier/repository/mongo"
	tgModel "fun-coice/internal/domain/commands/tg"
	tgbotapi "fun-coice/pkg/telegram-bot-api"
)

const defaultRandomEmojiPercent = 2

type data struct {
	list      tgModel.Commands
	reactions map[string]model.BotData
	repo      model.Repository
	randEmoji map[string]int
}

//TODO: translate

// 2024.07.19
var supportedEmoji = "👍👎❤️🔥🥰👏😁🤔🤯😱🤬😢🎉🤩🤮💩🙏👌🕊🤡🥱🥴😍🐳❤‍🔥🌚🌭💯🤣⚡️🍌🏆💔🤨😐🍓🍾💋🖕😈😴😭🤓👻👨‍💻👀🎃🙈😇😨🤝✍️🤗\U0001FAE1🎅🎄☃️💅🤪🗿🆒💘🙉🦄😘💊🙊😎👾🤷‍♂🤷🤷‍♀😡"

var supportedEmojiList = []string{
	"👍", "👎", "❤️", "🔥", "🥰", "👏", "😁", "🤔", "🤯", "😱", "🤬", "😢", "🎉", "🤩", "🤮", "💩", "🙏", "👌", "🕊", "🤡", "🥱", "🥴",
	"😍", "🐳", "❤", "‍🔥", "🌚", "🌭", "💯", "🤣", "⚡️", "🍌", "🏆", "💔", "🤨", "😐", "🍓", "🍾", "💋", "🖕", "😈", "😴", "😭",
	"🤓", "👻", "👨‍💻", "👀", "🎃", "🙈", "😇", "😨", "🤝", "✍️", "🤗", "\U0001FAE1", "🎅", "🎄", "☃️", "💅", "🤪", "🗿", "🆒", "💘",
	"🙉", "🦄", "😘", "💊", "🙊", "😎", "👾", "🤷", "‍♂🤷", "🤷‍♀😡",
}

func New() tgModel.Service {
	serv := data{
		list:      tgModel.NewCommands(),
		reactions: make(map[string]model.BotData),
		randEmoji: make(map[string]int),
	}
	tgModel.NewCommand().
		Simple("repeatEmoji", "Add repeat emoji Example: /repeatEmoji 👌:👌", serv.addRepeatReaction).
		Push(serv.list)
	tgModel.NewCommand().
		Simple("delRepeatEmoji", "Del repeat emoji Example: /delRepeatEmoji 👌", serv.delRepeatReaction).
		Push(serv.list)
	tgModel.NewCommand().
		Simple("addTextReaction", "Add text trigger Example: /addTextReaction love:❤️", serv.addTextReaction).
		Push(serv.list)
	tgModel.NewCommand().
		Simple("delTextReaction", "Del text trigger Example: /delTextReaction love", serv.delTextReaction).
		Push(serv.list)
	tgModel.NewCommand().
		Simple("textReactions", "Show emoji triggers", serv.TextReactions).
		Push(serv.list)
	tgModel.NewCommand().
		Simple("repeatReactions", "Show emoji repeats", serv.RepeatReactions).
		Push(serv.list)
	tgModel.NewCommand().
		Simple("start", "Start bot", serv.description).WithPerm(tgModel.PrivatePerms).
		Push(serv.list)
	tgModel.NewCommand().
		Simple("about", "About", serv.description, "help").
		Push(serv.list)
	tgModel.NewCommand().
		Simple("testReaction", "Try to set emoji Example: /testReaction 👌", serv.testReaction).
		Push(serv.list)
	tgModel.NewCommand().
		Simple("available", "Show available emoji", serv.available).
		Push(serv.list)
	tgModel.AdminCommand().
		Simple("emojierSave", "save emoji to db", serv.save).
		Push(serv.list)
	tgModel.AdminCommand().
		Simple("emojierClone", "Clone emoji from other bot", serv.clone).
		Push(serv.list)
	tgModel.NewCommand().
		Simple("emojiCommands", "emoji commands", serv.commandsList).
		Push(serv.list)
	tgModel.NewCommand().
		Simple("emojiSet", "Set emodji to repost message Example: /emojiSet 👌", serv.setReaction).
		Push(serv.list)
	tgModel.AdminCommand().
		Simple("emojiPercentRandom", "Set random emoji persent", serv.setPercent).
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
	return tgModel.ServiceDependsIs(tgModel.MongoDbDependency)
}

func (d *data) Configure(sc tgModel.ServiceConfig) error {
	var err error
	d.repo, err = repository.NewRepo(sc.MongoClient)
	if err != nil {
		return err
	}
	if d.repo != nil {
		botdata, err := d.repo.List(context.Background(), sc.MessageSender.BotName())
		if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
			log.Println("botdata load data err", err)
		} else {
			d.reactions[sc.MessageSender.BotName()] = botdata
		}
	}
	d.randEmoji[sc.MessageSender.BotName()] = defaultRandomEmojiPercent
	if len(d.reactions[sc.MessageSender.BotName()].TextReactions) == 0 &&
		len(d.reactions[sc.MessageSender.BotName()].TextReactions) == 0 {
		d.reactions[sc.MessageSender.BotName()] = model.BotData{
			BotName:         sc.MessageSender.BotName(),
			RepeatReactions: make(map[string]string),
			TextReactions:   make(map[string]string),
		}
		d.reactions[sc.MessageSender.BotName()].RepeatReactions["💩"] = "💩"
		d.reactions[sc.MessageSender.BotName()].TextReactions["говно"] = "💩"
		d.reactions[sc.MessageSender.BotName()].TextReactions["гавно"] = "💩"
		d.reactions[sc.MessageSender.BotName()].TextReactions["shit"] = "💩"
		d.reactions[sc.MessageSender.BotName()].TextReactions["shit"] = "💩"
		d.reactions[sc.MessageSender.BotName()].TextReactions["моряша"] = "❤️"
		d.reactions[sc.MessageSender.BotName()].TextReactions["кот"] = "❤️"
		d.reactions[sc.MessageSender.BotName()].TextReactions["котик"] = "❤️"
	}
	return nil
}

func (d *data) description(msg *tgbotapi.Message, command *tgModel.Command) *tgModel.HandlerResult {
	log.Println("============ description", command.Command)
	return tgModel.SimpleReply(msg.Chat.ID,
		"Bot can repeat users emoji and set reactions by words triggers, commands list: /emojiCommands",
		msg.MessageID)
}

func (d *data) commandsList(msg *tgbotapi.Message, command *tgModel.Command) *tgModel.HandlerResult {
	list := d.Commands().Available(msg, command.Bot.AdminId)
	commandsList := "Commands:\n"
	for key, item := range list {
		commandsList += "/" + key + " - " + item.Description + "\n"
	}
	return tgModel.Simple(msg.Chat.ID, commandsList)
}

func (d *data) save(msg *tgbotapi.Message, command *tgModel.Command) *tgModel.HandlerResult {
	err := d.addOrUpdate(command.Bot.Login, d.reactions[command.Bot.Login])
	if err != nil {
		return tgModel.SimpleReply(msg.Chat.ID, "Save err: "+err.Error(), msg.MessageID)
	}
	return tgModel.Reaction(msg.Chat.ID, msg.MessageID, "👌")
}

func (d *data) clone(msg *tgbotapi.Message, command *tgModel.Command) *tgModel.HandlerResult {
	oldBotName := command.Arguments.Raw
	dataForClone, err := d.repo.List(context.Background(), oldBotName)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return tgModel.SimpleReply(msg.Chat.ID, "Not found by"+oldBotName+": "+err.Error(), msg.MessageID)
	} else {
		dataForClone.BotName = command.Bot.Login
		d.reactions[command.Bot.Login] = dataForClone
		err := d.addOrUpdate(command.Bot.Login, d.reactions[command.Bot.Login])
		if err != nil {
			return tgModel.SimpleReply(msg.Chat.ID, "Save err: "+err.Error(), msg.MessageID)
		}
	}
	return tgModel.Reaction(msg.Chat.ID, msg.MessageID, "👌")
}

func (d *data) addOrUpdate(botName string, botData model.BotData) error {
	_, err := d.repo.List(context.Background(), botName)
	if errors.Is(err, mongo.ErrNoDocuments) {
		_, err = d.repo.Create(context.Background(), botData)
		if err != nil {
			log.Println("botdata Create data err", err)
			return err
		}
	} else {
		err = d.repo.Update(context.Background(), botData)
		if err != nil {
			log.Println("botdata Update data err", err)
			return err
		}
	}
	return nil
}

func (d *data) available(msg *tgbotapi.Message, command *tgModel.Command) *tgModel.HandlerResult {
	return tgModel.SimpleReply(msg.Chat.ID,
		"👍👎❤🔥🥰👏😁🤔🤯😱🤬😢🎉🤩🤮💩🙏👌🕊🤡🥱🥴😍🐳❤‍🔥🌚🌭💯🤣⚡🍌🏆💔🤨😐🍓🍾💋🖕😈😴😭🤓👻👨‍💻👀🎃🙈😇😨🤝✍🤗\U0001FAE1🎅🎄☃💅🤪🗿🆒💘🙉🦄😘💊🙊😎👾🤷‍♂🤷🤷‍♀😡",
		msg.MessageID)
}

func (d *data) reactionEvent(msg *tgbotapi.Message, command *tgModel.Command) *tgModel.HandlerResult {
	//👌😱💯🔥👎❤️👍💩
	for oldReaction, newReaction := range d.reactions[command.Bot.Login].RepeatReactions {
		if d.IsNewReaction(msg, oldReaction) {
			return tgModel.Reaction(msg.Chat.ID, msg.MessageID, newReaction)
		}
	}
	return tgModel.EmptyCommand()
}

func (d *data) textReactionEvent(msg *tgbotapi.Message, command *tgModel.Command) *tgModel.HandlerResult {
	if msg.Text != "" {
		for trigger, newReaction := range d.reactions[command.Bot.Login].TextReactions {
			if strings.Contains(strings.ToLower(msg.Text), trigger) {
				return tgModel.Reaction(msg.Chat.ID, msg.MessageID, newReaction)
			}
		}
	}
	if msg.Caption != "" {
		for trigger, newReaction := range d.reactions[command.Bot.Login].TextReactions {
			if strings.Contains(strings.ToLower(msg.Caption), trigger) {
				return tgModel.Reaction(msg.Chat.ID, msg.MessageID, newReaction)
			}
		}
	}
	if d.Random(command.Bot.Login) && d.randEmoji[command.Bot.Login] > 0 {
		return tgModel.Reaction(msg.Chat.ID, msg.MessageID, d.RandomEmpji())
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
	d.reactions[command.Bot.Login].RepeatReactions[items[0]] = strings.TrimSpace(items[1])
	_ = d.addOrUpdate(command.Bot.Login, d.reactions[command.Bot.Login])
	return tgModel.Reaction(msg.Chat.ID, msg.MessageID, strings.TrimSpace(items[1]))
}

func (d *data) delRepeatReaction(msg *tgbotapi.Message, command *tgModel.Command) *tgModel.HandlerResult {
	reaction := strings.TrimSpace(command.Arguments.Raw)
	_, ok := d.reactions[command.Bot.Login].RepeatReactions[reaction]
	if ok {
		delete(d.reactions[command.Bot.Login].RepeatReactions, reaction)
		return tgModel.Reaction(msg.Chat.ID, msg.MessageID, "👌")
	}
	_ = d.addOrUpdate(command.Bot.Login, d.reactions[command.Bot.Login])
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
	d.reactions[command.Bot.Login].TextReactions[items[0]] = strings.TrimSpace(items[1])
	_ = d.addOrUpdate(command.Bot.Login, d.reactions[command.Bot.Login])
	return tgModel.Reaction(msg.Chat.ID, msg.MessageID, strings.TrimSpace(items[1]))
}

func (d *data) delTextReaction(msg *tgbotapi.Message, command *tgModel.Command) *tgModel.HandlerResult {
	reaction := strings.TrimSpace(command.Arguments.Raw)
	_, ok := d.reactions[command.Bot.Login].TextReactions[reaction]
	if ok {
		delete(d.reactions[command.Bot.Login].TextReactions, reaction)
		return tgModel.Reaction(msg.Chat.ID, msg.MessageID, "👌")
	}
	_ = d.addOrUpdate(command.Bot.Login, d.reactions[command.Bot.Login])
	return tgModel.SimpleReply(msg.Chat.ID, "Not found!!", msg.MessageID)
}

func (d *data) setReaction(msg *tgbotapi.Message, command *tgModel.Command) *tgModel.HandlerResult {
	reaction := strings.TrimSpace(command.Arguments.Raw)
	if msg.ReplyToMessage == nil {
		return tgModel.SimpleReply(msg.Chat.ID, "Use reply, please", msg.MessageID)
	}
	return tgModel.Reaction(msg.Chat.ID, msg.ReplyToMessage.MessageID, reaction)
}

func (d *data) TextReactions(msg *tgbotapi.Message, command *tgModel.Command) *tgModel.HandlerResult {
	result := ""
	for trigger, newReaction := range d.reactions[command.Bot.Login].TextReactions {
		result += fmt.Sprintf("[%s|=>|%s] ", trigger, newReaction)
	}
	return tgModel.SimpleReply(msg.Chat.ID, result, msg.MessageID)
}

func (d *data) RepeatReactions(msg *tgbotapi.Message, command *tgModel.Command) *tgModel.HandlerResult {
	result := ""
	for trigger, newReaction := range d.reactions[command.BotName].RepeatReactions {
		result += fmt.Sprintf("[%s|=>|%s] ", trigger, newReaction)
	}
	return tgModel.SimpleReply(msg.Chat.ID, result, msg.MessageID)
}

func (d *data) testReaction(msg *tgbotapi.Message, command *tgModel.Command) *tgModel.HandlerResult {
	reaction := strings.TrimSpace(command.Arguments.Raw)
	return tgModel.Reaction(msg.Chat.ID, msg.MessageID, reaction)
}

func (d *data) IsNewReaction(msg *tgbotapi.Message, reaction string) bool {
	fmt.Println(msg.Text, reaction)
	eventReaction := msg.Text
	if eventReaction == "❤" {
		eventReaction = "❤️" // Telegram WTF?
	}
	return strings.Contains(eventReaction, reaction) && !strings.Contains(msg.Caption, reaction)
}

func (d *data) IsSupport(reaction string) bool {
	fmt.Println(fmt.Sprintf("IsSupport[%s]", reaction))
	return strings.Contains(supportedEmoji, reaction)
}

func (d *data) setPercent(msg *tgbotapi.Message, command *tgModel.Command) *tgModel.HandlerResult {
	percent := strings.TrimSpace(command.Arguments.Raw)
	percentInt, _ := strconv.Atoi(percent)
	if percentInt > 100 || percentInt < 0 {
		return tgModel.SimpleReply(msg.Chat.ID, "U invalid, set value of 0..100", msg.MessageID)
	}
	d.randEmoji[command.Bot.Login] = percentInt
	return tgModel.Reaction(msg.Chat.ID, msg.MessageID, "👌")
}

func (d *data) Random(param string) bool {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	return r.Intn(100) <= d.randEmoji[param]
}

func (d *data) RandomEmpji() string {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	randNumber := r.Intn(len(supportedEmojiList) - 1)
	return supportedEmojiList[randNumber] //fmt.Sprintf("__%d[%s]", randNumber, supportedEmojiList[randNumber])
}
