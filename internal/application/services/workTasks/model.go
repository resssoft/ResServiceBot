package workTasks

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
)

const (
	taskTitleTmp = "[%s] Трэкинг GMT %s.  \n\nНачало: %s\nКонец: %s \nОбщее время: %s\n\n%s"

	timeFormat = "15:04" // "15:04:05"

	BreakName = "Перерыв" //"Break"

	timeTrackTitle         = "Выберете действие"
	addTaskButtonTitle     = "Начать трэкинг"
	settingsButtonTitle    = "Настройки"
	takeBreakButtonTitle   = "Перерыв"
	stopBreakButtonTitle   = "Возобновить"
	StoppedTaskButtonTitle = "Завершить"

	startTaskButtonAction = "event:timeTraker_startTask"
	startTaskButtonEvent  = "timeTraker_startTask"

	settingsButtonAction = "event:timeTraker_settings"
	settingsButtonEvent  = "timeTraker_settings"

	takeBreakButtonAction = "event:timeTraker_break"
	takeBreakButtonEvent  = "timeTraker_break"

	stopBreakButtonAction = "event:timeTraker_stop_break"
	stopBreakButtonEvent  = "timeTraker_stop_break"

	StoppedTaskButtonAction = "event:timeTraker_stop"
	StoppedTaskButtonEvent  = "timeTraker_stop"
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
	End    time.Time
	Break  time.Time
	Pause  bool
	Close  bool
	Title  string
	UserId int64
	MsgId  int
	Items  []timeItem
	//GMT    string use for time show
}

type Tasks map[int64]Task
