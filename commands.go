package main

var commands = TGCommands{
	"calc": {
		Command:     "/calc",
		Synonyms:    []string{"calc", "калк"},
		Triggers:    []string{"calc", "калк", "сколько будет"},
		Templates:   []string{`\d+\s*.\s*\d+`},
		Description: "(параметры - строка для продсчета данных, пример 2+2 или (2.5 - 1.35) * 2.0",
		CommandType: "text",
		Permissions: TGCommandPermissions{
			ChatPermissions: "all",
			UserPermissions: "all",
		},
	},
	"addCheckItem": {
		Command:     "/addCheckItem",
		Description: "(параметры - имя чеклиста, =1 - если публичный, =1 если уже установлен) - создание элемента чеклиста в указанную группу",
		CommandType: "text",
		Permissions: TGCommandPermissions{
			ChatPermissions: "all",
			UserPermissions: "all",
		},
	},
	"updateCheckItem": {
		Command:     "/updateCheckItem",
		Description: "(параметр - имя чеклиста, =1 или =0 для статуса, полный текст элемента для обновления) - вывод указанной группы чеклиста",
		CommandType: "text",
		Permissions: TGCommandPermissions{
			ChatPermissions: "all",
			UserPermissions: "all",
		},
	},
	"сheckList": {
		Command:     "/сheckList",
		Description: "(параметр - имя чеклиста) - вывод указанной группы чеклиста",
		CommandType: "text",
		Permissions: freePerms,
	},
	"start": {
		Command:     "/start",
		Description: "Service registration, only private",
		CommandType: "text",
		Permissions: freePerms,
	},
	"myInfo": {
		Command:     "/myInfo",
		Description: "Write GT user info",
		CommandType: "text",
		Permissions: freePerms,
		Handler:     myInfo,
	},
	"appInfo": {
		Command:     "/appInfo",
		Description: "Write app info",
		CommandType: "text",
		Permissions: adminPerms,
		Handler:     appInfo,
	},
	"member": {
		Command:     "/member",
		Description: "Write GT user info and member status",
		CommandType: "text",
		Permissions: freePerms,
	},
	"getUserList": {
		Command:     "/getUserList",
		Description: "-",
		CommandType: "text",
		Permissions: TGCommandPermissions{
			ChatPermissions: "admin",
			UserPermissions: "admin",
		},
	},
	"addSaveCommand": {
		Command:     "/addSaveCommand",
		Description: "Создать комманду сохранения коротких текстовых сообщений, чтобы потом ею сохранять текстовые строки. например. '/addSaveCommand whatToDo' и потом 'whatToDo вымыть посуду'",
		CommandType: "text",
		Permissions: TGCommandPermissions{
			ChatPermissions: "moder",
			UserPermissions: "moder",
		},
	},
	"addFeature": {
		Command:     "/addFeature",
		Description: "Создание описание фичи",
		CommandType: "text",
		Permissions: TGCommandPermissions{
			ChatPermissions: "all",
			UserPermissions: "all",
		},
	},
	"getFeatures": {
		Command:     "/getFeatures",
		Description: "Список фич приложения",
		CommandType: "text",
		Permissions: TGCommandPermissions{
			ChatPermissions: "all",
			UserPermissions: "all",
		},
	},
	"cat": {
		Command:     "/cat",
		Description: "Какой ты кот",
		CommandType: "text",
		Permissions: TGCommandPermissions{
			ChatPermissions: "all",
			UserPermissions: "all",
		},
	},
	"SaveCommandsList": {
		Command:     "/SaveCommandsList",
		Description: "Список комманд для сохранения текстовых строк",
		CommandType: "text",
		Permissions: TGCommandPermissions{
			ChatPermissions: "all",
			UserPermissions: "all",
		},
	},
	"listOf": {
		Command:     "/listOf",
		Description: "(+ аргумент) Список сохраненных сообщений по указанной комманде",
		CommandType: "text",
		Permissions: TGCommandPermissions{
			ChatPermissions: "all",
			UserPermissions: "all",
		},
	},
	"admin": {
		Command:     "/admin",
		Description: "Вывод логина админа",
		CommandType: "text",
		Permissions: TGCommandPermissions{
			ChatPermissions: "all",
			UserPermissions: "all",
		},
	},
	"version": {
		Command:     "/version",
		Description: "Вывод версии",
		CommandType: "text",
		Permissions: TGCommandPermissions{
			ChatPermissions: "all",
			UserPermissions: "all",
		},
	},
	"appVersion": {
		Command:     "/appVersion",
		Description: "синоним version",
		CommandType: "text",
		Permissions: TGCommandPermissions{
			ChatPermissions: "all",
			UserPermissions: "all",
		},
	},
	"версия": {
		Command:     "/версия",
		Description: "синоним version",
		CommandType: "text",
		Permissions: TGCommandPermissions{
			ChatPermissions: "all",
			UserPermissions: "all",
		},
	},
	"commands": {
		Command:     "/commands",
		Description: "Список комманд",
		CommandType: "text",
		Permissions: TGCommandPermissions{
			ChatPermissions: "all",
			UserPermissions: "all",
		},
	},
	"rebuild": {
		Command:     "/rebuild",
		Description: "rebuild",
		CommandType: "text",
		Permissions: TGCommandPermissions{
			ChatPermissions: "admin",
			UserPermissions: "admin",
		},
	},
	"homeweb": {
		Command:     "/homeweb",
		Description: "get image link from cam1",
		CommandType: "text",
		Permissions: TGCommandPermissions{
			ChatPermissions: "admin",
			UserPermissions: "admin",
		},
	},
	"games": {
		Command:     "/games",
		Description: "games list",
		CommandType: "text",
		Permissions: TGCommandPermissions{
			ChatPermissions: "all",
			UserPermissions: "all",
		},
	},
}
