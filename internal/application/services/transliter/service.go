package transliter

import (
	"fmt"
	tgCommands "fun-coice/internal/domain/commands/tg"
	gt "github.com/bas24/googletranslatefree"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

type data struct {
	list tgCommands.Commands
}

var LatToAm map[string]string

var AmToLat map[string]string = map[string]string{
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

func New() tgCommands.Service {
	result := data{}
	commandsList := make(tgCommands.Commands)
	commandsList["transit"] = tgCommands.Command{
		Command:     "/transit",
		Description: "Encode string to base64",
		CommandType: "text",
		Permissions: tgCommands.FreePerms,
		Handler:     result.transit,
	}
	LatToAm = make(map[string]string)
	for key, val := range AmToLat {
		LatToAm[val] = key
	}

	result.list = commandsList
	return &result
}

func (d data) Commands() tgCommands.Commands {
	return d.list
}

func (d data) transit(msg *tgbotapi.Message, commandName string, param string, params []string) tgCommands.HandlerResult {
	if param == "" {
		fmt.Println("EMPTY TRANSLITE COMMAND")
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
