package tgbot

import (
	"fmt"
	"fun-coice/config"
	tgCommands "fun-coice/internal/domain/commands/tg"
	"fun-coice/pkg/appStat"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	zlog "github.com/rs/zerolog/log"
	"log"
	"sync"
)

var defaultWorkersCount = 10

type data struct {
	WebUri       string
	Token        string
	StartMsg     bool
	WebMode      bool
	Commands     tgCommands.Commands
	Bot          *tgbotapi.BotAPI
	WorkersCount int
	Deferred     map[uint64]tgCommands.Command
	mutex        *sync.Mutex
}

func New(token, webUri string) data {
	bot, err := tgbotapi.NewBotAPI(config.TelegramToken())
	if err != nil {
		log.Panic(err)
	}
	return data{
		Token:        token,
		WebUri:       webUri,
		Commands:     defaultCommands,
		Bot:          bot,
		WorkersCount: defaultWorkersCount,
		StartMsg:     true,
		Deferred:     make(map[uint64]tgCommands.Command),
		mutex:        &sync.Mutex{},
	}
}

func (d data) Run() error {
	//d.Bot.Debug = true
	//TODO: d.Bot.GetMyCommands() AND SET THEM
	if d.StartMsg {
		msg := tgbotapi.NewMessage(config.TelegramAdminId(), "Bot Started with version "+appStat.Version)
		d.Bot.Send(msg)
	}

	log.Printf("Authorized on account %s", d.Bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 10

	if d.WebMode {
		fmt.Println("tg bot WebMode")
		webUpdates := d.Bot.ListenForWebhook(d.WebUri)
		for i := 0; i < d.WorkersCount; i++ {
			go d.CommandsHandler(webUpdates)
		}
	} else {
		fmt.Println("tg bot UpdateMode")
		updates := d.Bot.GetUpdatesChan(u)
		for i := 0; i < d.WorkersCount; i++ {
			go d.CommandsHandler(updates)
		}
	}
	return nil
}

func (d data) CommandsHandler(updates tgbotapi.UpdatesChannel) {
	var err error
	for update := range updates {
		/*
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
				//d.Bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data))
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
					d.Bot.Send(msg)

				case "lovelyGameJoin":
					isRegisteredUser := SaveUserToChannelList(
						contentType,
						chat.ID,
						chat.Title,
						from.ID,
						from.String(),
					)
					if !isRegisteredUser {
						d.Bot.Send(tgbotapi.NewMessage(chat.ID, from.String()+", write me to private for register"))
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
					d.Bot.Send(msg)

				case "lovelyGameStart":
					messageText := ""
					unregisteredUsers := unregisteredChannelUsers(contentType, chat.ID)
					if unregisteredUsers != "" {
						messageText = "I can`t start, unregistered users: " + unregisteredUsers
						d.Bot.Send(tgbotapi.NewMessage(chat.ID, messageText))
					} else {
						messageText = "Start lovely Game with: \n" +
							getChannelUsers(contentType, chat.ID) +
							"\n Wait for the killer to choose a player..."
						go sendRoleToUser(bot, chat.ID, contentType)
						msg := tgbotapi.NewEditMessageText(
							chat.ID,
							messageID,
							messageText)
						d.Bot.Send(msg)
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
					d.Bot.Send(tgbotapi.NewMessage(mainChatIDInt64, messageText))
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
		*/

		zlog.Info().Any("msg", update.Message).Any("InlineQuery", update.InlineQuery).Send()

		if update.Message == nil || (update.Message == nil && update.InlineQuery != nil) {
			zlog.Info().Any("update", update).Send()
			continue
		}

		//TODO: MOVE TO image processing service with wait messages types
		/*
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
				d.Bot.Send(message)

			}
		*/
		//fmt.Println(update.Message.Text)

		for _, command := range d.Commands {
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
			if !command.IsImplemented(commandName, d.Bot.Self.UserName) {
				if command.IsMatched(update.Message.Text, d.Bot.Self.UserName) {
					commandValue = update.Message.Text
				} else {
					continue
				}
			}
			result := command.Handler(update.Message, command.Command, commandValue, splitedCommands)
			if result.Prepared {
				_, err = d.Bot.Send(result.ChatEvent)
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
