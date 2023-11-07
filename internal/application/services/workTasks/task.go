package workTasks

import (
	"fmt"
	tgModel "fun-coice/internal/domain/commands/tg"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/hako/durafmt"
	"github.com/rs/zerolog/log"
	"sort"
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
	activeTask, exist := t.Tasks[t.ActiveTask]
	if exist {
		activeTask.End = time.Now()
		activeTask.Duration = activeTask.Duration + activeTask.End.Sub(activeTask.Start)
		//activeTask.Duration = activeTask.Duration + activeTask.End.Sub(t.Start)
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
	if !exist {
		return Track{}, false
	}
	{ //debug
		tasksInfo := ""
		for i, t := range userTrack.Tasks {
			tasksInfo += fmt.Sprintf("\n[%v]%s-%s/%s %s",
				i, t.Start.Format(timeFormatS), t.End.Format(timeFormatS), Duration(t.Duration), t.Name)
		}
		log.Info().
			Str("SetTrackBreak before", tasksInfo).
			Send()
		d.tracks[uid] = userTrack.SetBreak()
	}
	d.tracks[uid] = userTrack
	userTrack.Title = userTrack.GetTitle()
	{ //debug
		tasksInfo := ""
		for i, t := range userTrack.Tasks {
			tasksInfo += fmt.Sprintf("\n[%v]%s-%s/%s %s",
				i, t.Start.Format(timeFormatS), t.End.Format(timeFormatS), Duration(t.Duration), t.Name)
		}
		log.Info().
			Str("SetTrackBreak After", tasksInfo).
			Send()
		log.Info().Any("SetTrackBreak", d.tracks[uid]).Send()

	}
	return userTrack, exist
}

func (d *data) StopTrackBreak(uid int64) (Track, bool) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	userTrack, exist := d.tracks[uid]
	if exist {
		tasksInfo := ""
		for i, t := range userTrack.Tasks {
			tasksInfo += fmt.Sprintf("\n[%v]%s-%s/%s %s",
				i, t.Start.Format(timeFormatS), t.End.Format(timeFormatS), Duration(t.Duration), t.Name)
		}
		log.Info().
			Str("StopTrackBreak before", tasksInfo).
			Send()
		d.tracks[uid] = userTrack.StopBreak()
	}
	d.tracks[uid] = userTrack
	userTrack.Title = userTrack.GetTitle()
	tasksInfo := ""
	for i, t := range userTrack.Tasks {
		tasksInfo += fmt.Sprintf("\n[%v]%s-%s/%s %s",
			i, t.Start.Format(timeFormatS), t.End.Format(timeFormatS), Duration(t.Duration), t.Name)
	}
	log.Info().
		Str("StopTrackBreak After", tasksInfo).
		Send()
	log.Info().Any("SetTrackBreak", d.tracks[uid]).Send()
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

func (t *Track) getTasks(inactive bool) (map[int]timeItem, []int) {
	tasks := make(map[int]timeItem)
	for index, task := range t.Tasks {
		if inactive && index == t.ActiveTask {
			continue
		}
		tasks[index] = task
	}
	keys := make([]int, 0, len(tasks))
	for k := range tasks {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	return tasks, keys
}

func (d *data) activeTrackButtons(uid int64) *tgbotapi.InlineKeyboardMarkup {
	userTrack, exist := d.tracks[uid]
	if !exist {
		return tgModel.GetTGButtons(tgModel.KeyBoardTG{})
	}
	tasks, keys := userTrack.getTasks(true)
	var taskRows []tgModel.KeyBoardRowTG
	taskRows = append(taskRows,
		d.ButtonRow(startTaskEvent, takeBreakEvent, StoppedTaskEvent, settingsEvent),
		d.ButtonRow(setTaskNameEvent))
	for _, taskIndex := range keys {
		taskRows = append(
			taskRows,
			tgModel.KBButs(
				tgModel.KeyBoardButtonTG{
					Text: fmt.Sprintf(taskIcon + " " + tasks[taskIndex].Name),
					Data: fmt.Sprintf("%s:%v", SetTaskAction, taskIndex),
				}))
	}
	//taskRows = append(taskRows, d.ButtonRow(startTaskEvent))

	return tgModel.GetTGButtons(tgModel.KBRows(taskRows...))
}

func (d *data) breakTrackButtons(_ int64) *tgbotapi.InlineKeyboardMarkup {
	return tgModel.GetTGButtons(tgModel.KBRows(
		d.ButtonRow(stopBreakEvent, StoppedTaskEvent, settingsEvent),
		d.ButtonRow(setBreakNameEvent)))
}

func (d *data) trackButtons(_ int64) *tgbotapi.InlineKeyboardMarkup {
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
	//t.Pause = true
	t.Status = StatusPause
	activeTask, exist := t.Tasks[t.ActiveTask]
	if exist {
		activeTask.End = time.Now()
		activeTask.Duration = activeTask.Duration + activeTask.End.Sub(activeTask.Start)
		t.Tasks[t.ActiveTask] = activeTask
	}
	log.Info().Any("task SetBreak", t).Send()
	return *t
}

func (t *Track) StopBreak() Track {
	if t == nil {
		return Track{}
	}
	breakStopTime := time.Now()
	t.add(DefaultBreakName, t.Break, breakStopTime)
	//t.Pause = false
	t.Title = t.GetTitle()
	t.Status = StatusProgress
	activeTask, exist := t.Tasks[t.ActiveTask]
	if exist {
		activeTask.Start = time.Now()
		t.Tasks[t.ActiveTask] = activeTask
	}
	log.Info().Any("task StopBreak", t).Send()
	return *t
}

func (t *Track) StopTrack() Track {
	log.Info().Any("task STARTFUN StopTask", t).Send()
	if t == nil {
		return Track{}
	}
	if t.IsPaused() {
		withoutBreak := t.StopBreak()
		t = &withoutBreak
		log.Info().Any("task StopTask withoutBreak", t).Send()
	}

	stopTime := time.Now()
	t.End = stopTime
	t.Title = t.GetTitle()
	t.Status = StatusStopped
	log.Info().Any("task StopTask", t).Send()
	return *t
}

func (t *Track) GetTitle() string {
	tasksInfo := ""
	breaks := ""
	fullDuration := time.Now().Sub(t.Start)
	for _, item := range t.Breaks {
		breaks += fmt.Sprintf("\n %s %s: %s-%s [%s]",
			breakIcon,
			item.Name,
			item.Start.Format(timeFormat),
			item.End.Format(timeFormat),
			Duration(item.Duration))
		fullDuration -= item.Duration
	}

	//sorted tasks
	icon := activeTaskIcon
	duration := ""
	tasks, keys := t.getTasks(false)
	for _, tIndex := range keys {
		task := tasks[tIndex]
		icon = taskIcon
		if t.ActiveTask == tIndex {
			icon = activeTaskIcon
			if t.IsPaused() {
				icon = taskPauseIcon
				duration = Duration(task.Duration)
			} else {
				duration = Duration(task.Duration + time.Now().Sub(task.Start))
			}
			if t.IsStopped() {
				icon = taskIcon
				duration = Duration(task.Duration)
			}
		} else {
			duration = Duration(task.Duration)
		}
		debug := fmt.Sprintf(" (%s %s %s) ",
			task.Start.Format(timeFormatS),
			task.Start.Format(timeFormatS),
			Duration(task.Duration))

		tasksInfo += fmt.Sprintf("\n%s %s : %s", icon, task.Name+debug, duration)
	}
	if t.IsPaused() {
		breaks += fmt.Sprintf(
			"\n%s %s : %s - [%s]",
			activeBreakIcon,
			DefaultBreakName,
			t.Break.Format(timeFormat),
			Duration(time.Now().Sub(t.Break)))
		fullDuration -= time.Now().Sub(t.Break)
	}
	return fmt.Sprintf(
		"[%s GMT %s] %s - %s \n⏱ %s\n%s:%s\n\n%s",
		t.Start.Format("2006-01-02"),
		t.Start.Format("-0700"),
		t.Start.Format(timeFormat),
		t.End.Format(timeFormat),
		Duration(fullDuration),
		TasksText,
		tasksInfo,
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
	activeTask, exist := userTrack.Tasks[userTrack.ActiveTask]
	if exist {
		activeTask.End = time.Now()
		activeTask.Duration = activeTask.Duration + activeTask.End.Sub(activeTask.Start)
		userTrack.Tasks[userTrack.ActiveTask] = activeTask
	}
	nextTask, exist := userTrack.Tasks[id]
	if !exist {
		return false
	} else {
		nextTask.Start = time.Now()
		userTrack.Tasks[id] = nextTask
	}
	userTrack.ActiveTask = id
	userTrack.Title = userTrack.GetTitle()
	d.tracks[uid] = userTrack
	return true
}

func (d *data) keyboard(t Track) *tgbotapi.InlineKeyboardMarkup {
	var keyboard *tgbotapi.InlineKeyboardMarkup
	switch t.Status {
	case StatusProgress:
		keyboard = d.activeTrackButtons(t.UserId)
	case StatusPause:
		keyboard = d.breakTrackButtons(t.UserId)
	default:
		keyboard = d.activeTrackButtons(t.UserId)
	}
	return keyboard
}

func (t *Track) IsPaused() bool {
	return t.Status == StatusPause
}

func (t *Track) IsStopped() bool {
	return t.Status == StatusStopped
}
