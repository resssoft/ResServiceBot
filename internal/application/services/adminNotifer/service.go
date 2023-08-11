package adminNotifer

import (
	tgCommands "fun-coice/internal/domain/commands/tg"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type data struct {
	events  tgCommands.Commands
	adminId int64
}

var _ = (tgCommands.Service)(&data{})

func New(adminId int64) tgCommands.Service {
	result := data{
		adminId: adminId,
	}
	commandsList := tgCommands.NewCommands()
	commandsList["event:"+tgCommands.StartBotEvent.String()] = tgCommands.Command{
		CommandType: "event",
		Handler:     result.startEvent,
	}
	commandsList["event:"+tgCommands.UserLeaveChantEvent.String()] = tgCommands.Command{
		Command:     "/event:" + tgCommands.UserLeaveChantEvent.String(),
		CommandType: "event",
		Handler:     result.UserLeaveChantEvent,
	}
	commandsList["event:"+tgCommands.UserJoinedChantEvent.String()] = tgCommands.Command{
		CommandType: "event",
		Handler:     result.UserJoinedChantEvent,
	}

	result.events = commandsList
	return &result
}

func (d data) Commands() tgCommands.Commands {
	return d.events
}

func (d data) startEvent(msg *tgbotapi.Message, commandName string, param string, params []string) tgCommands.HandlerResult {
	return tgCommands.Simple(d.adminId, "New bot start:\n"+tgCommands.UserInfo(msg.From))
}

func (d data) UserLeaveChantEvent(msg *tgbotapi.Message, commandName string, param string, params []string) tgCommands.HandlerResult {
	if msg == nil {
		return tgCommands.EmptyCommand()
	}
	if msg.LeftChatMember == nil {
		return tgCommands.EmptyCommand()
	}
	var info string

	info += tgCommands.UserAndChatInfo(msg.LeftChatMember, msg.Chat)
	if msg.From.ID != msg.LeftChatMember.ID {
		info += "\nBy " + tgCommands.UserInfo(msg.From)
	}
	return tgCommands.Simple(d.adminId, "User Leave Chant:\n"+info)
}

func (d data) UserJoinedChantEvent(msg *tgbotapi.Message, commandName string, param string, params []string) tgCommands.HandlerResult {
	if msg == nil {
		return tgCommands.EmptyCommand()
	}
	if len(msg.NewChatMembers) == 0 {
		return tgCommands.EmptyCommand()
	}
	var info string
	for _, user := range msg.NewChatMembers {
		info += tgCommands.UserAndChatInfo(&user, msg.Chat)
		if msg.From.ID != user.ID {
			info += "\nBy " + tgCommands.UserInfo(msg.From)
		}
		info += "\n"
	}
	return tgCommands.Simple(d.adminId, "Users Joined Chant:\n"+info)
}
