package tgModel

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"regexp"
)

type Command struct {
	Command       string //TODO: check for needles field
	Synonyms      []string
	Triggers      []string
	Templates     []string
	Description   string
	CommandType   string //deprecated
	IsEvent       bool
	ListExclude   bool
	Permissions   CommandPermissions
	Handler       HandlerFunc
	Arguments     CommandArguments
	Service       string // set in the bot only
	FileTypes     FileTypes
	BotName       string //command author
	Deferred      bool   // send by Deferred method
	FilesCallback FileHandlerFunc
	//State       string //offline or online, service can be down
	//WithFiles   bool // Files need prepare
}

func NewCommand() *Command {
	return &Command{}
}

func FreeCommand() *Command {
	return &Command{
		Permissions: FreePerms,
	}
}

func AdminCommand() *Command {
	return &Command{
		Permissions: FreePerms,
	}
}

func (t *Command) Simple(
	name, description string,
	handler HandlerFunc,
	synonyms ...string) *Command {
	return &Command{
		Command:     "/" + name,
		Description: description,
		Permissions: FreePerms,
		Handler:     handler,
		Synonyms:    synonyms,
	}
}

func (t *Command) WithPerm(perm CommandPermissions) *Command {
	t.Permissions = perm
	return t
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

func (t *Command) IsMatched(msg, botName string) bool {
	if len(t.Templates) > 0 {
		for _, template := range t.Templates {
			templateMatched, _ := regexp.MatchString(template, msg)
			if templateMatched {
				return true
			}
		}
	}
	return false
}

func (t *Command) WithHandler(handler HandlerFunc) *Command {
	t.Handler = handler
	return t
}

func (t *Command) Permission(messageItem *tgbotapi.Message, adminId int64) bool {
	if messageItem != nil {
		if messageItem == nil {
			return false
		}
		switch messageItem.Chat.Type {
		case "private":
			if t.Permissions.Check(messageItem.From, adminId) {
				return true
			}
		case "chat":
			if t.Permissions.Check(messageItem.From, adminId) {
				return true
			}
		}
	}
	return false
}

func (t *Command) SetArgs(args string) *Command {
	t.Arguments = CommandArguments{
		Raw: args,
	}
	return t
}

func (t *Command) ParsedArgs() []string {
	return t.Arguments.Parse()
}

func IsCommand(command, msg, botName string) bool {
	return msg == ("/"+command) || msg == fmt.Sprintf("%s@%s", command, botName)
}
