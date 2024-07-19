package users

import (
	"fmt"
	"strconv"

	tgModel "fun-coice/internal/domain/commands/tg"
	"fun-coice/pkg/scribble"
	tgbotapi "fun-coice/pkg/telegram-bot-api"
)

type data struct {
	list tgModel.Commands
	DB   *scribble.Driver
}

func New() tgModel.Service {
	result := data{}
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

func (d *data) Commands() tgModel.Commands {
	return d.list
}

func (d *data) Name() string {
	return "users"
}

func (d *data) Destroy() {}

func (d *data) Dependency() *tgModel.ServiceDepends {
	return tgModel.ServiceDependsIs(tgModel.FileDbDependency)
}

func (d *data) Configure(sc tgModel.ServiceConfig) error {
	if sc.FileDb == nil {
		return fmt.Errorf("file db is nil")
	}
	d.DB = sc.FileDb
	return nil
}

func (d *data) startBot(msg *tgbotapi.Message, command *tgModel.Command) *tgModel.HandlerResult {
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
