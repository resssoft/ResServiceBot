package workTasks

import (
	"database/sql"
	tgModel "fun-coice/internal/domain/commands/tg"
	"github.com/doug-martin/goqu/v9"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
	"sync"
)

const (
	linkTmp         = "https://t.me/%s?start=%v"
	eventDelMsg     = "event:p2p_remove_this_msg"
	eventSendMsg    = "event:p2p_send_message_to:"
	eventPrepareMsg = "event:p2p_prepare_message_to:"
)

type data struct {
	list    tgModel.Commands
	users   map[int64]User //temporary
	storage *sql.DB
	builder goqu.DialectWrapper
	tasks   Tasks
	mutex   *sync.Mutex
}

func New(DB *sql.DB) tgModel.Service {
	result := data{
		storage: DB,
		users:   make(map[int64]User), // temporary
		builder: goqu.Dialect("sqlite3"),
		mutex:   &sync.Mutex{},
		tasks:   make(Tasks),
	}

	commandsList := tgModel.NewCommands()
	commandsList.AddSimple("timeTrack", "", result.timeTrack)
	//commandsList.AddSimple("start", "", result.start) // replace default start event - show start buttons

	commandsList.AddEvent(startTaskButtonEvent, result.startTaskButtonEventHandler)
	commandsList.AddEvent(settingsButtonEvent, result.settingsButtonEventHandler)
	commandsList.AddEvent(takeBreakButtonEvent, result.takeBreakButtonEventHandler)
	commandsList.AddEvent(StoppedTaskButtonEvent, result.StoppedTaskButtonEventHandler)
	commandsList.AddEvent(stopBreakButtonEvent, result.stopBreakButtonEventHandler)

	//TODO: add workers for active trackers for update time info (check type buttons before edit message)
	//TODO set tasks text labels and duration (like as breaks) - some tasks by active tracker
	//TODO some trackers per day by user - feature: set random or user traker name
	//TODO: add user break type buttons (coffe break for example) - settings
	//TODO: add set user GMT - settings
	//TODO save info to db
	result.list = commandsList

	//TODO: read from db to RAM active tasks(rename task to traker)
	go result.tracking()

	return &result
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

func (d *data) startTaskButtonEventHandler(msg *tgbotapi.Message, c *tgModel.Command) *tgModel.HandlerResult {
	log.Info().Msg("startTaskButtonEventHandler")
	task := d.AddTask(msg.Chat.ID, msg.MessageID)
	return tgModel.SimpleEditWithButtons(msg.Chat.ID, msg.MessageID, task.Title, d.activeTaskButtons())
}

func (d *data) settingsButtonEventHandler(msg *tgbotapi.Message, c *tgModel.Command) *tgModel.HandlerResult {
	log.Info().Msg("settingsButtonEventHandler")
	return tgModel.SimpleReply(msg.Chat.ID, "Not implement", msg.MessageID)
}

func (d *data) takeBreakButtonEventHandler(msg *tgbotapi.Message, c *tgModel.Command) *tgModel.HandlerResult {
	log.Info().Msg("takeBreakButtonEventHandler")
	task, exist := d.SetTaskBreak(msg.Chat.ID)
	if !exist {
		return tgModel.SimpleReply(msg.Chat.ID, "Task not found, sorry, create new by /timeTrack", msg.MessageID)
	}
	return tgModel.SimpleEditWithButtons(msg.Chat.ID, msg.MessageID, task.Title, d.breakTaskButtons())
}

func (d *data) stopBreakButtonEventHandler(msg *tgbotapi.Message, c *tgModel.Command) *tgModel.HandlerResult {
	log.Info().Msg("stopBreakButtonEventHandler")
	task, exist := d.StopTaskBreak(msg.Chat.ID)
	if !exist {
		return tgModel.SimpleReply(msg.Chat.ID, "Task not found, sorry, create new by /timeTrack", msg.MessageID)
	}
	return tgModel.SimpleEditWithButtons(msg.Chat.ID, msg.MessageID, task.Title, d.activeTaskButtons())
}

func (d *data) StoppedTaskButtonEventHandler(msg *tgbotapi.Message, c *tgModel.Command) *tgModel.HandlerResult {
	log.Info().Msg("StoppedTaskButtonEventHandler")
	task, exist := d.StopTask(msg.Chat.ID)
	if !exist {
		return tgModel.SimpleReply(msg.Chat.ID, "Task not found, sorry, create new by /timeTrack", msg.MessageID)
	}
	return tgModel.SimpleEdit(msg.Chat.ID, msg.MessageID, task.Title)
}
