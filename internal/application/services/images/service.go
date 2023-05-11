package images

import (
	tgCommands "fun-coice/internal/domain/commands/tg"
)

type data struct {
	list tgCommands.Commands
}

func New() tgCommands.Service {
	commandsList := make(tgCommands.Commands)
	result := data{
		list: commandsList,
	}
	commandsList["imageHelp"] = tgCommands.Command{
		Command:     "/imageHelp",
		Description: "image commands info",
		CommandType: "text",
		Permissions: tgCommands.FreePerms,
		Handler:     result.help,
	}
	commandsList["resize"] = tgCommands.Command{
		Command:     "/resize",
		Description: "resize image",
		CommandType: "text",
		Permissions: tgCommands.FreePerms,
		Handler:     result.resize,
	}
	commandsList["resizeImage"] = tgCommands.Command{
		Command:     "/resizeImage",
		Description: "resize image",
		CommandType: "text",
		ListExclude: true,
		Permissions: tgCommands.FreePerms,
		Handler:     result.resizeImage,
	}
	commandsList["rotate"] = tgCommands.Command{
		Command:     "/rotate",
		Description: "rotate image",
		CommandType: "text",
		Permissions: tgCommands.FreePerms,
		Handler:     result.rotate,
	}
	commandsList["rotateImage"] = tgCommands.Command{
		Command:     "/rotateImage",
		Description: "rotate image",
		CommandType: "text",
		ListExclude: true,
		Permissions: tgCommands.FreePerms,
		Handler:     result.rotateImage,
	}
	return &result
}

func (d data) Commands() tgCommands.Commands {
	return d.list
}
