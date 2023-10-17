package tgModel

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

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

type Event struct {
	Name string
	Msg  *tgbotapi.Message
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

func GetTGButtons(params KeyBoardTG) tgbotapi.InlineKeyboardMarkup {
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
