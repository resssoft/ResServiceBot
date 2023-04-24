package main

import tgCommands "fun-coice/internal/domain/commands/tg"

var commands = tgCommands.Commands{
	"calc": {
		Command:     "/calc",
		Synonyms:    []string{"calc", "калк"},
		Triggers:    []string{"calc", "калк", "сколько будет"},
		Templates:   []string{`\d+\s*.\s*\d+`},
		Description: "(параметры - строка для продсчета данных, пример 2+2 или (2.5 - 1.35) * 2.0",
		CommandType: "text",
		Permissions: tgCommands.FreePerms,
	},
	"addCheckItem": {
		Command:     "/addCheckItem",
		Description: "(параметры - имя чеклиста, =1 - если публичный, =1 если уже установлен) - создание элемента чеклиста в указанную группу",
		CommandType: "text",
		Permissions: tgCommands.FreePerms,
	},
	"updateCheckItem": {
		Command:     "/updateCheckItem",
		Description: "(параметр - имя чеклиста, =1 или =0 для статуса, полный текст элемента для обновления) - вывод указанной группы чеклиста",
		CommandType: "text",
		Permissions: tgCommands.FreePerms,
	},
	"сheckList": {
		Command:     "/сheckList",
		Description: "(параметр - имя чеклиста) - вывод указанной группы чеклиста",
		CommandType: "text",
		Permissions: tgCommands.FreePerms,
	},
	"start": {
		Command:     "/start",
		Description: "Service registration, only private",
		CommandType: "text",
		Permissions: tgCommands.FreePerms,
	},
	"myInfo": {
		Command:     "/myInfo",
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
		Description: "Write user id",
		CommandType: "text",
		Permissions: tgCommands.FreePerms,
		Handler:     userId,
	},
	"getUserList": {
		Command:     "/getUserList",
		Description: "-",
		CommandType: "text",
		Permissions: tgCommands.AdminPerms,
	},
	"addSaveCommand": {
		Command:     "/addSaveCommand",
		Description: "Создать комманду сохранения коротких текстовых сообщений, чтобы потом ею сохранять текстовые строки. например. '/addSaveCommand whatToDo' и потом 'whatToDo вымыть посуду'",
		CommandType: "text",
		Permissions: tgCommands.ModerPerms,
	},
	"addFeature": {
		Command:     "/addFeature",
		Description: "Создание описание фичи",
		CommandType: "text",
		Permissions: tgCommands.FreePerms,
	},
	"getFeatures": {
		Command:     "/getFeatures",
		Description: "Список фич приложения",
		CommandType: "text",
		Permissions: tgCommands.FreePerms,
	},
	"SaveCommandsList": {
		Command:     "/SaveCommandsList",
		Description: "Список комманд для сохранения текстовых строк",
		CommandType: "text",
		Permissions: tgCommands.FreePerms,
	},
	"listOf": {
		Command:     "/listOf",
		Description: "(+ аргумент) Список сохраненных сообщений по указанной комманде",
		CommandType: "text",
		Permissions: tgCommands.FreePerms,
	},
	"version": {
		Command:     "/version",
		Synonyms:    []string{"appVersion", "ver", "версия"},
		Description: "Write bot version",
		CommandType: "text",
		Permissions: tgCommands.FreePerms,
		Handler:     appVersion,
	},
	"commands": {
		Command:     "/commands",
		Description: "Список комманд",
		CommandType: "text",
		Permissions: tgCommands.FreePerms,
	},
	"rebuild": {
		Command:     "/rebuild",
		Description: "rebuild",
		CommandType: "text",
		Permissions: tgCommands.AdminPerms,
	},
	"games": {
		Command:     "/games",
		Description: "games list",
		CommandType: "text",
		Permissions: tgCommands.FreePerms,
	},
}
