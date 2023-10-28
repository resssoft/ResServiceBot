package workTasks

import (
	"fmt"
	tgModel "fun-coice/internal/domain/commands/tg"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/hako/durafmt"
	"github.com/rs/zerolog/log"
	"strings"
	"time"
)

func (d *data) AddTrack(uid int64, msgId int) Track {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	userTrack := Track{
		Start:  time.Now(),
		UserId: uid,
		MsgId:  msgId,
	}
	userTrack.Title = userTrack.GetTitle()
	log.Info().Any("AddTrack", userTrack).Send()
	d.tracks[uid] = userTrack
	return userTrack
}

func (d *data) GetTrack(uid int64) (Track, bool) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	userTrack, exist := d.tracks[uid]
	return userTrack, exist
}

func (d *data) SetTrackBreak(uid int64) (Track, bool) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	userTrack, exist := d.tracks[uid]
	if exist {
		d.tracks[uid] = userTrack.SetBreak()
	}
	userTrack.Title = userTrack.GetTitle()
	log.Info().Any("SetTrackBreak", d.tracks[uid]).Send()
	return userTrack, exist
}

func (d *data) StopTrackBreak(uid int64) (Track, bool) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	userTrack, exist := d.tracks[uid]
	if exist {
		d.tracks[uid] = userTrack.StopBreak()
	}
	log.Info().Any("StopTrackBreak", d.tracks[uid]).Send()
	return userTrack, exist
}

func (d *data) StopTrack(uid int64) (Track, bool) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	userTrack, exist := d.tracks[uid]
	if exist {
		d.tracks[uid] = userTrack.StopTask()
	}
	log.Info().Any("StopTrack", d.tracks[uid]).Send()
	return userTrack, exist
}

func (d *data) activeTrackButtons() *tgbotapi.InlineKeyboardMarkup {
	return tgModel.GetTGButtons(tgModel.KBRows(
		d.ButtonRow(takeBreakEvent, StoppedTaskEvent, settingsEvent),
		d.ButtonRow(setTaskNameEvent),
		d.ButtonRow(startTaskEvent)))
}

func (d *data) breakTrackButtons() *tgbotapi.InlineKeyboardMarkup {
	return tgModel.GetTGButtons(tgModel.KBRows(
		d.ButtonRow(stopBreakEvent, StoppedTaskEvent, settingsEvent),
		d.ButtonRow(setBreakNameEvent)))
}

func (d *data) trackButtons() *tgbotapi.InlineKeyboardMarkup {
	return tgModel.GetTGButtons(tgModel.KBRows(d.ButtonRow(startTrackEvent, showProfileEvent)))
}

func (t *Track) add(name string, start, end time.Time) timeItem {
	breakItem := timeItem{
		Name:     name,
		Start:    start,
		End:      end,
		Duration: end.Sub(start),
	}
	t.Breaks = append(t.Breaks, breakItem)
	return breakItem
}

func (t *Track) SetBreak() Track {
	if t == nil {
		return Track{}
	}
	breakTime := time.Now()
	t.Break = breakTime
	t.Pause = true
	log.Info().Any("task SetBreak", t).Send()
	return *t
}

func (t *Track) StopBreak() Track {
	if t == nil {
		return Track{}
	}
	breakStopTime := time.Now()
	t.add(BreakName, t.Break, breakStopTime)
	t.Pause = false
	t.Title = t.GetTitle()
	log.Info().Any("task StopBreak", t).Send()
	return *t
}

func (t *Track) StopTask() Track {
	log.Info().Any("task STARTFUN StopTask", t).Send()
	if t == nil {
		return Track{}
	}
	if t.Pause {
		withoutBreak := t.StopBreak()
		t = &withoutBreak

		log.Info().Any("task StopTask withoutBreak", t).Send()
	}

	stopTime := time.Now()
	t.End = stopTime
	t.Title = t.GetTitle()
	t.Close = true
	log.Info().Any("task StopTask", t).Send()
	return *t
}

func (t *Track) GetTitle() string {
	breaks := ""
	fullDuration := time.Now().Sub(t.Start)
	for _, item := range t.Breaks {
		breaks += fmt.Sprintf("\n %s: [%s]",
			item.Name,
			Duration(item.Duration))
		fullDuration -= item.Duration
	}
	if t.Pause {
		breaks += fmt.Sprintf("\n %s : %s - ", BreakName, t.Break.Format(timeFormat))
	}
	return fmt.Sprintf(
		taskTitleTmp,
		t.Start.Format("2006-01-02"),
		t.Start.Format("-0700"),
		t.Start.Format(timeFormat),
		t.End.Format(timeFormat),
		Duration(fullDuration),
		breaks)
}

func Duration(dt time.Duration) string {
	//TODO: translates
	formatted := ""
	if dt.Milliseconds() < time.Minute.Milliseconds() {
		formatted = fmt.Sprintf("%.0f секунд", dt.Seconds())
	} else {
		formatted = durafmt.Parse(dt).LimitFirstN(2).String()
	}
	formatted = strings.ReplaceAll(formatted, "milliseconds", "милимсекунд")
	formatted = strings.ReplaceAll(formatted, "seconds", "секунд")
	formatted = strings.ReplaceAll(formatted, "minutes", "минут")
	formatted = strings.ReplaceAll(formatted, "hours", "часов")
	return formatted
}
