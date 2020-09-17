package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Syfaro/telegram-bot-api"
	"github.com/nanobox-io/golang-scribble"
	"github.com/patrickmn/go-cache"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

const appVersion = "2.0.014dg86"
const doneMessage = "Done"
const telegramSingleMessageLengthLimit = 4096
const HWCSURLEvent = "go"
const HWCSURLImage = "result/"

var HWCSURL = ""

type homeWebCamServiceURLImageData struct {
	Result string `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}

func getHomeWebCamImage() (string, error) {
	HWCSData := homeWebCamServiceURLImageData{}
	resp, err := http.Get(HWCSURL + HWCSURLEvent)
	if err != nil {
		log.Printf("Status: %v Error: %v \n", err.Error())
	}
	defer resp.Body.Close()
	jsonData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Read body error: %v \n", err.Error())
		return "", err
	}

	if err = json.Unmarshal(jsonData, &HWCSData); err != nil {
		log.Printf("Unmarshal error: %s \n", err.Error())
		return "", err
	}
	return HWCSURL + HWCSURLImage + HWCSData.Result, nil
}

type TGUser struct {
	UserID  int
	ChatId  int64
	Login   string
	Name    string
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
	HWCSURL  string
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

type KeyBoardTG struct {
	Rows []KeyBoardRowTG
}

type KeyBoardRowTG struct {
	Buttons []KeyBoardButtonTG
}

type KeyBoardButtonTG struct {
	Text string
	Data string
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
		Description: "Service registration, only private",
		CommandType: "tg",
		Permissions: TGCommandPermissions{
			ChatPermissions: "all",
			UserPermissions: "all",
		},
	},
	"myInfo": {
		Command:     "/myInfo",
		Description: "Write GT user info",
		CommandType: "tg",
		Permissions: TGCommandPermissions{
			ChatPermissions: "all",
			UserPermissions: "all",
		},
	},
	"getUserList": {
		Command:     "/getUserList",
		Description: "-",
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
	"homeweb": {
		Command:     "/homeweb",
		Description: "get image link from cam1",
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

type ChatUser struct {
	ChatId      int64
	ChatName    string
	ContentType string
	CustomRole  string
	VoteCount   int
	User        TGUser
}

var ChatUserList = make([]ChatUser, 1)

var gamesListKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("🧡 Lovely game", "lovelyGame"),
		tgbotapi.NewInlineKeyboardButtonURL("Rules", "http://1073.ru/games/lovely/rules/"),
	),
)

func getChannelUserCount(contentType string, chatId int64) int {
	userCount := 0
	for _, item := range ChatUserList {
		if item.ChatId == chatId && item.ContentType == contentType && item.CustomRole != "dead" {
			userCount++
		}
	}
	return userCount
}

func getChannelUserMaxVoted(contentType string, chatId int64) (ChatUser, int, []ChatUser) {
	maxVote := 0
	maxVotedUser := ChatUser{}
	var maxVotedUsers []ChatUser
	for _, item := range ChatUserList {
		if item.ChatId == chatId &&
			item.ContentType == contentType &&
			item.CustomRole != "dead" {
			if item.VoteCount > maxVote {
				maxVote = item.VoteCount
				maxVotedUser = item
			}
		}
	}
	if maxVote == 0 {
		return ChatUser{}, 0, make([]ChatUser, 0)
	} else {
		maxVotedUsers = append(maxVotedUsers, maxVotedUser)
	}
	// get more players with max voteCount
	maxVoteCount := 1
	for _, item := range ChatUserList {
		if item.ChatId == chatId &&
			item.ContentType == contentType &&
			item.CustomRole != "dead" &&
			item.User.UserID != maxVotedUser.User.UserID {
			if item.VoteCount == maxVote {
				maxVotedUser = item
				maxVoteCount++
				maxVotedUsers = append(maxVotedUsers, maxVotedUser)
			}
		}
	}
	return maxVotedUser, maxVoteCount, maxVotedUsers
}

func getChannelUsers(contentType string, chatId int64) string {
	users := ""
	var userList []string
	for _, item := range ChatUserList {
		if item.ChatId == chatId && item.ContentType == contentType {
			userList = append(userList, item.User.Name)
		}
	}
	if len(userList) > 0 {
		users = strings.Join(userList, "\n")
	}
	return users
}

func removeChannelUsers(contentType string, chatId int64) {
	var ChatUserListNew = make([]ChatUser, 1)
	for _, item := range ChatUserList {
		if !(item.ChatId == chatId && item.ContentType == contentType) {
			ChatUserListNew = append(ChatUserListNew, item)
		}
	}
	ChatUserList = ChatUserListNew
}

func setZeroCountsChannelUsersList(contentType string, chatId int64) {
	for itemIndex, item := range ChatUserList {
		if item.ChatId == chatId && item.ContentType == contentType {
			ChatUserList[itemIndex].VoteCount = 0
		}
	}
}

func getUsersVoteMessageConfig(contentType string, chatID int64, messageText string) tgbotapi.MessageConfig {
	activeChatUsers := getChannelUsersList(contentType, chatID, false)
	buttons := getUsersButtons(activeChatUsers, chatID, "lovelyGamePlayerVoteChoice")
	msg := tgbotapi.NewMessage(
		chatID,
		"Voting")
	msg.ReplyMarkup = buttons
	return msg
}

func updateUsersVoteMessageConfig(contentType string, chatID int64, messageText string, messageID int) tgbotapi.EditMessageTextConfig {
	activeChatUsers := getChannelUsersList(contentType, chatID, false)
	buttons := getUsersButtons(activeChatUsers, chatID, "lovelyGamePlayerVoteChoice")
	msg := tgbotapi.NewEditMessageText(
		chatID,
		messageID,
		messageText)
	msg.ReplyMarkup = tgbotapi.NewEditMessageReplyMarkup(
		chatID,
		messageID,
		buttons,
	).ReplyMarkup
	return msg
}

func incCountsChannelUsersList(contentType string, chatId int64, userId int) {
	for itemIndex, item := range ChatUserList {
		if item.ChatId == chatId && item.ContentType == contentType && item.User.UserID == userId {
			ChatUserList[itemIndex].VoteCount += 1
		}
	}
}

func getCountsChannelUsersList(contentType string, chatId int64) int {
	sum := 0
	for _, item := range ChatUserList {
		if item.ChatId == chatId && item.ContentType == contentType {
			sum += item.VoteCount
		}
	}
	return sum
}

func getChannelUsersList(contentType string, chatId int64, excludeActiveRole bool) []ChatUser {
	var userList []ChatUser
	for _, item := range ChatUserList {
		var userIncluded = true
		if item.CustomRole == "dead" {
			userIncluded = false
		}
		if excludeActiveRole && item.CustomRole == "killer" {
			userIncluded = false
		}
		if item.ChatId == chatId && item.ContentType == contentType && userIncluded {
			userList = append(userList, item)
		}
	}
	return userList
}

func getChannelUser(contentType string, chatId int64, userId int) (ChatUser, error) {
	for _, item := range ChatUserList {
		if item.ChatId == chatId && item.ContentType == contentType && item.User.UserID == userId {
			return item, nil
		}
	}
	return ChatUser{}, errors.New("user not found")
}

func getUsersButtons(chatUsers []ChatUser, chatID int64, code string) tgbotapi.InlineKeyboardMarkup {
	var rows []KeyBoardRowTG
	for _, chatUser := range chatUsers {
		rows = append(rows, KBButs(KeyBoardButtonTG{
			Text: chatUser.User.Name + " (" + strconv.Itoa(chatUser.VoteCount) + ")",
			Data: strconv.Itoa(chatUser.User.UserID) + "|" + strconv.FormatInt(chatID, 10) + "#" + code,
		}))
	}
	return getTGButtons(KeyBoardTG{rows})
}

func sendRoleToUser(bot *tgbotapi.BotAPI, chatID int64, contentType string) {
	chatUsers := getChannelUsersList(contentType, chatID, true)
	random := rand.New(rand.NewSource(time.Now().Unix()))
	user := chatUsers[random.Intn(len(chatUsers))]
	SetUserRoleToChannelList(contentType, chatID, user.User.UserID, "killer")
	time.Sleep(5 * time.Second)
	msg := tgbotapi.NewMessage(int64(user.User.UserID), "Please, choice:")
	msg.ReplyMarkup = getUsersButtons(chatUsers, chatID, "lovelyGamePlayerChoice")
	messageID, _ := bot.Send(msg)

	fmt.Printf("messageID %+v\n", messageID)
}

func SaveUserToChannelList(contentType string, chatId int64, chatName string, userId int, userName string) bool {
	isNewUser := true
	for _, item := range ChatUserList {
		if item.ChatId == chatId && item.ContentType == contentType && item.User.UserID == userId {
			isNewUser = false
		}
	}
	_, isAdmin := checkPermission("admin", userId)
	if isNewUser {
		ChatUserList = append(
			ChatUserList,
			ChatUser{
				ChatId:      chatId,
				ChatName:    chatName,
				ContentType: contentType,
				CustomRole:  "",
				VoteCount:   0,
				User: TGUser{
					UserID:  userId,
					ChatId:  0,
					Name:    userName,
					Login:   userName,
					IsAdmin: isAdmin,
				},
			},
		)
	}
	return checkUserRegister(userId)
}

func SetUserRoleToChannelList(contentType string, chatId int64, userId int, userRole string) {
	for itemIndex, item := range ChatUserList {
		if item.ChatId == chatId && item.ContentType == contentType && item.User.UserID == userId {
			ChatUserList[itemIndex].CustomRole = userRole
		}
	}
}

func checkUserRegister(userId int) bool {
	// check - bot can write to user
	isRegistered := false
	var existUser = TGUser{}
	err := DB.Read("user", strconv.Itoa(userId), &existUser)
	if err == nil {
		if existUser.ChatId != 0 {
			isRegistered = true
		}
	}
	return isRegistered
}

func unregisteredChannelUsers(contentType string, chatId int64) string {
	users := ""
	var userList []string
	for _, item := range ChatUserList {
		if item.ChatId == chatId && item.ContentType == contentType {
			if !checkUserRegister(item.User.UserID) {
				userList = append(userList, item.User.Name)
			}
		}
	}
	if len(userList) > 0 {
		users = strings.Join(userList, "\n")
	}
	return users
}

func splitCommand(command string, separate string) ([]string, string) {
	if command == "" {
		return []string{}, ""
	}
	if separate == "" {
		separate = " "
	}
	result := strings.Split(command, separate)
	return result, strings.Replace(command, result[0]+separate, "", -1)
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

func KBRows(KBrows ...KeyBoardRowTG) KeyBoardTG {
	var rows []KeyBoardRowTG
	rows = append(rows, KBrows...)
	return KeyBoardTG{rows}
}

func KBButs(KBrows ...KeyBoardButtonTG) KeyBoardRowTG {
	var rows []KeyBoardButtonTG
	rows = append(rows, KBrows...)
	return KeyBoardRowTG{rows}
}

func getSimpleTGButton(text, data string) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(text, data),
		),
	)
}

func getTGButtons(params KeyBoardTG) tgbotapi.InlineKeyboardMarkup {
	var row []tgbotapi.InlineKeyboardButton
	var rows [][]tgbotapi.InlineKeyboardButton
	for _, rowsData := range params.Rows {
		for _, button := range rowsData.Buttons {
			row = append(row, tgbotapi.NewInlineKeyboardButtonData(button.Text, button.Data))
		}
		rows = append(rows, row)
		row = nil
	}
	return tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: rows,
	}
}

var existAdmin = TGUser{}
var DB *scribble.Driver

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
	HWCSURL = configuration.HWCSURL
	log.Printf("Admin is ..." + configuration.Telegram.AdminLogin)

	log.Printf("Work with DB...")
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	//TODO: remove this block, duplicate DB - CONFIG - when use db
	// read admin info from DB or write it to db
	DB, err = scribble.New(dir+"/data", nil)
	if err != nil {
		fmt.Println("Error", err)
	}
	if err := DB.Read("user", configuration.Telegram.AdminId, &existAdmin); err != nil {
		fmt.Println("admin not found error", err)
		adminIdInt, err := strconv.Atoi(configuration.Telegram.AdminId)
		if err != nil {
			fmt.Println("error getting admin ID", err)
		} else {
			existAdmin = TGUser{
				UserID:  adminIdInt,
				ChatId:  0,
				Login:   "",
				Name:    "",
				IsAdmin: false,
			}
			if err := DB.Write("user", configuration.Telegram.AdminId, existAdmin); err != nil {
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
			from := update.CallbackQuery.From
			fromName := update.CallbackQuery.From.String()
			chat := update.CallbackQuery.Message.Chat
			messageID := update.CallbackQuery.Message.MessageID
			contentType := "lovelyGame"
			//debug
			fmt.Printf("update.CallbackQuery %+v\n", update.CallbackQuery)
			fmt.Printf("update.CallbackQuery.Message %+v\n", update.CallbackQuery.Message)
			fmt.Printf("update.CallbackQuery.Message.Chat %+v\n", chat)
			fmt.Printf("update.CallbackQuery.From %+v %+v\n", from.ID, from.UserName)

			bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data))
			splitedCallbackQuery, clearCallbackQuery := splitCommand(update.CallbackQuery.Data, "#")
			commandsCount := len(splitedCallbackQuery)

			fmt.Printf("clearCallbackQuery %+v\n", clearCallbackQuery)
			switch clearCallbackQuery {
			case "lovelyGame":
				removeChannelUsers(contentType, chat.ID)
				buttonText := "Join (" + strconv.Itoa(getChannelUserCount(contentType, chat.ID)) + ")"
				msg := tgbotapi.NewEditMessageText(
					chat.ID,
					messageID,
					"Please, join to game.")
				msg.ReplyMarkup = tgbotapi.NewEditMessageReplyMarkup(
					chat.ID,
					messageID,
					getSimpleTGButton(buttonText, "lovelyGameJoin"),
				).ReplyMarkup
				bot.Send(msg)

			case "lovelyGameJoin":
				isRegisteredUser := SaveUserToChannelList(
					contentType,
					chat.ID,
					chat.Title,
					from.ID,
					from.String(),
				)
				if !isRegisteredUser {
					bot.Send(tgbotapi.NewMessage(chat.ID, from.String()+", write me to private for register"))
				}
				buttonText := "Join (" +
					strconv.Itoa(getChannelUserCount(
						contentType,
						chat.ID)) + ")"
				msg := tgbotapi.NewEditMessageText(
					chat.ID,
					messageID,
					"Please, join to game. After team complete, click to end joins")
				msg.ReplyMarkup = tgbotapi.NewEditMessageReplyMarkup(
					chat.ID,
					messageID,
					getTGButtons(KBRows(KBButs(
						KeyBoardButtonTG{buttonText, "lovelyGameJoin"},
						KeyBoardButtonTG{"End joins and start", "lovelyGameStart"},
					))),
				).ReplyMarkup
				bot.Send(msg)

			case "lovelyGameStart":
				messageText := ""
				unregisteredUsers := unregisteredChannelUsers(contentType, chat.ID)
				if unregisteredUsers != "" {
					messageText = "I can`t start, unregistered users: " + unregisteredUsers
					bot.Send(tgbotapi.NewMessage(chat.ID, messageText))
				} else {
					messageText = "Start lovely Game with: \n" +
						getChannelUsers(contentType, chat.ID) +
						"\n Wait for the killer to choose a player..."
					go sendRoleToUser(bot, chat.ID, contentType)
					msg := tgbotapi.NewEditMessageText(
						chat.ID,
						messageID,
						messageText)
					bot.Send(msg)
				}

			case "lovelyGameVoting":
				setZeroCountsChannelUsersList(contentType, chat.ID)
				bot.Send(getUsersVoteMessageConfig(contentType, chat.ID, "Start voting"))

			case "lovelyGamePlayerVoteChoice":
				messageText := ""
				if commandsCount <= 1 {
					continue
				}
				customDataItems, _ := splitCommand(splitedCallbackQuery[0], "|")
				customDataItemsCount := len(customDataItems)
				if customDataItemsCount <= 1 {
				}
				choicedUserID, _ := strconv.Atoi(customDataItems[0])
				mainChatID := customDataItems[1]
				mainChatIDInt64, _ := strconv.ParseInt(mainChatID, 10, 64)
				chatUser, _ := getChannelUser(contentType, mainChatIDInt64, choicedUserID)
				incCountsChannelUsersList(contentType, mainChatIDInt64, choicedUserID)
				voteSum := getCountsChannelUsersList(contentType, mainChatIDInt64)
				usersCount := getChannelUserCount(contentType, mainChatIDInt64)
				messageText = fromName + " voted for: " + chatUser.User.Name
				bot.Send(tgbotapi.NewMessage(mainChatIDInt64, messageText))
				if voteSum == usersCount {
					votedUser, voteUsersCount, votedUsers := getChannelUserMaxVoted(contentType, mainChatIDInt64)
					if 1 == voteUsersCount {
						SetUserRoleToChannelList(contentType, mainChatIDInt64, choicedUserID, "dead")
						if votedUser.CustomRole == "killer" {
							messageText = "Killer is dead and game of ending"
						} else if usersCount <= 2 { //TODO: check minimal users to 3
							messageText = "Game of ending. Killer won"
						} else {
							messageText = "Wait for the killer to choose a player..."
							go sendRoleToUser(bot, chat.ID, contentType)
						}
						msg := tgbotapi.NewEditMessageText(
							chat.ID,
							messageID,
							messageText)
						bot.Send(msg)
						// vote again
						//bot.Send(getUsersVoteMessageConfig(contentType, chat.ID, "Voting"))

					} else {
						messageText = "Multiple voting: "
						for _, votedUsersItem := range votedUsers {
							messageText += "\n" + votedUsersItem.User.Name
						}
						msg := tgbotapi.NewEditMessageText(
							chat.ID,
							messageID,
							messageText)
						msg.ReplyMarkup = tgbotapi.NewEditMessageReplyMarkup(
							chat.ID,
							messageID,
							getTGButtons(KBRows(KBButs(
								KeyBoardButtonTG{"Vote again", "lovelyGameVoting"},
								KeyBoardButtonTG{"Vote all", "lovelyGameVotingAll"},
							))),
						).ReplyMarkup
						bot.Send(msg)
					}
				} else {
					bot.Send(updateUsersVoteMessageConfig(contentType, mainChatIDInt64, "Voting", messageID))
				}

			case "lovelyGameVotingAll":
				continue

			case "lovelyGamePlayerChoice":
				if commandsCount <= 1 {
					continue
				}
				customDataItems, _ := splitCommand(splitedCallbackQuery[0], "|")
				customDataItemsCount := len(customDataItems)
				if customDataItemsCount > 1 {
					choicedUserID, _ := strconv.Atoi(customDataItems[0])
					mainChatID := customDataItems[1]
					mainChatIDInt64, _ := strconv.ParseInt(mainChatID, 10, 64)
					chatUser, _ := getChannelUser(contentType, mainChatIDInt64, choicedUserID)
					bot.Send(tgbotapi.NewMessage(mainChatIDInt64, "Killer choice: "+chatUser.User.Name))
					SetUserRoleToChannelList(contentType, mainChatIDInt64, choicedUserID, "dead")
					bot.Send(getUsersVoteMessageConfig(contentType, mainChatIDInt64, "Voting"))

					msg := tgbotapi.NewEditMessageText(
						chat.ID,
						messageID,
						"Your choice: "+chatUser.User.Name)
					bot.Send(msg)

					fmt.Printf("Private chat %+v\n", chat.ID)
					fmt.Printf("messageID edit %+v\n", messageID)
				}

			default:
				bot.Send(tgbotapi.NewMessage(chat.ID, "Data: "+update.CallbackQuery.Data))
			}

		} //update.CallbackQuery != nil
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
			_, isAdmin := checkPermission("admin", update.Message.From.ID)
			user := TGUser{
				UserID:  update.Message.From.ID,
				ChatId:  update.Message.Chat.ID,
				Login:   update.Message.From.UserName,
				Name:    update.Message.From.String(),
				IsAdmin: isAdmin,
			}
			if err := DB.Write("user", strconv.Itoa(update.Message.From.ID), user); err != nil {
				fmt.Println("add command error", err)
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Hi "+update.Message.From.String()+", you are registered!")
			bot.Send(msg)
		case "/myInfo":
			from := update.Message.From
			chat := update.Message.Chat

			chatMember, _ := bot.GetChatMember(tgbotapi.ChatConfigWithUser{
				ChatID:             chat.ID,
				SuperGroupUsername: "",
				UserID:             from.ID,
			})

			userInfo := "--== UserInfo==-- \n" +
				"ID: " + strconv.Itoa(from.ID) + "\n" +
				"UserName: " + from.UserName + "\n" +
				"FirstName: " + from.FirstName + "\n" +
				"LastName: " + from.LastName + "\n" +
				"LanguageCode: " + from.LanguageCode + "\n" +
				"--==ChatInfo==-- \n" +
				"ID: " + strconv.FormatInt(chat.ID, 10) + "\n" +
				"Title: " + chat.Title + "\n" +
				"Type: " + chat.Type + "\n" +
				"--==MemberInfo==-- \n" +
				"Status: " + chatMember.Status + "\n" +
				"ID: " + strconv.Itoa(chatMember.User.ID) + "\n" +
				"UserName: " + chatMember.User.UserName + "\n" +
				"FirstName: " + chatMember.User.FirstName + "\n" +
				"LastName: " + chatMember.User.LastName + "\n"
			msg := tgbotapi.NewMessage(chat.ID, userInfo)
			bot.Send(msg)
			//chat.ID,"",update.Message.From.ID

		case "/getUserList":
			err, permission := checkPermission("rebuild", update.Message.From.ID)
			if err != nil {
				log.Printf("Failed permissions: %v", err)
			}
			if permission {
				records, err := DB.ReadAll("user")
				if err != nil {
					fmt.Println("Error", err)
				}

				userList := []string{}
				for _, f := range records {
					userFound := TGUser{}
					if err := json.Unmarshal([]byte(f), &userFound); err != nil {
						fmt.Println("Error", err)
					}
					userList = append(userList, "["+strconv.Itoa(userFound.UserID)+"] "+userFound.Name)
				}
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, strings.Join(userList, "\n"))
				bot.Send(msg)
			} else {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Permission denied")
				bot.Send(msg)
			}

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
			} else {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Permission denied")
				bot.Send(msg)
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

			if err := DB.Write("command", commandValue, command); err != nil {
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
			//TODO: why it doesnt work
			err, messages := readLines("./features.txt", telegramSingleMessageLengthLimit)
			if err != nil {
				fmt.Println("write command error", err)
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, messages)
			bot.Send(msg)

		case "/SaveCommandsList":
			records, err := DB.ReadAll("command")
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
			records, err := DB.ReadAll("saved")
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

		case "/homeweb":
			homeWebImageURL, _ := getHomeWebCamImage()
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, homeWebImageURL)
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

			if err := DB.Write("checkList", itemCode, checkListItem); err != nil {
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

			records, err := DB.ReadAll("checkList")
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
						if err := DB.Write("checkList", checkListGroup, commandFound); err != nil {
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

			records, err := DB.ReadAll("сheckList")
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
			records, err := DB.ReadAll("command")
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
					if err := DB.Write(
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
