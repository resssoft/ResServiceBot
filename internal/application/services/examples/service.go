package examples

import (
	tgCommands "fun-coice/internal/domain/commands/tg"
	"sync"
)

type data struct {
	list    tgCommands.Commands
	counter int
	mutex   *sync.Mutex
}

func New() tgCommands.Service {
	commandsList := make(tgCommands.Commands)
	result := data{
		list:  commandsList,
		mutex: &sync.Mutex{},
	}
	commandsList["examples"] = tgCommands.Command{
		Command:     "/examples",
		Description: "Show examples info",
		CommandType: "text",
		ListExclude: true, // do not show in the commands list
		Permissions: tgCommands.FreePerms,
		Handler:     result.help,
	}
	commandsList["example_text"] = tgCommands.Command{
		Command:     "/example_text",
		Description: "Show text",
		CommandType: "text",
		ListExclude: true, // do not show in the commands list
		Permissions: tgCommands.FreePerms,
		Handler:     result.exampleText,
	}
	commandsList["example_reply"] = tgCommands.Command{
		Command:     "/example_reply",
		Synonyms:    []string{"example_text_synonym1", "example_text_synonym2"},
		Description: "Show text with reply",
		CommandType: "text",
		ListExclude: true, // do not show in the commands list
		Permissions: tgCommands.FreePerms,
		Handler:     result.exampleText,
	}
	commandsList["example_buttons"] = tgCommands.Command{
		Command:     "/example_buttons",
		Description: "Show buttons",
		CommandType: "text",
		ListExclude: true, // do not show in the commands list
		Permissions: tgCommands.FreePerms,
		Handler:     result.exampleShowInlineButtons,
	}
	commandsList["example_remove_buttons_trigger"] = tgCommands.Command{
		Command:     "/example_remove_buttons_trigger",
		Description: "inner command for buttons event",
		CommandType: "text",
		ListExclude: true, // do not show in the commands list
		Permissions: tgCommands.FreePerms,
		Handler:     result.exampleRemoveButtons,
	}
	commandsList["example_button_counter"] = tgCommands.Command{
		Command:     "/example_button_counter",
		Description: "inner command for buttons event",
		CommandType: "text",
		ListExclude: true, // do not show in the commands list
		Permissions: tgCommands.FreePerms,
		Handler:     result.exampleCounterIncrement,
	}
	commandsList["example_buttons_edit"] = tgCommands.Command{
		Command:     "/example_buttons_edit",
		Description: "inner command for buttons event",
		CommandType: "text",
		ListExclude: true, // do not show in the commands list
		Permissions: tgCommands.FreePerms,
		Handler:     result.exampleEditInlineButtons,
	}
	commandsList["example_notify"] = tgCommands.Command{
		Command:     "/example_notify",
		Description: "user notify",
		CommandType: "event",
		ListExclude: true, // do not show in the commands list
		Permissions: tgCommands.FreePerms,
		Handler:     result.exampleNotify,
	}

	result.list = commandsList
	return &result
}

func (d *data) Commands() tgCommands.Commands {
	return d.list
}

func (d *data) Counter() int {
	//safe with goroutines counter
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.counter += 1
	return d.counter
}
