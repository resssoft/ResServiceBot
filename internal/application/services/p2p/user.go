package p2p

import (
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
)

func (d *data) userInfo(u *tgbotapi.User) (User, error) {
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

func (d *data) userSave(u User) error {
	//TODO: save to DB
	d.users[u.tgUser.ID] = u
	return nil
}

func (d *data) userInfoByID(id int64) (User, error) {
	//TODO: get from db
	existed, exist := d.users[id]
	if exist {
		return existed, nil
	}
	return User{}, errors.New(fmt.Sprintf("user not found by id %v", id))
}

func (d *data) userRegistered(by string) bool {
	intval, err := strconv.ParseInt(by, 10, 64)
	if err != nil {
		return false
	}
	user, err := d.userInfoByID(intval)
	if err == nil {
		if !user.IsNew {
			return true
		}
	}
	return false
}
