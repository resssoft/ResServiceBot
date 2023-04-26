package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"fun-coice/config"
	"fun-coice/funs"
	"fun-coice/internal/application/admins"
	"fun-coice/internal/application/b64"
	"fun-coice/internal/application/calculator"
	"fun-coice/internal/application/datatimes"
	"fun-coice/internal/application/lists"
	financy "fun-coice/internal/application/money"
	qrcodes "fun-coice/internal/application/qrcodes"
	"fun-coice/internal/application/translate"
	"fun-coice/pkg/appStat"
	"fun-coice/pkg/scribble"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
	"log"
	"net/http"
	"os"
	"strconv"
)

const doneMessage = "Done"
const telegramSingleMessageLengthLimit = 4096
const dbDateFormatMonth = "2006-01-02"

var HWCSURL = ""

var ChatUserList = make([]ChatUser, 1)

func main() {
	zlog.Level(zerolog.DebugLevel)
	fmt.Print("Load configuration... ")
	config.Configure()

	fmt.Println(fmt.Sprintf("apilayer[%s]", config.Str("plugins.apilayer.token")))
	fmt.Println(fmt.Sprintf("Telegram[%s]", config.TelegramToken()))
	fmt.Println(fmt.Sprintf("Admin[%v]", config.TelegramAdminId()))

	bot, err := tgbotapi.NewBotAPI(config.TelegramToken())
	if err != nil {
		log.Panic(err)
	}
	botName := bot.Self.UserName
	log.Printf("Admin is ..." + config.TelegramAdminLogin())
	log.Printf("Work with DB...")
	appPath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	//TODO: TRANSLATES

	//TODO: moved simple DB implement to pkg
	// read admin info from DB or write it to db
	DB, err = scribble.New(appPath + "/data")
	if err != nil {
		fmt.Println("Error", err)
	}
	if err := DB.Read("user", strconv.FormatInt(int64(config.TelegramAdminId()), 10), &existAdmin); err != nil {
		fmt.Println("admin not found error", err)
		existAdmin = TGUser{
			UserID:  config.TelegramAdminId(),
			ChatId:  0,
			Login:   "",
			Name:    "",
			IsAdmin: false,
		}
		if err := DB.Write("user", strconv.FormatInt(int64(config.TelegramAdminId()), 10), existAdmin); err != nil {
			fmt.Println("Error", err)
		}
	}

	funCommandsService := funs.New(DB)
	commands = commands.Merge(funCommandsService.Commands())

	b64Service := b64.New()
	commands = commands.Merge(b64Service.Commands())

	QrCodesService := qrcodes.New()
	commands = commands.Merge(QrCodesService.Commands())

	dataTimesService := datatimes.New()
	commands = commands.Merge(dataTimesService.Commands())

	trService := translate.New()
	commands = commands.Merge(trService.Commands())

	calculatorService := calculator.New()
	commands = commands.Merge(calculatorService.Commands())

	financeService := financy.New(config.Str("plugins.apilayer.token"))
	commands = commands.Merge(financeService.Commands())

	listService := lists.New(DB)
	commands = commands.Merge(listService.Commands())

	usersService := lists.New(DB)
	commands = commands.Merge(usersService.Commands())

	adminService := admins.New(bot, DB, commands)
	commands = commands.Merge(adminService.Commands())

	//fmt.Println("funCommandsService", funCommandsService.Commands())

	//fmt.Println("commands", commands)

	//bot.Debug = true
	msg := tgbotapi.NewMessage(config.TelegramAdminId(), "Bot Started with version "+appStat.Version)
	bot.Send(msg)
	bot.GetMyCommands()

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 10

	updates := bot.GetUpdatesChan(u)

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

			//TODO: FIX MIGRATE FROM v4 to v5
			//bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data))
			splitedCallbackQuery, clearCallbackQuery := splitCommand(update.CallbackQuery.Data, "#")
			commandsCount := len(splitedCallbackQuery)

			zlog.Info().Interface("update", update).Send()
			zlog.Info().Interface("chat", chat).Send()
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
				choicedUserID, _ := strconv.ParseInt(customDataItems[0], 10, 64)
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
					choicedUserID, _ := strconv.ParseInt(customDataItems[0], 10, 64)
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

		zlog.Info().Any("msg", update.Message).Any("InlineQuery", update.InlineQuery).Send()

		if update.Message == nil || (update.Message == nil && update.InlineQuery != nil) {
			zlog.Info().Any("update", update).Send()
			continue
		}

		if update.Message.Photo != nil {
			fileId := ""
			for _, photoItem := range update.Message.Photo {
				fileId = photoItem.FileID
			}
			//response, err := http.Get(fmt.Sprintf("https://api.telegram.org/bot%s/getFile?file_id=%s", config.TelegramToken(), fileId))
			response, err := http.Get(fmt.Sprintf("https://api.telegram.org/bot%s/getFile?file_id=%s", config.TelegramToken(), fileId))
			if err != nil {
				log.Println("download TG photo error")
				continue
			}
			buf := new(bytes.Buffer)
			buf.ReadFrom(response.Body)
			result := buf.String()
			//log.Println("tg fileInfo unparsed")
			fileInfo := TgFileInfo{}
			err = json.Unmarshal([]byte(result), &fileInfo)
			if err != nil {
				log.Println("Decode fileInfo err")
				continue
			}
			fileUrl := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s",
				config.TelegramToken(), fileInfo.Result.FilePath)

			response, err = http.Get(fileUrl)
			if err != nil {
				log.Println("download TG import file error")
				continue
			}
			buf = new(bytes.Buffer)
			buf.ReadFrom(response.Body)

			newImage, err := getMagic(buf.Bytes())
			tgNewfile := tgbotapi.FileBytes{
				Name:  "photo.jpg",
				Bytes: newImage,
			}
			var message tgbotapi.Chattable
			message = tgbotapi.NewPhoto(update.Message.Chat.ID, tgNewfile)
			bot.Send(message)

		}
		//fmt.Println(update.Message.Text)

		for _, command := range commands {
			if !command.Permission(update.Message) || command.Handler == nil {
				continue
			}
			splitedCommands, commandValue := splitCommand(update.Message.Text, " ")
			if len(splitedCommands) == 0 {
				continue
			}
			commandName := splitedCommands[0]
			commandsCount := len(splitedCommands)
			if commandsCount == 0 {
				continue
			}
			if !command.IsImplemented(commandName, botName) {
				if command.IsMatched(commandName, botName) {
					commandValue = update.Message.Text
				} else {
					continue
				}
			}
			botMsg, prepared := command.Handler(update.Message, command.Command, commandValue, splitedCommands)
			if prepared {
				_, err = bot.Send(botMsg)
				if err != nil {
					fmt.Println(err.Error())
				}
			}
		}

		//TODO:: add service bot informer for /member - admin service and SAVER service
		//TODO:  calc service, fiat service
		//TODO:: defaults to services
		//TODO:: photo and other file handlers to services (USE WAIT LIST)
		//TODO: /commands - show with perms
		//TODO: added wait answer commands (or files-images combiner)

		/*
			switch commandName {
			case "/games":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Games list")
				msg.ReplyMarkup = gamesListKeyboard
				bot.Send(msg)

			default:
			}
		*/

		//log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		if update.Message.Chat.Type == "private" && config.Str("logLevel") == "private" || config.Str("logLevel") == "chat" {
			log.Printf("INNER MESSAGE %s[%d]: %s",
				update.Message.From.UserName,
				update.Message.From.ID,
				update.Message.Text)
			fmt.Printf("inline query %+v\n", update.InlineQuery)
		}

		//msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		//msg.ReplyToMessageID = update.Message.MessageID

		//bot.Send(msg)
	}

}

type TgFileInfo struct {
	Ok     bool `json:"ok,omitempty"`
	Result struct {
		FileId       string `json:"file_id,omitempty"`
		FileUniqueId string `json:"file_unique_id,omitempty"`
		FileSize     int    `json:"file_size,omitempty"`
		FilePath     string `json:"file_path,omitempty"`
	} `json:"result,omitempty"`
}

type AmoCrmMMessageMedia struct {
	Type     string
	Media    string
	FileName string
	FileSize int
}

func splitInt(n int) []int {
	slc := []int{}
	for n > 0 {
		slc = append(slc, n%10)
		n = n / 10
	}
	return slc
}

func splitInt64(n int64) []int64 {
	slc := []int64{}
	for n > 0 {
		slc = append(slc, n%10)
		n = n / 10
	}
	return slc
}

func lastDigits(n int) (int, int) {
	result := splitInt(n)
	if len(result) < 2 {
		return 0, 0
	}
	return result[0], result[1]
}

func lastDigits64(n int64) (int64, int64) {
	result := splitInt64(n)
	if len(result) < 2 {
		return 0, 0
	}
	return result[len(result)-1], result[len(result)-2]
}

var catNameSet1 = map[int]string{
	0: "Немытый",
	1: "Жирный",
	2: "Горячий",
	3: "Лысый",
	4: "Всратый",
	5: "Забивной",
	6: "Пушистый",
	7: "Бешенный",
	8: "Депресивный",
	9: "Отбитый",
}

var catNameSet2 = map[int]string{
	0: "Гей",
	1: "Тигр",
	2: "Даун",
	3: "Кошак",
	4: "Чмо",
	5: "Красавчик",
	6: "Уебан",
	7: "Пидорас",
	8: "Кiт",
	9: "Чухан",
}

var putinSpeech = map[int]string{
	0:  "Глотаю пыль",
	1:  "Аграрий",
	2:  "Простой человек с кухни",
	3:  "Гендерно нейтральный бог",
	4:  "Ядерные объедки",
	5:  "Ветеран",
	6:  "Оступившийся человек",
	7:  "Второсортный чужак",
	8:  "Джинн из бутылки",
	9:  "Кхе кхе",
	10: "Крапленые карты",
	11: "Лысый черт",
	12: "Пособник Киевского режима",
}

//TODO: implement check command
//t.IsCommand(commandName, "/setLeadStatus")
//func (t *tgConfig) IsCommand(msg, command string) bool { return msg == command || msg == fmt.Sprintf("%s@%s", command, t.BotName)}
