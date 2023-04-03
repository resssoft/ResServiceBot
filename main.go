package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"fun-coice/config"
	"fun-coice/funs"
	tgCommands "fun-coice/internal/domain/commands/tg"
	"fun-coice/pkg/appStat"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/hako/durafmt"
	scribble "github.com/nanobox-io/golang-scribble"
	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
	qrcode "github.com/skip2/go-qrcode"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const doneMessage = "Done"
const telegramSingleMessageLengthLimit = 4096
const HWCSURLEvent = "go"
const HWCSURLImage = "result/"
const dbDateFormatMonth = "2006-01-02"

var HWCSURL = ""

var ChatUserList = make([]ChatUser, 1)

func main() {
	zlog.Level(zerolog.DebugLevel)
	fmt.Print("Load configuration... ")
	config.Configure()
	configureConverter(config.Str("plugins.apilayer.token"))

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

	// read admin info from DB or write it to db
	DB, err = scribble.New(appPath+"/data", nil)
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

	//initFunCommand()

	funCommandsService := funs.New(DB)
	commands = commands.Merge(funCommandsService.Commands())

	fmt.Println("funCommandsService", funCommandsService.Commands())

	fmt.Println("commands", commands)

	//bot.Debug = true
	msg := tgbotapi.NewMessage(config.TelegramAdminId(), "!Bot Started with version "+appStat.Version)
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
		splitedCommands, commandValue := splitCommand(update.Message.Text, " ")
		commandsCount := len(splitedCommands)
		if commandsCount == 0 {
			continue
		}
		commandName := splitedCommands[0]
		//fmt.Println("splitedCommands", splitedCommands)

		for _, command := range commands {
			if !command.IsImplemented(commandName, botName) {
				continue
			}
			if !command.Permission(update.Message) {
				continue
			}
			if command.Handler != nil {
				botMsg, prepared := command.Handler(update.Message, commandName, commandValue, splitedCommands)
				if prepared {
					_, err = bot.Send(botMsg)
					if err != nil {
						fmt.Println(err.Error())
					}
				}
			}
		}

		/*
			if commandData, exist := isFunCommand(commandName); exist {
				s1 := rand.NewSource(time.Now().UnixNano())
				r1 := rand.New(s1)
				time.Sleep(time.Millisecond * time.Duration(r1.Int63n(600)))
				r2 := rand.New(s1)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, commandData.List1[r1.Intn(len(commandData.List1))]+" "+commandData.List2[r2.Intn(len(commandData.List2))])
				msg.ReplyToMessageID = update.Message.MessageID
				bot.Send(msg)
			} else {
				//fmt.Println("NOT isFunCommand")
			}
		*/

		//TODO:: add service bot informer for /member and SAVER service
		//TODO:: defaults to services
		//TODO:: photo and other file handlers to services (USE WAIT LIST)

		//TODO: set permissions for default commands
		switch commandName {
		case "/start":
			_, isAdmin := checkPermission("admin", update.Message.From.ID)
			user := TGUser{
				UserID:  int64(update.Message.From.ID),
				ChatId:  update.Message.Chat.ID,
				Login:   update.Message.From.UserName,
				Name:    update.Message.From.String(),
				IsAdmin: isAdmin,
			}
			if err := DB.Write("user", strconv.FormatInt(update.Message.From.ID, 10), user); err != nil {
				fmt.Println("add command error", err)
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Hi "+update.Message.From.String()+", you are registered!")
			bot.Send(msg)
			/*
				case "/addfan":
					fmt.Println(splitedCommands)
					text := ""
					if len(splitedCommands) != 4 {
						text = "format: /addfan newcommandname list1_item1,list1_item2 list2_item1,list2_item2"
						text += "\nExample: cats cute,funny,fluffy Molly,Charlie,Oscar"
						text += "\n no more than 3 spase in the string"
					} else {
						text = addFunCommand(splitedCommands[1], splitedCommands[2], splitedCommands[3])
					}
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
					bot.Send(msg)
			*/
		case "/member":
			from := update.Message.From
			chat := update.Message.Chat

			chatConfigWithUser := tgbotapi.ChatConfigWithUser{
				ChatID:             chat.ID,
				SuperGroupUsername: "",
				UserID:             from.ID,
			}

			chatMember, _ := bot.GetChatMember(tgbotapi.GetChatMemberConfig{chatConfigWithUser})

			userInfo := fmt.Sprintf(
				"--== UserInfo==--\n"+
					"ID: %v\nUserName: %s\nFirstName: %s\nLastName: %s\nLanguageCode: %s"+
					"\n--==ChatInfo==--\n"+
					"ID: %v\nTitle: %s\nType: %s"+
					"\n--== MemberInfo==--\n"+
					"Status: %s",
				from.ID,
				from.UserName,
				from.FirstName,
				from.LastName,
				from.LanguageCode,
				chat.ID,
				chat.Title,
				chat.Type,
				chatMember.Status,
			)
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
					userList = append(userList, "["+strconv.FormatInt(config.TelegramAdminId(), 10)+"] "+userFound.Name)
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
			command := tgCommands.Command{
				Command:     commandValue,
				CommandType: "SaveCommand",
				Permissions: tgCommands.CommandPermissions{
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
			formattedMessage := currentTime + "[" + appStat.Version + "]: " + commandValue
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

		case "/catrandom", "/catrandom@FunChoiceBot":
			s1 := rand.NewSource(time.Now().UnixNano())
			r1 := rand.New(s1)
			time.Sleep(time.Millisecond * time.Duration(r1.Int63n(600)))
			r2 := rand.New(s1)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, catNameSet1[r1.Intn(10)]+" "+catNameSet2[r2.Intn(10)])
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)

		case "/putinSpeech", "/putinSpeech@FunChoiceBot":
			s1 := rand.NewSource(time.Now().UnixNano())
			r1 := rand.New(s1)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, putinSpeech[r1.Intn(12)])
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)

		case "/calc", "/calc@FunChoiceBot":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, calcFromStr(commandValue))
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)

		case "calc", "калк", "сколько будет":
			if update.Message.Chat.Type != "private" {
				continue
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, calcFromStr(commandValue))
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)

		case "/fiat", "/fiat@FunChoiceBot":
			convertFrom := "AMD"
			convertTo1 := "RUB"
			convertTo2 := "USD"
			if len(splitedCommands) > 2 {
				switch splitedCommands[2] {
				case "a", "amd", "am", "ам", "амд", "дпам", "драм", "др":
					convertFrom = "AMD"
					convertTo1 = "RUB"
					convertTo2 = "USD"
				case "r", "ru", "rub", "rur", "ру", "р", "руб", "рублей":
					convertFrom = "RUB"
					convertTo1 = "AMD"
					convertTo2 = "USD"
				case "s", "us", "usd", "$", "дол", "до", "доларов", "долларов":
					convertFrom = "USD"
					convertTo1 = "AMD"
					convertTo2 = "RUB"
				}
			}
			_, err = strconv.Atoi(splitedCommands[1])
			msgText := "-"
			if err != nil {
				msgText = "digit err"
			} else {
				msgText = fmt.Sprintf("Convert result from %s %s = \n%s %s \n%s %s \n[%s]",
					splitedCommands[1], convertFrom,
					fiat(convertFrom, convertTo1, splitedCommands[1]), convertTo1,
					fiat(convertFrom, convertTo2, splitedCommands[1]), convertTo2,
					time.Now().Format(dbDateFormatMonth))
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)

		case "fiat", "convert", "конверт", "кон", "из", "from":
			if update.Message.Chat.Type != "private" {
				continue
			}
			convertFrom := "AMD"
			convertTo1 := "RUB"
			convertTo2 := "USD"
			if len(splitedCommands) > 2 {
				switch splitedCommands[2] {
				case "a", "amd", "am", "ам", "амд", "дпам", "драм", "др":
					convertFrom = "AMD"
					convertTo1 = "RUB"
					convertTo2 = "USD"
				case "r", "ru", "rub", "rur", "ру", "р", "руб", "рублей":
					convertFrom = "RUB"
					convertTo1 = "AMD"
					convertTo2 = "USD"
				case "s", "us", "usd", "$", "дол", "до", "доларов", "долларов":
					convertFrom = "USD"
					convertTo1 = "AMD"
					convertTo2 = "RUB"
				}
			}
			_, err = strconv.Atoi(splitedCommands[1])
			msgText := "-"
			if err != nil {
				msgText = "digit err"
			} else {
				msgText = fmt.Sprintf("Convert result from %s %s = \n%s %s \n%s %s \n[%s]",
					splitedCommands[1], convertFrom,
					fiat(convertFrom, convertTo1, splitedCommands[1]), convertTo1,
					fiat(convertFrom, convertTo2, splitedCommands[1]), convertTo2,
					time.Now().Format(dbDateFormatMonth))
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)

		case "/SaveCommandsList":
			records, err := DB.ReadAll("command")
			if err != nil {
				fmt.Println("Error", err)
			}

			commands := []string{}
			for _, f := range records {
				commandFound := tgCommands.Command{}
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

		case "/tm", "/timestamp":
			msgText := fmt.Sprintf("Input: %s\n\n", commandValue)
			timestampConvert, _ := regexp.MatchString(`^\d\d\d\d\d\d+$`, commandValue)
			timestampCompare, _ := regexp.MatchString(`^\d\d\d\d\d\d+\s\d\d\d\d\d\d+$`, commandValue)
			switch {
			case timestampConvert:
				timeInt, err := strconv.ParseInt(commandValue, 10, 64)
				if err != nil {
					msgText += "Error parse int timestamp: " + err.Error()
				} else {
					parsedTime := time.Unix(timeInt, 0)
					msgText += fmt.Sprintf("Date: %s \n", parsedTime.Format("2006-01-02 15:04:05 -0700"))
					msgText += fmt.Sprintf("Date: %s \n", parsedTime.Format("2006-01-02T15:04:05Z07:00"))
					msgText += fmt.Sprintf("Date: %s \n", parsedTime.Format("Monday, 02-Jan-06 15:04:05 MST"))
					msgText += fmt.Sprintf("Timestamp: %s \n", fmt.Sprintf("%v", parsedTime.Unix()))
					msgText += fmt.Sprintf("Diff: %s \n", durafmt.Parse(time.Now().Sub(parsedTime)).LimitFirstN(2).String())
					msgText += fmt.Sprintf("\nNow: %s \n", time.Now().Format("2006-01-02 15:04:05 -0700"))
				}
			case timestampCompare:
				r, _ := regexp.Compile(`^(\d\d\d\d\d\d+)\s(\d\d\d\d\d\d+)$`)
				parsedItems := r.FindStringSubmatch(commandValue)
				fmt.Println(parsedItems)
				if len(parsedItems) != 3 {
					msgText += fmt.Sprintf("Error: parsedItems != 3 It= %v", len(parsedItems))
				} else {
					timeInt1, err1 := strconv.ParseInt(parsedItems[1], 10, 64)
					if err != nil {
						msgText += "\nError: " + err.Error()
					}
					timeInt2, err2 := strconv.ParseInt(parsedItems[2], 10, 64)
					if err != nil {
						msgText += "\nError: " + err.Error()
					}
					if err1 == nil && err2 == nil {
						parsedTime1 := time.Unix(timeInt1, 0)
						parsedTime2 := time.Unix(timeInt2, 0)
						msgText += fmt.Sprintf("Diff: %s \n", durafmt.Parse(parsedTime1.Sub(parsedTime2)).LimitFirstN(2).String())
					}
				}
			default:
				parsedTime, err := time.Parse("2006-01-02 15:04:05", commandValue)
				if err != nil {
					parsedTime, err = time.Parse("2006-01-02 15:04", commandValue)
				}
				if err != nil {
					parsedTime, err = time.Parse("2006-01-02", commandValue)
				}
				if err != nil {
					parsedTime, err = time.Parse("02-01-2006", commandValue)
				}
				if err != nil {
					parsedTime, err = time.Parse("02.01.2006", commandValue)
				}
				if err != nil {
					parsedTime, err = time.Parse("02/01/2006", commandValue)
				}
				if err != nil {
					parsedTime, err = time.Parse("2006-01-02T15:04:05Z07:00", commandValue)
				}
				if err != nil {
					parsedTime, err = time.Parse("Mon Jan _2 15:04:05 2006", commandValue)
				}
				if err != nil {
					parsedTime, err = time.Parse("Mon Jan _2 15:04:05 MST 2006", commandValue)
				}
				if err != nil {
					parsedTime, err = time.Parse("02 Jan 06 15:04 MST", commandValue)
				}
				if err != nil {
					parsedTime, err = time.Parse("Mon, 02 Jan 2006 15:04:05 MST", commandValue)
				}
				if err != nil {
					parsedTime, err = time.Parse("Mon 02 Jan 2006 15:04:05 MST", commandValue)
				}
				if err != nil {
					parsedTime, err = time.Parse("Mon 02 Jan 2006 15:04:05 07:00", commandValue)
				}
				if err != nil {
					parsedTime, err = time.Parse("January 2, 2006", commandValue)
				}
				var parsedDuration time.Duration
				if err != nil {
					parsedDuration, err = time.ParseDuration(commandValue)
					parsedTime = time.Now().Add(parsedDuration)
				}
				if err != nil {
					msgText += "Error parse time"
				} else {
					msgText += fmt.Sprintf("Date: %s \n", parsedTime.Format("2006-01-02 15:04:05 -0700"))
					msgText += fmt.Sprintf("Date: %s \n", parsedTime.Format("2006-01-02T15:04:05Z07:00"))
					msgText += fmt.Sprintf("Date: %s \n", parsedTime.Format("Monday, 02-Jan-06 15:04:05 MST"))
					msgText += fmt.Sprintf("Timestamp: %s \n", fmt.Sprintf("%v", parsedTime.Unix()))
					msgText += fmt.Sprintf("Diff by current: %s \n", durafmt.Parse(time.Now().Sub(parsedTime)).LimitFirstN(2).String())
					msgText += fmt.Sprintf("\nNow: %s \n", time.Now().Format("2006-01-02 15:04:05 -0700"))
				}
			}
			time.Parse(commandValue, commandValue)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
			bot.Send(msg)

		case "/set":
			config.Set(splitedCommands[1], splitedCommands[2])
		case "/get":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("%v=%v", splitedCommands[1], config.Str(splitedCommands[1])))
			bot.Send(msg)
			//logLevel
		case "/vars":
			//TODO:fix
			config.Set(splitedCommands[1], splitedCommands[2])

		case "/admin":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Admin is @"+config.TelegramAdminLogin())
			bot.Send(msg)
		case "/qr256":
			var qr []byte
			qr, err := qrcode.Encode(commandValue, qrcode.Medium, 256)
			if err != nil {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "failed to encode qrcode")
				bot.Send(msg)
				continue
			}
			file := tgbotapi.FileBytes{Name: "qr.png", Bytes: qr}
			message := tgbotapi.NewPhoto(update.Message.Chat.ID, file)
			message.Caption = commandValue
			_, err = bot.Send(message)
			if err != nil {
				fmt.Println(err.Error())
			}
		case "/b64", "/base64":
			b64result := base64.StdEncoding.EncodeToString([]byte(commandValue))
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, b64result)
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)

		case "/qr":
			var qr []byte
			qr, err := qrcode.Encode(commandValue, qrcode.Medium, 1024)
			if err != nil {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "failed to encode qrcode")
				bot.Send(msg)
				continue
			}
			file := tgbotapi.FileBytes{Name: "qr.png", Bytes: qr}
			message := tgbotapi.NewPhoto(update.Message.Chat.ID, file)
			message.Caption = commandValue
			_, err = bot.Send(message)
			if err != nil {
				fmt.Println(err.Error())
			}
		case "/version", "/appVersion", "/версия":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, appStat.Version)
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
				fmt.Println("Error DB.ReadAl", err)
			}

			commandContain := false
			var commands []tgCommands.Command
			for _, f := range records {
				commandFound := tgCommands.Command{}
				if err := json.Unmarshal([]byte(f), &commandFound); err != nil {
					fmt.Println("Error commandFound Unmarshal", err)
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
			matchedCalc, _ := regexp.MatchString(`^[\d\s\+\\\-\*\(\)\.]+$`, update.Message.Text)
			matchedCalc2, _ := regexp.MatchString(`^\d+$`, update.Message.Text)
			//fmt.Println(matchedCalc)
			if matchedCalc && !matchedCalc2 {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, calcFromStr(update.Message.Text))
				msg.ReplyToMessageID = update.Message.MessageID
				bot.Send(msg)
			}

			if !commandContain {
				///////log.Println("This is unsupport command.")
				//msg := tgbotapi.NewMessage(update.Message.Chat.ID, "This is unsupport command.")
				//msg.ReplyToMessageID = update.Message.MessageID
				//bot.Send(msg)
			}
		}

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
