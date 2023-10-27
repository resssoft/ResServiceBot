package workTasks

import (
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (d data) userInfo(u *tgbotapi.User) (User, error) {
	if u == nil {
		return User{}, errors.New("user is nil")
	}
	existed, exist := d.users[u.ID]
	if exist {
		return existed, nil
	}
	newUser := User{
		IDStr:  fmt.Sprintf("%v", u.ID),
		tgUser: *u,
	}
	d.userSave(newUser)
	newUser.IsNew = true
	return newUser, nil
}

func (d data) userSave(u User) error {
	//TODO: save to DB
	d.users[u.tgUser.ID] = u
	return nil
}

func (d data) userInfoByID(id int64) (User, error) {
	//TODO: get from db
	existed, exist := d.users[id]
	if exist {
		return existed, nil
	}
	return User{}, errors.New(fmt.Sprintf("user not found by id %v", id))
}
