package model

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BotData struct {
	MongoID         primitive.ObjectID `bson:"_id"`
	BotName         string             `bson:"bot_name"`
	RepeatReactions map[string]string  `bson:"repeat_reactions"`
	TextReactions   map[string]string  `bson:"text_reactions"`
}

type BotsData map[string]BotData

type Repository interface {
	Create(context.Context, BotData) (BotData, error)
	Update(context.Context, BotData) error
	List(context.Context, string) (BotData, error)
	Delete(context.Context, int64) error
}
