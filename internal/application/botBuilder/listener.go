package botBuilder

import (
	"github.com/rs/zerolog/log"

	"fun-coice/internal/mediator"
)

type Listener struct {
	Client *BuilderService
}

func (u Listener) Listen(_ mediator.EventName, event interface{}) {
	switch event := event.(type) {
	case mediator.BotBuilderServiceChangeEvent:
		u.Client.AddService()
	default:
		log.Printf("registered an invalid botBuilder event: %T\n", event)
	}
}
