package games

import (
	tgCommands "fun-coice/internal/domain/commands/tg"
)

type data struct {
	list tgCommands.Commands
}

func New() tgCommands.Service {
	result := data{}
	commandsList := tgCommands.NewCommands()
	commandsList["games"] = tgCommands.Command{
		Command:     "/games",
		Description: "Games list",
		CommandType: "text",
		Permissions: tgCommands.FreePerms,
		Handler:     result.games,
	}

	result.list = commandsList
	return &result
}

func (d data) Commands() tgCommands.Commands {
	return d.list
}
