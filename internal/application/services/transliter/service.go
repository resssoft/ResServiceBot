package transliter

import (
	"fmt"
	tgCommands "fun-coice/internal/domain/commands/tg"
	gt "github.com/bas24/googletranslatefree"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"sort"
	"strings"
)

type data struct {
	list        tgCommands.Commands
	userStrings map[int64]string
}

var LatToAm map[string]string

var AmToLat = map[string]string{
	"ա":  "a",
	"բ":  "b",
	"գ":  "g",
	"դ":  "d",
	"ե":  "e",
	"զ":  "z",
	"է":  "e'",
	"ը":  "y'",
	"թ":  "t'",
	"ժ":  "jh",
	"ի":  "i",
	"լ":  "l",
	"խ":  "x",
	"ծ":  "c'",
	"կ":  "k",
	"հ":  "h",
	"ձ":  "d'",
	"ղ":  "gh",
	"ճ":  "tw",
	"մ":  "m",
	"յ":  "y",
	"ն":  "n",
	"շ":  "sh",
	"ո":  "o",
	"չ":  "ch",
	"պ":  "p",
	"ջ":  "j",
	"ռ":  "r'",
	"ս":  "s",
	"վ":  "v",
	"տ":  "t",
	"ր":  "r",
	"ց":  "c",
	"ւ":  "w",
	"փ":  "p'",
	"ք":  "q",
	"օ":  "o'",
	"ֆ":  "f",
	"ու": "u",
	"և":  "&",
}

var alphabetKeyboard tgbotapi.InlineKeyboardMarkup
var alphabetTrigger = "alphabetKey"

func New() tgCommands.Service {
	result := data{
		userStrings: make(map[int64]string),
	}
	commandsList := tgCommands.NewCommands()
	commandsList["translit"] = tgCommands.Command{
		Command:     "/translit",
		Description: "Encode string to base64",
		CommandType: "text",
		Permissions: tgCommands.FreePerms,
		Handler:     result.transit,
	}
	commandsList["alphabet"] = tgCommands.Command{
		Command:     "/alphabet",
		Synonyms:    []string{"af"},
		Description: "Encode string to base64",
		CommandType: "text",
		Permissions: tgCommands.FreePerms,
		Handler:     result.alphabet,
	}
	commandsList[alphabetTrigger] = tgCommands.Command{
		Command:     "/alphabetKey",
		Description: "alphabet notify by key",
		CommandType: "event",
		ListExclude: true, // do not show in the commands list
		Permissions: tgCommands.FreePerms,
		Handler:     result.alphabetEvent,
	}

	LatToAm = make(map[string]string)
	var amLetters []string
	for key, val := range AmToLat {
		LatToAm[val] = key
		amLetters = append(amLetters, key)
	}

	//set alphabet
	var row []tgbotapi.InlineKeyboardButton
	var rows [][]tgbotapi.InlineKeyboardButton
	itemsByRow := 6
	index := 0
	itemCount := 0

	sort.Strings(amLetters)
	for _, val := range amLetters {
		itemCount++
		index++
		if index <= itemsByRow && itemCount != len(amLetters) {
			row = append(row, tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf("%s (%s)", val, AmToLat[val]),
				alphabetTrigger+":"+val))
			continue
		}
		index = 0
		rows = append(rows, row)
		row = nil
	}
	row = append(row, tgbotapi.NewInlineKeyboardButtonData(" ", alphabetTrigger+":"+" "))
	row = append(row, tgbotapi.NewInlineKeyboardButtonData("←", alphabetTrigger+":"+"backspace"))
	row = append(row, tgbotapi.NewInlineKeyboardButtonData("↯", alphabetTrigger+":"+"translate"))
	rows = append(rows, row)
	alphabetKeyboard = tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: rows,
	}

	result.list = commandsList
	return &result
}

func (d *data) Commands() tgCommands.Commands {
	return d.list
}

func (d *data) transit(msg *tgbotapi.Message, commandName string, param string, params []string) tgCommands.HandlerResult {
	fmt.Println("transit command")
	//TODO: TRANSLATE FROM RUSSIAN TO AM (DETECT RUS)
	if param == "" {
		return tgCommands.EmptyCommand()
	}
	for latin, An := range LatToAm {
		if len([]rune(latin)) > 1 {
			param = strings.Replace(param, latin, An, -1)
			param = strings.Replace(param, strings.ToUpper(latin), strings.ToUpper(An), -1)
		}
	}
	for latin, An := range LatToAm {
		param = strings.Replace(param, latin, An, -1)
		param = strings.Replace(param, strings.ToUpper(latin), strings.ToUpper(An), -1)
	}
	result, err := gt.Translate(param, "hy", "ru")
	if err != nil {
		result = err.Error()
		fmt.Println(err)
	} else {
		param += "\n\n" + result
	}
	return tgCommands.SimpleReply(msg.Chat.ID, param, msg.MessageID)
}

//translit site https://www.hayastan.com/translit/

func (d *data) alphabet(msg *tgbotapi.Message, commandName string, param string, params []string) tgCommands.HandlerResult {
	//fmt.Println("alphabet") //TODO: send to statistic service

	newMsg := tgbotapi.NewMessage(msg.Chat.ID, "_")
	newMsg.ReplyMarkup = tgbotapi.NewEditMessageReplyMarkup(
		msg.Chat.ID,
		msg.MessageID,
		alphabetKeyboard).ReplyMarkup
	return tgCommands.PreparedCommand(newMsg)
}

func (d *data) alphabetEvent(msg *tgbotapi.Message, commandName string, param string, params []string) tgCommands.HandlerResult {
	//TODO: fix saved data
	from := msg.From.ID
	_, ok := d.userStrings[from]
	//TODO:DATA RACE FIX
	if !ok {
		d.userStrings[from] = ""
	}
	switch param {
	case "backspace":
		if len(d.userStrings[from]) > 0 {
			d.userStrings[from] = string([]rune(d.userStrings[from])[:len([]rune(d.userStrings[from]))-1])
		}
	case "translate":
		return d.transit(msg, commandName, d.userStrings[from], []string{})
	default:
		d.userStrings[from] += param
	}

	//fmt.Println("alphabet event") // send to statistic service
	newMsg := tgbotapi.NewEditMessageText(msg.Chat.ID, msg.MessageID, d.userStrings[from]+"_")
	newMsg.ReplyMarkup = tgbotapi.NewEditMessageReplyMarkup(
		msg.Chat.ID,
		msg.MessageID,
		alphabetKeyboard).ReplyMarkup
	return tgCommands.PreparedCommand(newMsg)
}
