package synccounter

import (
	"context"
	"strconv"
	"sync"

	"github.com/craiggwilson/go-streamdeck-sdk/streamdeckevent"

	"github.com/craiggwilson/go-streamdeck-sdk"
)

const uuid = "com.craiggwilson.streamdeck.example.synccounter"

func New() *streamdeck.DefaultAction {
	var l sync.Mutex
	var instances []*ActionInstance
	var count int

	increment := func() {
		l.Lock()
		defer l.Unlock()
		count++
		for _, instance := range instances {
			instance.display(count)
		}
	}

	return streamdeck.NewDefaultAction(
		uuid,
		func(eventContext streamdeck.EventContext, publisher streamdeck.ActionInstancePublisher) streamdeck.ActionInstance {
			instance := &ActionInstance{
				eventContext: eventContext,
				publisher:    publisher,
				inc: increment,
			}

			l.Lock()
			defer l.Unlock()
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

func (a *ActionInstance) UUID() streamdeck.ActionUUID {
	return uuid
}

func (a *ActionInstance) Context() streamdeck.EventContext {
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
