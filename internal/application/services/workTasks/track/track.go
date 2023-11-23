package track

import (
	"fmt"
	"github.com/hako/durafmt"
	"github.com/rs/zerolog/log"
	"sort"
	"strings"
	"time"
)

func (t *Track) AddTask(name string) TimeItem {
	id := len(t.Tasks)
	newTaskItem := TimeItem{
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

func (t *Track) GetTasks(inactive bool) (map[int]TimeItem, []int) {
	tasks := make(map[int]TimeItem)
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

func (t *Track) add(name string, start, end time.Time) TimeItem {
	breakItem := TimeItem{
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
	tasks, keys := t.GetTasks(false)
	for _, tIndex := range keys {
		task := tasks[tIndex]
		icon = TaskIcon
		if t.ActiveTask == tIndex {
			icon = activeTaskIcon
			if t.IsPaused() {
				icon = taskPauseIcon
				duration = Duration(task.Duration)
			} else {
				duration = Duration(task.Duration + time.Now().Sub(task.Start))
			}
			if t.IsStopped() {
				icon = TaskIcon
				duration = Duration(task.Duration)
			}
		} else {
			duration = Duration(task.Duration)
		}
		debug := fmt.Sprintf(" (%s %s %s) ",
			task.Start.Format(TimeFormatS),
			task.Start.Format(TimeFormatS),
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
	formatted = strings.ReplaceAll(formatted, "second", "секунда")
	formatted = strings.ReplaceAll(formatted, "minutes", "минут")
	formatted = strings.ReplaceAll(formatted, "minute", "минута")
	formatted = strings.ReplaceAll(formatted, "hours", "часов")
	formatted = strings.ReplaceAll(formatted, "hour", "час")
	return formatted
}

func (t *Track) UpdateTask(task TimeItem) bool {
	_, exist := t.Tasks[task.Id]
	if !exist {
		return false
	}
	t.Tasks[task.Id] = task
	return true
}

func (t *Track) IsPaused() bool {
	return t.Status == StatusPause
}

func (t *Track) IsStopped() bool {
	return t.Status == StatusStopped
}
