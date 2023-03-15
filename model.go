package main

import (
	"fmt"
	"fun-coice/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TGUser struct {
	UserID  int64
	ChatId  int64
	Login   string
	Name    string
	IsAdmin bool
}

type TGCommand struct {
	Command     string
	Synonyms    []string
	Triggers    []string
	Templates   []string
	Description string
	CommandType string
	Permissions TGCommandPermissions
	Handler     func(*tgbotapi.Message, string, string, []string) (tgbotapi.Chattable, bool)
}

func (t *TGCommand) isImplemented(msg, botName string) bool {
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

func (t *TGCommand) Permission(messageItem *tgbotapi.Message) bool {
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

func (tgp *TGCommandPermissions) Check(user *tgbotapi.User) bool {
	if tgp.UserPermissions == "all" {
		return true
	}
	if tgp.UserPermissions == "admin" && int64(user.ID) == config.TelegramAdminId() {
		return true
	}
	return false
}

type TGCommands map[string]TGCommand

type TGCommandPermissions struct {
	UserPermissions string
	ChatPermissions string
}

var freePerms = TGCommandPermissions{
	ChatPermissions: "all",
	UserPermissions: "all",
}

var adminPerms = TGCommandPermissions{
	ChatPermissions: "admin",
	UserPermissions: "admin",
}

type Configuration struct {
	HWCSURL  string
	Telegram TelegramConfig
}

type TelegramConfig struct {
	Bot        TgBot
	AdminId    string
	AdminLogin string
}

type TgBot struct {
	Token string
}

type KeyBoardTG struct {
	Rows []KeyBoardRowTG
}

type KeyBoardRowTG struct {
	Buttons []KeyBoardButtonTG
}

type KeyBoardButtonTG struct {
	Text string
	Data string
}

type SavedBlock struct {
	Group string
	User  string
	Text  string
}

type CheckList struct {
	Group  string
	ChatID int64
	Text   string
	Status bool
	Public bool
}

type ChatUser struct {
	ChatId      int64
	ChatName    string
	ContentType string
	CustomRole  string
	VoteCount   int
	User        TGUser
}
