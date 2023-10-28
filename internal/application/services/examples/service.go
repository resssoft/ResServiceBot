package examples

import (
	tgModel "fun-coice/internal/domain/commands/tg"
	"sync"
)

type data struct {
	list    tgModel.Commands
	counter int
	mutex   *sync.Mutex
}

func New() tgModel.Service {
	commandsList := tgModel.NewCommands()
	result := data{
		list:  commandsList,
		mutex: &sync.Mutex{},
	}
	commandsList["examples"] = tgModel.Command{
		Command:     "/examples",
		Description: "Show examples info",
		CommandType: "text",
		ListExclude: true, // do not show in the commands list
		Permissions: tgModel.FreePerms,
		Handler:     result.help,
	}
	commandsList["example_text"] = tgModel.Command{
		Command:     "/example_text",
		Description: "Show text",
		CommandType: "text",
		ListExclude: true, // do not show in the commands list
		Permissions: tgModel.FreePerms,
		Handler:     result.exampleText,
	}
	commandsList["example_reply"] = tgModel.Command{
		Command:     "/example_reply",
		Synonyms:    []string{"example_text_synonym1", "example_text_synonym2"},
		Description: "Show text with reply",
		CommandType: "text",
		ListExclude: true, // do not show in the commands list
		Permissions: tgModel.FreePerms,
		Handler:     result.exampleText,
	}
	commandsList["example_buttons"] = tgModel.Command{
		Command:     "/example_buttons",
		Description: "Show buttons",
		CommandType: "text",
		ListExclude: true, // do not show in the commands list
		Permissions: tgModel.FreePerms,
		Handler:     result.exampleShowInlineButtons,
	}
	commandsList["example_remove_buttons_trigger"] = tgModel.Command{
		Command:     "/example_remove_buttons_trigger",
		Description: "inner command for buttons event",
		CommandType: "text",
		ListExclude: true, // do not show in the commands list
		Permissions: tgModel.FreePerms,
		Handler:     result.exampleRemoveButtons,
	}
	commandsList["example_button_counter"] = tgModel.Command{
		Command:     "/example_button_counter",
		Description: "inner command for buttons event",
		CommandType: "text",
		ListExclude: true, // do not show in the commands list
		Permissions: tgModel.FreePerms,
		Handler:     result.exampleCounterIncrement,
	}
	commandsList["example_buttons_edit"] = tgModel.Command{
		Command:     "/example_buttons_edit",
		Description: "inner command for buttons event",
		CommandType: "text",
		ListExclude: true, // do not show in the commands list
		Permissions: tgModel.FreePerms,
		Handler:     result.exampleEditInlineButtons,
	}
	commandsList["example_notify"] = tgModel.Command{
		Command:     "/example_notify",
		Description: "user notify",
		CommandType: "event",
		ListExclude: true, // do not show in the commands list
		Permissions: tgModel.FreePerms,
		Handler:     result.exampleNotify,
	}

	result.list = commandsList
	return &result
}

func (d *data) Commands() tgModel.Commands {
	return d.list
}

func (d *data) Name() string {
	return "examples"
}

func (d *data) Configure(_ tgModel.ServiceConfig) {

}

func (d *data) Counter() int {
	//safe with goroutines counter
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.counter += 1
	return d.counter
}
