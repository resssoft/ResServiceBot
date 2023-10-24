package p2p

import (
	"database/sql"
	"fmt"
	tgModel "fun-coice/internal/domain/commands/tg"
	"github.com/doug-martin/goqu/v9"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
	zlog "github.com/rs/zerolog/log"
	"strings"
)

const (
	linkTmp         = "https://t.me/%s?start=%v"
	eventDelMsg     = "event:p2p_remove_this_msg"
	eventSendMsg    = "event:p2p_send_message_to:"
	eventPrepareMsg = "event:p2p_prepare_message_to:"
)

type data struct {
	list       tgModel.Commands
	users      map[int64]User //temporary
	storage    *sql.DB
	builder    goqu.DialectWrapper
	promoCodes promoCodes
}

func New(DB *sql.DB) tgModel.Service {
	result := data{
		storage: DB,
		users:   make(map[int64]User), // temporary
		builder: goqu.Dialect("sqlite3"),
		promoCodes: addPromoCodes(
			newPromoCode("resager", "premium1"),
			newPromoCode("norl", "no_rl"),
		),
	}
	commandsList := tgModel.NewCommands()
	commandsList.AddSimple("start", "", result.start) // replace default start event
	commandsList.AddSimple("p2p_default", "", result.defaultHandler)
	commandsList.AddEvent("p2p_prepare_message_to", result.prepareMessage)
	commandsList.AddEvent("p2p_send_message_to", result.sendAnonMessage)

	//TODO: add buttons (quickly) and events

	result.list = commandsList
	return &result
}

func (d data) Commands() tgModel.Commands {
	return d.list
}

func (d data) Name() string {
	return "p2p"
}

func (d data) start(msg *tgbotapi.Message, command *tgModel.Command) *tgModel.HandlerResult {
	command.Arguments.Parse()
	log.Info().Any("args", command.Arguments).Send()
	user, err := d.userInfo(msg.From)
	if err != nil {
		zlog.Error().Str("service error", d.Name()).Err(err).Send()
		return tgModel.EmptyCommand()
	}

	if len(command.Arguments.List) > 0 {
		param := command.Arguments.List[0]
		switch {
		case d.promoCodes.exist(param):
		case Number(param) > 0:
			user, err := d.userInfoByID(Number(param))
			if err == nil {
				name := ""
				if len(command.Arguments.List) >= 2 {
					name = command.Arguments.List[1]
				}
				return d.sendButtonsForMessage(msg.Chat.ID, user, name)
			} else {
				return tgModel.SimpleReply(msg.Chat.ID, "User not registered in the Bot", msg.MessageID) //TODO: translate
			}
		}
	}
	if user.IsNew {
		return tgModel.Simple(
			msg.Chat.ID,
			fmt.Sprintf(
				"Ваша ссылка: %s, поделитесь ею, чтобы вам могли присылать сообщения приватно. В настройках можно указать имя, которое будет указано у отправителя (только для сообщения к вам)",
				d.getLink(msg.From.ID, command.ParamCallback(tgModel.BotNameParam).Str())))
	}
	return tgModel.EmptyCommand()
}

func (d data) defaultHandler(msg *tgbotapi.Message, _ *tgModel.Command) *tgModel.HandlerResult {
	return tgModel.SimpleReply(msg.Chat.ID, "xe-xe", msg.MessageID)
}

func (d data) prepareMessage(msg *tgbotapi.Message, c *tgModel.Command) *tgModel.HandlerResult {
	fmt.Println("data", c.Data, msg.Chat.ID)
	return tgModel.DeferredWithText(
		msg.Chat.ID,
		"Напишите анонимное сообщение для отправки:",
		"p2p_send_message_to",
		c.Data,
		nil)
}

func (d data) sendAnonMessage(msg *tgbotapi.Message, c *tgModel.Command) *tgModel.HandlerResult {
	anonId := int64(0)
	separated := strings.Split(c.Data, ":")
	if len(separated) == 3 {
		anonId = Number(separated[2])
	}
	fmt.Println("data2", c.Data, anonId, msg.Chat.ID)
	//TODO: create new anon contact name
	return tgModel.SimpleWithButtons(anonId, "У тебя новое анонимное сообщение:\n\n"+msg.Text, d.replyButtons(msg.Chat.ID)).
		AddSimple(msg.Chat.ID, "✅ Сообщение отправлено!") // translate
}

func (d data) delSelf(msg *tgbotapi.Message, с *tgModel.Command) *tgModel.HandlerResult {
	return tgModel.SimpleReply(msg.Chat.ID, "Not implemented", msg.MessageID)
}

func (d data) getLink(uid int64, botName string) string {
	//TODO: create and save hash
	return fmt.Sprintf(linkTmp, botName, uid)
}

func (d data) sendButtonsForMessage(chat int64, to User, name string) *tgModel.HandlerResult {
	if name == "" {
		name = "this user" // translate
	}
	newMsg := tgbotapi.NewMessage(chat, "Новый контакт для "+name) //New dialog
	newMsg.ReplyMarkup = tgModel.GetTGButtons(tgModel.KBRows(tgModel.KBButs(
		tgModel.KeyBoardButtonTG{Text: "Написать анонимно", Data: eventPrepareMsg + to.IDStr}, //Write anonymously to
		tgModel.KeyBoardButtonTG{Text: "Удалить сообщение", Data: eventDelMsg},                //Remove this message
	)))
	return tgModel.PreparedCommand(newMsg)
}

func (d data) isPromocode(uid int64, botName string) string {
	//TODO: create and save hash
	return fmt.Sprintf(linkTmp, botName, uid)
}

func (d data) replyButtons(to int64) *tgbotapi.InlineKeyboardMarkup {
	return tgModel.GetTGButtons(tgModel.KBRows(tgModel.KBButs(
		tgModel.KeyBoardButtonTG{Text: "Ответить анонимно", Data: eventPrepareMsg + fmt.Sprintf("%v", to)},
	)))
}

//TODO: anon names
//TODO: save users
// provide names from hash codes, no IDS - id is deanon
