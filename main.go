package main

import (
	"fmt"
	"fun-coice/config"
	"fun-coice/internal/application/services/adminNotifer"
	"fun-coice/internal/application/services/admins"
	"fun-coice/internal/application/services/b64"
	"fun-coice/internal/application/services/calculator"
	"fun-coice/internal/application/services/datatimes"
	"fun-coice/internal/application/services/examples"
	"fun-coice/internal/application/services/funs"
	"fun-coice/internal/application/services/images"
	"fun-coice/internal/application/services/lists"
	"fun-coice/internal/application/services/money"
	"fun-coice/internal/application/services/qrcodes"
	"fun-coice/internal/application/services/text"
	"fun-coice/internal/application/services/translate"
	"fun-coice/internal/application/services/transliter"
	"fun-coice/internal/application/services/weather"
	"fun-coice/internal/application/tgbot"
	"fun-coice/pkg/scribble"
	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
	"log"
	"net/http"
	"os"
)

var DB *scribble.Driver

func main() {
	var err error
	zlog.Level(zerolog.DebugLevel)
	fmt.Print("Load configuration... ")
	config.Configure()

	fmt.Println(fmt.Sprintf("apilayer[%s]", config.Str("plugins.apilayer.token")))

	multiBot, err := tgbot.New("multybot")
	if err != nil {
		log.Fatal(err)
	}
	//multiBot.ADmin = config.TelegramAdminLogin()
	//log.Printf("Admin is ..." + config.TelegramAdminLogin())
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

	if config.Bool("telegram.bots.multybot.active") {
		fmt.Println("commands count", len(multiBot.Commands))
		funCommandsService := funs.New(DB)
		multiBot.AddCommands(funCommandsService.Commands())

		b64Service := b64.New()
		multiBot.AddCommands(b64Service.Commands())
		fmt.Println("commands count", len(multiBot.Commands))

		QrCodesService := qrcodes.New()
		multiBot.AddCommands(QrCodesService.Commands())

		dataTimesService := datatimes.New()
		multiBot.AddCommands(dataTimesService.Commands())

		trService := translate.New()
		multiBot.AddCommands(trService.Commands())

		calculatorService := calculator.New()
		multiBot.AddCommands(calculatorService.Commands())

		financeService := financy.New(config.Str("plugins.apilayer.token"))
		multiBot.AddCommands(financeService.Commands())

		listService := lists.New(DB)
		multiBot.AddCommands(listService.Commands())

		usersService := lists.New(DB)
		multiBot.AddCommands(usersService.Commands())

		textService := text.New()
		multiBot.AddCommands(textService.Commands())

		imageService := images.New("multybot")
		multiBot.AddCommands(imageService.Commands())

		exampleService := examples.New()
		multiBot.AddCommands(exampleService.Commands())

		//last init for command list
		adminService := admins.New(multiBot.Bot, DB, multiBot.Commands, "multybot") // TODO: provide Bot var to commandHandler
		multiBot.AddCommands(adminService.Commands())

		weatherTokens := map[string]string{
			"yandex":   config.Str("plugins.yandex_weather.token"),
			"gismeteo": config.Str("plugins.gismeteo.token"),
		}
		if false { //TODO: FIX SERVICE AND REMOVE CONDITION
			weatherService := weather.New(multiBot.GetSentMessages(), DB, weatherTokens)
			multiBot.AddCommands(weatherService.Commands())
		}
		err = multiBot.Run()
		if err != nil {
			log.Panic(err)
		}
	}

	if config.Bool("telegram.bots.translitbot.active") {
		translitBot, err := tgbot.New("translitbot")
		if err != nil {
			log.Fatal(err)
		}
		translitService := transliter.New()
		translitBot.AddCommands(translitService.Commands())
		translitBot.DefaultCommand = "translit"

		adminService2 := admins.New(translitBot.Bot, DB, translitBot.Commands, "translitbot")
		translitBot.AddCommands(adminService2.Commands())

		notiferEventsService := adminNotifer.New(config.TelegramAdminId("translitbot"))
		translitBot.AddCommands(notiferEventsService.Commands())
		err = translitBot.Run()
		if err != nil {
			log.Panic(err)
		}
	}

	fmt.Println("Start web server by" + config.WebServerAddr())
	err = http.ListenAndServe(config.WebServerAddr(), nil)
	if err != nil {
		fmt.Println("Error", err)
	}
}
