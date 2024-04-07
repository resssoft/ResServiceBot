package tgModel

import tgbotapi "fun-coice/pkg/telegram-bot-api"

const (
	StartBotEvent        = "start" //triggered by /start command from the bot
	UserLeaveChantEvent  = "user_leave_chat"
	UserJoinedChantEvent = "user_joined_chat"
	TextMsgBotEvent      = "text_msg" //triggered by /start command from the bot
	MessageReactionEvent = "message_reaction"
)

type ChatEvent string

type Redirect struct {
	CommandName string
	Message     *tgbotapi.Message
	Step        int
}
