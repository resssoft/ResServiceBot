package database

import (
	"context"
	"errors"
	"fmt"
	"fun-coice/internal/mediator"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
)

var (
	onceMongoAction sync.Once
	mongoClient     *mongo.Client
	mongoContext    context.Context
)

type MongoClientApplication interface {
	GetCollection(string) *mongo.Collection
	GetContext() context.Context
	CreateUniqueIndex(*mongo.Collection, string)
	CreateIndexWithTimeout(*mongo.Collection, string, int32)
}

type mongoClientOriginal struct {
	client     *mongo.Client
	context    context.Context
	dbName     string
	url        string
	dispatcher *mediator.Dispatcher
}

func ProvideMongo(url, dbName string, dispatcher *mediator.Dispatcher) (MongoClientApplication, error) {
	onceMongoAction.Do(func() {
		configureMongo(url, dispatcher)
	})
	if mongoClient == nil || mongoContext == nil {
		return &mongoClientOriginal{}, errors.New("mongo client or context is empty")
	}
	return &mongoClientOriginal{
		client:     mongoClient,
		context:    mongoContext,
		dbName:     dbName,
		url:        url,
		dispatcher: dispatcher,
	}, nil
}

func (r *mongoClientOriginal) GetCollection(collection string) *mongo.Collection {
	return mongoClient.Database(r.url).Collection(collection)
}

func (r *mongoClientOriginal) GetContext() context.Context {
	return r.context
}

func (r *mongoClientOriginal) CreateUniqueIndex(collection *mongo.Collection, key string) {
	indexName, err := collection.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys: bson.M{
				key: 1,
			},
			Options: options.Index().SetUnique(true),
		},
	)
	if err != nil {
		log.Error().AnErr("create index", err).Send()
	}
	log.Info().Interface("indexName", indexName).Send()
}

func (r *mongoClientOriginal) CreateIndexWithTimeout(collection *mongo.Collection, key string, expireTimeSeconds int32) {
	indexName, err := collection.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys: bson.M{
				key: 1,
			},
			Options: options.Index().SetExpireAfterSeconds(expireTimeSeconds),
		},
	)
	if err != nil {
		log.Error().AnErr("create index", err).Send()
	}
	log.Info().Interface("indexName", indexName).Send()
}

func configureMongo(address string, dispatcher *mediator.Dispatcher) {
	var err error
	mongoContext = context.Background()
	clientOptions := options.Client().ApplyURI(address)
	mongoClient, err = mongo.Connect(mongoContext, clientOptions)
	if err != nil {
		dispatcher.Dispatch(mediator.LogToFile, mediator.FileLoggerEvent{
			Src: mediator.FileLogFatal,
			Data: fmt.Sprintf("Cannot connect to mongo %s",
				err.Error()),
		})
		log.Fatal().
			Err(err).
			Msg("cannot connect to mongo")
	}
	err = mongoClient.Ping(mongoContext, nil)
	if err != nil {
		dispatcher.Dispatch(mediator.LogToFile, mediator.FileLoggerEvent{
			Src: mediator.FileLogFatal,
			Data: fmt.Sprintf("Cannot connect to mongo: %s",
				err.Error()),
		})
		log.Fatal().
			Err(err).
			Msg("cannot connect to mongo")
	}
}
