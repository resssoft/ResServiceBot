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

func (d *data) AddTask(uid int64, msgId int) Task {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	userTask := Task{
		Start:  time.Now(),
		UserId: uid,
		MsgId:  msgId,
	}
	userTask.Title = userTask.GetTitle()
	log.Info().Any("task add", userTask).Send()
	d.tasks[uid] = userTask
	return userTask
}

func (d *data) GetTask(uid int64) (Task, bool) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	userTask, exist := d.tasks[uid]
	return userTask, exist
}

func (d *data) SetTaskBreak(uid int64) (Task, bool) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	userTask, exist := d.tasks[uid]
	if exist {
		d.tasks[uid] = userTask.SetBreak()
	}
	userTask.Title = userTask.GetTitle()
	log.Info().Any("task SetTaskBreak", d.tasks[uid]).Send()
	return userTask, exist
}

func (d *data) StopTaskBreak(uid int64) (Task, bool) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	userTask, exist := d.tasks[uid]
	if exist {
		d.tasks[uid] = userTask.StopBreak()
	}
	log.Info().Any("task stopTaskBreak", d.tasks[uid]).Send()
	return userTask, exist
}

func (d *data) StopTask(uid int64) (Task, bool) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	userTask, exist := d.tasks[uid]
	if exist {
		d.tasks[uid] = userTask.StopTask()
	}
	log.Info().Any("task stopTaskBreak", d.tasks[uid]).Send()
	return userTask, exist
}

func (d *data) activeTaskButtons() *tgbotapi.InlineKeyboardMarkup {
	return tgModel.GetTGButtons(tgModel.KBRows(tgModel.KBButs(
		tgModel.KeyBoardButtonTG{Text: takeBreakButtonTitle, Data: takeBreakButtonAction},
		tgModel.KeyBoardButtonTG{Text: StoppedTaskButtonTitle, Data: StoppedTaskButtonAction},
	)))
}

func (d *data) breakTaskButtons() *tgbotapi.InlineKeyboardMarkup {
	return tgModel.GetTGButtons(tgModel.KBRows(tgModel.KBButs(
		tgModel.KeyBoardButtonTG{Text: stopBreakButtonTitle, Data: stopBreakButtonAction},
		tgModel.KeyBoardButtonTG{Text: StoppedTaskButtonTitle, Data: StoppedTaskButtonAction},
	)))
}

func (d *data) trackButtons() *tgbotapi.InlineKeyboardMarkup {
	return tgModel.GetTGButtons(tgModel.KBRows(tgModel.KBButs(
		tgModel.KeyBoardButtonTG{Text: addTaskButtonTitle, Data: startTaskButtonAction},
		tgModel.KeyBoardButtonTG{Text: settingsButtonTitle, Data: settingsButtonAction},
	)))
}

func (t *Task) add(name string, start, end time.Time) timeItem {
	breakItem := timeItem{
		Name:     name,
		Start:    start,
		End:      end,
		Duration: end.Sub(start),
	}
	t.Items = append(t.Items, breakItem)
	return breakItem
}

func (t *Task) SetBreak() Task {
	if t == nil {
		return Task{}
	}
	breakTime := time.Now()
	t.Break = breakTime
	t.Pause = true
	log.Info().Any("task SetBreak", t).Send()
	return *t
}

func (t *Task) StopBreak() Task {
	if t == nil {
		return Task{}
	}
	breakStopTime := time.Now()
	t.add(BreakName, t.Break, breakStopTime)
	t.Pause = false
	t.Title = t.GetTitle()
	log.Info().Any("task StopBreak", t).Send()
	return *t
}

func (t *Task) StopTask() Task {
	log.Info().Any("task STARTFUN StopTask", t).Send()
	if t == nil {
		return Task{}
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

func (t *Task) GetTitle() string {
	breaks := ""
	fullDuration := time.Now().Sub(t.Start)
	for _, item := range t.Items {
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
