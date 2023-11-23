package fileLogger

import (
	"fun-coice/internal/mediator"
	"github.com/rs/zerolog/log"
)

type Listener struct {
	Client *Client
}

func (u Listener) Listen(_ mediator.EventName, event interface{}) {
	switch event := event.(type) {
	case mediator.FileLoggerEvent:
		u.Client.Log(event.Src, event.Data, event.WithoutTime, event.ToDebug)
	default:
		log.Printf("registered an invalid fileLogger event: %T\n", event)
	}
}
