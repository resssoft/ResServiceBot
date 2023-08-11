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
	"time"
)

var (
	defaultWorkersCount = 10
	messagesChanLimit   = 300
	messageTimeLimit    = time.Millisecond * 100
)

type Data struct {
	WebUri         string
	Token          string
	StartMsg       bool
	WebMode        bool
	Commands       tgCommands.Commands
	Bot            *tgbotapi.BotAPI
	WorkersCount   int
	Deferred       map[int64]string
	mutexDeferred  *sync.Mutex
	DefaultCommand string
	messagesChan   chan tgbotapi.Chattable
	lastMsgTime    time.Time
	Name           string
	AdminId        int64
	mutexCommands  *sync.Mutex
}

func New(name string) (*Data, error) {
	bot, err := tgbotapi.NewBotAPI(config.TelegramToken(name))
	if err != nil {
		return nil, err
	}

	return &Data{
		Token:          config.TelegramToken(name),
		WebUri:         config.TelegramBotUrl(name),
		Commands:       defaultCommands,
		Bot:            bot,
		WorkersCount:   defaultWorkersCount,
		StartMsg:       true,
		WebMode:        config.TelegramIsWebMode(name),
		Name:           name,
		Deferred:       make(map[int64]string),
		mutexDeferred:  &sync.Mutex{},
		DefaultCommand: config.TelegramBotCommand(name),
		AdminId:        config.TelegramAdminId(name),
		messagesChan:   make(chan tgbotapi.Chattable, messagesChanLimit),
		mutexCommands:  &sync.Mutex{},
	}, nil
}

func (d *Data) Run() error {
	//d.Bot.Debug = true
	//TODO: d.Bot.GetMyCommands() AND SET THEM
	startMsg := fmt.Sprintf(
		"Bot %s Started with version %s \n webMode [%v]\nDefault Command: %s\nUri: %s",
		d.Name,
		appStat.Version,
		d.WebMode,
		d.DefaultCommand,
		d.WebUri)
	if d.StartMsg && d.AdminId != 0 { //ANOTHER PARAM
		msg := tgbotapi.NewMessage(d.AdminId, startMsg)
		d.Bot.Send(msg)
	}

	log.Printf("Authorized bot on account %s \n %s", d.Bot.Self.UserName, startMsg)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 10

	for i := 0; i < d.WorkersCount; i++ {
		go d.MessagesHandler()
	}
	if d.WebMode {
		fmt.Println("tg bot WebMode", d.Bot.Self.UserName, d.WebUri)
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

func (d *Data) SetWebMode(val bool) {
	d.WebMode = val
}

func (d *Data) GetSentMessages() tgCommands.SentMessages {
	return d.messagesChan
}

func (d *Data) AppendDeferred(user int64, command string) {
	d.mutexDeferred.Lock()
	defer d.mutexDeferred.Unlock()
	d.Deferred[user] = command
}

func (d *Data) CheckDeferred(user int64) string {
	d.mutexDeferred.Lock()
	defer d.mutexDeferred.Unlock()
	//fmt.Println(user, len(d.Deferred)) //DEVMODE
	founded := d.Deferred[user]
	delete(d.Deferred, user)
	return founded
}

func (d *Data) RunCommand(command tgCommands.Command, msg *tgbotapi.Message, commandName string, param string, params []string) bool {
	result := command.Handler(msg, command.Command, param, params)
	if result.Prepared {
		//fmt.Println("COMMAND PREPAERD") //DEVMODE
		log.Println("result.Messages", len(result.Messages))
		for _, chantEvent := range result.Messages {
			log.Println("chatEvent", chantEvent)
			msgRes, err := d.Bot.Send(chantEvent)
			if err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Println("Send message", msgRes.MessageID)
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
func (d *Data) MessagesHandler() {
	for msg := range d.messagesChan {
		msgRes, err := d.Bot.Send(msg)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println("Send message by", msgRes.MessageID)
		}
	}
}

func (d *Data) CommandsHandler(updates tgbotapi.UpdatesChannel, workerID string) {
	//log.Println("start worker CommandsHandler", workerID) // TODO: to debug
	isCommand := false
	commandName := ""
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
			if update.CallbackQuery.Message.Chat.Type == "private" {
				zlog.Info().Any("callback query data", update.CallbackQuery.Data).Send()
			}
			eventName := update.CallbackQuery.Data
			eventData := update.CallbackQuery.Data
			separatedData := strings.Split(eventName, ":")
			if len(separatedData) > 1 {
				eventName = separatedData[0]
				eventData = separatedData[1]
			}
			if commandDeferred, ok := d.GetCommand(eventName); ok {
				if d.RunCommand(
					commandDeferred,
					update.CallbackQuery.Message,
					commandDeferred.Command,
					eventData,
					strings.Split(update.CallbackQuery.Data, ":")) {
					break
				}
			}
			break
		}
		zlog.Info().
			Any("is command", update.Message.IsCommand()).
			Any("command", update.Message.Command()).
			Any("CommandWithAt", update.Message.CommandWithAt()).
			Any("Arguments", update.Message.CommandArguments()).
			Send()

		if update.Message.LeftChatMember != nil {
			go d.RunEvents(tgCommands.UserLeaveChantEvent.String(), update.Message, "", "")
		}

		if update.Message.NewChatMembers != nil {
			go d.RunEvents(tgCommands.UserJoinedChantEvent.String(), update.Message, "", "")
		}

		if update.EditedMessage != nil {
			zlog.Info().Any("update.EditedMessage", update.EditedMessage).Send()
		}
		msg := update.Message
		if msg == nil || (msg == nil && update.InlineQuery != nil) {
			zlog.Info().Any("nill MSG update", update).Send()
			continue
		}

		isCommand = msg.IsCommand()
		if isCommand {
			commandName = msg.Command()
		}
		if commandName == "start" {
			go d.RunEvents(tgCommands.StartBotEvent.String(), msg, "start", msg.CommandArguments())
		}
		if msg.Text != "" {
			go d.RunEvents(tgCommands.TextMsgBotEvent.String(), msg, "", "")
		}

		zlog.Info().Any("update FULL", update).Send() //TODO: MOVE TO DEBUG MODE

		if founded := d.CheckDeferred(update.Message.From.ID); founded != "" {
			if commandDeferred, ok := d.GetCommand(founded); ok {
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
		command, founded := d.GetCommand(commandName)

		if founded {
			sended = true
			log.Println("run founded Command")
			d.RunCommand(command, update.Message, command.Command, msg.CommandArguments(), nil)
		} else {
			for _, command := range d.Commands {
				if !command.Permission(update.Message, d.AdminId) || command.Handler == nil {
					continue
				}
				splitCommands, commandValue := splitCommand(update.Message.Text, " ")
				if len(splitCommands) == 0 {
					continue
				}
				//commandName := splitCommands[0]
				commandsCount := len(splitCommands)
				if commandsCount == 0 {
					continue
				}
				if !command.IsImplemented(commandName, d.Bot.Self.UserName) {
					if command.IsMatched(update.Message.Text, d.Bot.Self.UserName) {
						commandValue = update.Message.Text
					} else {
						log.Println("!IsMatched")
						continue
					}
				}
				sended = true
				log.Println("just RunCommand")
				if d.RunCommand(command, update.Message, command.Command, commandValue, splitCommands) {
					break
				}
			}
		}

		if d.DefaultCommand != "" && !sended {
			fmt.Println("default command:", d.DefaultCommand)
			if commandDeferred, ok := d.GetCommand(d.DefaultCommand); ok {
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

func (d *Data) GetSubCommands(subName string) []tgCommands.Command {
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

func (d *Data) GetCommand(name string) (tgCommands.Command, bool) {
	d.mutexCommands.Lock()
	item, ok := d.Commands[name]
	d.mutexCommands.Unlock()
	return item, ok
}

func (d *Data) AddCommands(newItems tgCommands.Commands) {
	d.mutexCommands.Lock()
	d.Commands = d.Commands.Merge(newItems)
	d.mutexCommands.Unlock()
}

func (d *Data) RunEvents(event string, msg *tgbotapi.Message, commandName string, param string) {
	eventCommands := d.GetSubCommands(event)
	log.Println("RunEvents", eventCommands)
	for _, eventCommand := range eventCommands {
		log.Println("tg event RunCommand", eventCommand)
		d.RunCommand(eventCommand, msg, commandName, param, nil)
	}
}
