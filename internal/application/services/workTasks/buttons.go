package workTasks

import (
	"fmt"

	"github.com/rs/zerolog/log"

	"fun-coice/internal/application/services/workTasks/track"
	tgModel "fun-coice/internal/domain/commands/tg"
	tgbotapi "fun-coice/pkg/telegram-bot-api"
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

func (d *data) activeTrackButtons(uid int64) *tgbotapi.InlineKeyboardMarkup {
	userTrack, exist := d.tracks[uid]
	if !exist {
		return tgModel.GetTGButtons(tgModel.KeyBoardTG{})
	}
	tasks, keys := userTrack.GetTasks(true)
	var taskRows []tgModel.KeyBoardRowTG
	taskRows = append(taskRows,
		d.ButtonRow(track.StartTaskEvent, track.TakeBreakEvent, track.StoppedTaskEvent, track.SettingsEvent),
		d.ButtonRow(track.SetTaskNameEvent))
	for _, taskIndex := range keys {
		taskRows = append(
			taskRows,
			tgModel.KBButs(
				tgModel.KeyBoardButtonTG{
					Text: fmt.Sprintf(track.TaskIcon + " " + tasks[taskIndex].Name),
					Data: fmt.Sprintf("%s:%v", track.SetTaskAction, taskIndex),
				}))
	}
	//taskRows = append(taskRows, d.ButtonRow(startTaskEvent))

	return tgModel.GetTGButtons(tgModel.KBRows(taskRows...))
}

func (d *data) breakTrackButtons(_ int64) *tgbotapi.InlineKeyboardMarkup {
	return tgModel.GetTGButtons(tgModel.KBRows(
		d.ButtonRow(track.StopBreakEvent, track.StoppedTaskEvent, track.SettingsEvent),
		d.ButtonRow(track.SetBreakNameEvent)))
}

func (d *data) trackButtons(_ int64) *tgbotapi.InlineKeyboardMarkup {
	return tgModel.GetTGButtons(tgModel.KBRows(d.ButtonRow(track.StartTrackEvent, track.ShowProfileEvent)))
}

func (d *data) keyboard(t track.Track) *tgbotapi.InlineKeyboardMarkup {
	var keyboard *tgbotapi.InlineKeyboardMarkup
	switch t.Status {
	case track.StatusProgress:
		keyboard = d.activeTrackButtons(t.UserId)
	case track.StatusPause:
		keyboard = d.breakTrackButtons(t.UserId)
	default:
		keyboard = d.activeTrackButtons(t.UserId)
	}
	return keyboard
}

func (d *data) userSettingsMainButtons(_ int64) *tgbotapi.InlineKeyboardMarkup {
	return tgModel.GetTGButtons(tgModel.KBRows(
		d.ButtonRow(track.UserSettingsSetDefaultTasks),
		d.ButtonRow(track.SetBreakNameEvent)))
}
