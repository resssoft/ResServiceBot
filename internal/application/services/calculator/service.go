package calculator

import (
	"fmt"
	tgModel "fun-coice/internal/domain/commands/tg"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mnogu/go-calculator"
	"log"
	"math"
	"strings"
)

type data struct {
	list tgModel.Commands
}

func New() tgModel.Service {
	result := data{}
	commandsList := tgModel.NewCommands()
	commandsList["calc"] = tgModel.Command{
		Command:     "/calc",
		Synonyms:    []string{"calc", "калк", "сколько будет"},
		Triggers:    []string{"calc", "калк", "сколько будет"},
		Templates:   []string{`^\d[\d\s\+\\\-\*\(\)\.]+$`},
		Description: "(параметры - строка для продсчета данных, пример 2+2 или (2.5 - 1.35) * 2.0",
		CommandType: "text",
		Permissions: tgModel.FreePerms,
		Handler:     result.calcFromStr,
	}

	result.list = commandsList
	return &result
}

func (d data) Commands() tgModel.Commands {
	return d.list
}

func (d data) Name() string {
	return "calculator"
}

func (d data) calcFromStr(msg *tgbotapi.Message, command *tgModel.Command) tgModel.HandlerResult {
	params := strings.Split(command.Arguments.Raw, " ")
	log.Println("calcFromStr", command.Arguments.Raw)
	log.Println("params", params)
	log.Println("commandName", command.Command)
	val, err := calculator.Calculate(command.Arguments.Raw)
	if err != nil {
		log.Println(err.Error())
		//TODO: admin errors log
		return tgModel.EmptyCommand()
	}
	resultText := fmt.Sprintf("%.2f", val)
	intPart, floatPart := math.Modf(val)
	if floatPart == 0 {
		resultText = fmt.Sprintf("%.0f", intPart)
	}
	if val < 0.01 {
		resultText = fmt.Sprintf("%.5f", val)
	}
	return tgModel.SimpleReply(msg.Chat.ID, resultText, msg.MessageID)
}
