package repository

import (
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"fun-coice/internal/application/services/emodjier/repository/model"
	"fun-coice/internal/database"
)

const CollectionName = "emojier_botsdata"

var _ model.Repository = (*RepoMongo)(nil)

type RepoMongo struct {
	dbApp      database.MongoClientApplication
	collection *mongo.Collection
}

func NewRepo(db database.MongoClientApplication) (model.Repository, error) {
	collection := db.GetCollection(CollectionName)
	return &RepoMongo{
		dbApp:      db,
		collection: collection,
	}, nil
}

func (r *RepoMongo) NewItemLink() *model.BotData {
	newItem := new(model.BotData)
	return newItem
}

func (r *RepoMongo) NewItemsLink() model.BotsData {
	return make(map[string]model.BotData)
}

func (r *RepoMongo) Migrate() error {
	return nil
}

func (r *RepoMongo) GetAllByField(name string, value interface{}) ([]model.BotData, error) {
	item := model.BotData{}
	filter := bson.M{name: value}
	var items []model.BotData
	cursor, err := r.collection.Find(r.dbApp.GetContext(), filter)
	if err != nil {
		return items, err
	}
	defer cursor.Close(r.dbApp.GetContext())
	for cursor.Next(r.dbApp.GetContext()) {
		err := cursor.Decode(&item)
		if err != nil {
			log.Error().AnErr("item read error", err).Send()
			continue
		}
		items = append(items, item)
	}
	if err := cursor.Err(); err != nil {
		return items, err
	}
	return items, nil
}

func (r *RepoMongo) GetByField(name string, value interface{}) (model.BotData, error) {
	item := model.BotData{}
	filter := bson.M{name: value}
	err := r.collection.FindOne(r.dbApp.GetContext(), filter).Decode(&item)
	if err != nil {
		log.Error().AnErr("user read error", err).Interface(name, value).Send()
		return item, err
	}
	return item, nil
}
