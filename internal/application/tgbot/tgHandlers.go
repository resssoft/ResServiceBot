package tgbot

import (
	"encoding/json"
	tgModel "fun-coice/internal/domain/commands/tg"
	"fun-coice/pkg/appStat"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	zlog "github.com/rs/zerolog/log"
	"strconv"
)

func myInfo(msg *tgbotapi.Message, _ *tgModel.Command) *tgModel.HandlerResult {
	if msg == nil {
		return tgModel.EmptyCommand()
	}
	return tgModel.Simple(msg.Chat.ID, tgModel.UserAndChatInfo(msg.From, msg.Chat))
}

func appInfo(msg *tgbotapi.Message, _ *tgModel.Command) *tgModel.HandlerResult {
	info := appStat.Info()
	infoJson, _ := json.MarshalIndent(info, "", "  ")
	return tgModel.Simple(msg.Chat.ID, string(infoJson))
}

func userId(msg *tgbotapi.Message, _ *tgModel.Command) *tgModel.HandlerResult {
	return tgModel.Simple(msg.Chat.ID, strconv.FormatInt(msg.Chat.ID, 10))
}

func appVersion(msg *tgbotapi.Message, _ *tgModel.Command) *tgModel.HandlerResult {
	return tgModel.Simple(msg.Chat.ID, appStat.Version)
}

func startDefault(msg *tgbotapi.Message, _ *tgModel.Command) *tgModel.HandlerResult {
	return tgModel.Simple(msg.Chat.ID, "Hi "+msg.From.String()+" and welcome")
}

func setNextCommand(msg *tgbotapi.Message, command *tgModel.Command) *tgModel.HandlerResult {
	//fmt.Println("setNextCommand", msg.Text, command.Command)
	//return tgModel.DeferredWithText(msg.Chat.ID, mstText, commandChoicer, message)
	return tgModel.EmptyCommand().WithDeferred(commandRedirect, msg)
}

func setRedirectByCommand(msg *tgbotapi.Message, command *tgModel.Command) *tgModel.HandlerResult {
	//fmt.Println("setNextCommand", msg.Text, command.Command)
	redirectToStr := []rune(command.Arguments.Raw)
	if len(redirectToStr) == 0 {
		return tgModel.EmptyCommand()
	}
	if redirectToStr[0] == rune('/') {
		redirectToStr = redirectToStr[1:]
	}
	zlog.Info().
		Any("setRedirectByCommand text", msg.Text).
		Any("msgid", msg.MessageID).
		Any("command.Command", string(redirectToStr)).
		Send()
	//return tgModel.DeferredWithText(msg.Chat.ID, mstText, commandChoicer, message)
	return tgModel.EmptyCommand().WithRedirect(string(redirectToStr), msg)
}
