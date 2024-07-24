package markovka

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	tgModel "fun-coice/internal/domain/commands/tg"
	tgbotapi "fun-coice/pkg/telegram-bot-api"
)

type data struct {
	list          tgModel.Commands
	messageSender tgModel.MessageSender
	paramsByBot   map[string]Params
	states        map[string]States
	durations     map[int64]int
	firstWords    []string
}

func New() tgModel.Service {
	serv := data{
		list:        tgModel.NewCommands(),
		paramsByBot: make(map[string]Params),
		states:      make(map[string]States),
		durations:   make(map[int64]int),
	}

	tgModel.NewCommand().
		Simple("markovka", "emoji commands", serv.SendGenerated).
		Push(serv.list)

	tgModel.AdminCommand().
		Simple("markovkaStat", "Show service stats", serv.stat).
		Push(serv.list)
	tgModel.NewCommand().
		Simple("markovkaCommands", "emoji commands", serv.commandsList).
		Push(serv.list)

	return &serv
}

func (d *data) Commands() tgModel.Commands {
	return d.list
}

func (d *data) Name() string {
	return "markovka"
}

func (d *data) Destroy() {}

func (d *data) Dependency() *tgModel.ServiceDepends {
	return tgModel.ServiceDependsIs(tgModel.FileDbDependency, tgModel.PluginsDependency)
}

func (d *data) Configure(sc tgModel.ServiceConfig) error {
	if sc.FileDb == nil {
		return fmt.Errorf("file db is nil")
	}
	if sc.Params == nil {
		return fmt.Errorf("params is nil")
	}
	params := Params{
		BotName: sc.MessageSender.BotName(),
		Chats:   make(map[int64]ChatData),
	}
	for key, value := range sc.Params {
		switch key {
		case "path":
			params.DataPath = value
		case "chats":
			chatsSettings := strings.Split(value, ",")
			for _, chatsSetting := range chatsSettings {
				chatsSetting = strings.TrimSpace(chatsSetting)
				chatSettings := strings.Split(chatsSetting, ":")
				if len(chatSettings) == 2 {
					chatId, _ := strconv.ParseInt(chatSettings[0], 10, 64)
					param, _ := strconv.Atoi(chatSettings[1])
					params.Chats[chatId] = ChatData{
						Id:    chatId,
						Param: param,
					}
					d.durations[chatId] = param
				}
			}
			params.DataPath = value
		}
	}
	d.paramsByBot[sc.MessageSender.BotName()] = params
	d.messageSender = sc.MessageSender
	if params.DataPath == "" {
		fmt.Println("markovka: no path set")
		return nil
	}
	fmt.Println("markovka: load states in the apth" + params.DataPath)
	dirEntry, err := os.ReadDir(params.DataPath)
	if err != nil {
		fmt.Println("markovka: error reading directory", err.Error())
		return nil
	}
	for _, entry := range dirEntry {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}
		fmt.Println("markovka: load states", params.DataPath+entry.Name())
		statesData, err := os.ReadFile(params.DataPath + entry.Name())
		if err != nil || len(statesData) == 0 {
			fmt.Println("markovka: error read states", err)
			return nil
		}
		var states []State
		err = json.Unmarshal(statesData, &states)
		if err != nil {
			fmt.Println("markovka: error unmarshalling states", err.Error())
			return nil
		}
		d.states[sc.MessageSender.BotName()] = States{
			States: states,
			Name:   entry.Name(),
		}
	}
	if len(d.states) == 0 {
		fmt.Println("markovka: empty states")
	}
	fmt.Println("markovka: load states", "first_words.txt")
	fWordsByte, err := os.ReadFile(params.DataPath + "first_words.txt")
	if err == nil || len(fWordsByte) != 0 {
		d.firstWords = strings.Split(string(fWordsByte), " ")
	}
	if len(d.firstWords) == 0 {
		fmt.Println("markovka: empty firstWords")
		return nil
	}
	for chatId, _ := range params.Chats {
		fmt.Println("markovka: start worker for chat", chatId)
		go d.worker(sc.MessageSender.BotName(), chatId)
	}
	return nil
}

func (d *data) worker(botName string, chatId int64) {
	for {
		select {
		case <-time.After(time.Duration(d.durations[chatId]) * time.Second):
			d.messageSender.PushHandleResult() <- tgModel.Simple(chatId, d.generateText(botName, chatId))
		}
	}
}

func (d *data) randomFirstWord() string {
	if d.firstWords == nil {
		return ""
	}
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	randomNum := r.Intn(len(d.firstWords))
	result := d.firstWords[randomNum]
	return result
}

func (d *data) stat(msg *tgbotapi.Message, command *tgModel.Command) *tgModel.HandlerResult {
	b64result := base64.StdEncoding.EncodeToString([]byte(command.Arguments.Raw))
	return tgModel.SimpleReply(msg.Chat.ID, b64result, msg.MessageID)
}

func (d *data) commandsList(msg *tgbotapi.Message, command *tgModel.Command) *tgModel.HandlerResult {
	list := d.Commands().Available(msg, command.Bot.AdminId)
	commandsList := "Commands:\n"
	for key, item := range list {
		commandsList += "/" + key + " - " + item.Description + "\n"
	}
	return tgModel.Simple(msg.Chat.ID, commandsList)
}

func (d *data) SendGenerated(msg *tgbotapi.Message, command *tgModel.Command) *tgModel.HandlerResult {
	generated := d.generateText(command.Bot.Login, msg.Chat.ID)
	if generated == "" {
		generated = d.generateText(command.Bot.Login, msg.Chat.ID)
	}
	if generated == "" {
		return tgModel.EmptyCommand()
	}
	return tgModel.SimpleReply(msg.Chat.ID, generated, msg.MessageID) //"["+randomWord+"]"+
}

func (d *data) generateText(botName string, chatId int64) string {
	if len(d.states[botName].States) == 0 {
		fmt.Println("start train")
		//states = markov.train(fullData)
	}
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	randomNum := r.Intn(15) + 6
	randomWord := d.randomFirstWord()
	generatedText := markov.generateText(d.states[botName].States, randomWord, randomNum)
	generatedText = strings.ReplaceAll(generatedText, "\n", ", ")
	return generatedText //"["+randomWord+"]"+
}
