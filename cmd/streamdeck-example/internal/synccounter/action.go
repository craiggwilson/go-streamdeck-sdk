package synccounter

import (
	"context"
	"strconv"

	"github.com/craiggwilson/go-streamdeck-sdk"
	"github.com/craiggwilson/go-streamdeck-sdk/streamdeckevent"
)

const actionUUID = "com.craiggwilson.streamdeck.example.synccounter"

func New() *streamdeck.InstancedAction {
	var instances []*ActionInstance
	var count int

	increment := func() {
		count++
		for _, instance := range instances {
			instance.display(count)
		}
	}

	return streamdeck.NewInstancedAction(
		actionUUID,
		func(eventContext streamdeck.EventContext, publisher streamdeck.ActionInstancePublisher) streamdeck.ActionInstance {
			instance := &ActionInstance{
				eventContext: eventContext,
				publisher:    publisher,
				inc: increment,
			}

			instances = append(instances, instance)
			return instance
		},
	)
}

type ActionInstance struct {
	eventContext streamdeck.EventContext
	publisher    streamdeck.ActionInstancePublisher
	inc func()
}

func (a *ActionInstance) ActionUUID() streamdeck.ActionUUID {
	return actionUUID
}

func (a *ActionInstance) EventContext() streamdeck.EventContext {
	return a.eventContext
}

func (a *ActionInstance) HandleKeyDown(_ context.Context, _ streamdeckevent.KeyDown) error {
	a.inc()
	return nil
}

func (a *ActionInstance) display(count int) {
	_ = a.publisher.SetTitle(streamdeckevent.SetTitlePayload{
		Title:  strconv.Itoa(count),
		Target: streamdeckevent.HardwareAndSoftware,
	})
}
