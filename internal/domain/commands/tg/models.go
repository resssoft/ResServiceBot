package tgCommands

import (
	"fmt"
	"fun-coice/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type User struct {
	UserID  int64
	ChatId  int64
	Login   string
	Name    string
	IsAdmin bool
}

type Command struct {
	Command     string
	Synonyms    []string
	Triggers    []string
	Templates   []string
	Description string
	CommandType string
	Permissions CommandPermissions
	Handler     func(*tgbotapi.Message, string, string, []string) (tgbotapi.Chattable, bool)
	//Bots        []string //multybots
}

func (t *Command) IsImplemented(msg, botName string) bool {
	if IsCommand(t.Command, msg, botName) {
		return true
	}
	for _, synonym := range t.Synonyms {
		if IsCommand(synonym, msg, botName) {
			return true
		}
	}
	return false
}

func (t *Command) Permission(messageItem *tgbotapi.Message) bool {
	if messageItem != nil {
		if messageItem == nil {
			return false
		}
		switch messageItem.Chat.Type {
		case "private":
			if t.Permissions.Check(messageItem.From) {
				return true
			}
		case "chat":
			if t.Permissions.Check(messageItem.From) {
				return true
			}
		}
	}
	return false
}

func IsCommand(command, msg, botName string) bool {
	return msg == command || msg == fmt.Sprintf("%s@%s", command, botName)
}

func (tgp *CommandPermissions) Check(user *tgbotapi.User) bool {
	if tgp.UserPermissions == "all" {
		return true
	}
	if tgp.UserPermissions == "admin" && int64(user.ID) == config.TelegramAdminId() {
		return true
	}
	return false
}

type Commands map[string]Command

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

func (cs Commands) Merge(list Commands) Commands {
	merged := make(Commands)
	for key, value := range cs {
		merged[key] = value
	}
	for key, value := range list {
		merged[key] = value
	}
	return merged
}
