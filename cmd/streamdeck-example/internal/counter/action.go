package counter

import (
	"context"
	"strconv"

	"github.com/craiggwilson/go-streamdeck-sdk/streamdeckevent"

	"github.com/craiggwilson/go-streamdeck-sdk"
)

const uuid = "com.craiggwilson.streamdeck.example.counter"

func New() *streamdeck.DefaultAction {
	return streamdeck.NewDefaultAction(
		uuid,
		func(eventContext streamdeck.EventContext, publisher streamdeck.ActionInstancePublisher) streamdeck.ActionInstance {
			return &ActionInstance{
				eventContext: eventContext,
				publisher: publisher,
			}
		},
	)
}

type ActionInstance struct {
	eventContext streamdeck.EventContext
	publisher    streamdeck.ActionInstancePublisher
	count        int
}

func (a *ActionInstance) UUID() streamdeck.ActionUUID {
	return uuid
}

func (a *ActionInstance) Context() streamdeck.EventContext {
	return a.eventContext
}

func (a *ActionInstance) HandleKeyDown(_ context.Context, _ streamdeckevent.KeyDown) error {
	a.count++
	return a.publisher.SetTitle(streamdeckevent.SetTitlePayload{
		Title: strconv.Itoa(a.count),
		Target: streamdeckevent.HardwareAndSoftware,
	})
}
