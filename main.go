package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"fun-coice/config"
	"fun-coice/internal/application/services/adminNotifer"
	"fun-coice/internal/application/services/admins"
	"fun-coice/internal/application/services/b64"
	"fun-coice/internal/application/services/calculator"
	"fun-coice/internal/application/services/chatAdmin"
	"fun-coice/internal/application/services/datatimes"
	"fun-coice/internal/application/services/examples"
	"fun-coice/internal/application/services/funs"
	"fun-coice/internal/application/services/images"
	"fun-coice/internal/application/services/lists"
	"fun-coice/internal/application/services/money"
	"fun-coice/internal/application/services/msgStore"
	"fun-coice/internal/application/services/p2p"
	"fun-coice/internal/application/services/qrcodes"
	"fun-coice/internal/application/services/text"
	"fun-coice/internal/application/services/translate"
	"fun-coice/internal/application/services/transliter"
	"fun-coice/internal/application/services/users"
	"fun-coice/internal/application/services/workTasks"
	"fun-coice/internal/application/tgbot"
	tgModel "fun-coice/internal/domain/commands/tg"
	tgmessage "fun-coice/internal/repositories/telegram/message"
	"fun-coice/pkg/appStat"
	"fun-coice/pkg/scribble"
	"fun-coice/pkg/version"
	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
	"log"
	"net/http"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

var DB *scribble.Driver

var dbFile = "./tg-sqlite3.db"

func main() {
	showVer := flag.Bool("v", false, "show version")
	checkConfig := flag.Bool("c", false, "check config")
	flag.Parse()
	if *showVer {
		fmt.Println("version", appStat.Version, "build date", version.Get())
		return
	}

	var err error
	zlog.Level(zerolog.DebugLevel)
	fmt.Print("Load configuration... ")
	config.Configure()
	if *checkConfig {
		botsConfigJson, err := json.MarshalIndent(config.TgBots(), "", "    ")
		fmt.Println(err, string(botsConfigJson))
		return
	}

	//TODO: add falgs for version and config test
	///fmt.Println(fmt.Sprintf("apilayer[%s]", config.Str("plugins.apilayer.token")))

	if _, err = os.Stat(dbFile); err != nil {
		log.Println("Creating sqlite-database.db...")
		file, err := os.Create(dbFile)
		if err != nil {
			log.Fatal(fmt.Errorf("cant create db sql3 file %w", err))
		}
		file.Close()
		dbFilePath, _ := filepath.Abs(dbFile)
		log.Println("sqlite-database.db created", dbFilePath)
	}
	dbFilePath, _ := filepath.Abs(dbFile)
	log.Println("sqlite-database.db ", dbFilePath)

	db, err := sql.Open("sqlite3", dbFile) // or file::memory:?cache=shared //:memory:
	if err != nil {
		log.Fatal(fmt.Errorf("cant open db sql3 file %w", err))
	}
	defer db.Close()

	msgRepo, err := tgmessage.New(db)
	if err != nil {
		log.Fatal(fmt.Errorf("cant create msg repo %w", err))
	}

	log.Printf("Work with DB...")
	appPath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	//TODO: TRANSLATES
	//TODO: DebugMode

	// read admin info from DB or write it to db
	//TODO: moved simple DB implement to pkg
	//TODO: create db interface layer + MOVE TO APPLICATION FOLDER
	DB, err = scribble.New(appPath + "/data")
	if err != nil {
		fmt.Println("Error", err)
	}

	_ = map[string]string{ //weatherTokens
		"yandex":   config.Str("plugins.yandex_weather.token"),
		"gismeteo": config.Str("plugins.gismeteo.token"),
	}

	//add configure or bot register for

	services := []tgModel.Service{
		funs.New(DB),
		b64.New(),
		qrcodes.New(),
		datatimes.New(),
		translate.New(),
		calculator.New(),
		financy.New(config.Str("plugins.apilayer.token")), // TODO: plugins tokens to settings (send admin notify for set token from TG
		lists.New(DB),
		users.New(DB),
		text.New(),
		images.New(), // TODO: provide Bot var to commandHandler
		examples.New(),
		adminNotifer.New(config.TelegramAdminId("multybot")), // TODO: provide Bot var to commandHandler OR from configure method
		chatAdmin.New(config.TelegramAdminId("multybot")),    // TODO: provide Bot var to commandHandler
		admins.New(DB), // TODO: provide Bot var to commandHandler use middleware channels
		msgStore.New(msgRepo),
		//weather.New(multiBot.GetSentMessages(), DB, weatherTokens),  // TODO: plugins tokens to settings (send admin notify for set token from TG
		transliter.New(),
		p2p.New(db),
		workTasks.New(db), // TODO: plan
	}

	for botName, tgBotConfig := range config.TgBots() {
		fmt.Print("Found bot [" + botName + "]\n")
		if tgBotConfig.Active {
			fmt.Print("\nPrepare bot " + botName + " with services: ")
			tgBot, err := tgbot.New(botName, tgBotConfig)
			if err != nil {
				log.Println("Error: bot cant be started: ", botName, err)
			}
			for _, botService := range tgBotConfig.Services {
				for _, serviceItem := range services {
					if botService == serviceItem.Name() {
						log.Print(", [" + serviceItem.Name() + "]")
						tgBot.AddCommands(serviceItem.Commands(), serviceItem.Name())
					}
				}
			}
			tgBot.DefaultCommand = tgBotConfig.DefaultCommand //TODO: set method
			log.Print(" Staring...\n")
			err = tgBot.Run()
			if err != nil {
				log.Println(err.Error())
			}
		} else {
			log.Println("Inactive bot: " + botName)
		}
	}

	fmt.Println("Start web server by " + config.WebServerAddr())
	err = http.ListenAndServe(config.WebServerAddr(), nil)
	if err != nil {
		fmt.Println("Error", err)
	}
}
