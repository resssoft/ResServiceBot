package workTasks

import (
	tgModel "fun-coice/internal/domain/commands/tg"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
)

const (
	taskTitleTmp = "[%s] Трэкинг GMT %s.  \n\nНачало: %s\nКонец: %s \nОбщее время: %s\n\n%s"

	timeFormat = "15:04" // "15:04:05"

	BreakName = "Перерыв" //"Break"

	timeTrackTitle = "Выберете действие"

	startTrackEvent   = "startTrack"
	settingsEvent     = "settings"
	takeBreakEvent    = "pause_track"
	stopBreakEvent    = "stop_break"
	StoppedTaskEvent  = "stop_task"
	setTaskNameEvent  = "setTaskName"
	startTaskEvent    = "startTask"
	showProfileEvent  = "showProfile"
	setBreakNameEvent = "setBreakName"
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

type Track struct {
	Start  time.Time
	End    time.Time
	Break  time.Time
	Pause  bool
	Close  bool
	Title  string
	UserId int64
	MsgId  int
	Breaks []timeItem
	Tasks  []timeItem
	//GMT    string use for time show
}

type Tracks map[int64]Track

type Button struct {
	Text   string
	Action string
	Event  string
	Data   tgModel.KeyBoardButtonTG
}
