package tgModel

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type MsgFilter struct {
	ID            int64  `json:"id"`
	ChatType      string `json:"chat_type"`
	ChatID        int    `json:"chat_id"`
	MerchantID    string `json:"merchant_id"`
	Limit         int    `json:"limit"`
	Sort          string `json:"sort"`
	SortDirection string `json:"sort_dir"`
}

type Message struct {
	TgMsg        *tgbotapi.Message
	BotName      string
	ID           int
	MsgType      string //text, file, images, new_member...
	MsgDirection int    // 0 from user, 1 from bot
	MsgJson      string
}

type MsgRepository interface {
	Create(ctx context.Context, item *Message) error
	Update(ctx context.Context, item *Message) error
	Get(ctx context.Context, ID int64) (*Message, error)
	List(ctx context.Context, filter *MsgFilter) (Messages, error)
	Delete(ctx context.Context, ID int64) error
}

type Messages []*Message

//TODO: remove TgMsg and provide fields
