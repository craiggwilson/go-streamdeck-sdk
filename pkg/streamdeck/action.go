package streamdeck

import (
	"github.com/craiggwilson/streamdeck-plugins/pkg/streamdeck/streamdeckcore"
)

// ActionInstanceFactory creates instances of actions.
type ActionInstanceFactory interface {
	UUID() streamdeckcore.ActionUUID
	CreateActionInstance(eventContext streamdeckcore.EventContext, publisher Publisher) ActionInstance
}

// NewActionInstanceFactory makes an implementation of ActionInstanceFactory.
func NewActionInstanceFactory(uuid streamdeckcore.ActionUUID, create func(eventContext streamdeckcore.EventContext, publisher Publisher) ActionInstance) ActionInstanceFactory {
	return &simpleActionInstanceFactory{
		uuid: uuid,
		create: create,
	}
}

type simpleActionInstanceFactory struct {
	uuid streamdeckcore.ActionUUID
	create func(eventContext streamdeckcore.EventContext, publisher Publisher) ActionInstance
}

func (f *simpleActionInstanceFactory) UUID() streamdeckcore.ActionUUID {
	return f.uuid
}

func (f *simpleActionInstanceFactory) CreateActionInstance(eventContext streamdeckcore.EventContext, publisher Publisher) ActionInstance {
	return f.create(eventContext, publisher)
}

// ActionInstance represents an instance of an action. It should also implement some of the event handlers
// in order to receive important events from the Streamdeck.
type ActionInstance interface {
	UUID() streamdeckcore.ActionUUID
	Context() streamdeckcore.EventContext
}
