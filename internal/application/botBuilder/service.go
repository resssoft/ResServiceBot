package botBuilder

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"

	"fun-coice/config"
	"fun-coice/internal/application/services/adminNotifer"
	"fun-coice/internal/application/services/admins"
	"fun-coice/internal/application/services/b64"
	"fun-coice/internal/application/services/calculator"
	"fun-coice/internal/application/services/chatAdmin"
	"fun-coice/internal/application/services/citiesTime"
	"fun-coice/internal/application/services/datatimes"
	"fun-coice/internal/application/services/emojiTaskTracker"
	"fun-coice/internal/application/services/examples"
	"fun-coice/internal/application/services/funs"
	"fun-coice/internal/application/services/images"
	"fun-coice/internal/application/services/lists"
	"fun-coice/internal/application/services/msgStore"
	"fun-coice/internal/application/services/p2p"
	"fun-coice/internal/application/services/qrcodes"
	"fun-coice/internal/application/services/testManager"
	"fun-coice/internal/application/services/text"
	"fun-coice/internal/application/services/translate"
	"fun-coice/internal/application/services/transliter"
	"fun-coice/internal/application/services/users"
	"fun-coice/internal/application/services/workTasks"
	"fun-coice/internal/application/tgbot"
	"fun-coice/internal/database"
	tgModel "fun-coice/internal/domain/commands/tg"
	"fun-coice/internal/mediator"
	tgmessage "fun-coice/internal/repositories/telegram/message"
	"fun-coice/pkg/scribble"
)

var (
	dbFile = "./tg-sqlite3.db"
)

type BuilderService struct {
	dispatcher    *mediator.Dispatcher
	sqliteDb      *sql.DB
	serviceConfig tgModel.ServiceConfig
}

func New(dispatcher *mediator.Dispatcher) *BuilderService {
	serv := &BuilderService{
		dispatcher: dispatcher,
	}
	if err := dispatcher.Register(
		Listener{
			Client: serv,
		},
		mediator.BotBuilderEvents...); err != nil {
		log.Info().Err(err).Send()
	}
	return serv
}

func (bs *BuilderService) Build() {
	//TODO: restart only one bot - read config again
	//TODO: turn off or on services by bot
	//weatherTokens := map[string]string{ //weatherTokens
	//	"yandex":   config.Str("plugins.yandex_weather.token"),
	//	"gismeteo": config.Str("plugins.gismeteo.token"),
	//}

	services := []tgModel.Service{
		funs.New(),
		b64.New(),
		qrcodes.New(),
		datatimes.New(),
		translate.New(),
		calculator.New(),
		//financy.New(config.Str("plugins.apilayer.token")), // TODO: plugins tokens to settings (send admin notify for set token from TG
		lists.New(),
		users.New(),
		text.New(),
		images.New(), // TODO: provide Bot var to commandHandler
		examples.New(),
		adminNotifer.New(),
		chatAdmin.New(),
		admins.New(), // TODO: provide Bot var to commandHandler use middleware channels
		msgStore.New(),
		//weather.New(), // TODO: plugins tokens to settings (send admin notify for set token from TG
		transliter.New(),
		p2p.New(),
		workTasks.New(), // TODO: plan
		emojiTaskTracker.New(),
		testManager.New(),
		citiesTime.New(),
	}

	for botName, tgBotConfig := range config.TgBots() {
		fmt.Print("Found bot [" + botName + "]\n")
		if tgBotConfig.Active {
			fmt.Print("\nPrepare bot " + botName + " with services: ")
			tgBot, err := tgbot.New(botName, tgBotConfig)
			if err != nil {
				log.Printf("Error: bot cant be started: ", botName, err)
				continue
			}
			installedServices := ""
			for _, botService := range tgBotConfig.Services {
				for _, serviceItem := range services {
					if botService == serviceItem.Name() {
						err = serviceItem.Configure(bs.SetDependency(serviceItem, tgBot))
						if err != nil {
							log.Print(", [SKIP(" + err.Error() + "):" + serviceItem.Name() + "]")
							continue
						}
						installedServices += fmt.Sprintf(" [%s]", serviceItem.Name())
						log.Print(" [" + serviceItem.Name() + "]")
						tgBot.AddCommands(serviceItem.Commands(), serviceItem.Name())
					}
				}
			}
			fmt.Print(installedServices)
			tgBot.DefaultCommand = tgBotConfig.DefaultCommand //TODO: set method
			log.Print(" Staring...\n")
			err = tgBot.Run()
			if err != nil {
				log.Info().Err(err).Send()
			} else {
				tgBot.SendMsg("Services: " + installedServices)
			}
		} else {
			log.Info().Msgf("Inactive bot: %s", botName)
		}
	}
}

func (bs *BuilderService) SetDependency(service tgModel.Service, tgBot *tgbot.Data) tgModel.ServiceConfig {
	var err error
	newConfig := tgModel.ServiceConfig{
		MessageSender: tgBot,
		OwnerId:       tgBot.AdminId,
	}
	if service.Dependency() == nil {
		return newConfig
	}
	if service.Dependency().Needed(tgModel.FileDbDependency) {
		if bs.serviceConfig.FileDb == nil {
			log.Printf("Work with DB...")
			appPath, err := os.Getwd()
			if err != nil {
				log.Fatal().Err(err).Send()
			}
			//TODO: moved simple DB implement to pkg
			//TODO: create db interface layer + MOVE TO APPLICATION FOLDER
			FileDb, err := scribble.New(appPath + "/data")
			if err != nil {
				fmt.Println("Error", err)
			}
			bs.serviceConfig.FileDb = FileDb
			newConfig.FileDb = FileDb
		} else {
			newConfig.FileDb = bs.serviceConfig.FileDb
		}
	}
	if service.Dependency().Needed(tgModel.MongoDbDependency) {
		if bs.serviceConfig.MongoClient == nil {
			MongoClient, err := database.ProvideMongo(config.DbMongoUrl(), config.DbMongoDbName(), bs.dispatcher)
			if err != nil {
				log.Error().Err(err).Send()
			} else {
				bs.serviceConfig.MongoClient = MongoClient
				newConfig.MongoClient = MongoClient
			}
		} else {
			newConfig.MongoClient = bs.serviceConfig.MongoClient
		}
	}
	if service.Dependency().Needed(tgModel.SqliteDbDependency) || service.Dependency().Needed(tgModel.MsgRepoDependency) {
		if bs.serviceConfig.SqliteDb == nil {
			if _, err = os.Stat(dbFile); err != nil {
				log.Info().Msg("Creating sqlite-database.db...")
				file, err := os.Create(dbFile)
				if err != nil {
					log.Fatal().Err(err).Msg("cant create db sql3 file")
				}
				file.Close()
				dbFilePath, _ := filepath.Abs(dbFile)
				log.Printf("\nsqlite-database.db created %s", dbFilePath)
			}
			dbFilePath, _ := filepath.Abs(dbFile)
			log.Printf("\nsqlite-database.db %s", dbFilePath)

			sqliteDb, err := sql.Open("sqlite3", dbFile) // or file::memory:?cache=shared //:memory:
			if err != nil {
				log.Error().Err(err).Msg("cant open db sql3 file")
			} else {
				bs.sqliteDb = sqliteDb
				bs.serviceConfig.SqliteDb = sqliteDb
				newConfig.SqliteDb = sqliteDb
			}
		} else {
			newConfig.SqliteDb = bs.serviceConfig.SqliteDb
		}
	}
	if service.Dependency().Needed(tgModel.MsgRepoDependency) {
		if bs.serviceConfig.MsgRepo == nil && bs.serviceConfig.SqliteDb != nil {
			msgRepo, err := tgmessage.New(bs.serviceConfig.SqliteDb)
			if err != nil {
				log.Fatal().Err(err).Msg("cant create msg repo")
			} else {
				bs.serviceConfig.MsgRepo = msgRepo
				newConfig.MsgRepo = msgRepo
			}
		} else {
			if bs.serviceConfig.MsgRepo != nil {
				newConfig.MsgRepo = bs.serviceConfig.MsgRepo
			}
		}
	}
	if service.Dependency().Needed(tgModel.PluginsDependency) {
		params := config.Params("plugins." + service.Name())
		newConfig.Params = params
	}
	return newConfig
}

func (bs *BuilderService) Destroy() {
	if bs.sqliteDb != nil {
		bs.sqliteDb.Close()
	}
}

func (bs *BuilderService) AddService() {
	//Not implemented
}
