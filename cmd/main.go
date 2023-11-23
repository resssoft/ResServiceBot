package main2

import (
	"fun-coice/internal/database"
	"fun-coice/internal/fileLogger"
	"fun-coice/internal/mediator"
	"fun-coice/internal/repository"
	"github.com/robfig/cron"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	//"gitlab.com/AppsgeyserGroup/servers/SalesBot/internal/amoCRM"
	//"gitlab.com/AppsgeyserGroup/servers/SalesBot/internal/database"
	//"gitlab.com/AppsgeyserGroup/servers/SalesBot/internal/fileLogger"
	//"gitlab.com/AppsgeyserGroup/servers/SalesBot/internal/mediator"
	//messenger "gitlab.com/AppsgeyserGroup/servers/SalesBot/internal/messengers"
	//"gitlab.com/AppsgeyserGroup/servers/SalesBot/internal/models"
	//"gitlab.com/AppsgeyserGroup/servers/SalesBot/internal/repository"
	//pipeline "gitlab.com/AppsgeyserGroup/servers/SalesBot/internal/triggersHandler"
	//routing "gitlab.com/AppsgeyserGroup/servers/SalesBot/internal/webServer"
	"os"
	"time"
)

type SystemListener struct{}

var onExit chan int

func main() {
	var err error
	onExit = make(chan int)
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	dispatcher := mediator.NewDispatcher()
	if err := dispatcher.Register(
		SystemListener{},
		mediator.AppExit,
		mediator.SetLogDebugMode,
		mediator.SetLogInfoMode); err != nil {
		log.Info().Err(err).Send()
	}

	logFiles := map[string]string{
		"fatal.txt":      mediator.FileLogFatal,
		"errors.txt":     mediator.FileLogErrors,
		"contacts.txt":   mediator.FileLogContacts,
		"webHooks.txt":   mediator.FileLogWebHooks,
		"requests.txt":   mediator.FileLogRequests,
		"messenger.txt":  mediator.FileLogMessenger,
		"amoCRM.txt":     mediator.FileLogAmoCRM,
		"amoLatency.txt": mediator.FileLogAmoLatency,
	}

	loggerClient := fileLogger.Provide(dispatcher)
	for filename, logName := range logFiles {
		err = loggerClient.AddSource(filename, logName)
		if err != nil {
			log.Info().Err(err).Msgf("Error open log file %s", filename)
		}
	}
	time.Sleep(time.Second)
	defer loggerClient.CloseAll()

	mongoDbApp, err := database.ProvideMongo("", "", dispatcher)
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	_ = repository.NewUserRepo(mongoDbApp)
	//leadRep := repository.NewLeadRepo(mongoDbApp)
	//tokenRep := repository.NewTokenRepo(mongoDbApp)
	//pipeline.Provide(userRep, leadRep, dispatcher)

	//go messenger.Initialize(dispatcher)
	//go routing.NewRouter(dispatcher)
	//go amocrmClient.Provide(dispatcher, leadRep, tokenRep)

	log.Info().Msg("Prepare cron jobs")
	cronJobs := cron.New()
	// Every 6 hours
	err = cronJobs.AddFunc("0 0 */6 * * *", func() {
		log.Debug().Msg("========= START CRON ========= TASK AmoCrmRefreshToken")
		//log.Info().Err(dispatcher.Dispatch(mediator.AmoCrmRefreshToken, models.AmoCrmRefreshTokenEvent{})).Send()
	})
	err = cronJobs.AddFunc("0 */30 * * * *", func() {
		log.Debug().Msg("========= START CRON ========= TASK AmoCrmCronSleepersEvent")
		//log.Info().Err(dispatcher.Dispatch(models.AmoCrmCronSleepers, models.AmoCrmCronSleepersEvent{})).Send()
	})
	if err != nil {
		log.Err(err).Msg("cron err")
	}
	go cronJobs.Start()

	for code := range onExit {
		os.Exit(code)
	}
}

func (u SystemListener) Listen(eventName mediator.EventName, _ interface{}) {
	switch eventName {
	case mediator.AppExit:
		onExit <- 0
	case mediator.SetLogDebugMode:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case mediator.SetLogInfoMode:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}
