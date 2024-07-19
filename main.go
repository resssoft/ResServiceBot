package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"fun-coice/config"
	"fun-coice/internal/application/botBuilder"
	"fun-coice/internal/mediator"
	"fun-coice/pkg/appStat"
	"fun-coice/pkg/version"

	_ "github.com/mattn/go-sqlite3"
)

var (
	onExit chan int
)

type SystemListener struct{}

func main() {
	onExit = make(chan int)

	showVer := flag.Bool("v", false, "show version")
	checkConfig := flag.Bool("c", false, "check config")
	flag.Parse()
	if *showVer {
		fmt.Println("version", appStat.Version, "build date", version.Get())
		return
	}

	var err error
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	fmt.Print("Load configuration... ")
	config.Configure()
	if *checkConfig {
		botsConfigJson, err := json.MarshalIndent(config.TgBots(), "", "    ")
		fmt.Println(err, string(botsConfigJson))
		return
	}

	//logs, crons
	dispatcher := mediator.NewDispatcher()
	if err := dispatcher.Register(
		SystemListener{},
		mediator.AppExit,
		mediator.SetLogDebugMode,
		mediator.SetLogInfoMode); err != nil {
		log.Info().Err(err).Send()
	}

	//TODO: TRANSLATES
	//TODO: DebugMode

	bots := botBuilder.New(dispatcher)
	bots.Build()

	fmt.Println("Start web server by " + config.WebServerAddr())
	err = http.ListenAndServe(config.WebServerAddr(), nil)
	if err != nil {
		fmt.Println("Error", err)
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
