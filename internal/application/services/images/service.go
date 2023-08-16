package images

import (
	tgModel "fun-coice/internal/domain/commands/tg"
)

type data struct {
	list    tgModel.Commands
	botName string
}

func New(botName string) tgModel.Service {
	commandsList := tgModel.NewCommands()
	result := data{
		list:    commandsList,
		botName: botName,
	}
	commandsList["imageHelp"] = tgModel.Command{
		Command:     "/imageHelp",
		Description: "image commands info",
		CommandType: "text",
		Permissions: tgModel.FreePerms,
		Handler:     result.help,
	}
	commandsList["resize"] = tgModel.Command{
		Command:     "/resize",
		Description: "resize image",
		CommandType: "text",
		Permissions: tgModel.FreePerms,
		Handler:     result.resize,
	}
	commandsList["resizeImage"] = tgModel.Command{
		Command:     "/resizeImage",
		Description: "resize image",
		CommandType: "text",
		ListExclude: true,
		Permissions: tgModel.FreePerms,
		Handler:     result.resizeImage,
	}
	commandsList["rotate"] = tgModel.Command{
		Command:     "/rotate",
		Description: "rotate image",
		CommandType: "text",
		Permissions: tgModel.FreePerms,
		Handler:     result.rotate,
	}
	commandsList["rotateImage"] = tgModel.Command{
		Command:     "/rotateImage",
		Description: "rotate image",
		CommandType: "text",
		ListExclude: true,
		Permissions: tgModel.FreePerms,
		Handler:     result.rotateImage,
	}
	return &result
}

func (d data) Commands() tgModel.Commands {
	return d.list
}
