package counter

import (
	"context"
	"strconv"

	"github.com/craiggwilson/go-streamdeck-sdk"
	"github.com/craiggwilson/go-streamdeck-sdk/streamdeckevent"
)

const actionUUID = "com.craiggwilson.streamdeck.example.counter"

func New() *streamdeck.InstancedAction {
	return streamdeck.NewInstancedAction(
		actionUUID,
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

func (a *ActionInstance) ActionUUID() streamdeck.ActionUUID {
	return actionUUID
}

func (a *ActionInstance) EventContext() streamdeck.EventContext {
	return a.eventContext
}

func (a *ActionInstance) HandleKeyDown(_ context.Context, _ streamdeckevent.KeyDown) error {
	a.count++
	return a.publisher.SetTitle(streamdeckevent.SetTitlePayload{
		Title: strconv.Itoa(a.count),
		Target: streamdeckevent.HardwareAndSoftware,
	})
}
