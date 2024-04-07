package repository

import (
	"context"
	"fun-coice/internal/application/services/workTasks/track"
)

type TrackRepository interface {
	Migrate() error
	Create(context.Context, track.Track) (track.Track, error)
	Update(context.Context, *track.Track) error
	Get(context.Context, int64) (*track.Track, error)
	List(context.Context, track.TrackFilter) (track.Tracks, error)
	Delete(context.Context, int64) error
}

type UserRepository interface {
	Migrate() error
	Create(context.Context, track.User) (track.User, error)
	Update(context.Context, *track.User) error
	Get(context.Context, int64) (*track.User, error)
	List(context.Context, track.UserFilter) (track.Users, error)
	Delete(context.Context, int64) error
}
