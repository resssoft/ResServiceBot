package tgModel

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type HandlerResult struct {
	Prepared bool              // command is prepared for sending
	Deferred bool              // wait next message for handled by next command
	Resend   *tgbotapi.Message // message for resend
	Next     string            // next command
	Redirect *Redirect         //
	Messages []tgbotapi.Chattable
	Data     string
	Buttons  *tgbotapi.InlineKeyboardMarkup
	Events   []Event // run some events (or commands) after processing the current command
}

type HandlerFunc func(*tgbotapi.Message, *Command) *HandlerResult

func EmptyCommand() *HandlerResult {
	return &HandlerResult{
		Messages: nil,
	}
}

func PreparedCommand(chatEvents ...tgbotapi.Chattable) *HandlerResult {
	return &HandlerResult{
		Prepared: true,
		Messages: chatEvents,
	}
}

func Simple(chatId int64, text string) *HandlerResult {
	return PreparedCommand(tgbotapi.NewMessage(chatId, text))
}

func SimpleReply(chatId int64, text string, replyTo int) *HandlerResult {
	newMsg := tgbotapi.NewMessage(chatId, text)
	newMsg.ReplyToMessageID = replyTo
	return PreparedCommand(newMsg)
}

func SimpleWithButtons(chatId int64, text string, bts *tgbotapi.InlineKeyboardMarkup) *HandlerResult {
	mewMsg := tgbotapi.NewMessage(chatId, text)
	if bts != nil {
		mewMsg.ReplyMarkup = bts
	}
	return PreparedCommand(mewMsg)
}

func UnPreparedCommand(chatEvent tgbotapi.Chattable) *HandlerResult {
	return &HandlerResult{
		Messages: []tgbotapi.Chattable{chatEvent},
	}
}

func DeferredCommand(command, data string, msg *tgbotapi.Message) *HandlerResult {
	return &HandlerResult{
		Deferred: true,
		Next:     command,
		Data:     data,
		Resend:   msg,
	}
}

func DeferredWithText(chatId int64, text, command, data string, msg *tgbotapi.Message) *HandlerResult {
	return &HandlerResult{
		Deferred: true,
		Prepared: true,
		Messages: []tgbotapi.Chattable{tgbotapi.NewMessage(chatId, text)},
		Next:     command,
		Data:     data,
		Resend:   msg,
	}
}

func WaitingPreparedCommand(chatEvent tgbotapi.Chattable) *HandlerResult {
	return &HandlerResult{
		Deferred: true,
		Prepared: true,
		Messages: []tgbotapi.Chattable{chatEvent},
	}
}

func (hr *HandlerResult) SetEvent(newEvent Event) *HandlerResult {
	hr.Events = append(hr.Events, newEvent)
	return hr
}

func (hr *HandlerResult) WithEvent(name string, msg *tgbotapi.Message) *HandlerResult {
	hr.Events = append(hr.Events, Event{
		Name: "event:" + name,
		Msg:  msg,
	})
	return hr
}

func (hr *HandlerResult) WithRedirect(name string, msg *tgbotapi.Message) *HandlerResult {
	hr.Redirect = &Redirect{
		CommandName: name,
		Message:     msg,
	}
	return hr
}

func (hr *HandlerResult) WithDeferred(command string, msg *tgbotapi.Message) *HandlerResult {
	hr.Deferred = true
	hr.Next = command
	hr.Resend = msg
	return hr
}

func (hr *HandlerResult) WithText(chatId int64, text string) *HandlerResult {
	hr.Messages = []tgbotapi.Chattable{tgbotapi.NewMessage(chatId, text)}
	return hr
}

func (hr *HandlerResult) AddSimple(chatId int64, text string) *HandlerResult {
	hr.Messages = append(hr.Messages, tgbotapi.NewMessage(chatId, text))
	return hr
}

//OLD func(*tgbotapi.Message, string, string, []string) (tgbotapi.Chattable, bool) tgModel.HandlerResult tgModel.PreparedCommand( tgModel.PreparedCommand
// tgModel.PreparedCommand(tgbotapi.NewMessage ->  tgModel.Simple
//TODO: Handler     func(*tgbotapi.Message, string, string, []string) (tgbotapi.Chattable, HandlerResult)
