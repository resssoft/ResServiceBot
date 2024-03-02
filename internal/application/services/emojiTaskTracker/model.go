package emojiTaskTracker

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type userData struct {
	MongoID primitive.ObjectID `bson:"_id"`
	tasks   map[string]Task
}

type Task struct {
	MongoID    primitive.ObjectID `bson:"_id"`
	Start      time.Time          `bson:"start"`
	End        time.Time          `bson:"end"`
	Break      time.Time          `bson:"break"`
	Title      string             `bson:"title"`
	UserId     int64              `bson:"user_id"`
	MsgId      int                `bson:"message_id"`
	Breaks     []TimeItem         `bson:"breaks"`
	Status     Status             `bson:"status"`
	ActiveTask int                `bson:"active_task"`
	BotName    string             `bson:"not_name"`
	Code       string             `bson:"code"`
	//GMT string use for time show
}

type Status int

const (
	StatusCreated = iota
	StatusProgress
	StatusPause
	StatusStopped
	StatusProfile
	StatusSettings
)

type TimeItem struct {
	Id       int
	Name     string
	Start    time.Time
	End      time.Time
	Duration time.Duration
}

func (s Status) String() string {
	switch s {
	case StatusCreated:
		return "Created"
	case StatusProgress:
		return "Started"
	case StatusPause:
		return "Paused"
	case StatusStopped:
		return "Stopped"
	default:
		return "none"
	}
}
