package tgModel

import tgbotapi "fun-coice/pkg/telegram-bot-api"

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

var PrivatePerms = CommandPermissions{
	ChatPermissions: "",
	UserPermissions: "all",
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

//TODO: perms by bot
