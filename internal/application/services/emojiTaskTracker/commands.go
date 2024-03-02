package emojiTaskTracker

import (
	"fmt"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"

	tgModel "fun-coice/internal/domain/commands/tg"
)

//👌😱💯🔥👎❤️👍

func (d *data) reactionEvent(msg *tgbotapi.Message, command *tgModel.Command) *tgModel.HandlerResult {
	//msg.Text hav emoji only !!!
	fmt.Println("====================================== reactionEvent")
	switch {
	case strings.Contains(msg.Text, "🤔"): //pause
		foundedTask := d.search(msg.Chat.ID, msg.Text)
		if foundedTask != nil {
			foundedTask.Status = StatusPause
			d.save(msg.Chat.ID, *foundedTask)
		}
		return tgModel.SimpleEdit(msg.Chat.ID, msg.MessageID, foundedTask.Format())
	case strings.Contains(msg.Text, "💯"):
		foundedTask := d.search(msg.Chat.ID, msg.Text)
		if foundedTask != nil {
			foundedTask.Status = StatusStopped
			d.save(msg.Chat.ID, *foundedTask)
		}
		return tgModel.SimpleEdit(msg.Chat.ID, msg.MessageID, foundedTask.Format())
	case strings.Contains(msg.Text, "👌"): //start
		foundedTask := d.search(msg.Chat.ID, msg.Text)
		if foundedTask != nil {
			foundedTask.Status = StatusProgress
			d.save(msg.Chat.ID, *foundedTask)
		}
		fmt.Println("NOT FOUND", msg.Chat.ID, msg.Text)
		return tgModel.SimpleEdit(msg.Chat.ID, msg.MessageID, foundedTask.Format())
	default:
		return tgModel.Simple(msg.Chat.ID, "test")
	}
}

func (d *data) NewTask(msg *tgbotapi.Message, command *tgModel.Command) *tgModel.HandlerResult {
	fmt.Println("====================================== NewTask")
	code := fmt.Sprintf("task%v", msg.MessageID)
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
	d.mutexTask.Lock()
	if _, ok := d.userData[msg.Chat.ID]; !ok {
		d.userData[msg.Chat.ID] = userData{
			tasks: make(map[string]Task),
		}
	}
	d.userData[msg.Chat.ID].tasks[code] = newTask
	defer d.mutexTask.Unlock()
	return tgModel.Simple(msg.Chat.ID, newTask.Format()).WithDelete(msg.Chat.ID, msg.MessageID)
}
