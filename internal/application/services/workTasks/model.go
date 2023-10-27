package workTasks

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
)

type User struct {
	tgUser tgbotapi.User
	IsNew  bool
	IDStr  string
}

type timeItem struct {
	Name     string
	Start    time.Time
	End      time.Time
	Duration time.Duration
}

type Task struct {
	Start  time.Time
	Break  time.Time
	Active bool
	Title  string
	UserId int64
	MsgId  int
	Items  []timeItem
	//GMT    string use for time show
}

type Tasks map[int64]Task

const (
	taskTitleTmp = "[%s] Трэкинг GMT %s.  \n\nСтарт: %s \nОбщее время:%s\n\n%s"

	BreakName = "Перерыв" //"Break"

	timeTrackTitle         = "Выберете действие"
	addTaskButtonTitle     = "Начать трэкинг"
	settingsButtonTitle    = "Настройки"
	takeBreakButtonTitle   = "Перерыв"
	stopBreakButtonTitle   = "Возобновить"
	StoppedTaskButtonTitle = "Завершить"

	startTaskButtonAction = "event:workTrack_startTask"
	startTaskButtonEvent  = "workTrack_startTask"

	settingsButtonAction = "event:workTrack_settings"
	settingsButtonEvent  = "workTrack_settings"

	takeBreakButtonAction = "event:workTrack_break"
	takeBreakButtonEvent  = "workTrack_break"

	stopBreakButtonAction = "event:workTrack_stop_break"
	stopBreakButtonEvent  = "workTrack_stop_break"

	StoppedTaskButtonAction = "event:workTrack_stop"
	StoppedTaskButtonEvent  = "workTrack_stop"
)
