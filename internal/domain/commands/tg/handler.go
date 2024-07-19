package tgModel

import "fun-coice/pkg/telegram-bot-api"

//"fun-coice/pkg/telegram-bot-api"

type HandlerResult struct {
	Prepared bool              // command is prepared for sending
	Deferred bool              // wait next message for handled by next command
	Resend   *tgbotapi.Message // message for resend
	Next     string            // next command
	Redirect *Redirect         //
	Messages []MessageEvent
	Data     string
	Buttons  *tgbotapi.InlineKeyboardMarkup
	Events   []Event // run some events (or commands) after processing the current command
	Reaction *TgReaction
}

type MessageEvent struct {
	Event    tgbotapi.Chattable
	Tag      string
	Callback chan CallbackData `json:"-"`
}

type HandlerFunc func(*tgbotapi.Message, *Command) *HandlerResult

func EmptyCommand() *HandlerResult {
	return &HandlerResult{
		Messages: nil,
	}
}

type ReactionItem struct {
	Type  string `json:"type"`
	Emoji string `json:"emoji"`
}

type TgReaction struct {
	ChatID    int64
	MessageID int
	Emoji     string
}

func Delete(chatId int64, msgId int) *HandlerResult {
	return PreparedCommand(tgbotapi.NewDeleteMessage(chatId, msgId))
}

func SimpleMessageEvents(chatEvents ...tgbotapi.Chattable) []MessageEvent {
	var list []MessageEvent
	for _, event := range chatEvents {
		list = append(list, MessageEvent{Event: event})
	}
	return list
}

func PreparedCommand(chatEvents ...tgbotapi.Chattable) *HandlerResult {
	return &HandlerResult{
		Prepared: true,
		Messages: SimpleMessageEvents(chatEvents...),
	}
}

func Simple(chatId int64, text string) *HandlerResult {
	return PreparedCommand(tgbotapi.NewMessage(chatId, text))
}

func SimpleWIthCallback(chatId int64, text, tag string, callback chan CallbackData) *HandlerResult {
	return &HandlerResult{
		Prepared: true,
		Messages: []MessageEvent{
			{
				Event:    tgbotapi.NewMessage(chatId, text),
				Tag:      tag,
				Callback: callback,
			},
		},
	}
}

func SimpleEdit(chatId int64, msgId int, text string) *HandlerResult {
	return PreparedCommand(tgbotapi.NewEditMessageText(chatId, msgId, text))
}

func Reaction(chatId int64, msgId int, text string) *HandlerResult {
	return &HandlerResult{
		Prepared: true,
		Reaction: &TgReaction{
			ChatID:    chatId,
			MessageID: msgId,
			Emoji:     text,
		},
	}
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

func SimpleEditWithButtons(chatId int64, msgId int, text string, bts *tgbotapi.InlineKeyboardMarkup) *HandlerResult {
	if bts == nil {
		return EmptyCommand()
	}
	mewMsg := tgbotapi.NewEditMessageTextAndMarkup(chatId, msgId, text, *bts)
	return PreparedCommand(mewMsg)
}

func UnPreparedCommand(chatEvent tgbotapi.Chattable) *HandlerResult {
	return &HandlerResult{
		Messages: SimpleMessageEvents(chatEvent),
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
		Messages: SimpleMessageEvents(tgbotapi.NewMessage(chatId, text)),
		Next:     command,
		Data:     data,
		Resend:   msg,
	}
}

func WaitingPreparedCommand(chatEvent tgbotapi.Chattable) *HandlerResult {
	return &HandlerResult{
		Deferred: true,
		Prepared: true,
		Messages: SimpleMessageEvents(chatEvent),
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
	hr.Messages = SimpleMessageEvents(tgbotapi.NewMessage(chatId, text))
	return hr
}

func (hr *HandlerResult) AddSimple(chatId int64, text string) *HandlerResult {
	hr.Messages = append(hr.Messages, SimpleMessageEvents(tgbotapi.NewMessage(chatId, text))...)
	return hr
}

func (hr *HandlerResult) WithDelete(chatId int64, msgId int) *HandlerResult {
	hr.Messages = append(hr.Messages, SimpleMessageEvents(tgbotapi.NewDeleteMessage(chatId, msgId))...)
	return hr
}

//OLD func(*tgbotapi.Message, string, string, []string) (tgbotapi.Chattable, bool) tgModel.HandlerResult tgModel.PreparedCommand( tgModel.PreparedCommand
// tgModel.PreparedCommand(tgbotapi.NewMessage ->  tgModel.Simple
//TODO: Handler     func(*tgbotapi.Message, string, string, []string) (tgbotapi.Chattable, HandlerResult)
