package workTasks

import (
	"fmt"
	tgModel "fun-coice/internal/domain/commands/tg"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/hako/durafmt"
	"github.com/rs/zerolog/log"
	"time"
)

func (d data) AddTask(uid int64, msgId int) Task {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	startTime := time.Now()
	title := fmt.Sprintf(
		taskTitleTmp,
		startTime.Format("2006-01-02"),
		startTime.Format("-0700"),
		startTime.Format("15:04:05"),
		"0", "")
	userTask := Task{
		Start:  startTime,
		Title:  title,
		UserId: uid,
		Active: true,
		MsgId:  msgId,
	}
	log.Info().Any("task add", userTask).Send()
	d.tasks[uid] = userTask
	return userTask
}

func (d data) GetTask(uid int64) (Task, bool) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	userTask, exist := d.tasks[uid]
	return userTask, exist
}

func (d data) SetTaskBreak(uid int64) (Task, bool) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	userTask, exist := d.tasks[uid]
	breakTime := time.Now()
	if exist {
		userTask.Break = breakTime
		userTask.Active = false
		d.tasks[uid] = userTask
	}
	//userTask.Title = fmt.Sprintf("%s \n %s start: %s", userTask.Title, BreakName, breakTime.Format("15:04:05"))
	userTask.Title = userTask.GetTitle()
	log.Info().Any("task SetTaskBreak", userTask).Send()
	return userTask, exist
}

func (d data) StopTaskBreak(uid int64) (Task, bool) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	userTask, exist := d.tasks[uid]
	breakStopTime := time.Now()
	if exist {
		//userTask.Break = breakTime
		breakItem := userTask.add(BreakName, userTask.Break, breakStopTime)
		userTask.Active = true
		userTask.Title = fmt.Sprintf("%s \n %s: [%s]",
			userTask.Title,
			breakItem.Name,
			durafmt.Parse(breakItem.Duration).LimitFirstN(2).String())
		userTask.Title = userTask.GetTitle()
		d.tasks[uid] = userTask
	}
	log.Info().Any("task stopTaskBreak", userTask).Send()
	return userTask, exist
}

func (d data) activeTaskButtons() *tgbotapi.InlineKeyboardMarkup {
	return tgModel.GetTGButtons(tgModel.KBRows(tgModel.KBButs(
		tgModel.KeyBoardButtonTG{Text: takeBreakButtonTitle, Data: takeBreakButtonAction},
		tgModel.KeyBoardButtonTG{Text: StoppedTaskButtonTitle, Data: StoppedTaskButtonAction},
	)))
}

func (d data) breakTaskButtons() *tgbotapi.InlineKeyboardMarkup {
	return tgModel.GetTGButtons(tgModel.KBRows(tgModel.KBButs(
		tgModel.KeyBoardButtonTG{Text: stopBreakButtonTitle, Data: stopBreakButtonAction},
		tgModel.KeyBoardButtonTG{Text: StoppedTaskButtonTitle, Data: StoppedTaskButtonAction},
	)))
}

func (d data) trackButtons() *tgbotapi.InlineKeyboardMarkup {
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

func (t *Task) GetTitle() string {
	breaks := ""
	fullDuration := time.Now().Sub(t.Start)
	for _, item := range t.Items {
		breaks += fmt.Sprintf("\n %s: [%s]",
			item.Name,
			durafmt.Parse(item.Duration).LimitFirstN(2).String())
		fullDuration -= item.Duration
	}
	if !t.Active {
		breaks += fmt.Sprintf("\n %s : %s", BreakName, t.Break.Format("15:04:05"))
	}
	return fmt.Sprintf(
		taskTitleTmp,
		t.Start.Format("2006-01-02"),
		t.Start.Format("-0700"),
		t.Start.Format("15:04:05"),
		durafmt.Parse(fullDuration).LimitFirstN(2).String(),
		breaks)
}
