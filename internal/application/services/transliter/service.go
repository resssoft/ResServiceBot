package transliter

import (
	"fmt"
	"fun-coice/internal/domain/commands/tg"
	gt "github.com/bas24/googletranslatefree"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"sort"
	"strings"
)

type data struct {
	list        tgModel.Commands
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

func New() tgModel.Service {
	result := data{
		userStrings: make(map[int64]string),
	}
	commandsList := tgModel.NewCommands()
	commandsList["translit"] = tgModel.Command{
		Command:     "/translit",
		Description: "Encode string to base64",
		CommandType: "text",
		Permissions: tgModel.FreePerms,
		Handler:     result.transit,
	}
	commandsList["alphabet"] = tgModel.Command{
		Command:     "/alphabet",
		Synonyms:    []string{"af"},
		Description: "Encode string to base64",
		CommandType: "text",
		Permissions: tgModel.FreePerms,
		Handler:     result.alphabet,
	}
	commandsList[alphabetTrigger] = tgModel.Command{
		Command:     "/alphabetKey",
		Description: "alphabet notify by key",
		CommandType: "event",
		ListExclude: true, // do not show in the commands list
		Permissions: tgModel.FreePerms,
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

func (d *data) Commands() tgModel.Commands {
	return d.list
}

func (d *data) transit(msg *tgbotapi.Message, command *tgModel.Command) tgModel.HandlerResult {
	fmt.Println("transit command")
	//TODO: TRANSLATE FROM RUSSIAN TO AM (DETECT RUS)
	param := command.Arguments.Raw
	if param == "" {
		return tgModel.EmptyCommand()
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
	return tgModel.SimpleReply(msg.Chat.ID, param, msg.MessageID)
}

//translit site https://www.hayastan.com/translit/

func (d *data) alphabet(msg *tgbotapi.Message, _ *tgModel.Command) tgModel.HandlerResult {
	//fmt.Println("alphabet") //TODO: send to statistic service

	newMsg := tgbotapi.NewMessage(msg.Chat.ID, "_")
	newMsg.ReplyMarkup = tgbotapi.NewEditMessageReplyMarkup(
		msg.Chat.ID,
		msg.MessageID,
		alphabetKeyboard).ReplyMarkup
	return tgModel.PreparedCommand(newMsg)
}

func (d *data) alphabetEvent(msg *tgbotapi.Message, command *tgModel.Command) tgModel.HandlerResult {
	//TODO: fix saved data
	from := msg.From.ID
	_, ok := d.userStrings[from]
	//TODO:DATA RACE FIX
	if !ok {
		d.userStrings[from] = ""
	}
	switch command.Arguments.Raw {
	case "backspace":
		if len(d.userStrings[from]) > 0 {
			d.userStrings[from] = string([]rune(d.userStrings[from])[:len([]rune(d.userStrings[from]))-1])
		}
	case "translate":
		newCommand := tgModel.Command{
			Arguments: tgModel.CommandArguments{Raw: d.userStrings[from]},
		}
		return d.transit(msg, &newCommand)
	default:
		d.userStrings[from] += command.Arguments.Raw
	}

	//fmt.Println("alphabet event") // send to statistic service
	newMsg := tgbotapi.NewEditMessageText(msg.Chat.ID, msg.MessageID, d.userStrings[from]+"_")
	newMsg.ReplyMarkup = tgbotapi.NewEditMessageReplyMarkup(
		msg.Chat.ID,
		msg.MessageID,
		alphabetKeyboard).ReplyMarkup
	return tgModel.PreparedCommand(newMsg)
}
