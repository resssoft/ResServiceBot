package tgbot

import (
	tgModel "fun-coice/internal/domain/commands/tg"
)

const commandChoicer = "event:commandChoicer"
const commandChoicerEvent = "commandChoicer"
const commandRedirect = "commandRedirect"

var defaultCommands = tgModel.Commands{
	"start": { //TODO: use add simple command
		Command:     "/start",
		Description: "start bot",
		Permissions: tgModel.FreePerms,
		Handler:     startDefault,
	},
	"myInfo": {
		Command:     "/myInfo",
		Synonyms:    []string{"info", "me"},
		Description: "Write GT user info",
		Permissions: tgModel.FreePerms,
		Handler:     myInfo,
	},
	"appInfo": {
		Command:     "/appInfo",
		Description: "Write app info",
		Permissions: tgModel.FreePerms,
		Handler:     appInfo,
	},
	"id": {
		Command:     "/id",
		Synonyms:    []string{"userId", "userid", "myid", "ид", "мой ид", "мой айди", "идентификатор"},
		Description: "Write user id",
		Permissions: tgModel.FreePerms,
		Handler:     userId,
	},
	"version": {
		Command:     "/version",
		Synonyms:    []string{"appVersion", "ver", "версия"},
		Description: "Write bot version",
		Permissions: tgModel.FreePerms,
		Handler:     appVersion,
	},
	commandChoicer: {
		IsEvent: true,
		Handler: setNextCommand,
	},
	commandRedirect: {
		Command: commandRedirect,
		Handler: setRedirectByCommand,
	},
}
