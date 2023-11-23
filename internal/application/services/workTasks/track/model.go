package track

import (
	tgModel "fun-coice/internal/domain/commands/tg"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
)

const (
	TrackNotFoundErrMsg = "Track not found, sorry, create new by /timeTrack"
	activeTaskIcon      = "⏳"
	TaskIcon            = "🔸"
	taskPauseIcon       = "⏸"
	breakIcon           = "🔸"
	activeBreakIcon     = "⏳"

	timeFormat  = "15:04"    // "15:04:05"
	TimeFormatS = "15:04:05" // "15:04:05"
	TasksText   = "Задачи"   //"Break"

	DefaultBreakName = "Перерыв" //"Break"
	DefaultTaskName  = "Работа"  //"Break"

	TimeTrackTitle = "Выберите действие"

	StartTrackEvent   = "startTrack"
	SettingsEvent     = "settings"
	TakeBreakEvent    = "pause_track"
	StopBreakEvent    = "stop_break"
	StoppedTaskEvent  = "stop_task"
	SetTaskNameEvent  = "setTaskName"
	StartTaskEvent    = "startTask"
	ShowProfileEvent  = "showProfile"
	SetBreakNameEvent = "setBreakName"

	SetTaskEvent  = "timeTraker_set_task"
	SetTaskAction = "event:timeTraker_set_task"
)

//💳📝📝💬💬✏️💬
//📅➕➖➗✖️✔️🕐🏁
//🆕▶️⏸⏯⏹➡️⬅️⬆️⬇️🔙
//📝✏️🔎🗑🛠💾⏱⏰⏳🚩🏁➕➖➗✖️✔️🟠🟡🟢🔵🟣⚫️⚪️🔸🚧

type User struct {
	TgUser  tgbotapi.User
	IsNew   bool
	IDStr   string
	LangISO string
}

type TimeItem struct {
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
	Title      string
	UserId     int64
	MsgId      int
	Breaks     []TimeItem
	Tasks      map[int]TimeItem
	Status     Status
	ActiveTask int
	BotName    string
	Code       string
	//GMT    string use for time show
}

type TrackFilter struct {
	UserId  *int64
	MsgId   *int
	Status  Status
	BotName *string
	Code    *string
	//GMT    string use for time show
}

type TrackFields bool

func (tf *TrackFields) BotName() string {
	return "bot_name"
}
func (tf *TrackFields) MsgId() string {
	return "msg_id"
}
func (tf *TrackFields) TrackId() string {
	return "track_id"
}
func (tf *TrackFields) TrackJson() string {
	return "track_json"
}
func (tf *TrackFields) UserId() string {
	return "user_id"
}
func (tf *TrackFields) Status() string {
	return "status"
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
