package emojiTaskTracker

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hako/durafmt"
)

func (d *data) search(chatId int64, msgId int) *Task {
	d.mutexTask.Lock()
	defer d.mutexTask.Unlock()
	foundUserData, ok := d.userData[chatId]
	if !ok {
		return nil
	}
	val, ok := foundUserData.tasks[fmt.Sprintf("%v_%v", chatId, msgId)]
	if ok {
		return &val
	}
	return nil
}

func (d *data) searchByCode(code string) (*Task, int64) {
	d.mutexTask.Lock()
	defer d.mutexTask.Unlock()
	foundUid, ok := d.tasksIdx[code]
	if !ok {
		return nil, 0
	}
	foundUserData, ok := d.userData[foundUid]
	if !ok {
		return nil, foundUid
	}
	val, ok := foundUserData.tasks[code]
	if ok {
		return &val, foundUid
	}
	return nil, 0
}

func (d *data) save(chatId int64, task Task, code string) error {
	d.mutexTask.Lock()
	defer d.mutexTask.Unlock()
	foundUserData, ok := d.userData[chatId]
	if !ok {
		return errors.New("user tasks not found")
	}
	if code != task.Code {
		delete(foundUserData.tasks, code)
	}
	foundUserData.tasks[task.Code] = task
	d.tasksIdx[task.Code] = chatId
	return nil
}

func (t *Task) Format() string {
	taskTime := "-:-:-"
	if t.Status == StatusStarted || t.Status == StatusStopped {
		taskTime = Duration(time.Now().Sub(t.Start))
	}
	return fmt.Sprintf("%s[%s] %s\n",
		t.Status,
		taskTime,
		t.Title,
	)
}

func Duration(dt time.Duration) string {
	//TODO: translates
	formatted := ""
	if dt.Milliseconds() < time.Minute.Milliseconds() {
		formatted = fmt.Sprintf("%.0f секунд", dt.Seconds())
	} else {
		formatted = durafmt.Parse(dt).LimitFirstN(2).String()
	}
	formatted = strings.ReplaceAll(formatted, "milliseconds", "милимсекунд")
	formatted = strings.ReplaceAll(formatted, "seconds", "секунд")
	formatted = strings.ReplaceAll(formatted, "second", "секунда")
	formatted = strings.ReplaceAll(formatted, "minutes", "минут")
	formatted = strings.ReplaceAll(formatted, "minute", "минута")
	formatted = strings.ReplaceAll(formatted, "hours", "часов")
	formatted = strings.ReplaceAll(formatted, "hour", "час")
	return formatted
}
