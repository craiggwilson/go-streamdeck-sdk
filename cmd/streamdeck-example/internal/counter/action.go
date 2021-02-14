package counter

import (
	"context"
	"strconv"

	"github.com/craiggwilson/streamdeck-plugins/pkg/streamdeck"
	"github.com/craiggwilson/streamdeck-plugins/pkg/streamdeck/streamdeckcore"
)

const uuid = "com.craiggwilson.streamdeck.counter"

func New() streamdeck.ActionInstanceFactory {
	return streamdeck.NewActionInstanceFactory(
		uuid,
		func(eventContext streamdeckcore.EventContext, publisher streamdeck.Publisher) streamdeck.ActionInstance {
			return &ActionInstance{
				eventContext: eventContext,
				publisher: publisher,
			}
		},
	)
}

type ActionInstance struct {
	eventContext streamdeckcore.EventContext
	publisher streamdeck.Publisher
	count     int
}

func (a *ActionInstance) UUID() streamdeckcore.ActionUUID {
	return uuid
}

func (a *ActionInstance) Context() streamdeckcore.EventContext {
	return a.eventContext
}

func (a *ActionInstance) HandleKeyDown(ctx context.Context, event streamdeckcore.KeyDownEvent) error {
	a.count++
	return a.publisher.SetTitle(streamdeckcore.SetTitlePayload{
		Title: strconv.Itoa(a.count),
		Target: streamdeckcore.HardwareAndSoftware,
	})
}
