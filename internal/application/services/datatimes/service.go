package datatimes

import (
	"fmt"
	tgModel "fun-coice/internal/domain/commands/tg"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/hako/durafmt"
	"github.com/pawelszydlo/humanize"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type data struct {
	list       tgModel.Commands
	humanTimes map[string]*humanize.Humanizer
}

var languages = []string{"en"}

type timeType uint16

const (
	Nanosecond  timeType = 1
	Microsecond timeType = 2
	Millisecond timeType = 4
	Second      timeType = 16
	Minute      timeType = 32
	Hour        timeType = 64
	Day         timeType = 128
	Week        timeType = 256
	Month       timeType = 512
	Year        timeType = 1024
)

func New() tgModel.Service {
	var err error
	humanTimes := make(map[string]*humanize.Humanizer)
	for _, lang := range languages {
		humanTimes[lang], err = humanize.New(lang)
		if err != nil {
			fmt.Println("Error init humanize for lang", lang, err.Error())
		}
	}
	result := data{
		humanTimes: humanTimes,
	}
	commandsList := tgModel.NewCommands()
	commandsList["tm"] = tgModel.Command{
		Command:     "/tm",
		Synonyms:    []string{"time", "timestamp", "datetime", "date", "время"},
		Description: "Get date time info or convert",
		CommandType: "text",
		Permissions: tgModel.FreePerms,
		Handler:     result.getInfo,
	}
	commandsList["tmdu"] = tgModel.Command{
		Command:     "/tmdu",
		Synonyms:    []string{"timeDur", "timeDuration", "timeOf", "duration", "продолжительность"},
		Templates:   []string{`^\d[\d\s\+\\\-\*\(\)\.]+$`},
		Description: "Get date time info or convert",
		CommandType: "text",
		Permissions: tgModel.FreePerms,
		Handler:     result.getDuration,
	}
	commandsList["timeconvert"] = tgModel.Command{
		Command:     "/timeconvert",
		Templates:   []string{`^\d+\s[^\d\s]+\s[i]{0,1}[n]{0,1}[в]{0,1}\s[^\d\s]+$`},
		Description: "convert time, format: /timeconvert <digit> <timeType: seconds, days, years> <digit> <timeType: seconds, days, years>",
		CommandType: "text",
		Permissions: tgModel.FreePerms,
		Handler:     result.timeConvert,
	}

	result.list = commandsList
	return &result
}

func (d *data) Commands() tgModel.Commands {
	return d.list
}

func (d *data) Name() string {
	return "datatimes"
}

func (d *data) getDuration(msg *tgbotapi.Message, command *tgModel.Command) *tgModel.HandlerResult {
	msgText := fmt.Sprintf("Input: %s\n\n", command.Arguments.Raw)
	duration := time.Second
	var err error
	if ht, ok := d.humanTimes["en"]; ok {
		if ht != nil {
			duration, err = d.humanTimes["en"].ParseDuration(command.Arguments.Raw)
			if err != nil {
				fmt.Println("Error init humanize for lang", err.Error())
			}
		} else {
			fmt.Println("humanize for en is nil")
		}
	} else {
		fmt.Println("humanize for en is not exist")
	}
	msgText += dateFormat(time.Now().Add(duration))
	return tgModel.Simple(msg.Chat.ID, msgText)
}

func (d *data) getInfo(msg *tgbotapi.Message, command *tgModel.Command) *tgModel.HandlerResult {
	params := strings.Split(command.Arguments.Raw, " ")
	commandValue := command.Arguments.Raw
	var err error
	if len(params) == 0 {
		return tgModel.Simple(msg.Chat.ID, time.Now().Format("2006-01-02 15:04:05 -0700"))
	}
	msgText := fmt.Sprintf("Input: %s\n\n", commandValue)
	timestampConvert, _ := regexp.MatchString(`^\d\d\d\d\d\d+$`, commandValue)
	timestampCompare, _ := regexp.MatchString(`^\d\d\d\d\d\d+\s\d\d\d\d\d\d+$`, commandValue)
	timesGMT, _ := regexp.MatchString(`^GMT\s?[\-\+]+\d+$`, commandValue)
	timesGMTRu, _ := regexp.MatchString(`^гмт\s?[\-\+]+\d+$`, commandValue)
	switch {
	case timestampConvert:
		timeInt, err := strconv.ParseInt(commandValue, 10, 64)
		if err != nil {
			msgText += "Error parse int timestamp: " + err.Error()
		} else {
			parsedTime := time.Unix(timeInt, 0)
			msgText += dateFormat(parsedTime)
		}
	case timestampCompare:
		r, _ := regexp.Compile(`^(\d\d\d\d\d\d+)\s(\d\d\d\d\d\d+)$`)
		parsedItems := r.FindStringSubmatch(commandValue)
		fmt.Println(parsedItems)
		if len(parsedItems) != 3 {
			msgText += fmt.Sprintf("Error: parsedItems != 3 It= %v", len(parsedItems))
		} else {
			timeInt1, err1 := strconv.ParseInt(parsedItems[1], 10, 64)
			if err1 != nil {
				msgText += "\nError: " + err1.Error()
			}
			timeInt2, err2 := strconv.ParseInt(parsedItems[2], 10, 64)
			if err2 != nil {
				msgText += "\nError: " + err2.Error()
			}
			if err1 == nil && err2 == nil {
				parsedTime1 := time.Unix(timeInt1, 0)
				parsedTime2 := time.Unix(timeInt2, 0)
				msgText += fmt.Sprintf("Diff: %s \n", durafmt.Parse(parsedTime1.Sub(parsedTime2)).LimitFirstN(2).String())
			}
		}
	case timesGMT:
		r, _ := regexp.Compile(`^GMT\s?([\-\+]+)(\d+)$`)
		parsedItems := r.FindStringSubmatch(commandValue)
		if len(parsedItems) != 3 {
			msgText += fmt.Sprintf("Error: parsedItems != 3 It= %v", len(parsedItems))
		} else {
			offset, _ := strconv.Atoi(parsedItems[2])
			if parsedItems[2] == "-" {
				offset *= -1
			}
			parsedTime := time.Now().In(time.FixedZone("GMT", int((time.Hour * 3).Seconds())))
			msgText += fmt.Sprintf("Date: %s \n", parsedTime.Format("15:04:05 -0700"))
		}
	case timesGMTRu:
		r, _ := regexp.Compile(`^гмт\s?([\-\+]+)(\d+)$`)
		parsedItems := r.FindStringSubmatch(commandValue)
		if len(parsedItems) != 3 {
			msgText += fmt.Sprintf("Error: parsedItems != 3 It= %v", len(parsedItems))
		} else {
			offset, _ := strconv.Atoi(parsedItems[2])
			if parsedItems[1] == "-" {
				offset *= -1
			}
			parsedTime := time.Now().In(time.FixedZone("GMT", int((time.Hour * 3).Seconds())))
			msgText += fmt.Sprintf("Date: %s \n", parsedTime.Format("15:04:05 -0700"))
		}
	default:
		parsedTime := time.Now()
		byTz := byTimezone(commandValue)
		if byTz != nil {
			fmt.Println("byTz", byTz.String())
			parsedTime = time.Now().In(byTz)
		} else {
			parsedTime, err = time.Parse("2006-01-02 15:04:05", commandValue)
		}
		if err != nil {
			parsedTime, err = time.Parse("2006-01-02 15:04", commandValue)
		}
		if err != nil {
			parsedTime, err = time.Parse("2006-01-02", commandValue)
		}
		if err != nil {
			parsedTime, err = time.Parse("02-01-2006", commandValue)
		}
		if err != nil {
			parsedTime, err = time.Parse("02.01.2006", commandValue)
		}
		if err != nil {
			parsedTime, err = time.Parse("02/01/2006", commandValue)
		}
		if err != nil {
			parsedTime, err = time.Parse("2006-01-02T15:04:05Z07:00", commandValue)
		}
		if err != nil {
			parsedTime, err = time.Parse("Mon Jan _2 15:04:05 2006", commandValue)
		}
		if err != nil {
			parsedTime, err = time.Parse("Mon Jan _2 15:04:05 MST 2006", commandValue)
		}
		if err != nil {
			parsedTime, err = time.Parse("02 Jan 06 15:04 MST", commandValue)
		}
		if err != nil {
			parsedTime, err = time.Parse("Mon, 02 Jan 2006 15:04:05 MST", commandValue)
		}
		if err != nil {
			parsedTime, err = time.Parse("Mon 02 Jan 2006 15:04:05 MST", commandValue)
		}
		if err != nil {
			parsedTime, err = time.Parse("Mon 02 Jan 2006 15:04:05 07:00", commandValue)
		}
		if err != nil {
			parsedTime, err = time.Parse("January 2, 2006", commandValue)
		}
		var parsedDuration time.Duration
		if err != nil {
			parsedDuration, err = time.ParseDuration(commandValue)
			parsedTime = time.Now().Add(parsedDuration)
		}
		if err != nil {
			msgText += "Error parse time"
		} else {
			msgText += timeFormat(parsedTime)
		}
	}
	//time.Parse(commandValue, commandValue)
	return tgModel.Simple(msg.Chat.ID, msgText)
}

func byTimezone(value string) *time.Location {
	tz, err := time.LoadLocation(value) //example "America/New_York"
	if err != nil {
		tz = byCity(value)
	}
	fmt.Println("tz", tz)
	if tz != nil {
		return tz
	}
	return nil
}

func byCity(value string) *time.Location {
	switch strings.ToLower(value) {
	case "msk", "moscow", "мск", "москва", "москве":
		fmt.Println("msk", time.FixedZone("GMT", int((time.Hour*3).Seconds())).String())
		return time.FixedZone("GMT", int((time.Hour * 3).Seconds()))
	case "nsk", "novosibirsk", "нск", "новосибирск", "новосиб":
		fmt.Println("nsk", time.FixedZone("GMT", int((time.Hour*7).Seconds())).String())
		return time.FixedZone("GMT", int((time.Hour * 7).Seconds()))
	}
	return nil
}

func timeFormat(val time.Time) string {
	msgText := ""
	msgText += fmt.Sprintf("Time: %s \n", val.Format("15:04:05 -0700"))
	msgText += fmt.Sprintf("Timestamp: %s \n", fmt.Sprintf("%v", val.Unix()))
	msgText += fmt.Sprintf("Diff by current: %s \n", durafmt.Parse(time.Now().Sub(val)).LimitFirstN(2).String())
	msgText += fmt.Sprintf("\nNow: %s \n", time.Now().Format("15:04:05 -0700"))
	return msgText
}

func dateFormat(val time.Time) string {
	msgText := ""
	msgText += fmt.Sprintf("Date: %s \n", val.Format("2006-01-02 15:04:05 -0700"))
	msgText += fmt.Sprintf("Date: %s \n", val.Format("2006-01-02T15:04:05Z07:00"))
	msgText += fmt.Sprintf("Date: %s \n", val.Format("Monday, 02-Jan-06 15:04:05 MST"))
	msgText += fmt.Sprintf("Timestamp: %s \n", fmt.Sprintf("%v", val.Unix()))
	msgText += fmt.Sprintf("Diff by current: %s \n", durafmt.Parse(time.Now().Sub(val)).LimitFirstN(2).String())
	msgText += fmt.Sprintf("\nNow: %s \n", time.Now().Format("2006-01-02 15:04:05 -0700"))
	return msgText
}

func durationFormat(val string) string {
	msgText := ""
	parsed, err := durafmt.ParseString(val)
	if err == nil && parsed != nil {
		msgText += fmt.Sprintf("Time: %s \n", parsed.LimitFirstN(2).String())
	} else {
		fmt.Println("")
	}

	msgText += fmt.Sprintf("\nNow: %s \n", time.Now().Format("15:04:05 -0700"))
	return msgText
}

func (d *data) timeConvert(msg *tgbotapi.Message, command *tgModel.Command) *tgModel.HandlerResult {
	commandValue := strings.ToLower(command.Arguments.Raw)
	fromType := Nanosecond
	toType := Nanosecond
	from := int64(1)
	resultName := ""
	fmt.Println("timeConvert", commandValue) // DEBUG
	fromVal, fromName, toName := d.parseTimeCovert(commandValue)
	if fromName == "" || fromVal == "" {
		fmt.Println("not parsed", commandValue) // DEBUG
		return tgModel.EmptyCommand()
	}
	switch fromName {
	case "second", "seconds", "секунд", "секунда", "секунды", "секундах":
		fromType = Second
	case "minute", "minutes", "минут", "минута", "минуты", "минутах":
		fromType = Minute
	case "hour", "hours", "часов", "час", "часа", "часах":
		fromType = Hour
	case "days", "day", "дней", "дня", "день", "днях":
		fromType = Day
	case "Week", "Weeks", "недель", "неделя", "недели", "неделях":
		fromType = Week
	case "Month", "Months", "месяцев", "месяца", "месяц", "месяцах":
		fromType = Month
	case "Year", "Years", "лет", "годов", "года", "годах":
		fromType = Year
	}
	switch toName {
	case "second", "seconds", "секунд", "секунда", "секунды", "секундах":
		toType = Second
	case "minute", "minutes", "минут", "минута", "минуты", "минутах":
		toType = Minute
	case "hour", "hours", "часов", "час", "часа", "часах":
		toType = Hour
	case "days", "day", "дней", "дня", "день", "днях":
		toType = Day
	case "Week", "Weeks", "недель", "неделя", "недели", "неделях":
		toType = Week
	case "Month", "Months", "месяцев", "месяца", "месяц", "месяцах":
		toType = Month
	case "Year", "Years", "лет", "годов", "года", "годах":
		toType = Year
	}
	from, _ = strconv.ParseInt(fromVal, 10, 64)
	fmt.Println(fmt.Sprintf("fromName[%s]%v toName[%s]%v val: %v", fromName, fromType, toName, toType, from)) // DEBUG
	result := d.convert(from, fromType, toType)
	return tgModel.Simple(msg.Chat.ID, fmt.Sprintf("%v %s", result, resultName))
}

func (d *data) parseTimeCovert(val string) (string, string, string) {
	r, _ := regexp.Compile(`^\s{0,1}(\d+)\s([^\d\s]+)\s{0,1}in\s([^\d\s]+)\s{0,1}`)
	parsedItems := r.FindStringSubmatch(val)
	if len(parsedItems) == 4 {
		return parsedItems[1], parsedItems[2], parsedItems[3]
	}
	fmt.Println("not parsed 1") // DEBUG
	r, _ = regexp.Compile(`^\s{0,1}(\d+)\s([^\d\s]+)\s{0,1}в\s([^\d\s]+)\s{0,1}`)
	parsedItems = r.FindStringSubmatch(val)
	if len(parsedItems) == 4 {
		return parsedItems[1], parsedItems[2], parsedItems[3]
	}
	fmt.Println("not parsed 2") // DEBUG
	r, _ = regexp.Compile(`^/timeconvert\s{0,1}(\d+)\s([^\d\s]+)\s([^\d\s]+)\s{0,1}`)
	parsedItems = r.FindStringSubmatch(val)
	if len(parsedItems) == 4 {
		return parsedItems[1], parsedItems[2], parsedItems[3]
	}
	fmt.Println("not parsed 3") // DEBUG
	r, _ = regexp.Compile(`^\s{0,1}(\d+)\s([^\d\s]+)\s([^\d\s]+)\s{0,1}`)
	parsedItems = r.FindStringSubmatch(val)
	if len(parsedItems) == 4 {
		return parsedItems[1], parsedItems[2], parsedItems[3]
	}
	fmt.Println("not parsed 4") // DEBUG
	return "", "", ""
}

func (d *data) convert(fromVal int64, fromType, toType timeType) int64 {
	from := int64(1)
	to := int64(1)
	result := int64(1)
	switch fromType {
	case Nanosecond:
		from = fromVal * int64(time.Nanosecond)
	case Microsecond:
		from = fromVal * int64(time.Microsecond)
	case Millisecond:
		from = fromVal * int64(time.Millisecond)
	case Second:
		from = fromVal * int64(time.Second)
	case Minute:
		from = fromVal * int64(time.Minute)
	case Hour:
		from = fromVal * int64(time.Hour)
	case Day:
		from = fromVal * int64(time.Hour*24)
	case Week:
		from = fromVal * int64(time.Hour*24*7)
	case Month:
		from = fromVal * int64(time.Hour*24*30)
	case Year:
		from = fromVal * int64(time.Hour*24*30*12)
	}
	switch toType {
	case Nanosecond:
		to = int64(time.Nanosecond)
	case Microsecond:
		to = int64(time.Microsecond)
	case Millisecond:
		to = int64(time.Millisecond)
	case Second:
		to = int64(time.Second)
	case Minute:
		to = int64(time.Minute)
	case Hour:
		to = int64(time.Hour)
	case Day:
		to = int64(time.Hour * 24)
	case Week:
		to = int64(time.Hour * 24 * 7)
	case Month:
		to = int64(time.Hour * 24 * 30)
	case Year:
		to = int64(time.Hour * 24 * 30 * 12)
	}
	result = from / to
	return result
}
