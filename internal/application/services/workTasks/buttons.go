package workTasks

import (
	"fun-coice/internal/application/services/workTasks/track"
	tgModel "fun-coice/internal/domain/commands/tg"
	"github.com/rs/zerolog/log"
)

func (d *data) addButton(text, event string, handler tgModel.HandlerFunc) track.Button {
	publicEvent := d.Name() + "_" + event
	btn := track.Button{
		Text:   text,
		Action: "event:" + publicEvent,
		Event:  publicEvent,
	}
	btn.Data = tgModel.KeyBoardButtonTG{Text: btn.Text, Data: btn.Action}
	d.buttons[event] = btn
	itemEvent := tgModel.NewEvent(publicEvent, handler)
	log.Info().Any("btn", btn).Send()
	log.Info().Any("itemEvent", itemEvent).Send()
	d.list.Add(publicEvent, *itemEvent)
	//log.Info().Any("d.list", d.list).Send()
	return btn
}

func (d *data) Button(event string) track.Button {
	//TODO: move to tgModel
	btn, ok := d.buttons[event]
	if ok {
		//btn.Data.Text = "" //TODO: translates
		return btn
	}
	log.Warn().Msg("Empy button used")
	return track.Button{}
}

func (d *data) ButtonRow(events ...string) tgModel.KeyBoardRowTG {
	var rows []tgModel.KeyBoardButtonTG
	for _, event := range events {
		rows = append(rows, d.Button(event).Data)
	}
	return tgModel.KeyBoardRowTG{Buttons: rows}
}
