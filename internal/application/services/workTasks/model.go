package workTasks

import (
	tgModel "fun-coice/internal/domain/commands/tg"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
)

const (
	taskTitleTmp = "[%s GMT %s] %s - %s \n⏱ %s\nЗадачи:%s\n\n%s"
	//taskTitleTmp = "[%s] Трэкинг GMT %s.  \n\nНачало: %s Конец: %s \nОбщее время: %s\n\n%s"
	TrackNotFoundErrMsg = "Track not found, sorry, create new by /timeTrack"
	activeTaskIcon      = "⏳"
	taskIcon            = "🔸"
	breakIcon           = "🔸"

	timeFormat = "15:04" // "15:04:05"

	DefaultBreakName = "Перерыв" //"Break"
	DefaultTaskName  = "Работа"  //"Break"

	timeTrackTitle = "Выберите действие"

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

//💳📝📝💬💬✏️💬
//📅➕➖➗✖️✔️🕐🏁
//🆕▶️⏸⏯⏹➡️⬅️⬆️⬇️🔙
//📝✏️🔎🗑🛠💾⏱⏰⏳🚩🏁➕➖➗✖️✔️🟠🟡🟢🔵🟣⚫️⚪️🔸🚧

type User struct {
	tgUser  tgbotapi.User
	IsNew   bool
	IDStr   string
	LangISO string
}

type timeItem struct {
	Id       int
	Name     string
	Start    time.Time
	End      time.Time
	Duration time.Duration
}

type Track struct {
	Start      time.Time
	End        time.Time
	Break      time.Time
	Pause      bool
	Close      bool
	Title      string
	UserId     int64
	MsgId      int
	Breaks     []timeItem
	Tasks      map[int]timeItem
	Status     Status
	ActiveTask int
	//GMT    string use for time show
}

type Status int

const (
	StatusStart = iota
	StatusProgress
	StatusPause
	StatusStopped
	StatusProfile
	StatusSettings
)

func (s Status) Is(status Status) bool {
	return s == status
}

type Tracks map[int64]Track

type Button struct {
	Text   string
	Action string
	Event  string
	Data   tgModel.KeyBoardButtonTG
}
