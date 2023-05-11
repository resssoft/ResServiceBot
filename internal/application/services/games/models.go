package games

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

var ChatUserList = make([]ChatUser, 1)

var gamesListKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("🧡 Lovely game", "lovelyGame"),
		tgbotapi.NewInlineKeyboardButtonURL("Rules", ""),
	),
)

type ChatUser struct {
	ChatId      int64
	ChatName    string
	ContentType string
	CustomRole  string
	VoteCount   int
	User        TGUser
}

type TGUser struct {
	UserID  int64
	ChatId  int64
	Login   string
	Name    string
	IsAdmin bool
}
