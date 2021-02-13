package counter

import (
	"context"

	"github.com/craiggwilson/streamdeck-plugins/pkg/streamdeck"
	"github.com/craiggwilson/streamdeck-plugins/pkg/streamdeck/streamdeckcore"
)

func New() streamdeck.Action {
	return streamdeck.NewInstancedAction(
		streamdeck.ActionInstanceFactoryFunc(func(_ streamdeckcore.EventContext) streamdeck.Action {
			return &Action{}
		}),
	)
}

type Action struct {
	eventPublisher streamdeckcore.EventPublisher
	count int
}

func (a *Action) UUID() streamdeckcore.ActionUUID {
	return "com.craiggwilson.streamdeck.counter"
}

func (a *Action) Initialize(eventPublisher streamdeckcore.EventPublisher) {
	a.eventPublisher = eventPublisher
}

func (a *Action) HandleKeyDown(ctx context.Context, event streamdeckcore.KeyDownEvent) error {
	a.count++
	return nil
}
