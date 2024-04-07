package track

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	tgModel "fun-coice/internal/domain/commands/tg"
	tgbotapi "fun-coice/pkg/telegram-bot-api"
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

	UserSettingsSetDefaultTasks = "settingsSetDefaultTasks"
)

//💳📝📝💬💬✏️💬
//📅➕➖➗✖️✔️🕐🏁
//🆕▶️⏸⏯⏹➡️⬅️⬆️⬇️🔙
//📝✏️🔎🗑🛠💾⏱⏰⏳🚩🏁➕➖➗✖️✔️🟠🟡🟢🔵🟣⚫️⚪️🔸🚧

type User struct {
	MongoID  primitive.ObjectID `bson:"_id"`
	TgUser   tgbotapi.User
	Settings UserSettings
}

type UserSettings struct {
	ID               int64
	DefaultTaskNames []string
	LangISO          string
}

func (u *User) defaultSettings() {
	u.Settings = UserSettings{
		ID:               u.TgUser.ID,
		DefaultTaskNames: []string{DefaultTaskName}, //set by lang
		LangISO:          u.TgUser.LanguageCode,
	}
}

type TimeItem struct {
	Id       int
	Name     string
	Start    time.Time
	End      time.Time
	Duration time.Duration
}

type Track struct {
	MongoID    primitive.ObjectID `bson:"_id"`
	Start      time.Time          `bson:"start"`
	End        time.Time          `bson:"end"`
	Break      time.Time          `bson:"break"`
	Title      string             `bson:"title"`
	UserId     int64              `bson:"user_id"`
	MsgId      int                `bson:"message_id"`
	Breaks     []TimeItem         `bson:"breaks"`
	Tasks      map[int]TimeItem   `bson:"tasks"`
	Status     Status             `bson:"status"`
	ActiveTask int                `bson:"active_task"`
	BotName    string             `bson:"not_name"`
	Code       string             `bson:"code"`
	//GMT string use for time show
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

type Users map[int64]User

type UserFilter struct {
	UserId  *int64
	MsgId   *int
	Status  Status
	BotName *string
	Code    *string
	//GMT    string use for time show
}
