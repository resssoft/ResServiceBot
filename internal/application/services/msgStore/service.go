package msgStore

import (
	"context"
	"encoding/json"
	"fmt"
	tgModel "fun-coice/internal/domain/commands/tg"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	zlog "github.com/rs/zerolog/log"
)

type data struct {
	list    tgModel.Commands
	msgRepo tgModel.MsgRepository
}

func New(msgRepo tgModel.MsgRepository) tgModel.Service {
	result := data{
		msgRepo: msgRepo,
	}
	commandsList := tgModel.NewCommands()
	commandsList["event:"+tgModel.TextMsgBotEvent] = tgModel.Command{
		CommandType: "event",
		Handler:     result.msgEvent,
	}
	commandsList["msgCount"] = tgModel.Command{
		Command:     "/msgCount",
		Permissions: tgModel.FreePerms,
		Handler:     result.msgCount,
	}
	commandsList["msgDb"] = tgModel.Command{
		Command:     "/msgDb",
		Permissions: tgModel.FreePerms,
		Handler:     result.msgGetDB,
	}
	commandsList["msgImport"] = tgModel.Command{
		Command:     "/msgImport",
		Permissions: tgModel.FreePerms,
		Handler:     result.msgImport,
	}

	//TODO count msg of chat id or name, type (private or supergroup), by user, LAST X messages from privates
	//TODO: IMPORT
	//TODO: do not save some chats

	result.list = commandsList
	return &result
}

func (d *data) Commands() tgModel.Commands {
	return d.list
}

func (d data) Name() string {
	return "msgStore"
}

func (d *data) msgEvent(msg *tgbotapi.Message, command *tgModel.Command) *tgModel.HandlerResult {
	//fmt.Println("msgEvent")
	msgJson, err := json.Marshal(msg)
	if err != nil {
		fmt.Println(err.Error()) // to errors log or errors service
	}
	err = d.msgRepo.Create(context.Background(), &tgModel.Message{
		TgMsg:        msg,
		BotName:      command.BotName,
		MsgType:      "text",
		MsgDirection: 0,
		MsgJson:      string(msgJson),
	})
	if err != nil {
		fmt.Println(err.Error()) // to errors log or errors service
	}
	return tgModel.EmptyCommand()
}

func (d *data) msgCount(msg *tgbotapi.Message, _ *tgModel.Command) *tgModel.HandlerResult {
	items, err := d.msgRepo.List(context.Background(), nil)
	if err != nil {
		fmt.Println(err.Error())
	}
	return tgModel.Simple(msg.Chat.ID, fmt.Sprintf("Messages: %v", len(items)))
}

func (d *data) msgGetDB(msg *tgbotapi.Message, _ *tgModel.Command) *tgModel.HandlerResult {
	return tgModel.Simple(msg.Chat.ID, "Not implemented")
}

func (d *data) msgImport(msg *tgbotapi.Message, command *tgModel.Command) *tgModel.HandlerResult {
	if len(command.Arguments.Parse()) == 0 {
		return tgModel.DeferredWithText(msg.Chat.ID, "Send file or json text", "msgImport", "", nil)
	}
	var result = MsgData{}
	err := json.Unmarshal([]byte(msg.Text), &result)
	if err != nil {
		fmt.Println(err.Error())
	}
	zlog.Info().Any("json", result).Send()
	return tgModel.Simple(msg.Chat.ID, "Not implemented")
}
