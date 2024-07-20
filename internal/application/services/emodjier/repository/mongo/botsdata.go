package repository

import (
	"context"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"fun-coice/internal/application/services/emodjier/repository/model"
)

func (r *RepoMongo) Create(ctx context.Context, item model.BotData) (model.BotData, error) {
	item.MongoID = primitive.NewObjectID()
	log.Debug().Interface("new item", item).Send()
	_, err := r.collection.InsertOne(r.dbApp.GetContext(), item)
	if err != nil {
		log.Error().AnErr("Insert item error", err).Send()
		return item, err
	}
	return item, nil
}

func (r *RepoMongo) Update(ctx context.Context, item model.BotData) error {
	log.Info().Interface("upd item", item).Send()
	_, err := r.collection.UpdateOne(
		r.dbApp.GetContext(),
		bson.M{"_id": item.MongoID},
		bson.D{
			{"$set", &item},
		})
	if err != nil {
		log.Error().AnErr("Update item error", err).Send()
		return err
	}
	return nil
}

func (r *RepoMongo) Get(ctx context.Context, ID int64) (*model.BotData, error) {
	//Not implemented
	var err error
	return nil, err
}

func (r *RepoMongo) List(ctx context.Context, botName string) (model.BotData, error) {
	return r.GetByField("bot_name", botName)
}

func (r *RepoMongo) Delete(ctx context.Context, ID int64) error {
	//Not implemented
	var err error
	return err
}
