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

const appVersion = "2.0.0"

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

type Configuration struct {
	Telegram TelegramConfig
}

type TelegramConfig struct {
	Bot        TgBot
	AdminId    string
	AdminLogin string
}

type TgBot struct {
	Token string
}

type SavedBlock struct {
	Group string
	User  string
	Text  string
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

func main() {
	fmt.Print("Load configuration... ")
	configurationFile, _ := os.Open("configuration.json")
	defer configurationFile.Close()
	decoder := json.NewDecoder(configurationFile)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("load configuration error:", err)
	}
	fmt.Println(" telegram bot admin is " + configuration.Telegram.AdminLogin)

	bot, err := tgbotapi.NewBotAPI(configuration.Telegram.Bot.Token)
	if err != nil {
		log.Panic(err)
	}

	//TODO: remove this block, duplicate DB - CONFIG - when use cache
	log.Printf("Work with cache...")
	c := cache.New(95*time.Hour, 100*time.Hour)
	c.Set("admin", configuration.Telegram.AdminId, cache.DefaultExpiration)
	c.Set("adminLogin", configuration.Telegram.AdminLogin, cache.DefaultExpiration)

	log.Printf("Work with DB...")
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	//TODO: remove this block, duplicate DB - CONFIG - when use db
	// read admin info from DB or write it to db
	db, err := scribble.New(dir+"/data", nil)
	if err != nil {
		fmt.Println("Error", err)
	}
	existAdmin := TGUser{}
	if err := db.Read("user", configuration.Telegram.AdminId, &existAdmin); err != nil {
		fmt.Println("admin not found error", err)
		adminIdInt, err := strconv.Atoi(configuration.Telegram.AdminId)
		if err != nil {
			fmt.Println("error getting admin ID", err)
		} else {
			existAdmin = TGUser{
				UserID:  adminIdInt,
				Login:   "",
				IsAdmin: false,
			}
			if err := db.Write("user", configuration.Telegram.AdminId, existAdmin); err != nil {
				fmt.Println("Error", err)
			}
		}
	}

	//bot.Debug = true
	adminIdInt64, err := strconv.ParseInt(configuration.Telegram.AdminId, 10, 64)
	if err != nil {
		fmt.Println("error convert admin ID to int64", err)
	} else {
		msg := tgbotapi.NewMessage(adminIdInt64, "Bot Started")
		bot.Send(msg)
	}

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

		//TODO: set permissions for default commands
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
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Added "+commandValue)
			bot.Send(msg)

		case "/SaveCommandsList":
			records, err := db.ReadAll("command")
			if err != nil {
				fmt.Println("Error", err)
			}

			commands := []string{}
			for _, f := range records {
				commandFound := TGCommand{}
				if err := json.Unmarshal([]byte(f), &commandFound); err != nil {
					fmt.Println("Error", err)
				}
				commands = append(commands, commandFound.Command)
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, strings.Join(commands, ", "))
			bot.Send(msg)

		case "/listOf":
			records, err := db.ReadAll("saved")
			if err != nil {
				fmt.Println("Error", err)
			}

			commands := []string{}
			for _, f := range records {
				commandFound := TGCommand{}
				if err := json.Unmarshal([]byte(f), &commandFound); err != nil {
					fmt.Println("Error", err)
				}

				//commandValue
				commands = append(commands, commandFound.Command)
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, strings.Join(commands, ", "))
			bot.Send(msg)

		case "/admin":
			adminLogin, found := c.Get("adminLogin")
			if found {
				fmt.Println(adminLogin)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Admin is "+adminLogin.(string))
				bot.Send(msg)
			}

		case "/version", "/appVersion", "/версия":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, appVersion)
			bot.Send(msg)

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
					itemCode := commandName +
						"_" + strconv.FormatInt(update.Message.Chat.ID, 10) +
						"_" + strconv.FormatInt(time.Now().UnixNano(), 10)
					if err := db.Write(
						"saved",
						itemCode,
						SavedBlock{
							Text:  commandValue,
							Group: commandName,
							User:  strconv.FormatInt(update.Message.Chat.ID, 10),
						}); err != nil {
						fmt.Println("add command error", err)
					} else {
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Done")
						bot.Send(msg)
					}
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
