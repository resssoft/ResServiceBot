package tgModel

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type User struct {
	UserID  int64
	ChatId  int64
	Login   string
	Name    string
	IsAdmin bool
}

type TGUser struct {
	UserID  int64
	ChatId  int64
	Login   string
	Name    string
	IsAdmin bool
}

func UserAndChatInfo(user *tgbotapi.User, chat *tgbotapi.Chat) string {
	return UserInfo(user) + ChatInfo(chat)
}

func UserInfo(user *tgbotapi.User) string {
	if user == nil {
		return ""
	}
	userLogin := ""
	if user.UserName != "" {
		userLogin += fmt.Sprintf("(@%s)", user.UserName)
	}
	if user.LanguageCode != "" {
		userLogin += fmt.Sprintf("(%s)", user.LanguageCode)
	}
	userInfo := fmt.Sprintf("User: [%v] %s %s %s",
		user.ID,
		userLogin,
		user.FirstName,
		user.LastName,
	)
	return userInfo
}

func ChatInfo(chat *tgbotapi.Chat) string {
	if chat != nil {
		return fmt.Sprintf("\nChat: [%v] (%s): %s",
			chat.ID,
			chat.Type,
			chat.Title,
		)
	}
	return ""
}
