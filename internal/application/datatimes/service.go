package datatimes

import (
	"fmt"
	tgCommands "fun-coice/internal/domain/commands/tg"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/hako/durafmt"
	"regexp"
	"strconv"
	"time"
)

type data struct {
	list tgCommands.Commands
}

func New() tgCommands.Service {
	result := data{}
	commandsList := make(tgCommands.Commands)
	commandsList["tm"] = tgCommands.Command{
		Command:     "/tm",
		Synonyms:    []string{"time", "timestamp", "datetime", "date"},
		Description: "Get date time info or convert",
		CommandType: "text",
		Permissions: tgCommands.FreePerms,
		Handler:     result.getInfo,
	}

	result.list = commandsList
	return &result
}

func (d data) Commands() tgCommands.Commands {
	return d.list
}

func (d data) getInfo(msg *tgbotapi.Message, commandName string, commandValue string, params []string) (tgbotapi.Chattable, bool) {
	if len(params) == 0 {
		return tgbotapi.NewMessage(msg.Chat.ID, time.Now().Format("2006-01-02 15:04:05 -0700")), true
	}
	msgText := fmt.Sprintf("Input: %s\n\n", commandValue)
	timestampConvert, _ := regexp.MatchString(`^\d\d\d\d\d\d+$`, commandValue)
	timestampCompare, _ := regexp.MatchString(`^\d\d\d\d\d\d+\s\d\d\d\d\d\d+$`, commandValue)
	switch {
	case timestampConvert:
		timeInt, err := strconv.ParseInt(commandValue, 10, 64)
		if err != nil {
			msgText += "Error parse int timestamp: " + err.Error()
		} else {
			parsedTime := time.Unix(timeInt, 0)
			msgText += fmt.Sprintf("Date: %s \n", parsedTime.Format("2006-01-02 15:04:05 -0700"))
			msgText += fmt.Sprintf("Date: %s \n", parsedTime.Format("2006-01-02T15:04:05Z07:00"))
			msgText += fmt.Sprintf("Date: %s \n", parsedTime.Format("Monday, 02-Jan-06 15:04:05 MST"))
			msgText += fmt.Sprintf("Timestamp: %s \n", fmt.Sprintf("%v", parsedTime.Unix()))
			msgText += fmt.Sprintf("Diff: %s \n", durafmt.Parse(time.Now().Sub(parsedTime)).LimitFirstN(2).String())
			msgText += fmt.Sprintf("\nNow: %s \n", time.Now().Format("2006-01-02 15:04:05 -0700"))
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
	default:
		parsedTime, err := time.Parse("2006-01-02 15:04:05", commandValue)
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
			msgText += fmt.Sprintf("Date: %s \n", parsedTime.Format("2006-01-02 15:04:05 -0700"))
			msgText += fmt.Sprintf("Date: %s \n", parsedTime.Format("2006-01-02T15:04:05Z07:00"))
			msgText += fmt.Sprintf("Date: %s \n", parsedTime.Format("Monday, 02-Jan-06 15:04:05 MST"))
			msgText += fmt.Sprintf("Timestamp: %s \n", fmt.Sprintf("%v", parsedTime.Unix()))
			msgText += fmt.Sprintf("Diff by current: %s \n", durafmt.Parse(time.Now().Sub(parsedTime)).LimitFirstN(2).String())
			msgText += fmt.Sprintf("\nNow: %s \n", time.Now().Format("2006-01-02 15:04:05 -0700"))
		}
	}
	time.Parse(commandValue, commandValue)
	return tgbotapi.NewMessage(msg.Chat.ID, msgText), true
}
