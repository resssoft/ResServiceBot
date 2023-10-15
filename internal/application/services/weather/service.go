package weather

import (
	"context"
	"encoding/json"
	"fmt"
	tgModel "fun-coice/internal/domain/commands/tg"
	"fun-coice/pkg/scribble"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	zlog "github.com/rs/zerolog/log"
)

var _ = (tgModel.Service)(&data{})

type data struct {
	events  tgModel.Commands
	sentMsg tgModel.SentMessages
	DB      *scribble.Driver
	items   []item
	ctx     context.Context
	tokens  map[string]string
}

func New(sentMsg tgModel.SentMessages, DB *scribble.Driver, tokens map[string]string) tgModel.Service {
	result := data{
		sentMsg: sentMsg,
		DB:      DB,
		ctx:     context.Background(),
		tokens:  tokens,
	}
	commandsList := tgModel.NewCommands()
	commandsList.AddSimple("weather_add_chat", "Add weather notifier to chat", result.addWeatherChat)
	commandsList.AddSimple("weather_show", "Show weather", result.showWeatherToChat)
	commandsList.AddSimple("weather_test", "Show weather", result.showWeatherChatEvents)
	result.Configure()
	go result.worker()
	result.events = commandsList
	return &result
}

func (d *data) Commands() tgModel.Commands {
	return d.events
}

func (d data) Name() string {
	return "weather"
}

func (d *data) addWeatherChat(msg *tgbotapi.Message, command *tgModel.Command) *tgModel.HandlerResult {
	params := strings.Split(command.Arguments.Raw, " ")
	chatId := fmt.Sprintf("%v", msg.Chat.ID)
	hour := time.Now().Hour()
	min := time.Now().Minute()
	var err error
	if len(params) == 1 {
		partsTime := strings.Split(command.Arguments.Raw, ":")
		if len(partsTime) == 2 {
			hour, err = strconv.Atoi(partsTime[0])
			if err != nil {
				fmt.Println("cant parse hour", partsTime[0], err)
			}
			min, err = strconv.Atoi(partsTime[1])
			if err != nil {
				fmt.Println("cant parse min", partsTime[1], err)
			}
		}
	}
	weatherItem := fmt.Sprintf("%v#%v:%v", chatId, hour, min)
	fmt.Println("add weather item", weatherItem)
	d.items = append(d.items, item{
		chatId: msg.Chat.ID,
		hour:   hour,
		min:    min,
	})
	if err := d.DB.Write("weather_chats", chatId, weatherItem); err != nil {
		fmt.Println("add command error", err)
		return tgModel.Simple(msg.Chat.ID, "Cant do that, sorry")
	} else {
		return tgModel.Simple(msg.Chat.ID, "saved!")
	}
}

func (d *data) showWeatherChatEvents(msg *tgbotapi.Message, command *tgModel.Command) *tgModel.HandlerResult {
	return tgModel.Simple(msg.Chat.ID, "qq") //////////////////////////////
}

func (d *data) showWeatherToChat(msg *tgbotapi.Message, command *tgModel.Command) *tgModel.HandlerResult {
	return tgModel.Simple(msg.Chat.ID, d.GetGismeteoForecast()) //////////////////////////////
}

func (d *data) delWeatherChatEvents(msg *tgbotapi.Message, command *tgModel.Command) *tgModel.HandlerResult {
	return tgModel.Simple(msg.Chat.ID, "qq") //////////////////////////////
}

func (d *data) Configure() {
	records, err := d.DB.ReadAll("weather_chats")
	if err != nil {
		fmt.Println("db read error", err)
		return
	}
	if len(records) == 0 {
		fmt.Println("zero saved items found")
	}
	//parse item from simple db value
	var chatId int64
	hour := 0
	min := 0
	for _, record := range records {
		//item format: 3563463456#12:05 chatId#hour:minute
		chatId = 0
		hour = 0
		min = 0
		parts := strings.Split(record, "#")
		if len(parts) != 2 {
			continue
		}
		chatId, err = strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			fmt.Println("cant parse chatId", parts[0], err)
			continue
		}
		partsTime := strings.Split(parts[0], ":")
		if len(partsTime) != 2 {
			continue
		}
		hour, err = strconv.Atoi(partsTime[0])
		if err != nil {
			fmt.Println("cant parse hour", partsTime[0], err)
			continue
		}
		min, err = strconv.Atoi(partsTime[1])
		if err != nil {
			fmt.Println("cant parse min", partsTime[1], err)
			continue
		}
		d.items = append(d.items, item{
			chatId: chatId,
			hour:   hour,
			min:    min,
		})
	}
}

func (d *data) worker() {
	for {
		select {
		case <-d.ctx.Done():
			return
		case <-time.NewTimer(checkItemsDelay).C:
			fmt.Println("checkItemsDelay")
			hour := time.Now().Hour()
			min := time.Now().Minute()
			changed := false
			var removedItems []int
			for index, itemVal := range d.items {
				if itemVal.hour == hour && itemVal.min == min {
					fmt.Println("!!!!!!!!!!!!")
					changed = true
					removedItems = append(removedItems, index)
				}
			}
			if changed {
				var newItems []item
				for _, remIndex := range removedItems {
					for index, itemVal := range d.items {
						if index != remIndex {
							newItems = append(newItems, itemVal)
						}
					}
				}
				d.items = newItems
			}
		}
	}
}

func (d *data) SendWeather(chatId int64) {
	d.sentMsg <- tgModel.Simple(chatId, "User Leave Chant:\n").Messages[0]
}

func (d *data) GetGismeteoForecast() string {
	client := &http.Client{}
	req, err := http.NewRequest("GET", gismeteoUrl, nil)
	req.Header.Set("X-Gismeteo-Token", "6414437ae020d0.26060480")

	if err != nil {
		fmt.Println(err)
	}
	res, err := client.Do(req)
	if res.Body != nil {
		defer res.Body.Close()
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	fmt.Println(string(body))
	var resultObj gismeteoData
	err = json.Unmarshal(body, &resultObj)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	zlog.Info().Any("resultObj", resultObj).Send()
	result := ""
	for _, dayItem := range resultObj.Response {
		icon := ""
		if val, ok := gismeteoIcons[dayItem.Icon]; ok {
			icon = val
		}
		result += fmt.Sprintf("[%s] %s %.2f°C - %.2f°C, %s \n",
			dayItem.Date.UTC,
			icon,
			dayItem.Temperature.Air.Min.C,
			dayItem.Temperature.Air.Max.C,
			dayItem.Description.Full)
	}
	fmt.Println("======================6", result)
	return result
}
