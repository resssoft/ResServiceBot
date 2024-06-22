package testManager

import (
	"github.com/rs/zerolog/log"
	"github.com/sasha-s/go-deadlock"

	tgModel "fun-coice/internal/domain/commands/tg"
)

type data struct {
	list   tgModel.Commands
	tests  map[string]*testParent
	active map[int64]string
	users  map[int64]*testParent
	mutex  deadlock.Mutex
}

func New() tgModel.Service {
	result := data{
		tests:  make(map[string]*testParent),
		active: make(map[int64]string),
		users:  make(map[int64]*testParent),
	}
	commandsList := tgModel.NewCommands()

	tgModel.FreeCommand().
		Simple("test_add", "Added new test", result.add).
		AsMenu().
		Push(commandsList)

	tgModel.FreeCommand().
		Simple("test_run_", "Run test", result.run).
		WithTemplates([]string{`^/test_run_.+$`}).
		Push(commandsList)

	tgModel.FreeCommand().
		Simple("test_done", "Done edit test", result.done).
		Push(commandsList)

	tgModel.FreeCommand().
		Simple("test_import", "Import test from text or file", result.importTest).
		Push(commandsList)

	tgModel.NewEvent("test_set_name", result.setName).Push(commandsList)
	tgModel.NewEvent("test_append", result.append).Push(commandsList)
	tgModel.NewEvent("test_question", result.question).Push(commandsList)
	tgModel.NewEvent("test_answer", result.answer).Push(commandsList)
	tgModel.NewEvent("test_import", result.importTestEvent).Push(commandsList)

	log.Info().Interface("commandsList_", commandsList).Send() //DEBUG
	result.list = commandsList
	return &result
}

func (d *data) Commands() tgModel.Commands {
	return d.list
}

func (d *data) Name() string {
	return "testManager"
}

func (d *data) Configure(_ tgModel.ServiceConfig) {}
