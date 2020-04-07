package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/Syfaro/telegram-bot-api"
	"github.com/nanobox-io/golang-scribble"
	"github.com/patrickmn/go-cache"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

const appVersion = "2.0.009dg5"
const doneMessage = "Done"
const telegramSingleMessageLengthLimit = 4096

type TGUser struct {
	UserID  int
	Login   string
	IsAdmin bool
}

type TGCommand struct {
	Command     string
	Description string
	CommandType string
	Permissions TGCommandPermissions
}

type TGCommandPermissions struct {
	UserPermissions string
	ChatPermissions string
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

type CheckList struct {
	Group  string
	ChatID int64
	Text   string
	Status bool
	Public bool
}

var commands = map[string]TGCommand{
	"addCheckItem": {
		Command:     "/addCheckItem",
		Description: "(параметры - имя чеклиста, =1 - если публичный, =1 если уже установлен) - создание элемента чеклиста в указанную группу",
		CommandType: "tg",
		Permissions: TGCommandPermissions{
			ChatPermissions: "all",
			UserPermissions: "all",
		},
	},
	"updateCheckItem": {
		Command:     "/updateCheckItem",
		Description: "(параметр - имя чеклиста, =1 или =0 для статуса, полный текст элемента для обновления) - вывод указанной группы чеклиста",
		CommandType: "tg",
		Permissions: TGCommandPermissions{
			ChatPermissions: "all",
			UserPermissions: "all",
		},
	},
	"сheckList": {
		Command:     "/сheckList",
		Description: "(параметр - имя чеклиста) - вывод указанной группы чеклиста",
		CommandType: "tg",
		Permissions: TGCommandPermissions{
			ChatPermissions: "all",
			UserPermissions: "all",
		},
	},
	"start": {
		Command:     "/start",
		Description: "Регистрация в сервисе",
		CommandType: "tg",
		Permissions: TGCommandPermissions{
			ChatPermissions: "all",
			UserPermissions: "all",
		},
	},
	"addSaveCommand": {
		Command:     "/addSaveCommand",
		Description: "Создать комманду сохранения коротких текстовых сообщений, чтобы потом ею сохранять текстовые строки. например. '/addSaveCommand whatToDo' и потом 'whatToDo вымыть посуду'",
		CommandType: "tg",
		Permissions: TGCommandPermissions{
			ChatPermissions: "moder",
			UserPermissions: "moder",
		},
	},
	"addFeature": {
		Command:     "/addFeature",
		Description: "Создание описание фичи",
		CommandType: "tg",
		Permissions: TGCommandPermissions{
			ChatPermissions: "all",
			UserPermissions: "all",
		},
	},
	"getFeatures": {
		Command:     "/getFeatures",
		Description: "Список фич приложения",
		CommandType: "tg",
		Permissions: TGCommandPermissions{
			ChatPermissions: "all",
			UserPermissions: "all",
		},
	},
	"SaveCommandsList": {
		Command:     "/SaveCommandsList",
		Description: "Список комманд для сохранения текстовых строк",
		CommandType: "tg",
		Permissions: TGCommandPermissions{
			ChatPermissions: "all",
			UserPermissions: "all",
		},
	},
	"listOf": {
		Command:     "/listOf",
		Description: "(+ аргумент) Список сохраненных сообщений по указанной комманде",
		CommandType: "tg",
		Permissions: TGCommandPermissions{
			ChatPermissions: "all",
			UserPermissions: "all",
		},
	},
	"admin": {
		Command:     "/admin",
		Description: "Вывод логина админа",
		CommandType: "tg",
		Permissions: TGCommandPermissions{
			ChatPermissions: "all",
			UserPermissions: "all",
		},
	},
	"version": {
		Command:     "/version",
		Description: "Вывод версии",
		CommandType: "tg",
		Permissions: TGCommandPermissions{
			ChatPermissions: "all",
			UserPermissions: "all",
		},
	},
	"appVersion": {
		Command:     "/appVersion",
		Description: "синоним version",
		CommandType: "tg",
		Permissions: TGCommandPermissions{
			ChatPermissions: "all",
			UserPermissions: "all",
		},
	},
	"версия": {
		Command:     "/версия",
		Description: "синоним version",
		CommandType: "tg",
		Permissions: TGCommandPermissions{
			ChatPermissions: "all",
			UserPermissions: "all",
		},
	},
	"commands": {
		Command:     "/commands",
		Description: "Список комманд",
		CommandType: "tg",
		Permissions: TGCommandPermissions{
			ChatPermissions: "all",
			UserPermissions: "all",
		},
	},
	"rebuild": {
		Command:     "/rebuild",
		Description: "rebuild",
		CommandType: "tg",
		Permissions: TGCommandPermissions{
			ChatPermissions: "admin",
			UserPermissions: "admin",
		},
	},
	"games": {
		Command:     "/games",
		Description: "games list",
		CommandType: "tg",
		Permissions: TGCommandPermissions{
			ChatPermissions: "all",
			UserPermissions: "all",
		},
	},
}

type ChatUserCount struct {
	ChatId      int64
	ChatName    string
	ContentType string
	UserCount   int
}

var ChatUserCountList = make([]ChatUserCount, 1)

var gamesListKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("🧡 Lovely game", "lovelyGame"),
		tgbotapi.NewInlineKeyboardButtonURL("Rules", "http://1073.ru/games/lovely/rules/"),
	),
)

func getChannelUserCount(contentType string, chatId int64) int {
	for _, item := range ChatUserCountList {
		if item.ChatId == chatId && item.ContentType == contentType {
			return item.UserCount
		}
	}
	return 0
}

func IncreaseChannelUserCount(contentType string, chatId int64, chatName string) {
	founded := false
	for _, item := range ChatUserCountList {
		if item.ChatId == chatId && item.ContentType == contentType {
			item.UserCount = item.UserCount + 1
			founded = true
		}
	}
	if !founded {
		ChatUserCountList = append(
			ChatUserCountList,
			ChatUserCount{chatId, chatName, contentType, 1},
		)
	}
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

func writeLines(lines []string, path string) error {

	// overwrite file if it exists
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}

	defer file.Close()

	// new writer w/ default 4096 buffer size
	w := bufio.NewWriter(file)

	for _, line := range lines {
		_, err := w.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}

	// flush outstanding data
	return w.Flush()
}

func checkPermission(command string, userId int) (error, bool) {
	typeOfCommand := commands[command].Permissions.UserPermissions
	switch typeOfCommand {
	case "all":
		return nil, true
	case "admin":
		if userId == existAdmin.UserID {
			return nil, true
		} else {
			return nil, false
		}
	}
	return nil, true
}

func readLines(path string, resultLimit int) (error, string) {
	result := ""
	file, err := os.OpenFile(path, os.O_RDONLY, 0666)
	if err != nil {
		return err, ""
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	stringLen := 0
	for scanner.Scan() {
		result += scanner.Text() + "\n"
		fmt.Println(result)
		stringLen = utf8.RuneCountInString(result)
		if stringLen > resultLimit {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return err, ""
	}
	return nil, ""
}

var existAdmin = TGUser{}

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
		msg := tgbotapi.NewMessage(adminIdInt64, "Bot Started with version "+appVersion)
		bot.Send(msg)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 10

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.CallbackQuery != nil {
			fmt.Printf("CallbackQuery %+v\n", update.CallbackQuery)
			bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data))
			splitedCallbackQuery, clearCallbackQuery := splitCommand(update.Message.Text, " ")
			commandsCount := len(splitedCallbackQuery)
			callbackQueryMessageChatID := 0
			if commandsCount == 0 {
				continue
			}
			callbackQueryMessageChatID, _ = strconv.Atoi(splitedCallbackQuery[0])

			switch clearCallbackQuery {
			case "lovelyGame":
				messageID := strconv.Itoa(update.CallbackQuery.Message.MessageID)
				buttonText := "Join to Lovely game start (" +
					strconv.Itoa(getChannelUserCount(
						"lovelyGame",
						update.CallbackQuery.Message.Chat.ID)) + ")"
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "-")
				msg.ReplyMarkup = tgbotapi.NewEditMessageReplyMarkup(
					update.CallbackQuery.Message.Chat.ID,
					update.CallbackQuery.Message.MessageID,
					tgbotapi.NewInlineKeyboardMarkup(
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData(buttonText, messageID+"#lovelyGameJoin"),
						),
					),
				)
				bot.Send(msg)
			case "lovelyGameJoin":
				IncreaseChannelUserCount(
					"lovelyGame",
					update.CallbackQuery.Message.Chat.ID,
					update.CallbackQuery.Message.Chat.Title)
				messageID := strconv.Itoa(callbackQueryMessageChatID)
				buttonText := "Join to Lovely game (" +
					strconv.Itoa(getChannelUserCount(
						"lovelyGame",
						update.CallbackQuery.Message.Chat.ID)) + ")"
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "-")
				msg.ReplyMarkup = tgbotapi.NewEditMessageReplyMarkup(
					update.CallbackQuery.Message.Chat.ID,
					update.CallbackQuery.Message.MessageID,
					tgbotapi.NewInlineKeyboardMarkup(
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData(buttonText, messageID+"#lovelyGameJoin"),
							tgbotapi.NewInlineKeyboardButtonData("End joins and start", messageID+"#lovelyGameJoin"),
						),
					),
				)
				bot.Send(msg)
			default:
				bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Data: "+update.CallbackQuery.Data))
			}

		}
		fmt.Printf("inline query %+v\n", update.InlineQuery)
		if update.Message == nil || (update.Message == nil && update.InlineQuery != nil) {
			continue
		}
		//fmt.Println(update.Message.Text)
		splitedCommands, commandValue := splitCommand(update.Message.Text, " ")
		commandsCount := len(splitedCommands)
		if commandsCount == 0 {
			continue
		}
		commandName := splitedCommands[0]

		//TODO: set permissions for default commands
		switch commandName {
		case "/start":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Hi")
			bot.Send(msg)

		case "/rebuild":
			err, permission := checkPermission("rebuild", update.Message.From.ID)
			if err != nil {
				log.Printf("Failed permissions: %v", err)
			}
			if permission {
				dir, err := os.Getwd()
				if err != nil {
					log.Printf("Failed to get dir: %v", err)
				}
				cmd := exec.Command("/bin/sh", dir+"/rebuild.sh")
				if err := cmd.Start(); err != nil {
					log.Printf("Failed to start cmd: %v", err)
				}

				log.Println("Exit by command...")

				os.Exit(3)
			}

		case "/commands":
			commandsList := "Commands:\n"
			for _, commandsItem := range commands {
				commandsList += commandsItem.Command + " - " + commandsItem.Description + "\n"
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, commandsList)
			bot.Send(msg)

		case "/addSaveCommand":
			command := TGCommand{
				Command:     commandValue,
				CommandType: "SaveCommand",
				Permissions: TGCommandPermissions{
					UserPermissions: "",
					ChatPermissions: "",
				},
			}

			if err := db.Write("command", commandValue, command); err != nil {
				fmt.Println("add command error", err)
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Added "+commandValue)
			bot.Send(msg)

		case "/addFeature":
			currentTime := time.Now().Format(time.RFC3339)
			formattedMessage := currentTime + "[" + appVersion + "]: " + commandValue
			err := writeLines([]string{formattedMessage}, "./features.txt")
			if err != nil {
				fmt.Println("write command error", err)
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, doneMessage)
			bot.Send(msg)

		case "/games":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Games list")
			msg.ReplyMarkup = gamesListKeyboard
			bot.Send(msg)

		case "/getFeatures":
			err, messages := readLines("./features.txt", telegramSingleMessageLengthLimit)
			if err != nil {
				fmt.Println("write command error", err)
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, messages)
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
				commandFound := SavedBlock{}
				if err := json.Unmarshal([]byte(f), &commandFound); err != nil {
					fmt.Println("Error", err)
				}

				if commandFound.Group == commandValue && commandFound.User == strconv.FormatInt(update.Message.Chat.ID, 10) {
					commands = append(commands, commandFound.Text)
				}
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, commandValue+":\n-"+strings.Join(commands, "\n-"))
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

		case commands["addCheckItem"].Command:
			debugMessage := ""
			checkItemText := ""
			checkListGroup := splitedCommands[1]
			isPublic := false
			checkListStatus := false
			if checkListGroup == "" {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "need more info, read /commands")
				bot.Send(msg)
				break
			}
			checkItemText = strings.Replace(commandValue, checkListGroup+" ", "", -1)
			debugMessage += " [" + checkItemText + "] "
			if splitedCommands[2] == "=1" || splitedCommands[2] == "isPublic" {
				isPublic = true
				checkItemText = strings.Replace(commandValue, splitedCommands[2]+" ", "", -1)
				debugMessage += " isPublic "
			}
			if splitedCommands[3] == "=1" || splitedCommands[3] == "isCheck" {
				checkItemText = strings.Replace(commandValue, splitedCommands[3]+" ", "", -1)
				checkListStatus = true
				debugMessage += " checkListStatus "
			}
			debugMessage += " [" + checkItemText + "] "

			checkListItem := CheckList{
				Group:  checkListGroup,
				ChatID: update.Message.Chat.ID,
				Status: checkListStatus,
				Public: isPublic,
				Text:   checkItemText,
			}

			itemCode := checkListGroup +
				"_" + strconv.FormatInt(update.Message.Chat.ID, 10) +
				"_" + strconv.FormatInt(time.Now().UnixNano(), 10)

			if err := db.Write("checkList", itemCode, checkListItem); err != nil {
				fmt.Println("add command error", err)
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Added to "+checkListGroup+" debug:"+debugMessage)
			bot.Send(msg)

		case commands["updateCheckItem"].Command:
			checkListGroup := splitedCommands[1]
			if checkListGroup == "" {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "need more info, read /commands")
				bot.Send(msg)
				break
			}

			records, err := db.ReadAll("checkList")
			if err != nil {
				fmt.Println("db read error", err)
			}

			newStatus := false
			if splitedCommands[1] == "=1" {
				newStatus = true
			}

			checkItemText := strings.Replace(commandValue, splitedCommands[1]+" ", "", -1)
			updatedItems := 0

			for _, f := range records {
				commandFound := CheckList{}
				if err := json.Unmarshal([]byte(f), &commandFound); err != nil {
					fmt.Println("Error", err)
				}

				if commandFound.Group == checkListGroup && commandFound.ChatID == update.Message.Chat.ID {
					if commandFound.Text == checkItemText {
						commandFound.Status = newStatus
						if err := db.Write("checkList", checkListGroup, commandFound); err != nil {
							fmt.Println("add command error", err)
						} else {
							updatedItems++
						}
					}
				}
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "update "+strconv.Itoa(updatedItems)+"items")
			bot.Send(msg)

		case commands["сheckList"].Command:
			checkListGroup := splitedCommands[1]
			if checkListGroup == "" {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "need more info, read /commands")
				bot.Send(msg)
				break
			}

			records, err := db.ReadAll("сheckList")
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
				checkListFull += "[" + strconv.FormatInt(commandFound.ChatID, 10) + " == " + strconv.FormatInt(update.Message.Chat.ID, 10) + "] "
				if commandFound.Group == checkListGroup && commandFound.ChatID == update.Message.Chat.ID {
					if commandFound.Status == true {
						checkListFull += checkListStatusCheck
					} else {
						checkListFull += checkListStatusUnCheck
					}
					checkListFull += " " + commandFound.Text + "\n"
				}
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, checkListFull)
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
