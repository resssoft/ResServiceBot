package workTasks

import (
	"database/sql"
	tgModel "fun-coice/internal/domain/commands/tg"
	"github.com/doug-martin/goqu/v9"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
	"sync"
)

type data struct {
	list    tgModel.Commands
	users   map[int64]User //temporary
	storage *sql.DB
	builder goqu.DialectWrapper
	tracks  Tracks
	mutex   *sync.Mutex
	buttons map[string]Button
}

func New(DB *sql.DB) tgModel.Service {
	result := data{
		storage: DB,
		users:   make(map[int64]User), // temporary
		builder: goqu.Dialect("sqlite3"),
		mutex:   &sync.Mutex{},
		tracks:  make(Tracks),
		buttons: make(map[string]Button),
	}

	commandsList := tgModel.NewCommands()
	commandsList.AddSimple("timeTrack", "Show time track controls", result.timeTrack)
	result.list = commandsList

	result.addButton("🚗 Начать трэкинг", startTrackEvent, result.startTrackButtonEventHandler)
	result.addButton("⚙️", settingsEvent, result.settingsButtonEventHandler)
	result.addButton("⏸ перерыв", takeBreakEvent, result.takeBreakButtonEventHandler)
	result.addButton("⏹ возобновить", stopBreakEvent, result.stopBreakButtonEventHandler)
	result.addButton("▶️ остановить трэкинг", StoppedTaskEvent, result.StoppedTaskButtonEventHandler)

	result.addButton("💬 Задать имя активной задачи", setTaskNameEvent, result.NotImplementHandler)
	result.addButton("➕ Добавить задачу", startTaskEvent, result.NotImplementHandler)
	result.addButton("👤 Профиль", showProfileEvent, result.NotImplementHandler)
	result.addButton("💬 Задать имя перерыву", setBreakNameEvent, result.NotImplementHandler)

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
	go result.tracking()

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

func (d *data) tracking() {
	//TODO: implement
}

func (d *data) timeTrack(msg *tgbotapi.Message, _ *tgModel.Command) *tgModel.HandlerResult {
	return tgModel.SimpleWithButtons(msg.Chat.ID, timeTrackTitle, d.trackButtons())
}

func (d *data) startTrackButtonEventHandler(msg *tgbotapi.Message, c *tgModel.Command) *tgModel.HandlerResult {
	log.Info().Msg("startTaskButtonEventHandler")
	task := d.AddTrack(msg.Chat.ID, msg.MessageID)
	return tgModel.SimpleEditWithButtons(msg.Chat.ID, msg.MessageID, task.Title, d.activeTrackButtons())
}

func (d *data) settingsButtonEventHandler(msg *tgbotapi.Message, c *tgModel.Command) *tgModel.HandlerResult {
	log.Info().Msg("settingsButtonEventHandler")
	return tgModel.SimpleReply(msg.Chat.ID, "Not implement", msg.MessageID)
}

func (d *data) takeBreakButtonEventHandler(msg *tgbotapi.Message, c *tgModel.Command) *tgModel.HandlerResult {
	log.Info().Msg("takeBreakButtonEventHandler")
	task, exist := d.SetTrackBreak(msg.Chat.ID)
	if !exist {
		return tgModel.SimpleReply(msg.Chat.ID, "Track not found, sorry, create new by /timeTrack", msg.MessageID)
	}
	return tgModel.SimpleEditWithButtons(msg.Chat.ID, msg.MessageID, task.Title, d.breakTrackButtons())
}

func (d *data) stopBreakButtonEventHandler(msg *tgbotapi.Message, c *tgModel.Command) *tgModel.HandlerResult {
	log.Info().Msg("stopBreakButtonEventHandler")
	task, exist := d.StopTrackBreak(msg.Chat.ID)
	if !exist {
		return tgModel.SimpleReply(msg.Chat.ID, "Track not found, sorry, create new by /timeTrack", msg.MessageID)
	}
	return tgModel.SimpleEditWithButtons(msg.Chat.ID, msg.MessageID, task.Title, d.activeTrackButtons())
}

func (d *data) StoppedTaskButtonEventHandler(msg *tgbotapi.Message, c *tgModel.Command) *tgModel.HandlerResult {
	log.Info().Msg("StoppedTaskButtonEventHandler")
	task, exist := d.StopTrack(msg.Chat.ID)
	if !exist {
		return tgModel.SimpleReply(msg.Chat.ID, "Track not found, sorry, create new by /timeTrack", msg.MessageID)
	}
	return tgModel.SimpleEdit(msg.Chat.ID, msg.MessageID, task.Title)
}

func (d *data) setTaskNameButtonEventHandler(msg *tgbotapi.Message, c *tgModel.Command) *tgModel.HandlerResult {
	log.Info().Msg("setTaskNameButtonEventHandler")
	return tgModel.SimpleReply(msg.Chat.ID, "Not implement", msg.MessageID)
}

func (d *data) startTaskButtonEventHandler(msg *tgbotapi.Message, c *tgModel.Command) *tgModel.HandlerResult {
	log.Info().Msg("setTaskNameButtonEventHandler")
	return tgModel.SimpleReply(msg.Chat.ID, "Not implement", msg.MessageID)
}

func (d *data) NotImplementHandler(msg *tgbotapi.Message, c *tgModel.Command) *tgModel.HandlerResult {
	log.Info().Msg("setTaskNameButtonEventHandler")
	return tgModel.SimpleReply(msg.Chat.ID, "Not implement", msg.MessageID)
}
