package tgModel

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

const (
	StartBotEvent        = "start" //triggered by /start command from the bot
	UserLeaveChantEvent  = "user_leave_chat"
	UserJoinedChantEvent = "user_joined_chat"
	TextMsgBotEvent      = "text_msg" //triggered by /start command from the bot
)

type ChatEvent string

type Redirect struct {
	CommandName string
	Message     *tgbotapi.Message
	Step        int
}
