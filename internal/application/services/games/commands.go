package games

import (
	tgModel "fun-coice/internal/domain/commands/tg"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (d data) games(msg *tgbotapi.Message, command *tgModel.Command) *tgModel.HandlerResult {
	newMsg := tgbotapi.NewMessage(msg.Chat.ID, "Games list")
	newMsg.ReplyMarkup = gamesListKeyboard
	return tgModel.PreparedCommand(newMsg)
}

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
