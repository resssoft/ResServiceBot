package tgModel

var FreePerms = CommandPermissions{
	Chat:    true,
	Private: true,
}

var AdminPerms = CommandPermissions{
	AdminOnly: true,
}

var PrivatePerms = CommandPermissions{
	Private: true,
}

type CommandPermissions struct {
	Chat          bool
	Private       bool
	AdminOnly     bool
	CustomChat    []int64
	CustomPrivate []int64
}

//TODO: perms by bot
