package main

import (
	"fmt"
	"fun-coice/config"
	"fun-coice/funs"
	"fun-coice/internal/application/services/admins"
	"fun-coice/internal/application/services/b64"
	"fun-coice/internal/application/services/calculator"
	"fun-coice/internal/application/services/datatimes"
	"fun-coice/internal/application/services/examples"
	"fun-coice/internal/application/services/images"
	"fun-coice/internal/application/services/lists"
	"fun-coice/internal/application/services/money"
	"fun-coice/internal/application/services/qrcodes"
	"fun-coice/internal/application/services/text"
	"fun-coice/internal/application/services/translate"
	"fun-coice/internal/application/tgbot"
	tgCommands "fun-coice/internal/domain/commands/tg"
	"fun-coice/pkg/scribble"
	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
	"log"
	"net/http"
	"os"
	"strconv"
)

var DB *scribble.Driver

func main() {
	var err error
	zlog.Level(zerolog.DebugLevel)
	fmt.Print("Load configuration... ")
	config.Configure()

	fmt.Println(fmt.Sprintf("apilayer[%s]", config.Str("plugins.apilayer.token")))
	fmt.Println(fmt.Sprintf("Telegram[%s]", config.TelegramToken()))
	fmt.Println(fmt.Sprintf("Admin[%v]", config.TelegramAdminId()))

	multiBot := tgbot.New(config.TelegramToken(), "/tg/multybot/")
	//multiBot.ADmin = config.TelegramAdminLogin()
	log.Printf("Admin is ..." + config.TelegramAdminLogin())
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
	var existAdmin = tgCommands.TGUser{}
	if err := DB.Read("user", strconv.FormatInt(int64(config.TelegramAdminId()), 10), &existAdmin); err != nil {
		fmt.Println("admin not found error", err)
		existAdmin = tgCommands.TGUser{
			UserID:  config.TelegramAdminId(),
			ChatId:  0,
			Login:   "",
			Name:    "",
			IsAdmin: false,
		}
		if err := DB.Write("user", strconv.FormatInt(int64(config.TelegramAdminId()), 10), existAdmin); err != nil {
			fmt.Println("Error", err)
		}
	}
	fmt.Println("commands count", len(multiBot.Commands))
	funCommandsService := funs.New(DB)
	multiBot.Commands = multiBot.Commands.Merge(funCommandsService.Commands())

	b64Service := b64.New()
	multiBot.Commands = multiBot.Commands.Merge(b64Service.Commands())
	fmt.Println("commands count", len(multiBot.Commands))

	QrCodesService := qrcodes.New()
	multiBot.Commands = multiBot.Commands.Merge(QrCodesService.Commands())

	dataTimesService := datatimes.New()
	multiBot.Commands = multiBot.Commands.Merge(dataTimesService.Commands())

	trService := translate.New()
	multiBot.Commands = multiBot.Commands.Merge(trService.Commands())

	calculatorService := calculator.New()
	multiBot.Commands = multiBot.Commands.Merge(calculatorService.Commands())

	financeService := financy.New(config.Str("plugins.apilayer.token"))
	multiBot.Commands = multiBot.Commands.Merge(financeService.Commands())

	listService := lists.New(DB)
	multiBot.Commands = multiBot.Commands.Merge(listService.Commands())

	usersService := lists.New(DB)
	multiBot.Commands = multiBot.Commands.Merge(usersService.Commands())

	textService := text.New()
	multiBot.Commands = multiBot.Commands.Merge(textService.Commands())

	imageService := images.New()
	multiBot.Commands = multiBot.Commands.Merge(imageService.Commands())

	exampleService := examples.New()
	multiBot.Commands = multiBot.Commands.Merge(exampleService.Commands())

	//last init for command list
	adminService := admins.New(multiBot.Bot, DB, multiBot.Commands)
	multiBot.Commands = multiBot.Commands.Merge(adminService.Commands())

	err = multiBot.Run()
	if err != nil {
		log.Panic(err)
	}

	err = http.ListenAndServe("0.0.0.0:8080", nil)
	if err != nil {
		fmt.Println("Error", err)
	}
}
