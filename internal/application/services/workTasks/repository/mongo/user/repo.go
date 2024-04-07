package userMongoRepository

import (
	"fun-coice/internal/application/services/workTasks/repository"
	"fun-coice/internal/application/services/workTasks/track"
	"fun-coice/internal/database"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const trackCollectionName = "timeTraker_user"

var _ repository.UserRepository = (*RepoMongo)(nil)

type RepoMongo struct {
	dbApp      database.MongoClientApplication
	collection *mongo.Collection
}

func NewUserRepo(db database.MongoClientApplication) (repository.UserRepository, error) {
	collection := db.GetCollection(trackCollectionName)
	return &RepoMongo{
		dbApp:      db,
		collection: collection,
	}, nil
}

func (r *RepoMongo) NewItemLink() *track.User {
	newItem := new(track.User)
	return newItem
}

func (r *RepoMongo) NewItemsLink() track.Users {
	return make(map[int64]track.User)
}

func (r *RepoMongo) Migrate() error {
	return nil
}

func (r *RepoMongo) GetAllByField(name string, value interface{}) ([]track.User, error) {
	item := track.User{}
	filter := bson.M{name: value}
	var items []track.User
	cursor, err := r.collection.Find(r.dbApp.GetContext(), filter)
	if err != nil {
		return items, err
	}
	defer cursor.Close(r.dbApp.GetContext())
	for cursor.Next(r.dbApp.GetContext()) {
		err := cursor.Decode(&item)
		if err != nil {
			log.Error().AnErr("user read error", err).Send()
			continue
		}
		items = append(items, item)
	}
	if err := cursor.Err(); err != nil {
		return items, err
	}
	return items, nil
}

func (r *RepoMongo) GetByField(name string, value interface{}) (track.User, error) {
	item := track.User{}
	filter := bson.M{name: value}
	err := r.collection.FindOne(r.dbApp.GetContext(), filter).Decode(&item)
	if err != nil {
		log.Error().AnErr("user read error", err).Interface(name, value).Send()
		return item, err
	}
	return item, nil
}
