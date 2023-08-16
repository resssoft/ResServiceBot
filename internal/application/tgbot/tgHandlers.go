package tgbot

import (
	"encoding/json"
	tgModel "fun-coice/internal/domain/commands/tg"
	"fun-coice/pkg/appStat"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
)

func myInfo(msg *tgbotapi.Message, command *tgModel.Command) tgModel.HandlerResult {
	if msg == nil {
		return tgModel.EmptyCommand()
	}
	return tgModel.Simple(msg.Chat.ID, tgModel.UserAndChatInfo(msg.From, msg.Chat))
}

func appInfo(msg *tgbotapi.Message, command *tgModel.Command) tgModel.HandlerResult {
	info := appStat.Info()
	infoJson, _ := json.MarshalIndent(info, "", "  ")
	return tgModel.Simple(msg.Chat.ID, string(infoJson))
}

func userId(msg *tgbotapi.Message, command *tgModel.Command) tgModel.HandlerResult {
	return tgModel.Simple(msg.Chat.ID, strconv.FormatInt(msg.Chat.ID, 10))
}

func appVersion(msg *tgbotapi.Message, command *tgModel.Command) tgModel.HandlerResult {
	return tgModel.Simple(msg.Chat.ID, appStat.Version)
}

func startDefault(msg *tgbotapi.Message, command *tgModel.Command) tgModel.HandlerResult {
	return tgModel.Simple(msg.Chat.ID, "Hi "+msg.From.String()+" and welcome")
}
