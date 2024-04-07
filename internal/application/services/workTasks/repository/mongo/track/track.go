package trackMongoRepository

import (
	"context"
	"fun-coice/internal/application/services/workTasks/track"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (r *RepoMongo) Create(ctx context.Context, item track.Track) (track.Track, error) {
	item.MongoID = primitive.NewObjectID()
	log.Debug().Interface("new track", item).Send()
	_, err := r.collection.InsertOne(r.dbApp.GetContext(), item)
	if err != nil {
		log.Error().AnErr("Insert track error", err).Send()
		return item, err
	}
	return item, nil
}

func (r *RepoMongo) Update(ctx context.Context, item *track.Track) error {
	log.Info().Interface("upd track", item).Send()
	_, err := r.collection.UpdateOne(
		r.dbApp.GetContext(),
		bson.M{"_id": item.MongoID},
		bson.D{
			{"$set", item},
		})
	if err != nil {
		log.Error().AnErr("Update track error", err).Send()
		return err
	}
	return nil
}

func (r *RepoMongo) Get(ctx context.Context, ID int64) (*track.Track, error) {
	var err error
	return nil, err
}

func (r *RepoMongo) List(ctx context.Context, filter track.TrackFilter) (track.Tracks, error) {
	items := r.NewItemsLink()
	return items, nil
}

func (r *RepoMongo) Delete(ctx context.Context, ID int64) error {
	var err error
	return err
}
