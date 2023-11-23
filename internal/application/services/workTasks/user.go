package workTasks

import (
	"errors"
	"fmt"
	"fun-coice/internal/application/services/workTasks/track"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (d *data) userInfo(u *tgbotapi.User) (track.User, error) {
	if u == nil {
		return track.User{}, errors.New("user is nil")
	}
	existed, exist := d.users[u.ID]
	if exist {
		return existed, nil
	}
	newUser := track.User{
		IDStr:  fmt.Sprintf("%v", u.ID),
		TgUser: *u,
	}
	d.userSave(newUser)
	newUser.IsNew = true
	return newUser, nil
}

func (d *data) userSave(u track.User) error {
	//TODO: save to DB
	d.users[u.TgUser.ID] = u
	return nil
}

func (d *data) userInfoByID(id int64) (track.User, error) {
	//TODO: get from db
	existed, exist := d.users[id]
	if exist {
		return existed, nil
	}
	return track.User{}, errors.New(fmt.Sprintf("user not found by id %v", id))
}
