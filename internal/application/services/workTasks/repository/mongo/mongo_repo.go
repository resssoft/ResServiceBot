package mongo_repo

import (
	"context"
	"fun-coice/internal/application/services/workTasks/repository"
	"fun-coice/internal/application/services/workTasks/track"
	"fun-coice/internal/database"
	"go.mongodb.org/mongo-driver/mongo"
)

const LeadCollectionName = "leadData"

var _ repository.Repository = (*RepoMongo)(nil)

type RepoMongo struct {
	dbApp      database.MongoClientApplication
	collection *mongo.Collection
}

func NewMongoRepo(db database.MongoClientApplication) (repository.Repository, error) {
	collection := db.GetCollection(LeadCollectionName)
	return &RepoMongo{
		dbApp:      db,
		collection: collection,
	}, nil
}

func (r *RepoMongo) NewItemLink() *track.Track {
	newItem := new(track.Track)
	return newItem
}

func (r *RepoMongo) NewItemsLink() track.Tracks {
	return make(map[int64]track.Track)
}

func (r *RepoMongo) Migrate() error {
	return nil
}

func (r *RepoMongo) Create(ctx context.Context, track track.Track) error {
	var err error
	return err
}

func (r *RepoMongo) Update(ctx context.Context, item *track.Track) error {
	var err error
	return err
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
