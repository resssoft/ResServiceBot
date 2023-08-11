package users

import (
	"fmt"
	tgCommands "fun-coice/internal/domain/commands/tg"
	"fun-coice/pkg/scribble"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
)

type data struct {
	list tgCommands.Commands
	DB   *scribble.Driver
}

func New(DB *scribble.Driver) tgCommands.Service {
	result := data{
		DB: DB,
	}
	commandsList := tgCommands.NewCommands()
	commandsList["start"] = tgCommands.Command{
		Command:     "/start",
		Description: "start info",
		CommandType: "text",
		Permissions: tgCommands.FreePerms,
		Handler:     result.startBot,
	}

	//TODO: ADDED save commands

	result.list = commandsList
	return &result
}

func (d data) Commands() tgCommands.Commands {
	return d.list
}

func (d data) startBot(msg *tgbotapi.Message, commandName string, param string, params []string) tgCommands.HandlerResult {
	user := tgCommands.User{
		UserID: msg.From.ID,
		ChatId: msg.Chat.ID,
		Login:  msg.From.UserName,
		Name:   msg.From.String(),
		//IsAdmin: isAdmin, // TODO: implement
	}
	if err := d.DB.Write("user", strconv.FormatInt(msg.From.ID, 10), user); err != nil {
		fmt.Println("add command error", err)
	}

	return tgCommands.Simple(msg.Chat.ID, "Hi "+msg.From.String()+" and welcome! See by /commands")
}
