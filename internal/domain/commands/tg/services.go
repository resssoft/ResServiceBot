package tgModel

import tgbotapi "fun-coice/pkg/telegram-bot-api"

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
