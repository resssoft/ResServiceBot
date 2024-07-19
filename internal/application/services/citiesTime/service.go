package citiesTime

import (
	"fmt"
	"sort"
	"strings"
	"time"

	tgModel "fun-coice/internal/domain/commands/tg"
	tgbotapi "fun-coice/pkg/telegram-bot-api"
)

type data struct {
	list          tgModel.Commands
	messageSender tgModel.MessageSender
	cities        map[string]string
}

type formattedTime struct {
	UTCSort int
	UTC     string
	City    string
	Time    string
}

type formattedTimes []formattedTime

// Реализация интерфейса sort.Interface для типа ByAge
func (a formattedTimes) Len() int           { return len(a) }
func (a formattedTimes) Less(i, j int) bool { return a[i].UTCSort < a[j].UTCSort }
func (a formattedTimes) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func New() tgModel.Service {
	s := data{
		cities: make(map[string]string),
	}
	s.cities = map[string]string{
		"Сакраменто, США": "America/Los_Angeles",
		//"Варшава, Польша":       "Europe/Warsaw",
		//"Амстердам, Нидерланды": "Europe/Amsterdam",
		//"Белград, Сербия":       "Europe/Belgrade",
		//"Подгорица, Черногория": "Europe/Podgorica",
		"Европа":         "Europe/Podgorica",
		"Москва, Россия": "Europe/Moscow",
		//"Ижевск, Россия":        "Europe/Samara", // Ижевск находится в часовом поясе Самары
		//"Ереван, Армения":       "Asia/Yerevan",
		"Ереван, Ижевск": "Asia/Yerevan",
		"Пермь, Россия":  "Asia/Yekaterinburg",
		//"Новосибирск, Россия":   "Asia/Novosibirsk",
		//"Красноярск, Россия":    "Asia/Krasnoyarsk",
		"Красноярск, Новосибирск": "Asia/Krasnoyarsk",
	}
	commandsList := tgModel.NewCommands()
	tgModel.NewCommand().
		Simple("time", "Show cities time", s.printTimeList).
		WithTriggers("время").
		Push(commandsList)
	tgModel.NewCommand().
		Simple("timeAdd", "Show cities time", s.addCity).
		Push(commandsList)
	tgModel.NewCommand().
		Simple("timeDel", "Show cities time", s.delCity).
		Push(commandsList)

	s.list = commandsList
	return &s
}

func (d *data) Commands() tgModel.Commands {
	return d.list
}

func (d *data) Name() string {
	return "citiesTime"
}

func (d *data) Destroy() {}

func (d *data) Dependency() *tgModel.ServiceDepends {
	return nil
}

func (d *data) Configure(botData tgModel.ServiceConfig) error {
	d.messageSender = botData.MessageSender
	return nil
}

func (d *data) addCity(msg *tgbotapi.Message, command *tgModel.Command) *tgModel.HandlerResult {
	items := strings.Split(command.Arguments.Raw, ":")
	if len(items) < 2 {
		return tgModel.SimpleReply(msg.Chat.ID, "Incorrect! Format: Samara:Europe/Samara", msg.MessageID)
	}
	_, err := time.LoadLocation(items[1])
	if err != nil {
		fmt.Println("Ошибка при загрузке часового пояса:", items[1], err)
		return tgModel.SimpleReply(msg.Chat.ID, "Incorrect timezone:", msg.MessageID)
	}
	d.cities[items[0]] = items[1]
	return tgModel.SimpleReply(msg.Chat.ID, "Added!!", msg.MessageID)
}

func (d *data) delCity(msg *tgbotapi.Message, command *tgModel.Command) *tgModel.HandlerResult {
	_, ok := d.cities[command.Arguments.Raw]
	if ok {
		delete(d.cities, command.Arguments.Raw)
		return tgModel.SimpleReply(msg.Chat.ID, "Deleted by key:"+command.Arguments.Raw, msg.MessageID)
	}
	return tgModel.SimpleReply(msg.Chat.ID, "Not found!!", msg.MessageID)
}

func (d *data) printTimeList(msg *tgbotapi.Message, command *tgModel.Command) *tgModel.HandlerResult {
	result := ""
	var results formattedTimes

	for city, timezone := range d.cities {
		cityLocation, err := time.LoadLocation(timezone)
		if err != nil {
			fmt.Println("Ошибка при загрузке часового пояса:", timezone, err)
			continue
		}
		currentTime := time.Now().In(cityLocation)
		_, offset := currentTime.Zone()
		results = append(results, formattedTime{
			UTCSort: offset/3600 + (offset%3600)/60,
			UTC:     fmt.Sprintf("(UTC%+02d:%02d)", offset/3600, (offset%3600)/60),
			City:    city,
			Time:    currentTime.Format("15:04:05"),
		})
	}
	sort.Sort(results)
	for _, item := range results {
		//Выводим текущее время и смещение UTC
		result += fmt.Sprintf("\n%s: %s %s",
			item.Time,
			item.UTC,
			item.City,
		)
	}
	return tgModel.SimpleReply(msg.Chat.ID, result, msg.MessageID)
}

/*
05:25:58: (UTC-7:00) Сакраменто, США
14:25:58: (UTC+2:00) Варшава, Польша
14:25:58: (UTC+2:00) Подгорица, Черногория
14:25:58: (UTC+2:00) Амстердам, Нидерланды
14:25:58: (UTC+2:00) Белград, Сербия
15:25:58: (UTC+3:00) Москва, Россия
16:25:58: (UTC+4:00) Ижевск, Россия
16:25:58: (UTC+4:00) Ереван, Армения
17:25:58: (UTC+5:00) Пермь, Россия
19:25:58: (UTC+7:00) Красноярск, Россия
19:25:58: (UTC+7:00) Новосибирск, Россия
*/
