package workTasks

import (
	"fun-coice/internal/application/services/workTasks/track"
	tgModel "fun-coice/internal/domain/commands/tg"
)

func (d *data) initCommands() {
	commandsList := tgModel.NewCommands()
	commandsList.AddSimple("timeTrack", "Show time track controls", d.timeTrack)
	commandsList.AddSimple("timeTrack_add_task", "Add task to active track, need task name parameter", d.addTaskButtonEventHandler)
	commandsList.AddSimple("timeTrack_set_task_name", "Add task to active track, need task name parameter", d.setTaskNameButtonEventHandler)

	commandsList.AddEvent(track.SetTaskEvent, d.SetActiveTask)

	d.list = commandsList

	d.addButton("🚗 Начать трэкинг", track.StartTrackEvent, d.startTrackButtonEventHandler)
	d.addButton("⚙️", track.SettingsEvent, d.settingsButtonEventHandler)
	d.addButton("⏸", track.TakeBreakEvent, d.takeBreakButtonEventHandler)
	d.addButton("▶️", track.StopBreakEvent, d.stopBreakButtonEventHandler)
	d.addButton("🏁", track.StoppedTaskEvent, d.StoppedTrackButtonEventHandler)

	d.addButton("📝 Задать имя активной задачи", track.SetTaskNameEvent, d.setTaskNameButtonEventHandler)
	d.addButton("➕", track.StartTaskEvent, d.addTaskButtonEventHandler)
	d.addButton("👤 Профиль", track.ShowProfileEvent, d.NotImplementHandler)
	d.addButton("📝 Задать имя перерыву", track.SetBreakNameEvent, d.NotImplementHandler)

	//TODO edit time, duration, start, end

	//TODO some trackers per day by user - feature: set random or user traker name
	//TODO: add user break type buttons (coffe break for example) - tracker options
	//TODO: add set user GMT - settings user
	//TODO: show logs - settings tracker
	//TODO save info to db
	//TODO change/correct current time of tracker task or break

	//TODO: read from db to RAM active tasks(rename task to traker)
}
