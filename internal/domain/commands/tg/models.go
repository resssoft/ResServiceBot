package tgModel

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type SentMessages chan<- tgbotapi.Chattable
