package images

import (
	tgModel "fun-coice/internal/domain/commands/tg"
)

type data struct {
	list    tgModel.Commands
	botName string
}

func New() tgModel.Service {
	commandsList := tgModel.NewCommands()
	result := data{
		list: commandsList,
	}
	commandsList["imageHelp"] = tgModel.Command{
		Command:     "imageHelp",
		Description: "image commands info",
		CommandType: "text",
		Permissions: tgModel.FreePerms,
		Handler:     result.help,
	}
	commandsList["resize"] = tgModel.Command{
		Command:     "resize",
		Description: "resize image",
		CommandType: "text",
		Permissions: tgModel.FreePerms,
		Handler:     result.resize,
	}
	commandsList["resizeImage"] = tgModel.Command{
		Command:     "resizeImage",
		Description: "resize image",
		CommandType: "text",
		ListExclude: true,
		IsEvent:     true,
		Permissions: tgModel.FreePerms,
		Handler:     result.resizeImage,
		FileTypes:   tgModel.PhotoMediaTypes,
	}
	commandsList["rotate"] = tgModel.Command{
		Command:     "rotate",
		Description: "rotate image",
		CommandType: "text",
		Permissions: tgModel.FreePerms,
		Handler:     result.rotate,
	}
	commandsList["rotateImage"] = tgModel.Command{
		Command:     "rotateImage",
		Description: "rotate image",
		CommandType: "text",
		ListExclude: true,
		//IsEvent:     true,
		Permissions: tgModel.FreePerms,
		Handler:     result.rotateImage,
		FileTypes:   tgModel.PhotoMediaTypes,
	}
	return &result
}

func (d data) Commands() tgModel.Commands {
	return d.list
}

func (d data) Name() string {
	return "images"
}

// add ramka
// split to X parts
//with text
//whatermarks
//zoom in / zoom out
