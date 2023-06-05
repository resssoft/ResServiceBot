package tgbot

import (
	"encoding/json"
	tgCommands "fun-coice/internal/domain/commands/tg"
	"fun-coice/pkg/appStat"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
)

func myInfo(msg *tgbotapi.Message, commandName string, param string, params []string) tgCommands.HandlerResult {
	if msg == nil {
		return tgCommands.EmptyCommand()
	}
	return tgCommands.Simple(msg.Chat.ID, tgCommands.UserAndChatInfo(msg.From, msg.Chat))
}

func appInfo(msg *tgbotapi.Message, commandName string, param string, params []string) tgCommands.HandlerResult {
	info := appStat.Info()
	infoJson, _ := json.MarshalIndent(info, "", "  ")
	return tgCommands.Simple(msg.Chat.ID, string(infoJson))
}

func userId(msg *tgbotapi.Message, commandName string, param string, params []string) tgCommands.HandlerResult {
	return tgCommands.Simple(msg.Chat.ID, strconv.FormatInt(msg.Chat.ID, 10))
}

func appVersion(msg *tgbotapi.Message, commandName string, param string, params []string) tgCommands.HandlerResult {
	return tgCommands.Simple(msg.Chat.ID, appStat.Version)
}

func startDefault(msg *tgbotapi.Message, commandName string, param string, params []string) tgCommands.HandlerResult {
	return tgCommands.Simple(msg.Chat.ID, "Hi "+msg.From.String()+" and welcome").
		WithEvent(tgCommands.NewEvent(tgCommands.StartBotEvent, msg))
}
