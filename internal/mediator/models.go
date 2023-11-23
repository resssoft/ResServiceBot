package mediator

type Listener interface {
	Listen(eventName EventName, event interface{})
}

type Job struct {
	EventName EventName
	EventType interface{}
}

type EventName string

const AppExit EventName = "app.exit"
const LogToFile EventName = "fileLogger.log.data"
const SetLogDebugMode EventName = "log.mode.debug"
const SetLogInfoMode EventName = "log.mode.info"

type FileLoggerEvent struct {
	Src         string
	Data        string
	WithoutTime bool
	ToDebug     bool
}

var FileLoggerEvents = []EventName{
	LogToFile,
}

const FileLogFatal = "fatal"
const FileLogErrors = "errors"
const FileLogContacts = "contacts"
const FileLogWebHooks = "webHooks"
const FileLogRequests = "requests"
const FileLogMessenger = "messenger"
const FileLogAmoCRM = "amoCRM"
const FileLogAmoLatency = "amoLatency"
