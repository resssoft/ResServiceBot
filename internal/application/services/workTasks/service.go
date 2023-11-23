package workTasks

import (
	"context"
	"database/sql"
	"fmt"
	"fun-coice/internal/application/services/workTasks/repository"
	mongo_repo "fun-coice/internal/application/services/workTasks/repository/mongo"
	sqlRepo "fun-coice/internal/application/services/workTasks/repository/sql"
	"fun-coice/internal/application/services/workTasks/track"
	"fun-coice/internal/database"
	tgModel "fun-coice/internal/domain/commands/tg"
	"github.com/doug-martin/goqu/v9"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
	"github.com/sasha-s/go-deadlock"
	"time"
)

type data struct {
	list    tgModel.Commands
	users   map[int64]track.User //temporary
	builder goqu.DialectWrapper
	tracks  track.Tracks
	//mutex         *sync.Mutex
	mutex         deadlock.Mutex
	buttons       map[string]track.Button
	messageSender tgModel.MessageSender
	repo          repository.Repository
}

const trackingDuration = time.Second * 31

func New(dbSQL *sql.DB, mongoClient database.MongoClientApplication) tgModel.Service {
	var repo repository.Repository
	switch {
	case mongoClient != nil:
		repo, _ = mongo_repo.NewMongoRepo(mongoClient) //TODO: check errors for all services
	case dbSQL != nil:
		repo, _ = sqlRepo.NewSQLRepo(dbSQL) //TODO: check errors for all services
	default:
		//RAM repo
	}
	result := data{
		users:   make(map[int64]track.User), // temporary
		builder: goqu.Dialect("sqlite3"),
		//mutex:   &sync.Mutex{},
		tracks:  make(track.Tracks),
		buttons: make(map[string]track.Button),
		repo:    repo,
	}
	result.initCommands()
	go result.tracking(context.Background())

	return &result
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
				if track.IsStopped() {
					continue
				}
				d.updateTrackMessage(track)
			}
			d.mutex.Unlock()
		}
	}
}

func (d *data) updateTrackMessage(track track.Track) {
	newTitle := track.GetTitle()
	if newTitle == track.Title {
		return
	}
	if track.MsgId == 0 {
		return
	}
	track.Title = newTitle
	if d.messageSender != nil {
		d.messageSender.PushHandleResult() <- tgModel.SimpleEditWithButtons(track.UserId, track.MsgId, track.Title, d.keyboard(track))
	}
}

func (d *data) AddTrack(uid int64, msgId int) track.Track {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	userTrack := track.Track{
		Start:   time.Now(),
		UserId:  uid,
		MsgId:   msgId,
		Status:  track.StatusProgress,
		Tasks:   make(map[int]track.TimeItem),
		BotName: d.messageSender.BotName(),
		Code:    fmt.Sprintf("%v-%v-%s", uid, time.Now().Unix(), d.messageSender.BotName()),
	}
	userTrack.AddTask(track.DefaultTaskName)
	log.Info().Any("AddTrack", userTrack).Send()
	d.tracks[uid] = userTrack
	return userTrack
}

func (d *data) AddTask(uid int64, name string) (track.Track, bool) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	userTrack, exist := d.tracks[uid]
	if !exist {
		return userTrack, false
	}
	userTrack.AddTask(name)
	d.tracks[uid] = userTrack
	return userTrack, true
}

func (d *data) SetTrackBreak(uid int64) (track.Track, bool) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	userTrack, exist := d.tracks[uid]
	if !exist {
		return track.Track{}, false
	}
	{ //debug
		tasksInfo := ""
		for i, t := range userTrack.Tasks {
			tasksInfo += fmt.Sprintf("\n[%v]%s-%s/%s %s",
				i, t.Start.Format(track.TimeFormatS), t.End.Format(track.TimeFormatS), track.Duration(t.Duration), t.Name)
		}
		log.Info().
			Str("SetTrackBreak before", tasksInfo).
			Send()
		d.tracks[uid] = userTrack.SetBreak()
	}
	d.tracks[uid] = userTrack
	userTrack.Title = userTrack.GetTitle()
	{ //debug
		tasksInfo := ""
		for i, t := range userTrack.Tasks {
			tasksInfo += fmt.Sprintf("\n[%v]%s-%s/%s %s",
				i, t.Start.Format(track.TimeFormatS), t.End.Format(track.TimeFormatS), track.Duration(t.Duration), t.Name)
		}
		log.Info().
			Str("SetTrackBreak After", tasksInfo).
			Send()
		log.Info().Any("SetTrackBreak", d.tracks[uid]).Send()

	}
	return userTrack, exist
}

func (d *data) StopTrackBreak(uid int64) (track.Track, bool) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	userTrack, exist := d.tracks[uid]
	if exist {
		tasksInfo := ""
		for i, t := range userTrack.Tasks {
			tasksInfo += fmt.Sprintf("\n[%v]%s-%s/%s %s",
				i, t.Start.Format(track.TimeFormatS), t.End.Format(track.TimeFormatS), track.Duration(t.Duration), t.Name)
		}
		log.Info().
			Str("StopTrackBreak before", tasksInfo).
			Send()
		d.tracks[uid] = userTrack.StopBreak()
	}
	d.tracks[uid] = userTrack
	userTrack.Title = userTrack.GetTitle()
	tasksInfo := ""
	for i, t := range userTrack.Tasks {
		tasksInfo += fmt.Sprintf("\n[%v]%s-%s/%s %s",
			i, t.Start.Format(track.TimeFormatS), t.End.Format(track.TimeFormatS), track.Duration(t.Duration), t.Name)
	}
	log.Info().
		Str("StopTrackBreak After", tasksInfo).
		Send()
	log.Info().Any("SetTrackBreak", d.tracks[uid]).Send()
	log.Info().Any("StopTrackBreak", d.tracks[uid]).Send()
	return userTrack, exist
}

func (d *data) StopTrack(uid int64) (track.Track, bool) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	userTrack, exist := d.tracks[uid]
	if exist {
		d.tracks[uid] = userTrack.StopTrack()
	}
	log.Info().Any("StopTrack", d.tracks[uid]).Send()
	return userTrack, exist
}

func (d *data) activeTrackButtons(uid int64) *tgbotapi.InlineKeyboardMarkup {
	userTrack, exist := d.tracks[uid]
	if !exist {
		return tgModel.GetTGButtons(tgModel.KeyBoardTG{})
	}
	tasks, keys := userTrack.GetTasks(true)
	var taskRows []tgModel.KeyBoardRowTG
	taskRows = append(taskRows,
		d.ButtonRow(track.StartTaskEvent, track.TakeBreakEvent, track.StoppedTaskEvent, track.SettingsEvent),
		d.ButtonRow(track.SetTaskNameEvent))
	for _, taskIndex := range keys {
		taskRows = append(
			taskRows,
			tgModel.KBButs(
				tgModel.KeyBoardButtonTG{
					Text: fmt.Sprintf(track.TaskIcon + " " + tasks[taskIndex].Name),
					Data: fmt.Sprintf("%s:%v", track.SetTaskAction, taskIndex),
				}))
	}
	//taskRows = append(taskRows, d.ButtonRow(startTaskEvent))

	return tgModel.GetTGButtons(tgModel.KBRows(taskRows...))
}

func (d *data) breakTrackButtons(_ int64) *tgbotapi.InlineKeyboardMarkup {
	return tgModel.GetTGButtons(tgModel.KBRows(
		d.ButtonRow(track.StopBreakEvent, track.StoppedTaskEvent, track.SettingsEvent),
		d.ButtonRow(track.SetBreakNameEvent)))
}

func (d *data) trackButtons(_ int64) *tgbotapi.InlineKeyboardMarkup {
	return tgModel.GetTGButtons(tgModel.KBRows(d.ButtonRow(track.StartTrackEvent, track.ShowProfileEvent)))
}

func (d *data) updateActiveTaskName(uid int64, newName string) (track.Track, bool) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	userTrack, exist := d.tracks[uid]
	if !exist {
		return userTrack, false
	}
	activeTask, exist := userTrack.Tasks[userTrack.ActiveTask]
	if !exist {
		return userTrack, false
	}
	activeTask.Name = newName
	userTrack.UpdateTask(activeTask)
	return userTrack, true
}

func (d *data) setActiveTask(uid int64, id int) bool {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	userTrack, exist := d.tracks[uid]
	if !exist {
		return false
	}
	activeTask, exist := userTrack.Tasks[userTrack.ActiveTask]
	if exist {
		activeTask.End = time.Now()
		activeTask.Duration = activeTask.Duration + activeTask.End.Sub(activeTask.Start)
		userTrack.Tasks[userTrack.ActiveTask] = activeTask
	}
	nextTask, exist := userTrack.Tasks[id]
	if !exist {
		return false
	} else {
		nextTask.Start = time.Now()
		userTrack.Tasks[id] = nextTask
	}
	userTrack.ActiveTask = id
	userTrack.Title = userTrack.GetTitle()
	d.tracks[uid] = userTrack
	return true
}

func (d *data) keyboard(t track.Track) *tgbotapi.InlineKeyboardMarkup {
	var keyboard *tgbotapi.InlineKeyboardMarkup
	switch t.Status {
	case track.StatusProgress:
		keyboard = d.activeTrackButtons(t.UserId)
	case track.StatusPause:
		keyboard = d.breakTrackButtons(t.UserId)
	default:
		keyboard = d.activeTrackButtons(t.UserId)
	}
	return keyboard
}

func (d *data) GetTrack(uid int64) (track.Track, bool) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	userTrack, exist := d.tracks[uid]
	return userTrack, exist
}
