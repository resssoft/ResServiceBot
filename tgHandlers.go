package main

import (
	"encoding/json"
	"fmt"
	"fun-coice/pkg/appStat"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
)

func myInfo(msg *tgbotapi.Message, commandName string, param string, params []string) (tgbotapi.Chattable, bool) {
	from := msg.From
	chat := msg.Chat

	userInfo := fmt.Sprintf("--== UserInfo==--\nID: %v\nUserName: %s\nFirstName: %s\nLastName: %s\nLanguageCode: %s"+
		"\n--==ChatInfo==--\nID: %v\nTitle: %s\nType: %s",
		from.ID,
		from.UserName,
		from.FirstName,
		from.LastName,
		from.LanguageCode,
		chat.ID,
		chat.Title,
		chat.Type,
	)
	return tgbotapi.NewMessage(chat.ID, userInfo), true
}

func appInfo(msg *tgbotapi.Message, commandName string, param string, params []string) (tgbotapi.Chattable, bool) {
	info := appStat.Info()
	infoJson, _ := json.MarshalIndent(info, "", "  ")
	return tgbotapi.NewMessage(msg.Chat.ID, string(infoJson)), true
}

func userId(msg *tgbotapi.Message, commandName string, param string, params []string) (tgbotapi.Chattable, bool) {
	return tgbotapi.NewMessage(msg.Chat.ID, strconv.FormatInt(msg.Chat.ID, 10)), true
}

func appVersion(msg *tgbotapi.Message, commandName string, param string, params []string) (tgbotapi.Chattable, bool) {
	return tgbotapi.NewMessage(msg.Chat.ID, appStat.Version), true
}

func startDefault(msg *tgbotapi.Message, commandName string, param string, params []string) (tgbotapi.Chattable, bool) {
	return tgbotapi.NewMessage(msg.Chat.ID, "Hi "+msg.From.String()+" and welcome"), true
}
