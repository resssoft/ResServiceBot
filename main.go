package main

import (
	"encoding/json"
	"fmt"
	"github.com/Syfaro/telegram-bot-api"
	"github.com/nanobox-io/golang-scribble"
	"github.com/patrickmn/go-cache"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type TGUser struct {
	UserID  int
	Login   string
	IsAdmin bool
}

type TGCommand struct {
	Command     string
	CommandType string
	Permissions TGCommandPermissions
}

type TGCommandPermissions struct {
	UserPermissions  string
	ChantPermissions string
}

func splitCommand(command string, separate string) ([]string, string) {
	if command == "" {
		return []string{}, ""
	}
	if separate == "" {
		separate = " "
	}
	result := strings.Split(command, separate)
	return result, strings.Replace(command, result[0]+" ", "", -1)
}

func firstWords(value string, count int) string {
	// Loop over all indexes in the string.
	for i := range value {
		// If we encounter a space, reduce the count.
		if value[i] == ' ' {
			count -= 1
			// When no more words required, return a substring.
			if count == 0 {
				return value[0:i]
			}
		}
	}
	// Return the entire string.
	return value
}

func main() {
	adminId := 180564250
	bot, err := tgbotapi.NewBotAPI("1051149437:AAE9eQ7DZyjXhVnWciitMgypY2fW-SinRDw")
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Work with cache...")
	c := cache.New(95*time.Hour, 100*time.Hour)
	c.Set("admin", strconv.Itoa(adminId), cache.DefaultExpiration)
	c.Set("adminLogin", "@Resager", cache.DefaultExpiration)

	log.Printf("Work with DB...")
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	db, err := scribble.New(dir, nil)
	if err != nil {
		fmt.Println("Error", err)
	}
	existAdmin := TGUser{}
	if err := db.Read("user", strconv.Itoa(adminId), &existAdmin); err != nil {
		fmt.Println("admin not found error", err)
		admin := TGUser{
			UserID:  adminId,
			Login:   "",
			IsAdmin: false,
		}
		log.Printf("admin ID from DB = [%s] ", admin.UserID)
		if err := db.Write("user", strconv.Itoa(adminId), admin); err != nil {
			fmt.Println("Error", err)
		}
	}

	//bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 10

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil && update.InlineQuery != nil {
			continue
		}
		//fmt.Println(update.Message.Text)
		splitedCommands, commandValue := splitCommand(update.Message.Text, " ")
		commandName := splitedCommands[0]

		switch commandName {
		case "/start":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Hi")
			bot.Send(msg)
		case "/addSaveCommand":
			command := TGCommand{
				Command:     commandValue,
				CommandType: "SaveCommand",
				Permissions: TGCommandPermissions{
					UserPermissions:  "",
					ChantPermissions: "",
				},
			}

			if err := db.Write("command", commandValue, command); err != nil {
				fmt.Println("add command error", err)
			}
			//result1 := firstWords(update.Message.Text, 1)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Added"+commandValue)
			bot.Send(msg)
		case "/admin":
			adminLogin, found := c.Get("adminLogin")
			if found {
				fmt.Println(adminLogin)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Admin is "+adminLogin.(string))
				bot.Send(msg)
			}
		default:
			records, err := db.ReadAll("command")
			if err != nil {
				fmt.Println("Error", err)
			}

			commandContain := false
			commands := []TGCommand{}
			for _, f := range records {
				commandFound := TGCommand{}
				if err := json.Unmarshal([]byte(f), &commandFound); err != nil {
					fmt.Println("Error", err)
				}
				commands = append(commands, commandFound)
				if commandFound.Command == commandName {
					commandContain = true
					log.Printf("FOUND command in DB!")
				}
			}

			if !commandContain {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "This is unsupport command.")
				msg.ReplyToMessageID = update.Message.MessageID
				bot.Send(msg)
			}
		}

		//log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		log.Printf("INNER MESSAGE %s[%d]: %s",
			update.Message.From.UserName,
			update.Message.From.ID,
			update.Message.Text)

		//msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		//msg.ReplyToMessageID = update.Message.MessageID

		//bot.Send(msg)
	}

}
