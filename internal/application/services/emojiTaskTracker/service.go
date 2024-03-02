package emojiTaskTracker

import (
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
	mutexTask deadlock.Mutex
}

func New() tgModel.Service {
	result := data{
		events:   allowedUpdates,
		name:     name,
		list:     tgModel.NewCommands(),
		userData: make(map[int64]userData),
	}
	//commandsList := tgModel.NewCommands()
	result.list.AddSimple("NewTask", "Added task with emoji control", result.NewTask)
	result.list.AddEvent(tgModel.MessageReactionEvent, result.reactionEvent)
	result.list.AddEvent(tgModel.MessageReactionEvent, result.reactionEvent)
	result.list["event:"+tgModel.MessageReactionEvent] = tgModel.Command{
		Command: "/event:" + tgModel.MessageReactionEvent,
		IsEvent: true,
		Handler: result.reactionEvent,
	}
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

func (d *data) Configure(_ tgModel.ServiceConfig) {}
