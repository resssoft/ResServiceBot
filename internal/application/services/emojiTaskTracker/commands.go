package emojiTaskTracker

import (
	"fmt"
	"strings"
	"time"

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
		if foundedTask != nil {
			foundedTask.Status = StatusPause
			d.save(msg.Chat.ID, *foundedTask, foundedTask.Code)
			return tgModel.SimpleEdit(msg.Chat.ID, msg.MessageID, foundedTask.Format())
		}
		fmt.Println("NOT FOUND", msg.Chat.ID, msg.MessageID)
	case d.isFinished(msg.Text, msg.Chat.ID):
		foundedTask := d.search(msg.Chat.ID, msg.MessageID)
		if foundedTask != nil {
			foundedTask.Status = StatusStopped
			d.save(msg.Chat.ID, *foundedTask, foundedTask.Code)
			return tgModel.SimpleEdit(msg.Chat.ID, msg.MessageID, foundedTask.Format())
		}
		fmt.Println("NOT FOUND", msg.Chat.ID, msg.MessageID)
	case d.isStarted(msg.Text, msg.Chat.ID):
		foundedTask := d.search(msg.Chat.ID, msg.MessageID)
		if foundedTask != nil {
			if foundedTask.Status != StatusStarted {
				foundedTask.Status = StatusStarted
				foundedTask.Start = time.Now()
				d.save(msg.Chat.ID, *foundedTask, foundedTask.Code)
			}
			return tgModel.SimpleEdit(msg.Chat.ID, msg.MessageID, foundedTask.Format())
		}
		fmt.Println("NOT FOUND", msg.Chat.ID, msg.MessageID)
	default:
		foundedTask := d.search(msg.Chat.ID, msg.MessageID)
		if foundedTask != nil {
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
			tasks: make(map[string]Task),
		}
	}
	d.userData[msg.Chat.ID].tasks[code] = newTask
	d.tasksIdx[code] = msg.Chat.ID
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
