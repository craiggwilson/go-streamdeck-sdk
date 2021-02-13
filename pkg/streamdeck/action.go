package streamdeck

import (
	"github.com/craiggwilson/streamdeck-plugins/pkg/streamdeck/streamdeckcore"
)

// Action is an implementation of an Action.
type Action interface {
	UUID() streamdeckcore.ActionUUID
	Initialize(eventPublisher streamdeckcore.EventPublisher)
}
