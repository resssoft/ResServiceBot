package emojiTaskTracker

import (
	"fmt"

	"github.com/sasha-s/go-deadlock"

	"fun-coice/internal/domain/commands/tg"
)

const name = "emojiTaskTracker"

var allowedUpdates = []string{"message"}

type data struct {
	list      tgModel.Commands
	events    []string
	name      string
	userData  map[int64]userData
	tasks     map[int64]Task
	tasksIdx  map[string]int64
	mutexTask deadlock.Mutex
	callback  chan tgModel.CallbackData `json:"-"`
}

func New() tgModel.Service {
	result := data{
		events:   allowedUpdates,
		name:     name,
		list:     tgModel.NewCommands(),
		userData: make(map[int64]userData),
		tasksIdx: make(map[string]int64),
		callback: make(chan tgModel.CallbackData),
	}
	//commandsList := tgModel.NewCommands()
	result.list.AddSimple("NewTask", "Added task with emoji control", result.NewTask)
	//result.list.AddSimple("tasks", "Show active tasks", result.activeTasks)
	//result.list.AddSimple("history", "Show tasks history", result.historyTasks)
	result.list.AddSimple("emoji", "Show control emoji", result.emojiControls)
	result.list.AddSimple("help", "Show bot info", result.help)
	result.list.AddSimple("description", "Show bot info", result.help)
	result.list.AddSimple("about", "Show bot info", result.help)
	result.list.AddEvent(tgModel.MessageReactionEvent, result.reactionEvent)
	result.list["event:"+tgModel.MessageReactionEvent] = tgModel.Command{
		Command: "/event:" + tgModel.MessageReactionEvent,
		IsEvent: true,
		Handler: result.reactionEvent,
	}
	go result.CallbackHandler()
	//result.list = commandsList
	return &result
}

func (d *data) Commands() tgModel.Commands {
	return d.list
}

func (d *data) Name() string {
	return d.name
}

func (d *data) Events() []string {
	return d.events
}

func (d *data) Destroy() {}

func (d *data) Dependency() *tgModel.ServiceDepends {
	return nil
}

func (d *data) Configure(_ tgModel.ServiceConfig) error {
	return nil
}

func (d *data) CallbackHandler() {
	for callbackItem := range d.callback {
		fmt.Println("CALLBACK", callbackItem.Tag, callbackItem.Value)
		fmt.Println("CALLBACK", d.userData)
		task, chatId := d.searchByCode(callbackItem.Tag)
		if task != nil {
			task.Code = fmt.Sprintf("%v_%v", chatId, callbackItem.Value)
			task.MsgId = callbackItem.Value
			d.save(chatId, task, callbackItem.Tag)
		}
		fmt.Println("CALLBACK", d.userData)
	}
}

//TODO: save to db tasks
//TODO: timer for task time update to actial duration in tg msg
//TODO: changed user emoji
//TODO: add simple examples for first start
//TODO: tasks history command
//TODO: tasks controls - edit, delete
