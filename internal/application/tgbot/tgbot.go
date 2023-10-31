package tgbot

import (
	"fmt"
	"fun-coice/config"
	tgModel "fun-coice/internal/domain/commands/tg"
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
	defaultWorkersCount      = 10
	messagesChanLimit        = 300
	CommandsHandlerChanLimit = 100
	messageTimeLimit         = time.Millisecond * 100
)

type Data struct {
	WebUri         string
	Token          string
	StartMsg       bool
	WebMode        bool
	Commands       tgModel.Commands
	Bot            *tgbotapi.BotAPI
	WorkersCount   int
	Deferred       map[int64]Deferred
	mutexDeferred  *sync.Mutex
	DefaultCommand string
	messagesChan   chan tgbotapi.Chattable
	commandResults chan *tgModel.HandlerResult
	lastMsgTime    time.Time
	Name           string
	AdminId        int64
	AdminLogin     string
	mutexCommands  *sync.Mutex
	Ran            bool
	config         config.TgBotConfig
}

type Deferred struct {
	Command string
	Data    string
	Message *tgbotapi.Message
}

func New(name string, botConfig config.TgBotConfig) (*Data, error) {
	bot, err := tgbotapi.NewBotAPI(config.TelegramToken(name))
	if err != nil {
		return nil, err
	}
	//TODO: replace config.TelegramToken(name) to botConfig

	zlog.Info().
		Any("AdminLogin", botConfig.AdminLogin).
		Any("AdminId", botConfig.AdminId).
		Any("Login", botConfig.Login).
		Any("DefaultCommand", botConfig.DefaultCommand).
		Any("bot.Self.UserName", bot.Self.UserName).
		Any("bot.Self", bot.Self).
		Any("WebMode", botConfig.WebMode).
		Any("Uri", botConfig.Uri).
		Any("Token", botConfig.Token).
		Send() // TODO temporary

	return &Data{
		config:         botConfig,
		Token:          botConfig.Token,
		WebUri:         botConfig.Uri,
		Commands:       defaultCommands,
		Bot:            bot,
		WorkersCount:   defaultWorkersCount,
		StartMsg:       true,
		WebMode:        botConfig.WebMode,
		Name:           bot.Self.UserName,
		Deferred:       make(map[int64]Deferred),
		mutexDeferred:  &sync.Mutex{},
		DefaultCommand: botConfig.DefaultCommand,
		AdminId:        botConfig.AdminId,
		AdminLogin:     botConfig.AdminLogin,
		messagesChan:   make(chan tgbotapi.Chattable, messagesChanLimit),
		commandResults: make(chan *tgModel.HandlerResult, CommandsHandlerChanLimit),
		mutexCommands:  &sync.Mutex{},
		Ran:            false,
	}, nil
}

func (d *Data) setDefaults() {
	d.Commands.AddSimple("about", "About bot", d.about, "help")
	d.Commands.AddSimple("admin", "Bot admin info", d.admin, "админ", "кто админ")
	d.Commands.AddSimple("commands", "Show bot commands", d.commandsList, "список комманд", "команды")
}

func (d *Data) Run() error {
	//d.Bot.Debug = true
	//TODO: d.Bot.GetMyCommands() AND SET THEM
	d.setDefaults()

	startMsg := "-"
	defer func() {
		startMsg = fmt.Sprintf(
			"Bot %s Started with version %s \n webMode [%v]\nDefault Command: %s\nUri: %s",
			d.Bot.Self.UserName,
			appStat.Version,
			d.WebMode,
			d.DefaultCommand,
			d.WebUri)
		if d.StartMsg && d.AdminId != 0 { //ANOTHER PARAM
			msg := tgbotapi.NewMessage(d.AdminId, startMsg)
			d.Bot.Send(msg)
		}
	}()

	log.Printf("Authorized bot on account %s \n %s", d.Bot.Self.UserName, startMsg)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 10

	for i := 0; i < d.WorkersCount; i++ {
		go d.MessagesHandler()
	}

	whInfo, err := d.Bot.GetWebhookInfo()
	if err != nil {
		return err
	}
	if d.WebMode {
		if whInfo.URL == "" {
			if config.WebServerDomain() != "" {
				w, err := tgbotapi.NewWebhook("https://" + config.WebServerDomain() + d.WebUri)
				if err != nil {
					return err
				}
				apiResp, err := d.Bot.Request(w)
				if err != nil {
					zlog.Info().Any("apiResp", apiResp).Send()
					return err
				}
				zlog.Info().Msg("WebMode: set webHook: " + "https://" + config.WebServerDomain() + d.WebUri)
			} else {
				d.WebMode = false
				zlog.Info().Msg("WebMode set to false, URL is incorrect")
			}
		} else {
			//TODO: check for url from bot and from config is equal
		}
	}

	if d.WebMode {
		fmt.Println("tg bot WebMode", d.Bot.Self.UserName, d.WebUri)
		webUpdates := d.Bot.ListenForWebhook(d.WebUri)
		for i := 0; i < d.WorkersCount; i++ {
			go d.UpdatesHandler(webUpdates, strconv.Itoa(i)+" "+d.Bot.Self.UserName)
		}
	} else { //web updates mode
		if whInfo.URL != "" {
			zlog.Info().Any("found bot webhook will be removed", whInfo.URL).Send()
			w, err := tgbotapi.NewWebhook("")
			if err != nil {
				return err
			}
			apiResp, err := d.Bot.Request(w)
			if err != nil {
				zlog.Info().Any("apiResp reset webhook", apiResp).Send()
				return err
			}
			startMsg += "\n removed webHook"
		}
		fmt.Println("tg bot UpdateMode", d.Bot.Self.UserName)
		updates := d.Bot.GetUpdatesChan(u)
		for i := 0; i < d.WorkersCount; i++ {
			go d.UpdatesHandler(updates, strconv.Itoa(i)+" "+d.Bot.Self.UserName)
		}
	}
	for i := 0; i < d.WorkersCount; i++ {
		go d.commandsHandler()
	}
	d.Ran = true
	return nil
}

func (d *Data) commandsHandler() {
	for result := range d.commandResults {
		d.SendCommandResult(result, nil)
	}
}

func (d *Data) SetWebMode(val bool) {
	d.WebMode = val
}

func (d *Data) AppendDeferred(user int64, command, data string, msg *tgbotapi.Message) {
	d.mutexDeferred.Lock()
	defer d.mutexDeferred.Unlock()
	d.Deferred[user] = Deferred{
		Command: command,
		Data:    data,
		Message: msg,
	}
}

func (d *Data) CheckDeferred(user int64) Deferred {
	d.mutexDeferred.Lock()
	defer d.mutexDeferred.Unlock()
	//fmt.Println(user, len(d.Deferred)) //DEVMODE
	founded := d.Deferred[user]
	delete(d.Deferred, user)
	return founded
}

func (d *Data) RunCommand(command tgModel.Command, msg *tgbotapi.Message) bool {
	if command.Arguments.Raw == "" {
		command.SetArgs(msg.CommandArguments())
	}
	command.FilesCallback = d.getTgFile
	command.ParamCallback = d.getParam

	command.BotName = d.Name                                              // TODO: check if set is needle (bot local name or login)
	zlog.Info().Any("RunCommand command.Command", command.Command).Send() // WHy EMPTY?
	result := command.Handler(msg, &command)
	zlog.Info().Any("result handler", result).Send()
	return d.SendCommandResult(result, msg)
}

func (d *Data) SendCommandResult(result *tgModel.HandlerResult, msg *tgbotapi.Message) bool {
	fmt.Println("SendCommandResult")
	if result.Prepared {
		//fmt.Println("COMMAND PREPAERD") //DEVMODE
		log.Println("result.Messages", len(result.Messages))
		for _, chantEvent := range result.Messages {
			//log.Println("chatEvent", chantEvent)

			msgRes, err := d.Bot.Send(chantEvent) //TODO: check limits by this package method
			if err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Println("Send message", msgRes.MessageID)
			}
		}
		//set differ with event
		if len(result.Events) > 0 && msg != nil {
			log.Println("len(result.Events) > 0")
			for _, chatEvent := range result.Events {
				eventCommands := d.GetSubCommands(chatEvent.Name)
				zlog.Info().Any("result eventCommands", eventCommands).Send()
				for _, eventCommand := range eventCommands {
					log.Println("eventCommands", eventCommand)
					eventCommand.SetArgs(msg.CommandArguments())
					if chatEvent.Msg != nil {
						log.Println("event provide msg", chatEvent.Msg)
						msg = chatEvent.Msg
					}
					d.RunCommand(eventCommand, msg)
				}
			}
		}
		if !result.Deferred {
			return true
		}
	}
	if result.Redirect != nil && msg != nil {
		//TODO: check redirect step limit
		zlog.Info().Any("redirect", result.Redirect).Send()
		redirectCommand, founded := d.GetCommand(result.Redirect.CommandName)
		if founded {
			redirectMsg := msg
			if result.Redirect.Message != nil {
				redirectMsg = result.Redirect.Message
			}
			d.RunCommand(redirectCommand, redirectMsg)
			return true
		}
	} else {
		zlog.Info().Msg("Empty redirect")
	}
	if result.Deferred && msg != nil {
		defBy := msg.From.ID
		if msg.From.IsBot {
			defBy = msg.Chat.ID
		}
		fmt.Println("COMMAND Wait", defBy, result.Next) //DEVMODE
		zlog.Info().Any("msg resend", result.Resend).Send()
		d.AppendDeferred(defBy, result.Next, result.Data, result.Resend)
		return true
	}
	return false
}

func (d *Data) MessagesHandler() {
	for msg := range d.messagesChan {
		msgRes, err := d.Bot.Send(msg)
		if err != nil {
			fmt.Println("Send tg message error", err.Error())
		} else {
			fmt.Println("Sent message ", msgRes.MessageID)
		}
	}
}

func (d *Data) UpdatesHandler(updates tgbotapi.UpdatesChannel, workerID string) {
	//log.Println("start worker UpdatesHandler", workerID) // TODO: to debug
	isCommand := false
	commandName := ""
	for update := range updates {
		d.mutexDeferred.Lock() ////////////////////////
		d.mutexDeferred.Unlock()
		zlog.Info().Any("d.Deferred", d.Deferred).Send() //////////////////

		commandName = ""
		log.Println("update chan EVENT", update.UpdateID, workerID) // TODO: to debug
		zlog.Debug().Any("update", update).Send()
		//TODO: add sent messages history - save users by bot
		//TODO: save bot stat - new users, by date; new messages

		//only for debug level
		//zlog.Debug().Any("msg", update.Message).Any("InlineQuery", update.InlineQuery).Send()
		if update.Message != nil {
			if update.Message.Chat.Type == "private" {
				//zlog.Info().Any("msg", update.Message).Any("InlineQuery", update.InlineQuery).Send()
			}
		}
		//zlog.Info().Any("msg", update.Message).Any("InlineQuery", update.InlineQuery).Send() ///////

		if update.CallbackQuery != nil {
			//update.CallbackQuery.ID
			//http.Get("https://api.telegram.org/bot5761609803:AAFZT_tCSxnVmRmZ1qQXmUtFR4NiQTFecPE/answerCallbackQuery?callback_query_id=" + update.CallbackQuery.ID + "&text=done&show_alert=false")
			//d.Bot.StopPoll()

			//tgbotapi.NewCallback(update.CallbackQuery.ID, "Готово!")
			if update.CallbackQuery.Message.Chat.Type == "private" {
				zlog.Info().Any("callback query data", update.CallbackQuery.Data).Send()
			}
			eventName := update.CallbackQuery.Data
			//eventData := update.CallbackQuery.Data
			separatedData := strings.Split(eventName, ":")
			if len(separatedData) > 1 {
				eventName = separatedData[1]
				//eventData = separatedData[1] // ????????????????????????
			}
			if commandDeferred, ok := d.GetCommand(eventName); ok {

				commandDeferred.SetArgs(update.CallbackQuery.Message.CommandArguments())
				commandDeferred.Data = update.CallbackQuery.Data
				if d.RunCommand(
					commandDeferred,
					update.CallbackQuery.Message) {
					break
				}
			}
			break
		}
		/*
			zlog.Info().
				Any("is command", update.Message.IsCommand()). ////////////
				Any("command", update.Message.Command()).
				Any("CommandWithAt", update.Message.CommandWithAt()).
				Any("Arguments", update.Message.CommandArguments()).
				Send()
		*/

		if update.EditedMessage != nil {
			zlog.Info().Any("update.EditedMessage", update.EditedMessage).Send()
		}
		msg := update.Message
		if msg == nil || (msg == nil && update.InlineQuery != nil) {
			zlog.Info().Any("nil MSG update", update).Send()
			continue
		}

		if msg.LeftChatMember != nil {
			go d.RunEvents(tgModel.UserLeaveChantEvent, msg, new(tgModel.Command))
		}

		if msg.NewChatMembers != nil {
			go d.RunEvents(tgModel.UserJoinedChantEvent, msg, new(tgModel.Command))
		}

		zlog.Debug().
			Any("IsCommand", msg.IsCommand()).
			Any("Command", msg.Command()).
			Any("CommandArguments", msg.CommandArguments()).
			Any("msg", msg.Text).
			Any("IsCommand", msg.CommandWithAt()).
			Send()

		isCommand = msg.IsCommand()
		if isCommand {
			commandName = msg.Command()
		}
		if commandName == "start" {
			startCommand := tgModel.Command{
				Command: "start",
				Arguments: tgModel.CommandArguments{
					Raw: msg.CommandArguments(),
				},
			}
			go d.RunEvents(tgModel.StartBotEvent, msg, &startCommand)
		}
		if msg.Text != "" {
			go d.RunEvents(tgModel.TextMsgBotEvent, msg, new(tgModel.Command))
		}

		zlog.Info().Any("update FULL", update).Send() //TODO: MOVE TO DEBUG MODE

		//check waited commands
		if founded := d.CheckDeferred(msg.From.ID); founded.Command != "" {
			log.Println("Deferred RunCommand", founded.Command)
			if commandDeferred, ok := d.GetCommand(founded.Command); ok {
				deferredMsg := msg
				log.Println("Deferred RunCommand")
				commandDeferred.SetArgs(msg.Text)
				//zlog.Info().Any("commandDeferred", commandDeferred.Arguments).Send()
				zlog.Info().
					Any("commandDeferred msg", msg).
					Any("founded.Message", founded.Message).
					Any("command", commandDeferred.Command).
					Any("Arguments", commandDeferred.Arguments).
					Send()
				if founded.Message != nil {
					deferredMsg = founded.Message
				}
				commandDeferred.Data = founded.Data
				commandDeferred.Deferred = true
				if d.RunCommand(
					commandDeferred,
					deferredMsg) {
					continue
				}
			}
		}

		sent := false
		command, founded := d.GetCommand(commandName)
		if founded {
			sent = true
			log.Println("run founded Command by: " + commandName)
			command.SetArgs(msg.CommandArguments())
			d.RunCommand(command, msg)
		} else {
			for _, command := range d.Commands {
				if !command.Permission(msg, d.AdminId) || command.Handler == nil {
					continue
				}
				splitCommands, commandValue := splitCommand(msg.Text, " ")
				if len(splitCommands) == 0 {
					continue
				}
				//commandName := splitCommands[0]
				commandsCount := len(splitCommands)
				if commandsCount == 0 {
					continue
				}
				if !command.IsImplemented(commandName, d.Bot.Self.UserName) {
					if command.IsMatched(msg.Text, d.Bot.Self.UserName) {
						commandValue = msg.Text
					} else {
						//log.Println("!IsMatched")
						continue
					}
				}
				sent = true
				log.Println("just RunCommand")
				command.SetArgs(commandValue)
				if d.RunCommand(command, msg) {
					break
				}
			}
		}

		if d.DefaultCommand != "" && !sent {
			fmt.Println("default command:", d.DefaultCommand)
			if commandDeferred, ok := d.GetCommand(d.DefaultCommand); ok {
				log.Println("default RunCommand")
				commandDeferred.SetArgs(msg.Text)
				if d.RunCommand(commandDeferred, msg) {
					continue
				}
			}
			continue
		}

		var foundedCommands []tgModel.Command
		//check for sent files without commands
		if len(msg.Photo) > 0 {
			foundedCommands = d.GetFileCommands("Photo")
		}
		if msg.Audio != nil {
			foundedCommands = d.GetFileCommands("Audio")
		}
		if msg.Video != nil {
			foundedCommands = d.GetFileCommands("Video")
		}
		if msg.Venue != nil {
			foundedCommands = d.GetFileCommands("Venue")
		}
		if msg.Voice != nil {
			foundedCommands = d.GetFileCommands("Voice")
		}
		if msg.Sticker != nil {
			foundedCommands = d.GetFileCommands("Sticker")
		}
		if msg.Animation != nil {
			foundedCommands = d.GetFileCommands("Animation")
		}
		if msg.MediaGroupID != "" {
			foundedCommands = d.GetFileCommands("MediaGroupID")
		}
		if msg.VideoNote != nil {
			foundedCommands = d.GetFileCommands("VideoNote")
		}
		if msg.Poll != nil {
			foundedCommands = d.GetFileCommands("Poll")
		}

		if msg.Document != nil {
			foundedCommands = d.GetFileCommands(msg.Document.MimeType)
			if len(foundedCommands) == 0 {
				foundedCommands = d.GetFileCommands("Document")
			}
		}
		if len(foundedCommands) > 0 {
			if len(foundedCommands) == 1 {
				if d.RunCommand(foundedCommands[0], update.Message) {
					continue
				}
			}
			mstText := "Choice command:\n"
			for _, commandItem := range foundedCommands {
				mstText += "/" + commandItem.Command + "\n"
			}
			newCommand := tgModel.NewCommand().WithHandler(
				func(message *tgbotapi.Message, command *tgModel.Command) *tgModel.HandlerResult {
					return tgModel.DeferredWithText(msg.Chat.ID, mstText, commandRedirect, "", message)
				})
			if d.RunCommand(*newCommand, msg) {
				continue
			}
			//d.AppendDeferred(msg.Chat.ID, commandChoicer, msg)
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

func (d *Data) GetSubCommands(subName string) []tgModel.Command {
	var founded []tgModel.Command
	for key, command := range d.Commands { //TODO: separate events and commands by bot vars
		if command.IsEvent && subName == command.Command {
			founded = append(founded, command)
			continue
		}

		results := strings.Split(key, ":") //deprecated
		if len(results) < 2 {
			continue
		}
		if results[1] == subName {
			founded = append(founded, command)
		}
	}

	return founded
}

func (d *Data) GetFileCommands(fileType string) []tgModel.Command {
	if fileType == "" {
		return nil
	}
	fileType = strings.ToLower(fileType)
	var founded []tgModel.Command
	for _, command := range d.Commands {
		if command.FileTypes.Has(fileType) {
			founded = append(founded, command)
		}
	}
	return founded
}

func (d *Data) GetCommand(name string) (tgModel.Command, bool) {
	d.mutexCommands.Lock()
	item, ok := d.Commands[name]
	d.mutexCommands.Unlock()
	return item, ok
}

func (d *Data) AddCommands(newItems tgModel.Commands, serviceName string) {
	newItems.SetBotData(d.Bot.Self.UserName, serviceName)
	d.mutexCommands.Lock()
	d.Commands = d.Commands.Merge(newItems)
	d.mutexCommands.Unlock()
}

func (d *Data) RunEvents(event string, msg *tgbotapi.Message, command *tgModel.Command) {
	eventCommands := d.GetSubCommands(event)
	log.Println("RunEvents", eventCommands)
	for _, eventCommand := range eventCommands {
		log.Println("tg event", eventCommand.Command)
		eventCommand.SetArgs(command.Arguments.Raw)
		d.RunCommand(eventCommand, msg)
	}
}

func (d *Data) getParam(name tgModel.BotParamRequest) tgModel.BotParamResponse {
	switch name {
	case tgModel.BotNameParam:
		return tgModel.BotParamStr(d.Name)
	case tgModel.BotAdminLoginParam:
		return tgModel.BotParamStr(d.AdminLogin)
	case tgModel.BotNAdminIdParam:
		return tgModel.BotParamInt64(d.AdminId)
	default:
		return tgModel.BotParamNotFound()
	}
}

func (d *Data) PushMessage() chan<- tgbotapi.Chattable {
	return d.messagesChan
}

func (d *Data) PushHandleResult() chan<- *tgModel.HandlerResult {
	return d.commandResults
}

func (d *Data) BotName() string {
	return d.Name
}
