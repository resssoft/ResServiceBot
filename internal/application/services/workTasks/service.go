package workTasks

import (
	"context"
	"database/sql"
	tgModel "fun-coice/internal/domain/commands/tg"
	"github.com/doug-martin/goqu/v9"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
	"github.com/sasha-s/go-deadlock"
	"strconv"
	"strings"
	"time"
)

type data struct {
	list    tgModel.Commands
	users   map[int64]User //temporary
	storage *sql.DB
	builder goqu.DialectWrapper
	tracks  Tracks
	//mutex         *sync.Mutex
	mutex         deadlock.Mutex
	buttons       map[string]Button
	messageSender tgModel.MessageSender
}

const trackingDuration = time.Second * 131

func New(DB *sql.DB) tgModel.Service {
	result := data{
		storage: DB,
		users:   make(map[int64]User), // temporary
		builder: goqu.Dialect("sqlite3"),
		//mutex:   &sync.Mutex{},
		tracks:  make(Tracks),
		buttons: make(map[string]Button),
	}

	commandsList := tgModel.NewCommands()
	commandsList.AddSimple("timeTrack", "Show time track controls", result.timeTrack)
	commandsList.AddSimple("timeTrack_add_task", "Add task to active track, need task name parameter", result.addTaskButtonEventHandler)
	commandsList.AddSimple("timeTrack_set_task_name", "Add task to active track, need task name parameter", result.setTaskNameButtonEventHandler)

	commandsList.AddEvent(SetTaskEvent, result.SetActiveTask)

	result.list = commandsList

	result.addButton("🚗 Начать трэкинг", startTrackEvent, result.startTrackButtonEventHandler)
	result.addButton("⚙️", settingsEvent, result.settingsButtonEventHandler)
	result.addButton("⏸", takeBreakEvent, result.takeBreakButtonEventHandler)
	result.addButton("▶️", stopBreakEvent, result.stopBreakButtonEventHandler)
	result.addButton("🏁", StoppedTaskEvent, result.StoppedTaskButtonEventHandler)

	result.addButton("📝 Задать имя активной задачи", setTaskNameEvent, result.setTaskNameButtonEventHandler)
	result.addButton("➕ Добавить задачу", startTaskEvent, result.addTaskButtonEventHandler)
	result.addButton("👤 Профиль", showProfileEvent, result.NotImplementHandler)
	result.addButton("📝 Задать имя перерыву", setBreakNameEvent, result.NotImplementHandler)

	//TODO rename task name, edit time, duration, start, end
	//TODO buttons and task order

	//TODO: add workers for active trackers for update time info (check type buttons before edit message)
	//TODO set tasks text labels and duration (like as breaks) - some tasks by active tracker
	//TODO some trackers per day by user - feature: set random or user traker name
	//TODO button for add
	//TODO: add user break type buttons (coffe break for example) - tracker options
	//TODO: switching between active tasks
	//TODO: add set user GMT - settings user
	//TODO save info to db
	//TODO change/correct current time of tracker task or break

	//TODO: read from db to RAM active tasks(rename task to traker)
	go result.tracking(context.Background())

	return &result
}

func (d *data) addButton(text, event string, handler tgModel.HandlerFunc) Button {
	publicEvent := d.Name() + "_" + event
	btn := Button{
		Text:   text,
		Action: "event:" + publicEvent,
		Event:  publicEvent,
	}
	btn.Data = tgModel.KeyBoardButtonTG{Text: btn.Text, Data: btn.Action}
	d.buttons[event] = btn
	itemEvent := tgModel.NewEvent(publicEvent, handler)
	log.Info().Any("btn", btn).Send()
	log.Info().Any("itemEvent", itemEvent).Send()
	d.list.Add(publicEvent, *itemEvent)
	//log.Info().Any("d.list", d.list).Send()
	return btn
}

func (d *data) Button(event string) Button {
	//TODO: move to tgModel
	btn, ok := d.buttons[event]
	if ok {
		//btn.Data.Text = "" //TODO: translates
		return btn
	}
	log.Warn().Msg("Empy button used")
	return Button{}
}

func (d *data) ButtonRow(events ...string) tgModel.KeyBoardRowTG {
	var rows []tgModel.KeyBoardButtonTG
	for _, event := range events {
		rows = append(rows, d.Button(event).Data)
	}
	return tgModel.KeyBoardRowTG{Buttons: rows}
}

func (d *data) Commands() tgModel.Commands {
	return d.list
}

func (d *data) Name() string {
	return "timeTraker" //workTrack
}

func (d *data) Configure(botData tgModel.ServiceConfig) {
	d.messageSender = botData.MessageSender
}

func (d *data) tracking(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.NewTimer(trackingDuration).C:
			d.mutex.Lock()
			for _, track := range d.tracks {
				if track.GetTitle() == track.Title {
					continue
				}
				d.updateTrackMessage(track)
			}
			d.mutex.Unlock()
		}
	}
}

func (d *data) updateTrackMessage(track Track) {
	if track.Close && !(track.Status.Is(StatusProgress) || track.Status.Is(StatusPause)) {
	}
	var keyboard *tgbotapi.InlineKeyboardMarkup
	switch track.Status {
	case StatusProgress:
		keyboard = d.activeTrackButtons(track.UserId)
	case StatusPause:
		keyboard = d.breakTrackButtons(track.UserId)
	default:
	}
	newTitle := track.GetTitle()
	track.Title = newTitle
	if d.messageSender != nil {
		d.messageSender.PushHandleResult() <- tgModel.SimpleEditWithButtons(track.UserId, track.MsgId, track.Title, keyboard)
	}
}

func (d *data) timeTrack(msg *tgbotapi.Message, _ *tgModel.Command) *tgModel.HandlerResult {
	return tgModel.SimpleWithButtons(msg.Chat.ID, timeTrackTitle, d.trackButtons(msg.Chat.ID))
}

func (d *data) startTrackButtonEventHandler(msg *tgbotapi.Message, c *tgModel.Command) *tgModel.HandlerResult {
	log.Info().Msg("startTaskButtonEventHandler")
	task := d.AddTrack(msg.Chat.ID, msg.MessageID)
	return tgModel.SimpleEditWithButtons(msg.Chat.ID, msg.MessageID, task.Title, d.activeTrackButtons(msg.Chat.ID))
}

func (d *data) settingsButtonEventHandler(msg *tgbotapi.Message, c *tgModel.Command) *tgModel.HandlerResult {
	log.Info().Msg("settingsButtonEventHandler")
	return tgModel.SimpleReply(msg.Chat.ID, "Not implement", msg.MessageID)
}

func (d *data) takeBreakButtonEventHandler(msg *tgbotapi.Message, c *tgModel.Command) *tgModel.HandlerResult {
	log.Info().Msg("takeBreakButtonEventHandler")
	task, exist := d.SetTrackBreak(msg.Chat.ID)
	if !exist {
		return tgModel.SimpleReply(msg.Chat.ID, TrackNotFoundErrMsg, msg.MessageID)
	}
	return tgModel.SimpleEditWithButtons(msg.Chat.ID, msg.MessageID, task.Title, d.breakTrackButtons(msg.Chat.ID))
}

func (d *data) stopBreakButtonEventHandler(msg *tgbotapi.Message, c *tgModel.Command) *tgModel.HandlerResult {
	log.Info().Msg("stopBreakButtonEventHandler")
	task, exist := d.StopTrackBreak(msg.Chat.ID)
	if !exist {
		return tgModel.SimpleReply(msg.Chat.ID, TrackNotFoundErrMsg, msg.MessageID)
	}
	return tgModel.SimpleEditWithButtons(msg.Chat.ID, msg.MessageID, task.Title, d.activeTrackButtons(msg.Chat.ID))
}

func (d *data) StoppedTaskButtonEventHandler(msg *tgbotapi.Message, c *tgModel.Command) *tgModel.HandlerResult {
	log.Info().Msg("StoppedTaskButtonEventHandler")
	task, exist := d.StopTrack(msg.Chat.ID)
	if !exist {
		return tgModel.SimpleReply(msg.Chat.ID, TrackNotFoundErrMsg, msg.MessageID)
	}
	return tgModel.SimpleEdit(msg.Chat.ID, msg.MessageID, task.Title)
}

func (d *data) setTaskNameButtonEventHandler(msg *tgbotapi.Message, c *tgModel.Command) *tgModel.HandlerResult {
	if c.Arguments.Raw == "" {
		return tgModel.DeferredWithText(msg.Chat.ID, "Enter new task name", "timeTrack_set_task_name", "", nil)
	}
	foundedTrack, exist := d.updateActiveTaskName(msg.Chat.ID, c.Arguments.Raw)
	if !exist {
		return tgModel.Simple(msg.Chat.ID, TrackNotFoundErrMsg)
	}
	d.updateTrackMessage(foundedTrack)
	return tgModel.Simple(msg.Chat.ID, "Ok")
	//return tgModel.EmptyCommand()
}

func (d *data) addTaskButtonEventHandler(msg *tgbotapi.Message, c *tgModel.Command) *tgModel.HandlerResult {
	if c.Arguments.Raw == "" {
		return tgModel.DeferredWithText(msg.Chat.ID, "Enter task name", "timeTrack_add_task", "", nil)
	}
	foundedTrack, exist := d.AddTask(msg.Chat.ID, c.Arguments.Raw)
	if !exist {
		return tgModel.Simple(msg.Chat.ID, TrackNotFoundErrMsg)
	}
	d.updateTrackMessage(foundedTrack)
	return tgModel.Simple(msg.Chat.ID, "Ok")
	//TODO: EmptyCommand
	//return tgModel.EmptyCommand()
}

func (d *data) SetActiveTask(msg *tgbotapi.Message, c *tgModel.Command) *tgModel.HandlerResult {
	anonId := 0
	separated := strings.Split(c.Data, ":")
	if len(separated) == 3 {
		anonId, _ = strconv.Atoi(separated[2])
	}
	d.setActiveTask(msg.Chat.ID, anonId)
	userTrack, exist := d.GetTrack(msg.Chat.ID)
	if !exist {
		return tgModel.Simple(msg.Chat.ID, TrackNotFoundErrMsg)
	}
	d.updateTrackMessage(userTrack)
	return tgModel.EmptyCommand()
}

func (d *data) NotImplementHandler(msg *tgbotapi.Message, c *tgModel.Command) *tgModel.HandlerResult {
	log.Info().Msg("setTaskNameButtonEventHandler")
	return tgModel.SimpleReply(msg.Chat.ID, "Not implement", msg.MessageID)
}
