package tgbot

import (
	"fmt"
	"fun-coice/config"
	tgCommands "fun-coice/internal/domain/commands/tg"
	"fun-coice/pkg/appStat"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	zlog "github.com/rs/zerolog/log"
	"log"
	"strings"
	"sync"
)

var defaultWorkersCount = 10

type data struct {
	WebUri       string
	Token        string
	StartMsg     bool
	WebMode      bool
	Commands     tgCommands.Commands
	Bot          *tgbotapi.BotAPI
	WorkersCount int
	Deferred     map[int64]string
	mutex        *sync.Mutex
}

func New(token, webUri string) data {
	bot, err := tgbotapi.NewBotAPI(config.TelegramToken())
	if err != nil {
		log.Panic(err)
	}
	return data{
		Token:        token,
		WebUri:       webUri,
		Commands:     defaultCommands,
		Bot:          bot,
		WorkersCount: defaultWorkersCount,
		StartMsg:     true,
		Deferred:     make(map[int64]string),
		mutex:        &sync.Mutex{},
	}
}

func (d data) Run() error {
	//d.Bot.Debug = true
	//TODO: d.Bot.GetMyCommands() AND SET THEM
	if d.StartMsg {
		msg := tgbotapi.NewMessage(config.TelegramAdminId(), "Bot Started with version "+appStat.Version)
		d.Bot.Send(msg)
	}

	log.Printf("Authorized on account %s", d.Bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 10

	if d.WebMode {
		fmt.Println("tg bot WebMode")
		webUpdates := d.Bot.ListenForWebhook(d.WebUri)
		for i := 0; i < d.WorkersCount; i++ {
			go d.CommandsHandler(webUpdates)
		}
	} else {
		fmt.Println("tg bot UpdateMode")
		updates := d.Bot.GetUpdatesChan(u)
		for i := 0; i < d.WorkersCount; i++ {
			go d.CommandsHandler(updates)
		}
	}
	return nil
}

func (d data) AppendDeferred(user int64, command string) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.Deferred[user] = command
}

func (d data) CheckDeferred(user int64) string {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	//fmt.Println(user, len(d.Deferred)) //DEVMODE
	founded := d.Deferred[user]
	delete(d.Deferred, user)
	return founded
}

func (d data) RunCommand(command tgCommands.Command, msg *tgbotapi.Message, commandName string, param string, params []string) bool {
	result := command.Handler(msg, command.Command, param, params)
	if result.Prepared {
		//fmt.Println("COMMAND PREPAERD") //DEVMODE
		_, err := d.Bot.Send(result.ChatEvent)
		if err != nil {
			fmt.Println(err.Error())
		}
		if !result.Wait {
			return true
		}
	}
	if result.Wait {
		//fmt.Println("COMMAND Wait", msg.From.ID, result.Next)  //DEVMODE
		d.AppendDeferred(msg.From.ID, result.Next)
		return true
	}
	return false
}

func (d data) CommandsHandler(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		//TODO: add sended messages history

		zlog.Info().Any("msg", update.Message).Any("InlineQuery", update.InlineQuery).Send()

		if update.CallbackQuery != nil {
			if commandDeferred, ok := d.Commands[update.CallbackQuery.Data]; ok {
				if d.RunCommand(
					commandDeferred,
					update.CallbackQuery.Message,
					commandDeferred.Command,
					update.CallbackQuery.Data,
					strings.Split(update.CallbackQuery.Data, " ")) {
					break
				}
			}
			break
		}

		if update.Message == nil || (update.Message == nil && update.InlineQuery != nil) {
			zlog.Info().Any("update", update).Send()
			continue
		}

		if founded := d.CheckDeferred(update.Message.From.ID); founded != "" {
			if commandDeferred, ok := d.Commands[founded]; ok {
				if d.RunCommand(
					commandDeferred,
					update.Message,
					commandDeferred.Command,
					update.Message.Text,
					strings.Split(update.Message.Text, " ")) {
					break
				}
			}
		}
		for _, command := range d.Commands {
			if !command.Permission(update.Message) || command.Handler == nil {
				continue
			}
			splitCommands, commandValue := splitCommand(update.Message.Text, " ")
			if len(splitCommands) == 0 {
				continue
			}
			commandName := splitCommands[0]
			commandsCount := len(splitCommands)
			if commandsCount == 0 {
				continue
			}
			if !command.IsImplemented(commandName, d.Bot.Self.UserName) {
				if command.IsMatched(update.Message.Text, d.Bot.Self.UserName) {
					commandValue = update.Message.Text
				} else {
					continue
				}
			}
			if d.RunCommand(command, update.Message, command.Command, commandValue, splitCommands) {
				break
			}
		}

		if update.Message.Chat.Type == "private" && config.Str("logLevel") == "private" || config.Str("logLevel") == "chat" {
			log.Printf("INNER MESSAGE %s[%d]: %s",
				update.Message.From.UserName,
				update.Message.From.ID,
				update.Message.Text)
			fmt.Printf("inline query %+v\n", update.InlineQuery)
		}
	}
}
