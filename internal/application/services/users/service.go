package users

import (
	"fmt"
	tgModel "fun-coice/internal/domain/commands/tg"
	"fun-coice/pkg/scribble"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
)

type data struct {
	list tgModel.Commands
	DB   *scribble.Driver
}

func New(DB *scribble.Driver) tgModel.Service {
	result := data{
		DB: DB,
	}
	commandsList := tgModel.NewCommands()
	commandsList["start"] = tgModel.Command{
		Command:     "/start",
		Description: "start info",
		CommandType: "text",
		Permissions: tgModel.FreePerms,
		Handler:     result.startBot,
	}

	//TODO: ADDED save commands

	result.list = commandsList
	return &result
}

func (d data) Commands() tgModel.Commands {
	return d.list
}

func (d data) startBot(msg *tgbotapi.Message, command *tgModel.Command) tgModel.HandlerResult {
	user := tgModel.User{
		UserID: msg.From.ID,
		ChatId: msg.Chat.ID,
		Login:  msg.From.UserName,
		Name:   msg.From.String(),
		//IsAdmin: isAdmin, // TODO: implement
	}
	if err := d.DB.Write("user", strconv.FormatInt(msg.From.ID, 10), user); err != nil {
		fmt.Println("add command error", err)
	}

	return tgModel.Simple(msg.Chat.ID, "Hi "+msg.From.String()+" and welcome! See by /commands")
}
