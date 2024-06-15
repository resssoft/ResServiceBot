package emojiTaskTracker

import (
	"fmt"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson/primitive"

	tgbotapi "fun-coice/pkg/telegram-bot-api"

	tgModel "fun-coice/internal/domain/commands/tg"
)

const description = `Bot for create and tracking tasks and controls by emoji
emoji controls look by command /emoji
`

//👌😱💯🔥👎❤️👍

var emojiControlsDefaultPause = "🤔"
var emojiControlsDefaultFinish = "💯"
var emojiControlsDefaultStarted = "👌"

func (d *data) emoji(uid int64, name string) string {
	switch name {
	case "pause":
		return emojiControlsDefaultPause
	case "finish":
		return emojiControlsDefaultFinish
	case "start":
		return emojiControlsDefaultStarted
	}
	return "no-emoji"
}

func (d *data) isStarted(text string, uid int64) bool {
	return strings.Contains(text, d.emoji(uid, "start"))
}

func (d *data) isPaused(text string, uid int64) bool {
	return strings.Contains(text, d.emoji(uid, "pause"))
}

func (d *data) isFinished(text string, uid int64) bool {
	return strings.Contains(text, d.emoji(uid, "finish"))
}

func (d *data) reactionEvent(msg *tgbotapi.Message, command *tgModel.Command) *tgModel.HandlerResult {
	//msg.Text hav emoji only !!!
	fmt.Println("====================================== reactionEvent")
	switch {
	case d.isPaused(msg.Text, msg.Chat.ID):
		foundedTask := d.search(msg.Chat.ID, msg.MessageID)
		log.Info().Interface("task", foundedTask).Send()
		if foundedTask == nil {
			fmt.Println("NOT FOUND", msg.Chat.ID, msg.MessageID)
			return tgModel.EmptyCommand()
		}
		d.save(msg.Chat.ID, foundedTask.SetPaused(), foundedTask.Code)
		log.Info().Interface("task", foundedTask).Send()
		log.Info().Msg(foundedTask.Format())
		return tgModel.SimpleEdit(msg.Chat.ID, msg.MessageID, foundedTask.Format())
	case d.isFinished(msg.Text, msg.Chat.ID):
		foundedTask := d.search(msg.Chat.ID, msg.MessageID)
		log.Info().Interface("task", foundedTask).Send()
		if foundedTask == nil {
			fmt.Println("NOT FOUND", msg.Chat.ID, msg.MessageID)
			return tgModel.EmptyCommand()
		}
		d.save(msg.Chat.ID, foundedTask.SetStopped(), foundedTask.Code)
		log.Info().Interface("task", foundedTask).Send()
		log.Info().Msg(foundedTask.Format())
		return tgModel.SimpleEdit(msg.Chat.ID, msg.MessageID, foundedTask.Format())
	case d.isStarted(msg.Text, msg.Chat.ID):
		foundedTask := d.search(msg.Chat.ID, msg.MessageID)
		log.Info().Interface("task", foundedTask).Send()
		if foundedTask == nil {
			fmt.Println("NOT FOUND", msg.Chat.ID, msg.MessageID)
			return tgModel.EmptyCommand()
		}
		d.save(msg.Chat.ID, foundedTask.SetStarted(), foundedTask.Code)
		log.Info().Msg(foundedTask.Format())
		return tgModel.SimpleEdit(msg.Chat.ID, msg.MessageID, foundedTask.Format())
	default:
		fmt.Println("====================================== default")
		foundedTask := d.search(msg.Chat.ID, msg.MessageID)
		if foundedTask != nil {
			log.Info().Interface("task", foundedTask).Send()
			return tgModel.SimpleEdit(msg.Chat.ID, msg.MessageID, foundedTask.Format())
		}
		return tgModel.EmptyCommand()
	}
	return tgModel.EmptyCommand()
}

func (d *data) NewTask(msg *tgbotapi.Message, command *tgModel.Command) *tgModel.HandlerResult {
	d.mutexTask.Lock()
	defer d.mutexTask.Unlock()
	fmt.Println("====================================== NewTask")
	code := fmt.Sprintf("%v_%v", msg.Chat.ID, msg.MessageID)
	if msg.ForwardFromMessageID != 0 {

	}
	newTask := Task{
		MongoID: primitive.NewObjectID(),
		Start:   time.Time{},
		End:     time.Time{},
		Break:   time.Time{},
		Title:   msg.Text,
		UserId:  msg.Chat.ID,
		MsgId:   0,
		Breaks:  nil,
		Status:  StatusCreated,
		BotName: command.BotName,
		Code:    code,
	}
	if _, ok := d.userData[msg.Chat.ID]; !ok {
		d.userData[msg.Chat.ID] = userData{
			tasks: make(map[string]*Task),
		}
	}
	d.userData[msg.Chat.ID].tasks[code] = &newTask
	d.tasksIdx[code] = msg.Chat.ID
	if msg.ForwardFromMessageID != 0 {

	}
	return tgModel.SimpleWIthCallback(msg.Chat.ID, newTask.Format(), code, d.callback).WithDelete(msg.Chat.ID, msg.MessageID)
}

func (d *data) activeTasks(msg *tgbotapi.Message, command *tgModel.Command) *tgModel.HandlerResult {
	return tgModel.EmptyCommand() // TODO: later
}

func (d *data) historyTasks(msg *tgbotapi.Message, command *tgModel.Command) *tgModel.HandlerResult {
	return tgModel.EmptyCommand() // TODO: later
}

func (d *data) emojiControls(msg *tgbotapi.Message, command *tgModel.Command) *tgModel.HandlerResult {
	return tgModel.Simple(msg.Chat.ID, fmt.Sprintf("emoji controls:%s%s%s%s%s%s",
		"\nstart:", d.emoji(msg.Chat.ID, "start"),
		"\npause:", d.emoji(msg.Chat.ID, "pause"),
		"\nfinish:", d.emoji(msg.Chat.ID, "finish"),
	))
}

func (d *data) help(msg *tgbotapi.Message, command *tgModel.Command) *tgModel.HandlerResult {
	return tgModel.Simple(msg.Chat.ID, description)
}
