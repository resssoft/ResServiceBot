package workTasks

import (
	"fun-coice/internal/application/services/workTasks/track"
	tgModel "fun-coice/internal/domain/commands/tg"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
	"strconv"
	"strings"
)

func (d *data) timeTrack(msg *tgbotapi.Message, _ *tgModel.Command) *tgModel.HandlerResult {
	return tgModel.SimpleWithButtons(msg.Chat.ID, track.TimeTrackTitle, d.trackButtons(msg.Chat.ID))
}

func (d *data) startTrackButtonEventHandler(msg *tgbotapi.Message, c *tgModel.Command) *tgModel.HandlerResult {
	log.Info().Msg("startTaskButtonEventHandler")
	userTrack := d.AddTrack(msg.Chat.ID, msg.MessageID)
	return tgModel.SimpleEditWithButtons(msg.Chat.ID, msg.MessageID, userTrack.Title, d.activeTrackButtons(msg.Chat.ID))
}

func (d *data) settingsButtonEventHandler(msg *tgbotapi.Message, c *tgModel.Command) *tgModel.HandlerResult {
	log.Info().Msg("settingsButtonEventHandler")
	return tgModel.SimpleReply(msg.Chat.ID, "Not implement", msg.MessageID)
}

func (d *data) takeBreakButtonEventHandler(msg *tgbotapi.Message, c *tgModel.Command) *tgModel.HandlerResult {
	log.Info().Msg("takeBreakButtonEventHandler")
	userTrack, exist := d.SetTrackBreak(msg.Chat.ID)
	if !exist {
		return tgModel.SimpleReply(msg.Chat.ID, track.TrackNotFoundErrMsg, msg.MessageID)
	}
	return tgModel.SimpleEditWithButtons(msg.Chat.ID, msg.MessageID, userTrack.Title, d.breakTrackButtons(msg.Chat.ID))
}

func (d *data) stopBreakButtonEventHandler(msg *tgbotapi.Message, _ *tgModel.Command) *tgModel.HandlerResult {
	log.Info().Msg("stopBreakButtonEventHandler")
	userTrack, exist := d.StopTrackBreak(msg.Chat.ID)
	if !exist {
		return tgModel.SimpleReply(msg.Chat.ID, track.TrackNotFoundErrMsg, msg.MessageID)
	}
	return tgModel.SimpleEditWithButtons(msg.Chat.ID, msg.MessageID, userTrack.Title, d.activeTrackButtons(msg.Chat.ID))
}

func (d *data) StoppedTrackButtonEventHandler(msg *tgbotapi.Message, _ *tgModel.Command) *tgModel.HandlerResult {
	log.Info().Msg("StoppedTrackButtonEventHandler")
	userTrack, exist := d.StopTrack(msg.Chat.ID)
	if !exist {
		return tgModel.SimpleReply(msg.Chat.ID, track.TrackNotFoundErrMsg, msg.MessageID)
	}
	return tgModel.SimpleEdit(msg.Chat.ID, msg.MessageID, userTrack.Title)
}

func (d *data) setTaskNameButtonEventHandler(msg *tgbotapi.Message, c *tgModel.Command) *tgModel.HandlerResult {
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
	if c.Arguments.Raw == "" {
		return tgModel.DeferredWithText(msg.Chat.ID, "Enter task name", "timeTrack_add_task", "", nil)
	}
	foundedTrack, exist := d.AddTask(msg.Chat.ID, c.Arguments.Raw)
	if !exist {
		return tgModel.Simple(msg.Chat.ID, track.TrackNotFoundErrMsg)
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
		return tgModel.Simple(msg.Chat.ID, track.TrackNotFoundErrMsg)
	}
	d.updateTrackMessage(userTrack)
	return tgModel.EmptyCommand()
}

func (d *data) NotImplementHandler(msg *tgbotapi.Message, c *tgModel.Command) *tgModel.HandlerResult {
	log.Info().Msg("setTaskNameButtonEventHandler")
	return tgModel.SimpleReply(msg.Chat.ID, "Not implement", msg.MessageID)
}
