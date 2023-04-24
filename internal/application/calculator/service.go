package calculator

import (
	"fmt"
	tgCommands "fun-coice/internal/domain/commands/tg"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mnogu/go-calculator"
	"log"
	"math"
)

type data struct {
	list tgCommands.Commands
}

func New() tgCommands.Service {
	result := data{}
	commandsList := make(tgCommands.Commands)
	commandsList["calc"] = tgCommands.Command{
		Command:     "/calc",
		Synonyms:    []string{"calc", "калк", "сколько будет"},
		Templates:   []string{`^\d[\d\s\+\\\-\*\(\)\.]+$`},
		Description: "Encode string to base64",
		CommandType: "text",
		Permissions: tgCommands.FreePerms,
		Handler:     result.calcFromStr,
	}

	result.list = commandsList
	return &result
}

func (d data) Commands() tgCommands.Commands {
	return d.list
}

func (d data) calcFromStr(msg *tgbotapi.Message, commandName string, param string, params []string) (tgbotapi.Chattable, bool) {
	log.Println("calcFromStr", param)
	log.Println("params", params)
	log.Println("commandName", commandName)
	val, err := calculator.Calculate(param)
	if err != nil {
		log.Println(err.Error())
		//TODO: admin errors log
		return nil, false
	}
	resultText := fmt.Sprintf("%.2f", val)
	intPart, floatPart := math.Modf(val)
	if floatPart == 0 {
		resultText = fmt.Sprintf("%.0f", intPart)
	}
	if val < 0.01 {
		resultText = fmt.Sprintf("%.5f", val)
	}
	msgNew := tgbotapi.NewMessage(msg.Chat.ID, resultText)
	msgNew.ReplyToMessageID = msg.MessageID
	return msgNew, true
}
