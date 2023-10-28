package funs

import (
	"encoding/json"
	"fmt"
	tgModel "fun-coice/internal/domain/commands/tg"
	"fun-coice/pkg/scribble"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"
)

type data struct {
	list tgModel.Commands
	DB   *scribble.Driver
}

var errorCommandMsg = "Error, sorry! Write to admin or send command /bug with command name"

type FunCommand struct {
	Name      string
	TgCommand tgModel.Command
	List1     []string
	List2     []string
}

var FunCommands map[string]FunCommand
var syncMap *sync.Mutex
var funCommandDCollection = "funcommands"
var funCommandType = "funcommand"

// New TODO: move to aplication folder
// New TODO: add list commands and remove (by admin)
func New(DB *scribble.Driver) tgModel.Service {
	result := data{
		DB: DB,
	}
	commandsList := tgModel.NewCommands()
	commandsList["addfan"] = tgModel.Command{
		Command:     "/addfan",
		Synonyms:    []string{"addfan", "добавитьфан"},
		Description: "Добавить генератор фанов",
		CommandType: "text",
		Permissions: tgModel.FreePerms,
		Handler:     result.add,
	}

	FunCommands = make(map[string]FunCommand)
	syncMap = new(sync.Mutex)

	records, err := DB.ReadAll(funCommandDCollection)
	if err != nil {
		fmt.Println("Error DB.ReadAll", err)
	}
	appPath, err := os.Getwd()
	if err != nil {
		fmt.Println("Error os.Getwd", err)
		return nil
	}
	if err := os.Mkdir(appPath+string(os.PathSeparator), os.ModeDir); err != nil {
		fmt.Println("initFunCommand", err.Error())
		//return errorCommandMsg
	}
	for _, f := range records {
		funCommand := FunCommand{}
		if err := json.Unmarshal([]byte(f), &funCommand); err != nil {
			fmt.Println("Error Unmarshal", err)
			continue
		}
		funCommand.TgCommand.Handler = result.run
		appendFunCommand(funCommand.Name, funCommand)
		commandsList[funCommand.Name] = funCommand.TgCommand
		//fmt.Println("Add fun command", funCommand.Name)
	}
	//fmt.Println("FunCommands", FunCommands)
	result.list = commandsList
	return &result
}

func (d *data) Commands() tgModel.Commands {
	return d.list
}

func (d *data) Name() string {
	return "funs"
}

func (d *data) Configure(_ tgModel.ServiceConfig) {

}

func (d *data) add(msg *tgbotapi.Message, command *tgModel.Command) *tgModel.HandlerResult {
	params := strings.Split(command.Arguments.Raw, " ")
	fmt.Println(params)
	text := ""
	if len(params) != 4 {
		text = "format: /addfan newcommandname list1_item1,list1_item2 list2_item1,list2_item2"
		text += "\nExample: cats cute,funny,fluffy Molly,Charlie,Oscar"
		text += "\n no more than 3 spase in the string"
	} else {
		text = d.addFunCommand(params[1], params[2], params[3])
	}
	return tgModel.Simple(msg.Chat.ID, text)
}

func (d *data) run(msg *tgbotapi.Message, command *tgModel.Command) *tgModel.HandlerResult {
	if commandData, exist := isFunCommand(command.Command); exist {
		s1 := rand.NewSource(time.Now().UnixNano())
		r1 := rand.New(s1)
		time.Sleep(time.Millisecond * time.Duration(r1.Int63n(200)))
		r2 := rand.New(s1)
		return tgModel.SimpleReply(msg.Chat.ID,
			commandData.List1[r1.Intn(len(commandData.List1))]+" "+commandData.List2[r2.Intn(len(commandData.List2))],
			msg.MessageID)
	}
	return tgModel.Simple(msg.Chat.ID, "Something wrong! Write to admin")
}

func isFunCommand(name string) (FunCommand, bool) {
	name = strings.Replace(name, "/", "", 1)
	syncMap.Lock()
	defer syncMap.Unlock()
	if data, ok := FunCommands[name]; ok {
		return data, ok
	}
	if data, ok := FunCommands[name]; ok {
		return data, ok
	}
	return FunCommand{}, false
}

func appendFunCommand(name string, command FunCommand) {
	syncMap.Lock()
	if FunCommands != nil {
		FunCommands[name] = command
	}
	syncMap.Unlock()
}

func (d *data) addFunCommand(name, list1, list2 string) string {
	command := tgModel.Command{
		Command:     "/" + name,
		Synonyms:    nil,
		Description: "Get random fanny words!",
		CommandType: funCommandType,
		Permissions: tgModel.FreePerms,
		Handler:     d.run,
	}
	funCommand := FunCommand{
		Name:      name,
		TgCommand: command,
		List1:     strings.Split(list1, ","),
		List2:     strings.Split(list2, ","),
	}
	appendFunCommand(name, funCommand)

	if err := d.DB.Write(funCommandDCollection, name, funCommand); err != nil {
		fmt.Println("addFun", err.Error())
		return errorCommandMsg
	}
	return name + " added!"
}
