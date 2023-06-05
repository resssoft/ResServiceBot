package lists

import (
	"encoding/json"
	"fmt"
	tgCommands "fun-coice/internal/domain/commands/tg"
	"fun-coice/pkg/scribble"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"strings"
	"time"
)

type data struct {
	list tgCommands.Commands
	DB   *scribble.Driver
}

func New(DB *scribble.Driver) tgCommands.Service {
	result := data{
		DB: DB,
	}
	commandsList := make(tgCommands.Commands)
	commandsList["addCheckItem"] = tgCommands.Command{
		Command:     "/addCheckItem",
		Description: "(параметры - имя чеклиста, =1 - если публичный, =1 если уже установлен) - создание элемента чеклиста в указанную группу",
		CommandType: "text",
		Permissions: tgCommands.FreePerms,
		Handler:     result.addCheckItem,
	}
	commandsList["updateCheckItem"] = tgCommands.Command{
		Command:     "/updateCheckItem",
		Description: "(параметр - имя чеклиста, =1 или =0 для статуса, полный текст элемента для обновления) - вывод указанной группы чеклиста",
		CommandType: "text",
		Permissions: tgCommands.FreePerms,
		Handler:     result.updateCheckItem,
	}
	commandsList["сheckList"] = tgCommands.Command{
		Command:     "/сheckList",
		Description: "(параметр - имя чеклиста) - вывод указанной группы чеклиста",
		CommandType: "text",
		Permissions: tgCommands.FreePerms,
		Handler:     result.сheckList,
	}
	commandsList["addSaveCommand"] = tgCommands.Command{
		Command:     "/addSaveCommand",
		Description: "Создать комманду сохранения коротких текстовых сообщений, чтобы потом ею сохранять текстовые строки. например. '/addSaveCommand whatToDo' и потом 'whatToDo вымыть посуду'",
		CommandType: "text",
		Permissions: tgCommands.FreePerms,
		Handler:     result.addSaveCommand,
	}
	commandsList["SaveCommandsList"] = tgCommands.Command{
		Command:     "/SaveCommandsList",
		Description: "Список комманд для сохранения текстовых строк",
		CommandType: "text",
		Permissions: tgCommands.FreePerms,
		Handler:     result.SaveCommandsList,
	}
	commandsList["listOf"] = tgCommands.Command{
		Command:     "/listOf",
		Description: "(+ аргумент) Список сохраненных сообщений по указанной комманде",
		CommandType: "text",
		Permissions: tgCommands.FreePerms,
		Handler:     result.listOf,
	}

	//TODO: ADDED save commands

	result.list = commandsList
	return &result
}

func (d data) Commands() tgCommands.Commands {
	return d.list
}

func (d data) addCheckItem(msg *tgbotapi.Message, commandName string, param string, params []string) tgCommands.HandlerResult {

	if len(params) <= 1 {
		return tgCommands.Simple(msg.Chat.ID, "set list name")
	}
	debugMessage := ""
	checkItemText := ""
	checkListGroup := params[1]
	isPublic := false
	checkListStatus := false
	if checkListGroup == "" {
		return tgCommands.Simple(msg.Chat.ID, "need more info, read /commands")
	}
	checkItemText = strings.Replace(param, checkListGroup+" ", "", -1)
	debugMessage += " [" + checkItemText + "] "
	if params[2] == "=1" || params[2] == "isPublic" {
		isPublic = true
		checkItemText = strings.Replace(param, params[2]+" ", "", -1)
		debugMessage += " isPublic "
	}
	if params[3] == "=1" || params[3] == "isCheck" {
		checkItemText = strings.Replace(param, params[3]+" ", "", -1)
		checkListStatus = true
		debugMessage += " checkListStatus "
	}
	debugMessage += " [" + checkItemText + "] "

	checkListItem := CheckList{
		Group:  checkListGroup,
		ChatID: msg.Chat.ID,
		Status: checkListStatus,
		Public: isPublic,
		Text:   checkItemText,
	}

	itemCode := checkListGroup +
		"_" + strconv.FormatInt(msg.Chat.ID, 10) +
		"_" + strconv.FormatInt(time.Now().UnixNano(), 10)

	if err := d.DB.Write("checkList", itemCode, checkListItem); err != nil {
		fmt.Println("add command error", err)
		return tgCommands.EmptyCommand()
	}
	return tgCommands.Simple(msg.Chat.ID, "Added to "+checkListGroup+" debug:"+debugMessage)
}

func (d data) updateCheckItem(msg *tgbotapi.Message, commandName string, param string, params []string) tgCommands.HandlerResult {
	if len(params) <= 1 {
		return tgCommands.Simple(msg.Chat.ID, "set list name")
	}
	checkListGroup := params[1]
	if checkListGroup == "" {
		return tgCommands.Simple(msg.Chat.ID, "need more info, read /commands")
	}

	records, err := d.DB.ReadAll("checkList")
	if err != nil {
		fmt.Println("db read error", err)
	}

	newStatus := false
	if params[1] == "=1" {
		newStatus = true
	}

	checkItemText := strings.Replace(param, params[1]+" ", "", -1)
	updatedItems := 0

	for _, f := range records {
		commandFound := CheckList{}
		if err := json.Unmarshal([]byte(f), &commandFound); err != nil {
			fmt.Println("Error", err)
		}

		if commandFound.Group == checkListGroup && commandFound.ChatID == msg.Chat.ID {
			if commandFound.Text == checkItemText {
				commandFound.Status = newStatus
				if err := d.DB.Write("checkList", checkListGroup, commandFound); err != nil {
					fmt.Println("add command error", err)
				} else {
					updatedItems++
				}
			}
		}
	}
	return tgCommands.Simple(msg.Chat.ID, "update "+strconv.Itoa(updatedItems)+"items")
}

func (d data) сheckList(msg *tgbotapi.Message, commandName string, param string, params []string) tgCommands.HandlerResult {
	if len(params) <= 1 {
		return tgCommands.Simple(msg.Chat.ID, "set list name")
	}
	checkListGroup := params[1]
	if checkListGroup == "" {
		return tgCommands.Simple(msg.Chat.ID, "need more info, read /commands")
	}

	records, err := d.DB.ReadAll("сheckList")
	if err != nil {
		fmt.Println("db read error", err)
	}

	checkListStatusCheck := "✓"
	checkListStatusUnCheck := "○"
	checkListFull := checkListGroup + ":\n"
	for _, f := range records {
		checkListFull += "."
		commandFound := CheckList{}
		if err := json.Unmarshal([]byte(f), &commandFound); err != nil {
			fmt.Println("Error", err)
		}

		checkListFull += "[" + commandFound.Group + " == " + checkListGroup + "] "
		checkListFull += "[" + strconv.FormatInt(commandFound.ChatID, 10) + " == " + strconv.FormatInt(msg.Chat.ID, 10) + "] "
		if commandFound.Group == checkListGroup && commandFound.ChatID == msg.Chat.ID {
			if commandFound.Status == true {
				checkListFull += checkListStatusCheck
			} else {
				checkListFull += checkListStatusUnCheck
			}
			checkListFull += " " + commandFound.Text + "\n"
		}
	}
	return tgCommands.Simple(msg.Chat.ID, checkListFull)
}

func (d data) addSaveCommand(msg *tgbotapi.Message, commandName string, param string, params []string) tgCommands.HandlerResult {
	command := tgCommands.Command{
		Command:     param,
		CommandType: "SaveCommand",
		Permissions: tgCommands.CommandPermissions{
			UserPermissions: "",
			ChatPermissions: "",
		},
	}

	if err := d.DB.Write("command", param, command); err != nil {
		fmt.Println("add command error", err)
		return tgCommands.EmptyCommand()
	}
	return tgCommands.Simple(msg.Chat.ID, "Added "+param)
}

func (d data) SaveCommandsList(msg *tgbotapi.Message, commandName string, param string, params []string) tgCommands.HandlerResult {
	records, err := d.DB.ReadAll("command")
	if err != nil {
		fmt.Println("Error", err)
	}

	commands := []string{}
	for _, f := range records {
		commandFound := tgCommands.Command{}
		if err := json.Unmarshal([]byte(f), &commandFound); err != nil {
			fmt.Println("Error", err)
			return tgCommands.EmptyCommand()
		}
		commands = append(commands, commandFound.Command)
	}
	return tgCommands.SimpleReply(msg.Chat.ID, strings.Join(commands, ", "), msg.MessageID)
}

func (d data) listOf(msg *tgbotapi.Message, commandName string, param string, params []string) tgCommands.HandlerResult {
	records, err := d.DB.ReadAll("saved")
	if err != nil {
		fmt.Println("Error", err)
	}
	commands := []string{}
	for _, f := range records {
		commandFound := SavedBlock{}
		if err := json.Unmarshal([]byte(f), &commandFound); err != nil {
			fmt.Println("Error", err)
		}
		if commandFound.Group == param && commandFound.User == strconv.FormatInt(msg.Chat.ID, 10) {
			commands = append(commands, commandFound.Text)
		}
	}
	return tgCommands.Simple(msg.Chat.ID, param+":\n-"+strings.Join(commands, "\n-"))
}