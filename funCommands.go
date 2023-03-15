package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
)

var errorCommandMsg = "Error, sorry! Write to admin or send command /bug with command name"

type FunCommand struct {
	Name      string
	TgCommand TGCommand
	List1     []string
	List2     []string
}

var FunCommands map[string]FunCommand
var syncMap *sync.Mutex
var funCommandDCollection = "funcommands"
var funCommandType = "funcommand"

func initFunCommand() {
	FunCommands = make(map[string]FunCommand)
	syncMap = new(sync.Mutex)

	records, err := DB.ReadAll(funCommandDCollection)
	if err != nil {
		fmt.Println("Error DB.ReadAll", err)
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
		appendFunCommand(funCommand.Name, funCommand)
		fmt.Println("Add fun command", funCommand.Name)
	}
	fmt.Println("FunCommands", FunCommands)
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

func addFunCommand2(command, list1, list2 string) string {
	text := "Added!"
	err := os.Mkdir(appPath+string(os.PathSeparator), os.ModeDir)
	if err != nil {
		fmt.Println("addFun", err.Error())
		return errorCommandMsg
	}
	f, err := os.Create(command + ".txt")
	if err != nil {
		fmt.Println("addFun", err.Error())
		return errorCommandMsg
	}
	defer f.Close()
	_, err = f.WriteString(command + "\n")
	if err != nil {
		fmt.Println("addFun", err.Error())
		return errorCommandMsg
	}
	f.WriteString(list1 + "\n")
	f.WriteString(list2 + "\n")
	return text
}

func addFunCommand(name, list1, list2 string) string {
	command := TGCommand{
		Command:     name,
		Synonyms:    nil,
		Description: "Get random fanny words",
		CommandType: funCommandType,
		Permissions: freePerms,
	}
	funCommand := FunCommand{
		Name:      name,
		TgCommand: command,
		List1:     strings.Split(list1, ","),
		List2:     strings.Split(list2, ","),
	}
	appendFunCommand(name, funCommand)

	if err := DB.Write(funCommandDCollection, name, funCommand); err != nil {
		fmt.Println("addFun", err.Error())
		return errorCommandMsg
	}
	return name + " added!"
}
