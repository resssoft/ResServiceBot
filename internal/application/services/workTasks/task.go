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
		Status: StatusProgress,
		Tasks:  make(map[int]timeItem),
	}
	userTrack.AddTask(DefaultTaskName)
	log.Info().Any("AddTrack", userTrack).Send()
	d.tracks[uid] = userTrack
	return userTrack
}

func (d *data) AddTask(uid int64, name string) (Track, bool) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	userTrack, exist := d.tracks[uid]
	if !exist {
		return userTrack, false
	}
	userTrack.AddTask(name)
	d.tracks[uid] = userTrack
	return userTrack, true
}

func (t *Track) AddTask(name string) timeItem {
	id := len(t.Tasks)
	newTaskItem := timeItem{
		Name:  name,
		Start: time.Now(),
		Id:    id,
	}
	//stopped preview task
	if id > 0 {
		activeTask := t.Tasks[t.ActiveTask]
		activeTask.End = time.Now()
		activeTask.Duration = activeTask.End.Sub(t.Start)
		t.Tasks[t.ActiveTask] = activeTask
	}
	t.ActiveTask = id
	t.Tasks[id] = newTaskItem
	t.Title = t.GetTitle()
	return newTaskItem
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
		d.tracks[uid] = userTrack.StopTrack()
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
	t.Status = StatusPause
	log.Info().Any("task SetBreak", t).Send()
	return *t
}

func (t *Track) StopBreak() Track {
	if t == nil {
		return Track{}
	}
	breakStopTime := time.Now()
	t.add(DefaultBreakName, t.Break, breakStopTime)
	t.Pause = false
	t.Title = t.GetTitle()
	t.Status = StatusProgress
	log.Info().Any("task StopBreak", t).Send()
	return *t
}

func (t *Track) StopTrack() Track {
	log.Info().Any("task STARTFUN StopTask", t).Send()
	if t == nil {
		return Track{}
	}
	if t.Pause {
		withoutBreak := t.StopBreak()
		t = &withoutBreak
		t.Status = StatusStopped

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
	tasks := ""
	breaks := ""
	fullDuration := time.Now().Sub(t.Start)
	for _, item := range t.Breaks {
		breaks += fmt.Sprintf("\n %s: [%s]",
			item.Name,
			Duration(item.Duration))
		fullDuration -= item.Duration
	}
	for tIndex, task := range t.Tasks {
		if t.ActiveTask == tIndex {
			tasks += fmt.Sprintf(
				"\n%s %s : %s",
				activeTaskIcon,
				task.Name,
				Duration(task.Duration+time.Now().Sub(task.Start)))
		} else {
			tasks += fmt.Sprintf(
				"\n%s %s : %s",
				taskIcon,
				task.Name,
				Duration(task.Duration))
		}
	}
	if t.Pause {
		breaks += fmt.Sprintf(
			"\n%s %s : %s - ",
			breakIcon,
			DefaultBreakName,
			t.Break.Format(timeFormat))
		fullDuration -= time.Now().Sub(t.Break)
	}
	return fmt.Sprintf(
		taskTitleTmp,
		t.Start.Format("2006-01-02"),
		t.Start.Format("-0700"),
		t.Start.Format(timeFormat),
		t.End.Format(timeFormat),
		Duration(fullDuration),
		tasks,
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

func (t *Track) updateTask(task timeItem) bool {
	_, exist := t.Tasks[task.Id]
	if !exist {
		return false
	}
	t.Tasks[task.Id] = task
	return true
}

func (d *data) updateActiveTaskName(uid int64, newName string) (Track, bool) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	userTrack, exist := d.tracks[uid]
	if !exist {
		return userTrack, false
	}
	activeTask, exist := userTrack.Tasks[userTrack.ActiveTask]
	if !exist {
		return userTrack, false
	}
	activeTask.Name = newName
	userTrack.updateTask(activeTask)
	return userTrack, true
}

func (d *data) setActiveTask(uid int64, id int) bool {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	userTrack, exist := d.tracks[uid]
	if !exist {
		return false
	}
	userTrack.ActiveTask = id
	d.tracks[uid] = userTrack
	return true
}
