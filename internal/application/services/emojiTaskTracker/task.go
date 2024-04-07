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

func (d *data) save(chatId int64, task Task) error {
	d.mutexTask.Lock()
	defer d.mutexTask.Unlock()
	foundUserData, ok := d.userData[chatId]
	if !ok {
		return errors.New("user tasks not found")
	}
	foundUserData.tasks[task.Code] = task
	return nil
}

func (t *Task) Format() string {
	taskTime := "-:-:-"
	if t.Status == StatusStarted || t.Status == StatusStopped {
		taskTime = Duration(time.Now().Sub(t.Start))
	}
	return fmt.Sprintf("%s[%s]\n%s\n",
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
