package streamdeck

// ActionInstanceFactory creates instances of an action.
type ActionInstanceFactory func(eventContext EventContext, publisher ActionInstancePublisher) ActionInstance

// ActionInstance represents an instance of an action. It should also implement some of the event handlers
// in order to receive relevant events from a device.
type ActionInstance interface {
	ActionUUID() ActionUUID
	EventContext() EventContext
}
