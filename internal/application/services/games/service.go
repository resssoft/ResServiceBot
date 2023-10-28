package games

import (
	tgModel "fun-coice/internal/domain/commands/tg"
)

type data struct {
	list tgModel.Commands
}

func New() tgModel.Service {
	result := data{}
	commandsList := tgModel.NewCommands()
	commandsList["games"] = tgModel.Command{
		Command:     "/games",
		Description: "Games list",
		CommandType: "text",
		Permissions: tgModel.FreePerms,
		Handler:     result.games,
	}

	result.list = commandsList
	return &result
}

func (d *data) Commands() tgModel.Commands {
	return d.list
}
