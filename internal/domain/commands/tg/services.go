package tgModel

import (
	"database/sql"
	"slices"

	"fun-coice/internal/database"
	"fun-coice/pkg/scribble"
	tgbotapi "fun-coice/pkg/telegram-bot-api"
)

type Service interface {
	Commands() Commands
	Name() string // TODO: use in the tgBot after append commands
	Configure(ServiceConfig) error
	Dependency() *ServiceDepends
	Destroy()
}

type ServiceConfig struct {
	MessageSender MessageSender
	FileDb        *scribble.Driver
	SqliteDb      *sql.DB
	MongoClient   database.MongoClientApplication
	MsgRepo       MsgRepository
	Params        map[string]string
	OwnerId       int64
	AdminsIds     []int64
}

type MessageSender interface {
	BotName() string
	PushMessage() chan<- tgbotapi.Chattable
	PushHandleResult() chan<- *HandlerResult
}

type ServiceDepend string

const (
	FileDbDependency   ServiceDepend = "FileDbDependency"
	MongoDbDependency  ServiceDepend = "MongoDbDependency"
	SqliteDbDependency ServiceDepend = "SqliteDbDependency"
	MsgRepoDependency  ServiceDepend = "MsgRepoDependency"
	PluginsDependency  ServiceDepend = "PluginsDependency"
)

type ServiceDepends struct {
	Names []ServiceDepend
}

func (s ServiceDepends) Needed(depend ServiceDepend) bool {
	return slices.Contains(s.Names, depend)
}
func ServiceDependsIs(depends ...ServiceDepend) *ServiceDepends {
	return &ServiceDepends{
		Names: depends,
	}
}
