package emojiTaskTracker

import (
	"fmt"
	"regexp"
)

func (d *data) search(chatId int64, text string) *Task {
	r, _ := regexp.Compile(`^task(\d+)\n`)
	parsedItems := r.FindStringSubmatch(text)
	fmt.Println(parsedItems)
	if len(parsedItems) != 2 {
		return nil
	}
	d.mutexTask.Lock()
	val, ok := d.userData[chatId].tasks[parsedItems[1]]
	defer d.mutexTask.Unlock()
	if ok {
		return &val
	}
	return nil
}

func (d *data) save(chatId int64, task Task) {
	d.mutexTask.Lock()
	d.userData[chatId].tasks[task.Code] = task
	defer d.mutexTask.Unlock()
}

func (t *Task) Format() string {
	return fmt.Sprintf("%s\n%s\nStatus: %s",
		t.Code,
		t.Title,
		t.Status,
	)
}
