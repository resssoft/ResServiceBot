package workTasks

import (
	"context"
	"database/sql"
	"fmt"
	"fun-coice/internal/application/services/workTasks/repository"
	"fun-coice/internal/application/services/workTasks/repository/mongo/track"
	"fun-coice/internal/application/services/workTasks/repository/mongo/user"
	sqlRepo "fun-coice/internal/application/services/workTasks/repository/sql"
	"fun-coice/internal/application/services/workTasks/track"
	"fun-coice/internal/database"
	tgModel "fun-coice/internal/domain/commands/tg"
	"github.com/doug-martin/goqu/v9"
	"github.com/rs/zerolog/log"
	"github.com/sasha-s/go-deadlock"
	"time"
)

type data struct {
	list          tgModel.Commands
	users         map[int64]track.User
	mutexUser     deadlock.Mutex
	builder       goqu.DialectWrapper
	tracks        track.Tracks
	mutex         deadlock.Mutex
	buttons       map[string]track.Button
	messageSender tgModel.MessageSender
	trackRepo     repository.TrackRepository
	userRepo      repository.UserRepository
	//mutex         *sync.Mutex
}

const trackingDuration = time.Second * 31

func New(dbSQL *sql.DB, mongoClient database.MongoClientApplication) tgModel.Service {
	var trackRepo repository.TrackRepository
	var userRepo repository.UserRepository
	switch {
	case mongoClient != nil:
		trackRepo, _ = trackMongoRepository.NewTrackRepo(mongoClient)
		userRepo, _ = userMongoRepository.NewUserRepo(mongoClient)
	case dbSQL != nil:
		trackRepo, _ = sqlRepo.NewSQLRepo(dbSQL) //TODO: check errors for all services
	default:
		//RAM trackRepo
	}
	result := data{
		users:   make(map[int64]track.User), // temporary
		builder: goqu.Dialect("sqlite3"),
		//mutex:   &sync.Mutex{},
		tracks:    make(track.Tracks),
		buttons:   make(map[string]track.Button),
		trackRepo: trackRepo,
		userRepo:  userRepo,
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
	fmt.Println("===updateTrackMessage start")
	newTitle := track.GetTitle()
	if newTitle == track.Title {
		fmt.Println("===updateTrackMessage title no dif")
		return
	}
	if track.MsgId == 0 {
		fmt.Println("===updateTrackMessage MsgId = 0")
		return
	}
	track.Title = newTitle
	if d.messageSender != nil {
		fmt.Println("===updateTrackMessage push")
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
	userTrack.AddTasks(d.user(uid).Settings.DefaultTaskNames)
	log.Info().Any("AddTrack", userTrack).Send()
	log.Info().Msg("==================================!!=== AddTrack")
	userTrack, err := d.trackRepo.Create(context.Background(), userTrack) // handle error
	if err != nil {
		log.Info().Err(err).Any("Update ERR", userTrack).Send()
	}
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

	d.updateTrackMessage(userTrack)
	log.Info().Msg("===================================== AddTask")
	err := d.trackRepo.Update(context.Background(), &userTrack) // handle error
	if err != nil {
		log.Info().Err(err).Any("Update ERR", userTrack).Send()
	}
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
	log.Info().Msg("===================================== SetTrackBreak")
	err := d.trackRepo.Update(context.Background(), &userTrack) // handle error
	if err != nil {
		log.Info().Err(err).Any("Update ERR", userTrack).Send()
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
	log.Info().Msg("===================================== StopTrackBreak")
	err := d.trackRepo.Update(context.Background(), &userTrack) // handle error
	if err != nil {
		log.Info().Err(err).Any("Update ERR", userTrack).Send()
	}
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
	log.Info().Msg("===================================== StopTrack")
	err := d.trackRepo.Update(context.Background(), &userTrack) // handle error
	if err != nil {
		log.Info().Err(err).Any("Update ERR", userTrack).Send()
	}
	return userTrack, exist
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
	log.Info().Msg("===================================== updateActiveTaskName")
	err := d.trackRepo.Update(context.Background(), &userTrack) // handle error
	if err != nil {
		log.Info().Err(err).Any("Update ERR", userTrack).Send()
	}
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
	log.Info().Msg("===================================== setActiveTask")
	err := d.trackRepo.Update(context.Background(), &userTrack) // handle error
	if err != nil {
		log.Info().Err(err).Any("Update ERR", userTrack).Send()
	}
	return true
}

func (d *data) GetTrack(uid int64) (track.Track, bool) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	userTrack, exist := d.tracks[uid]
	return userTrack, exist
}

func (d *data) UpdateTrack(uid int64, userTrack track.Track) {
	err := d.trackRepo.Update(context.Background(), &userTrack) // handle error
	if err != nil {
		log.Info().Err(err).Any("Update ERR", userTrack).Send()
	}
	d.tracks[uid] = userTrack
}
