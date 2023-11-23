package repository

import (
	"fun-coice/internal/database"
	tgModel "fun-coice/internal/domain/commands/tg"
	"github.com/rs/zerolog/log"
	//"gitlab.com/AppsgeyserGroup/servers/SalesBot/internal/database"
	//"gitlab.com/AppsgeyserGroup/servers/SalesBot/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const UserCollectionName = "user"

type UserRepository interface {
	Add(tgModel.User) (tgModel.User, error)
	Update(primitive.ObjectID, tgModel.User) error
	GetAll() ([]tgModel.User, error)
	GetByMessengerID(string, string) (tgModel.User, error)
	GetByID(string) (tgModel.User, error)
	Count() (int64, error)
}

type userRepo struct {
	dbApp      database.MongoClientApplication
	collection *mongo.Collection
}

func NewUserRepo(db database.MongoClientApplication) UserRepository {
	collection := db.GetCollection(UserCollectionName)
	return &userRepo{
		dbApp:      db,
		collection: collection,
	}
}

func (u *userRepo) Add(user tgModel.User) (tgModel.User, error) {
	user.MongoID = primitive.NewObjectID()
	log.Debug().Interface("new user", user).Send()
	_, err := u.collection.InsertOne(u.dbApp.GetContext(), user)
	if err != nil {
		log.Error().AnErr("Insert user error", err).Send()
		return user, err
	}
	return user, nil
}

func (u *userRepo) Update(id primitive.ObjectID, user tgModel.User) error {
	log.Info().Interface("upd user", user).Send()
	_, err := u.collection.UpdateOne(
		u.dbApp.GetContext(),
		bson.M{"_id": id},
		bson.D{
			{"$set", user},
		})
	if err != nil {
		log.Error().AnErr("Insert user error", err).Send()
		return err
	}
	return nil
}

func (u *userRepo) GetByField(name string, value interface{}) (tgModel.User, error) {
	user := tgModel.User{}
	filter := bson.M{name: value}
	err := u.collection.FindOne(u.dbApp.GetContext(), filter).Decode(&user)
	if err != nil {
		log.Error().AnErr("user read error", err).Interface(name, value).Send()
		return user, err
	}
	return user, nil
}

func (u *userRepo) GetAllByField(name string, value interface{}) ([]tgModel.User, error) {
	lead := tgModel.User{}
	filter := bson.M{name: value}
	var leads []tgModel.User
	cursor, err := u.collection.Find(u.dbApp.GetContext(), filter)
	if err != nil {
		return leads, err
	}
	defer cursor.Close(u.dbApp.GetContext())
	for cursor.Next(u.dbApp.GetContext()) {
		err := cursor.Decode(&lead)
		if err != nil {
			log.Error().AnErr("lead read error", err).Send()
			continue
		}
		leads = append(leads, lead)
	}
	if err := cursor.Err(); err != nil {
		return leads, err
	}
	return leads, nil
}

func (u *userRepo) GetAll() ([]tgModel.User, error) {
	user := tgModel.User{}
	var users []tgModel.User
	cursor, err := u.collection.Find(u.dbApp.GetContext(), bson.D{})
	if err != nil {
		return users, err
	}
	defer cursor.Close(u.dbApp.GetContext())
	for cursor.Next(u.dbApp.GetContext()) {
		err := cursor.Decode(&user)
		if err != nil {
			log.Error().AnErr("user read error", err).Send()
			continue
		}
		users = append(users, user)
	}
	if err := cursor.Err(); err != nil {
		return users, err
	}
	return users, nil
}

func (u *userRepo) GetByID(id string) (tgModel.User, error) {
	user, err := u.GetByField("id", id)
	return user, err
}

func (u *userRepo) GetByMessengerID(id, messenger string) (tgModel.User, error) {
	user, err := u.GetByField(messenger+".id", id)
	return user, err
}

func (u *userRepo) Count() (int64, error) {
	return u.collection.CountDocuments(u.dbApp.GetContext(), bson.D{})
}
