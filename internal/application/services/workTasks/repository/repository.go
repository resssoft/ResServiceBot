package repository

import (
	"context"
	"fun-coice/internal/application/services/workTasks/track"
)

type Repository interface {
	Migrate() error
	Create(context.Context, track.Track) error
	Update(context.Context, *track.Track) error
	Get(context.Context, int64) (*track.Track, error)
	List(context.Context, track.TrackFilter) (track.Tracks, error)
	Delete(context.Context, int64) error
}
