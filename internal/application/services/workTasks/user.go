package workTasks

import (
	"context"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"fun-coice/internal/application/services/workTasks/track"
	tgbotapi "fun-coice/pkg/telegram-bot-api"
)

func (d *data) userDefault(id int64) track.User {
	return track.User{
		TgUser: tgbotapi.User{
			ID: id,
		},
		Settings: track.UserSettings{
			ID:               id,
			DefaultTaskNames: make([]string, 0),
			LangISO:          "en",
		},
	}
}

func (d *data) checkUser(u *tgbotapi.User) bool {
	if u == nil {
		return false
	}
	if u.IsBot {
		return false
	}
	d.mutexUser.Lock()
	_, exist := d.users[u.ID]
	d.mutexUser.Unlock()
	if exist {
		return true
	}
	ctx := context.Background()
	founded, err := d.userRepo.Get(ctx, u.ID)
	if founded == nil || err != nil {
		newUser := track.User{
			TgUser: *u,
			Settings: track.UserSettings{
				ID:               u.ID,
				DefaultTaskNames: []string{track.DefaultTaskName},
				LangISO:          u.LanguageCode,
			},
		}
		err = d.userSave(ctx, newUser)
		if err != nil {
			log.Info().Err(err).Msg("checkUser: userSave err")
			return false
		}
	}
	return true
}

func (d *data) userSave(ctx context.Context, u track.User) error {
	d.mutexUser.Lock()
	defer d.mutexUser.Unlock()
	d.users[u.TgUser.ID] = u
	var err error
	if u.MongoID == primitive.NilObjectID {
		_, err = d.userRepo.Create(ctx, u)
	} else {
		err = d.userRepo.Update(ctx, &u)
	}
	return err
}

func (d *data) user(id int64) track.User {
	d.mutexUser.Lock()
	defer d.mutexUser.Unlock()
	existed, exist := d.users[id]
	if exist {
		return existed
	}
	founded, err := d.userRepo.Get(context.Background(), id)
	if err != nil || founded == nil {
		log.Info().Err(err).Msg("user: userRepo.Get err")
		return d.userDefault(id)
	}
	return *founded
}
