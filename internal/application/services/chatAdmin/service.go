package chatAdmin

import (
	tgModel "fun-coice/internal/domain/commands/tg"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type data struct {
	events  tgModel.Commands
	adminId int64
}

//TODO: use messages service and other DB

var _ = (tgModel.Service)(&data{})

func New(adminId int64) tgModel.Service {
	result := data{
		adminId: adminId,
	}
	commandsList := tgModel.NewCommands()

	commandsList["getChatLink"] = tgModel.Command{
		CommandType: "/getChatLink",

		Handler: result.ChatLink,
	}
	commandsList["event:"+tgModel.StartBotEvent] = tgModel.Command{
		CommandType: "event",
		Handler:     result.startEvent,
	}
	commandsList["event:"+tgModel.UserLeaveChantEvent] = tgModel.Command{
		Command:     "/event:" + tgModel.UserLeaveChantEvent,
		CommandType: "event",
		Handler:     result.UserLeaveChantEvent,
	}
	commandsList["event:"+tgModel.UserJoinedChantEvent] = tgModel.Command{
		CommandType: "event",
		Handler:     result.UserJoinedChantEvent,
	}
	//TODO: configure notify chat (admin is default)

	result.events = commandsList
	return &result
}

func (d data) Commands() tgModel.Commands {
	return d.events
}

func (d data) Name() string {
	return "chatAdmin"
}

func (d data) startEvent(msg *tgbotapi.Message, command *tgModel.Command) *tgModel.HandlerResult {
	return tgModel.Simple(d.adminId, "New bot start:\n"+tgModel.UserInfo(msg.From))
}

func (d data) UserLeaveChantEvent(msg *tgbotapi.Message, _ *tgModel.Command) *tgModel.HandlerResult {
	if msg == nil {
		return tgModel.EmptyCommand()
	}
	if msg.LeftChatMember == nil {
		return tgModel.EmptyCommand()
	}
	var info string

	info += tgModel.UserAndChatInfo(msg.LeftChatMember, msg.Chat)
	if msg.From.ID != msg.LeftChatMember.ID {
		info += "\nBy " + tgModel.UserInfo(msg.From)
	}
	return tgModel.Simple(d.adminId, "User Leave Chat:\n"+info)
}

func (d data) UserJoinedChantEvent(msg *tgbotapi.Message, _ *tgModel.Command) *tgModel.HandlerResult {
	if msg == nil {
		return tgModel.EmptyCommand()
	}
	if len(msg.NewChatMembers) == 0 {
		return tgModel.EmptyCommand()
	}
	var info string
	for _, user := range msg.NewChatMembers {
		info += tgModel.UserAndChatInfo(&user, msg.Chat)
		if msg.From.ID != user.ID {
			info += "\nBy " + tgModel.UserInfo(msg.From)
		}
		info += "\n"
	}
	return tgModel.Simple(d.adminId, "Users Joined Chant:\n"+info)
}

func (d data) ChatLink(msg *tgbotapi.Message, _ *tgModel.Command) *tgModel.HandlerResult {
	return tgModel.Simple(d.adminId, "Users Joined Chant:\n") ////////////////////////////////////////////
}
