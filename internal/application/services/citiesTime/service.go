package citiesTime

import (
	"fmt"
	"time"

	tgModel "fun-coice/internal/domain/commands/tg"
	tgbotapi "fun-coice/pkg/telegram-bot-api"
)

type data struct {
	list   tgModel.Commands
	cities map[string]string
}

func New() tgModel.Service {
	s := data{
		cities: make(map[string]string),
	}
	s.cities = map[string]string{
		"Сакраменто, США":       "America/Los_Angeles",
		"Варшава, Польша":       "Europe/Warsaw",
		"Амстердам, Нидерланды": "Europe/Amsterdam",
		"Белград, Сербия":       "Europe/Belgrade",
		"Подгорица, Черногория": "Europe/Podgorica",
		"Москва, Россия":        "Europe/Moscow",
		"Ижевск, Россия":        "Europe/Samara", // Ижевск находится в часовом поясе Самары
		"Ереван, Армения":       "Asia/Yerevan",
		"Пермь, Россия":         "Asia/Yekaterinburg",
		"Новосибирск, Россия":   "Asia/Novosibirsk",
	}
	commandsList := tgModel.NewCommands()
	tgModel.NewCommand().
		Simple("time", "Show cities time", s.printTimeList).
		WithTriggers("время").
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

func (d *data) Configure(_ tgModel.ServiceConfig) {

}

func (d *data) printTimeList(msg *tgbotapi.Message, command *tgModel.Command) *tgModel.HandlerResult {
	result := ""
	for city, timezone := range d.cities {
		cityLocation, err := time.LoadLocation(timezone)
		if err != nil {
			fmt.Println("Ошибка при загрузке часового пояса:", timezone, err)
			continue
		}
		currentTime := time.Now().In(cityLocation)
		_, offset := currentTime.Zone()

		// Выводим текущее время и смещение UTC
		result += fmt.Sprintf("\n%s: (UTC%+02d:%02d) %s",
			currentTime.Format("15:04:05"),
			offset/3600,      // Часы
			(offset%3600)/60, // Минуты
			city,
		)
	}
	return tgModel.SimpleReply(msg.Chat.ID, result, msg.MessageID)
}
