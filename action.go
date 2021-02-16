package streamdeck

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/craiggwilson/go-streamdeck-sdk/streamdeckcore"
)

// Action represents a discrete action and acts as a action to createInstance instances.
type Action interface {
	UUID() streamdeckcore.ActionUUID
	Initialize(pluginUUID streamdeckcore.PluginUUID, publisher streamdeckcore.Publisher)
	HandleEvent(ctx context.Context, raw json.RawMessage) error
}

// NewDefaultAction makes an implementation of a DefaultAction.
func NewDefaultAction(actionUUID streamdeckcore.ActionUUID, createInstance func(eventContext streamdeckcore.EventContext, publisher ActionInstancePublisher) ActionInstance) *DefaultAction {
	return &DefaultAction{
		actionUUID:     actionUUID,
		createInstance: createInstance,
		instances:      make(map[streamdeckcore.EventContext]ActionInstance),
	}
}

// DefaultAction implements the Action interface and delegates to a func for action instance creation, passing along
// events to the appropriate instances.
type DefaultAction struct {
	actionUUID     streamdeckcore.ActionUUID
	createInstance func(eventContext streamdeckcore.EventContext, publisher ActionInstancePublisher) ActionInstance
	instances      map[streamdeckcore.EventContext]ActionInstance

	pluginUUID streamdeckcore.PluginUUID
	publisher streamdeckcore.Publisher
}

// UUID implements the Action interface.
func (a *DefaultAction) UUID() streamdeckcore.ActionUUID {
	return a.actionUUID
}

// Initialize implements the Action interface.
func (a *DefaultAction) Initialize(pluginUUID streamdeckcore.PluginUUID, publisher streamdeckcore.Publisher) {
	a.pluginUUID = pluginUUID
	a.publisher = publisher
}

// HandleEvent implements the streamdeckcore.Handler interface.
func (a *DefaultAction) HandleEvent(ctx context.Context, raw json.RawMessage) error {
	var eventHeader struct {
		Event streamdeckcore.EventName `json:"event"`
		Action streamdeckcore.ActionUUID `json:"action"`
		Context streamdeckcore.EventContext `json:"context"`
	}
	if err := json.Unmarshal(raw, &eventHeader); err != nil {
		return fmt.Errorf("unmarshalling action and eventHeader: %w", err)
	}

	if eventHeader.Action != "" && eventHeader.Action != a.actionUUID {
		return nil
	}

	if eventHeader.Context == "" {
		for _, instance := range a.instances {
			if err := dispatchEvent(ctx, instance, eventHeader.Event, raw); err != nil {
				return fmt.Errorf("dispatching event %q to action instance %q: %w", eventHeader.Event, eventHeader.Context, err)
			}
		}

		return nil
	}

	instance, ok := a.instances[eventHeader.Context]
	if !ok {
		publisher := newCoreActionInstancePublisher(a.pluginUUID, a.actionUUID, eventHeader.Context, a.publisher)
		instance = a.createInstance(eventHeader.Context, publisher)
		a.instances[eventHeader.Context] = instance
	}

	if err := dispatchEvent(ctx, instance, eventHeader.Event, raw); err != nil {
		return fmt.Errorf("dispatching event %q to action instance %q: %w", eventHeader.Event, eventHeader.Context, err)
	}

	return nil
}

// ActionInstance represents an instance of an action. It should also implement some of the event handlers
// in order to receive important events from the Streamdeck.
type ActionInstance interface {
	UUID() streamdeckcore.ActionUUID
	Context() streamdeckcore.EventContext
}
