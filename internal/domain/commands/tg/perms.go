package tgModel

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type CommandPermissions struct {
	UserPermissions string
	ChatPermissions string
}

var FreePerms = CommandPermissions{
	ChatPermissions: "all",
	UserPermissions: "all",
}

var AdminPerms = CommandPermissions{
	ChatPermissions: "admin",
	UserPermissions: "admin",
}

var ModerPerms = CommandPermissions{
	ChatPermissions: "moder",
	UserPermissions: "moder",
}

func (tgp *CommandPermissions) Check(user *tgbotapi.User, adminId int64) bool {
	if tgp.UserPermissions == "all" {
		return true
	}
	if tgp.UserPermissions == "admin" && user.ID == adminId {
		return true
	}
	return false
}
