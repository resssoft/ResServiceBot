package main

import tgCommands "fun-coice/internal/domain/commands/tg"

var commands = tgCommands.Commands{
	"start": {
		Command:     "/start",
		Description: "start bot",
		CommandType: "text",
		Permissions: tgCommands.FreePerms,
		Handler:     startDefault,
	},
	"myInfo": {
		Command:     "/myInfo",
		Synonyms:    []string{"info"},
		Description: "Write GT user info",
		CommandType: "text",
		Permissions: tgCommands.FreePerms,
		Handler:     myInfo,
	},
	"appInfo": {
		Command:     "/appInfo",
		Description: "Write app info",
		CommandType: "text",
		Permissions: tgCommands.FreePerms,
		Handler:     appInfo,
	},
	"id": {
		Command:     "/id",
		Synonyms:    []string{"userId", "userid", "myid", "ид", "мой ид", "мой айди", "идентификатор"},
		Description: "Write user id",
		CommandType: "text",
		Permissions: tgCommands.FreePerms,
		Handler:     userId,
	},
	"version": {
		Command:     "/version",
		Synonyms:    []string{"appVersion", "ver", "версия"},
		Description: "Write bot version",
		CommandType: "text",
		Permissions: tgCommands.FreePerms,
		Handler:     appVersion,
	},
}
