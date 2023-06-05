package tgbot

import (
	"fmt"
	"fun-coice/config"
	tgCommands "fun-coice/internal/domain/commands/tg"
	"fun-coice/pkg/appStat"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	zlog "github.com/rs/zerolog/log"
	"log"
	"strconv"
	"strings"
	"sync"
)

var defaultWorkersCount = 10

type Data struct {
	WebUri         string
	Token          string
	StartMsg       bool
	WebMode        bool
	Commands       tgCommands.Commands
	Bot            *tgbotapi.BotAPI
	WorkersCount   int
	Deferred       map[int64]string
	mutex          *sync.Mutex
	DefaultCommand string
}

func New(token, webUri string) Data { //TODO: add error
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}
	return Data{
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

func (d Data) Run() error {
	//d.Bot.Debug = true
	//TODO: d.Bot.GetMyCommands() AND SET THEM
	if d.StartMsg && config.TelegramAdminId() != 0 { //ANOTHER PARAM
		msg := tgbotapi.NewMessage(config.TelegramAdminId(), "Bot Started with version "+appStat.Version)
		d.Bot.Send(msg)
	}

	log.Printf("Authorized bot on account %s", d.Bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 10

	if d.WebMode {
		fmt.Println("tg bot WebMode", d.Bot.Self.UserName)
		webUpdates := d.Bot.ListenForWebhook(d.WebUri)
		for i := 0; i < d.WorkersCount; i++ {
			go d.CommandsHandler(webUpdates, strconv.Itoa(i)+" "+d.Bot.Self.UserName)
		}
	} else {
		fmt.Println("tg bot UpdateMode", d.Bot.Self.UserName)
		updates := d.Bot.GetUpdatesChan(u)
		for i := 0; i < d.WorkersCount; i++ {
			go d.CommandsHandler(updates, strconv.Itoa(i)+" "+d.Bot.Self.UserName)
		}
	}
	return nil
}

func (d Data) AppendDeferred(user int64, command string) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.Deferred[user] = command
}

func (d Data) CheckDeferred(user int64) string {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	//fmt.Println(user, len(d.Deferred)) //DEVMODE
	founded := d.Deferred[user]
	delete(d.Deferred, user)
	return founded
}

func (d Data) RunCommand(command tgCommands.Command, msg *tgbotapi.Message, commandName string, param string, params []string) bool {
	result := command.Handler(msg, command.Command, param, params)
	if result.Prepared {
		//fmt.Println("COMMAND PREPAERD") //DEVMODE
		log.Println("result.Messages", len(result.Messages))
		for _, chantEvent := range result.Messages {
			log.Println("chatEvent", chantEvent)
			_, err := d.Bot.Send(chantEvent)
			if err != nil {
				fmt.Println(err.Error())
			}
		}
		if len(result.Events) > 0 {
			for _, chatEvent := range result.Events {
				eventCommands := d.GetSubCommands(chatEvent.Name.String())
				log.Println("eventCommands", eventCommands)
				for _, eventCommand := range eventCommands {
					log.Println("eventCommands", eventCommand)
					d.RunCommand(eventCommand, msg, commandName, param, params)
				}
			}
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

func (d Data) CommandsHandler(updates tgbotapi.UpdatesChannel, workerID string) {
	//log.Println("start worker CommandsHandler", workerID) // TODO: to debug
	for update := range updates {
		log.Println("update chan EVENT", update.UpdateID, workerID) // TODO: to debug
		//TODO: add sended messages history - save users by bot
		//TODO: save bot stat - new users, by date; new messages

		//only for debug level
		//zlog.Debug().Any("msg", update.Message).Any("InlineQuery", update.InlineQuery).Send()
		if update.Message != nil {
			if update.Message.Chat.Type == "private" {
				zlog.Info().Any("msg", update.Message).Any("InlineQuery", update.InlineQuery).Send()
			}
		}

		if update.CallbackQuery != nil {
			if commandDeferred, ok := d.Commands[update.CallbackQuery.Data]; ok {
				log.Println("CallbackQuery RunCommand")
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

		if update.Message.LeftChatMember != nil {
			d.RunEvents(tgCommands.UserLeaveChantEvent.String(), update.Message, "", "", nil)
		}

		if update.Message.NewChatMembers != nil {
			d.RunEvents(tgCommands.UserJoinedChantEvent.String(), update.Message, "", "", nil)
		}

		zlog.Info().Any("update FULL", update).Send() //TODO: REMOVE

		if founded := d.CheckDeferred(update.Message.From.ID); founded != "" {
			if commandDeferred, ok := d.Commands[founded]; ok {
				log.Println("Deferred RunCommand")
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
		sended := false
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
			sended = true
			log.Println("just RunCommand")
			if d.RunCommand(command, update.Message, command.Command, commandValue, splitCommands) {
				break
			}
		}

		if d.DefaultCommand != "" && !sended {
			fmt.Println("default command:", d.DefaultCommand)
			if commandDeferred, ok := d.Commands[d.DefaultCommand]; ok {
				log.Println("default RunCommand")
				if d.RunCommand(
					commandDeferred,
					update.Message,
					commandDeferred.Command,
					update.Message.Text,
					strings.Split(update.Message.Text, " ")) {
					break
				}
			}
			break
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

func (d Data) GetSubCommands(subName string) []tgCommands.Command {
	var founded []tgCommands.Command
	for key, command := range d.Commands {
		results := strings.Split(key, ":")
		//fmt.Println(key, results, len(results))
		if len(results) < 2 {
			continue
		}
		if results[1] == subName {
			founded = append(founded, command)
		}
	}

	return founded
}

func (d Data) RunEvents(event string, msg *tgbotapi.Message, commandName string, param string, params []string) {
	eventCommands := d.GetSubCommands(event)
	log.Println("RunEvents", eventCommands)
	for _, eventCommand := range eventCommands {
		log.Println("tg event RunCommand", eventCommand)
		d.RunCommand(eventCommand, msg, "", "", nil)
	}
}
