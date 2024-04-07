package workTasks

import (
	"context"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"

	"fun-coice/internal/application/services/workTasks/track"
	tgModel "fun-coice/internal/domain/commands/tg"
	tgbotapi "fun-coice/pkg/telegram-bot-api"
)

func (d *data) timeTrack(msg *tgbotapi.Message, _ *tgModel.Command) *tgModel.HandlerResult {
	d.checkUser(msg.From)
	return tgModel.SimpleWithButtons(msg.Chat.ID, track.TimeTrackTitle, d.trackButtons(msg.Chat.ID))
}

func (d *data) startTrackButtonEventHandler(msg *tgbotapi.Message, c *tgModel.Command) *tgModel.HandlerResult {
	d.checkUser(msg.From)
	userTrack := d.AddTrack(msg.Chat.ID, msg.MessageID)
	return tgModel.SimpleEditWithButtons(msg.Chat.ID, msg.MessageID, userTrack.Title, d.activeTrackButtons(msg.Chat.ID))
}

func (d *data) settingsButtonEventHandler(msg *tgbotapi.Message, c *tgModel.Command) *tgModel.HandlerResult {
	d.checkUser(msg.From)
	return tgModel.SimpleReply(msg.Chat.ID, "Not implement", msg.MessageID)
}

func (d *data) profileButtonEventHandler(msg *tgbotapi.Message, c *tgModel.Command) *tgModel.HandlerResult {
	d.checkUser(msg.From)
	return tgModel.SimpleWithButtons(msg.Chat.ID, track.TimeTrackTitle, d.userSettingsMainButtons(msg.Chat.ID))
}

func (d *data) takeBreakButtonEventHandler(msg *tgbotapi.Message, c *tgModel.Command) *tgModel.HandlerResult {
	d.checkUser(msg.From)
	userTrack, exist := d.SetTrackBreak(msg.Chat.ID)
	if !exist {
		return tgModel.SimpleReply(msg.Chat.ID, track.TrackNotFoundErrMsg, msg.MessageID)
	}
	return tgModel.SimpleEditWithButtons(msg.Chat.ID, msg.MessageID, userTrack.Title, d.breakTrackButtons(msg.Chat.ID))
}

func (d *data) stopBreakButtonEventHandler(msg *tgbotapi.Message, _ *tgModel.Command) *tgModel.HandlerResult {
	d.checkUser(msg.From)
	userTrack, exist := d.StopTrackBreak(msg.Chat.ID)
	if !exist {
		return tgModel.SimpleReply(msg.Chat.ID, track.TrackNotFoundErrMsg, msg.MessageID)
	}
	return tgModel.SimpleEditWithButtons(msg.Chat.ID, msg.MessageID, userTrack.Title, d.activeTrackButtons(msg.Chat.ID))
}

func (d *data) StoppedTrackButtonEventHandler(msg *tgbotapi.Message, _ *tgModel.Command) *tgModel.HandlerResult {
	d.checkUser(msg.From)
	userTrack, exist := d.StopTrack(msg.Chat.ID)
	if !exist {
		return tgModel.SimpleReply(msg.Chat.ID, track.TrackNotFoundErrMsg, msg.MessageID)
	}
	return tgModel.SimpleEdit(msg.Chat.ID, msg.MessageID, userTrack.Title)
}

func (d *data) setTaskNameButtonEventHandler(msg *tgbotapi.Message, c *tgModel.Command) *tgModel.HandlerResult {
	d.checkUser(msg.From)
	if c.Arguments.Raw == "" {
		return tgModel.DeferredWithText(msg.Chat.ID, "Enter new task name", "timeTrack_set_task_name", "", nil)
	}
	foundedTrack, exist := d.updateActiveTaskName(msg.Chat.ID, c.Arguments.Raw)
	if !exist {
		return tgModel.Simple(msg.Chat.ID, track.TrackNotFoundErrMsg)
	}
	d.updateTrackMessage(foundedTrack)
	return tgModel.Simple(msg.Chat.ID, "Ok")
	//return tgModel.EmptyCommand()
}

func (d *data) addTaskButtonEventHandler(msg *tgbotapi.Message, c *tgModel.Command) *tgModel.HandlerResult {
	d.checkUser(msg.From)
	if c.Arguments.Raw == "" {
		return tgModel.DeferredWithText(msg.Chat.ID, "Enter task name", "timeTrack_add_task", "", nil)
	}
	_, exist := d.AddTask(msg.Chat.ID, c.Arguments.Raw)
	if !exist {
		return tgModel.Simple(msg.Chat.ID, track.TrackNotFoundErrMsg)
	}
	return tgModel.Simple(msg.Chat.ID, "Ok")
	//TODO: EmptyCommand
	//return tgModel.EmptyCommand()
}

func (d *data) SetActiveTask(msg *tgbotapi.Message, c *tgModel.Command) *tgModel.HandlerResult {
	d.checkUser(msg.From)
	anonId := 0
	separated := strings.Split(c.Data, ":")
	if len(separated) == 3 {
		anonId, _ = strconv.Atoi(separated[2])
	}
	d.setActiveTask(msg.Chat.ID, anonId)
	userTrack, exist := d.GetTrack(msg.Chat.ID)
	if !exist {
		return tgModel.Simple(msg.Chat.ID, track.TrackNotFoundErrMsg)
	}
	d.updateTrackMessage(userTrack)
	return tgModel.EmptyCommand()
}

func (d *data) NotImplementHandler(msg *tgbotapi.Message, c *tgModel.Command) *tgModel.HandlerResult {
	d.checkUser(msg.From)
	log.Info().Msg("setTaskNameButtonEventHandler")
	return tgModel.SimpleReply(msg.Chat.ID, "Not implement", msg.MessageID)
}

func (d *data) UserSettingsTasks(msg *tgbotapi.Message, c *tgModel.Command) *tgModel.HandlerResult {
	d.checkUser(msg.From)
	if c.Arguments.Raw == "" {
		return tgModel.DeferredWithText(msg.Chat.ID, "Напишите список задач, которые будут созданы автоматически при старте трэка, каждую на новой строке, последняя будет запущена при старте трэка", "timeTrack_set_defaultTasks", "", nil)
	}
	tasks := strings.Split(c.Arguments.Raw, "\n")
	currentUser := d.user(msg.From.ID)
	currentUser.Settings.DefaultTaskNames = tasks
	err := d.userSave(context.Background(), currentUser)
	if err != nil {
		log.Info().Err(err).
			Int64("tg user id", msg.From.ID).
			Str("c.Arguments.Raw", c.Arguments.Raw).
			Msg("UserSettingsTasks save user err")
		return tgModel.Simple(msg.Chat.ID, "Some problems to save, write admin")
	}
	return tgModel.Simple(msg.Chat.ID, "Ok")
}
