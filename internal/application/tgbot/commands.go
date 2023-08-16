package tgbot

import (
	tgModel "fun-coice/internal/domain/commands/tg"
)

var defaultCommands = tgModel.Commands{
	"start": {
		Command:     "/start",
		Description: "start bot",
		CommandType: "text",
		Permissions: tgModel.FreePerms,
		Handler:     startDefault,
	},
	"myInfo": {
		Command:     "/myInfo",
		Synonyms:    []string{"info", "me"},
		Description: "Write GT user info",
		CommandType: "text",
		Permissions: tgModel.FreePerms,
		Handler:     myInfo,
	},
	"appInfo": {
		Command:     "/appInfo",
		Description: "Write app info",
		CommandType: "text",
		Permissions: tgModel.FreePerms,
		Handler:     appInfo,
	},
	"id": {
		Command:     "/id",
		Synonyms:    []string{"userId", "userid", "myid", "ид", "мой ид", "мой айди", "идентификатор"},
		Description: "Write user id",
		CommandType: "text",
		Permissions: tgModel.FreePerms,
		Handler:     userId,
	},
	"version": {
		Command:     "/version",
		Synonyms:    []string{"appVersion", "ver", "версия"},
		Description: "Write bot version",
		CommandType: "text",
		Permissions: tgModel.FreePerms,
		Handler:     appVersion,
	},
}
