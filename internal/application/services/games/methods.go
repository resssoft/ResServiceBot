package games

import (
	"bufio"
	"errors"
	"fmt"
	tgCommands "fun-coice/internal/domain/commands/tg"
	"fun-coice/pkg/scribble"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

//TODO: CLEAR OLD FUNCTIONS OR MOVE TO OTHERS SERVICES

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

func incCountsChannelUsersList(contentType string, chatId int64, userId int64) {
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

func getChannelUser(contentType string, chatId int64, userId int64) (ChatUser, error) {
	for _, item := range ChatUserList {
		if item.ChatId == chatId && item.ContentType == contentType && item.User.UserID == userId {
			return item, nil
		}
	}
	return ChatUser{}, errors.New("user not found")
}

func getUsersButtons(chatUsers []ChatUser, chatID int64, code string) tgbotapi.InlineKeyboardMarkup {
	var rows []tgCommands.KeyBoardRowTG
	for _, chatUser := range chatUsers {
		rows = append(rows, KBButs(tgCommands.KeyBoardButtonTG{
			Text: chatUser.User.Name + " (" + strconv.Itoa(chatUser.VoteCount) + ")",
			Data: strconv.FormatInt(chatUser.User.UserID, 10) + "|" + strconv.FormatInt(chatID, 10) + "#" + code,
		}))
	}
	return getTGButtons(tgCommands.KeyBoardTG{rows})
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

func SaveUserToChannelList(contentType string, chatId int64, chatName string, userId int64, userName string) bool {
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

func SetUserRoleToChannelList(contentType string, chatId int64, userId int64, userRole string) {
	for itemIndex, item := range ChatUserList {
		if item.ChatId == chatId && item.ContentType == contentType && item.User.UserID == userId {
			ChatUserList[itemIndex].CustomRole = userRole
		}
	}
}

func checkUserRegister(userId int64) bool {
	// check - bot can write to user
	isRegistered := false
	var existUser = TGUser{}
	err := DB.Read("user", strconv.FormatInt(userId, 10), &existUser) // USE from service
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

func checkPermission(command string, userId int64) (error, bool) {

	/*
		typeOfCommand := tgbot.commands[command].Permissions.UserPermissions
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
	*/
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

func KBRows(KBrows ...tgCommands.KeyBoardRowTG) tgCommands.KeyBoardTG {
	var rows []tgCommands.KeyBoardRowTG
	rows = append(rows, KBrows...)
	return tgCommands.KeyBoardTG{rows}
}

func KBButs(KBrows ...tgCommands.KeyBoardButtonTG) tgCommands.KeyBoardRowTG {
	var rows []tgCommands.KeyBoardButtonTG
	rows = append(rows, KBrows...)
	return tgCommands.KeyBoardRowTG{rows}
}

func getSimpleTGButton(text, data string) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(text, data),
		),
	)
}

func getTGButtons(params tgCommands.KeyBoardTG) tgbotapi.InlineKeyboardMarkup {
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

var DB *scribble.Driver
