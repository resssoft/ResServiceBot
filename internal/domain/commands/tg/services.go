package tgModel

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type Service interface {
	Commands() Commands
	Name() string // TODO: use in the tgBot after append commands
	Configure(ServiceConfig)
}

type ServiceConfig struct {
	MessageSender MessageSender
}

type MessageSender interface {
	BotName() string
	PushMessage() chan<- tgbotapi.Chattable
	PushHandleResult() chan<- *HandlerResult
}
