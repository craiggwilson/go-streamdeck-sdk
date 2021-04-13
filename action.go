package streamdeck

import (
	"context"
	"encoding/json"
	"fmt"
)

// Action represents a discrete action and acts as a action to createInstance instances.
type Action interface {
	ActionUUID() ActionUUID
	InitializeAction(pluginUUID PluginUUID, publisher ActionPublisher)
}

// NewInstancedAction makes an implementation of a InstancedAction.
func NewInstancedAction(actionUUID ActionUUID, createInstance ActionInstanceFactory) *InstancedAction {
	return &InstancedAction{
		actionUUID:     actionUUID,
		createInstance: createInstance,
		instances:      make(map[EventContext]ActionInstance),
	}
}

// InstancedAction implements the Action interface and delegates to a func for action instance creation, passing along
// events to the appropriate instances.
type InstancedAction struct {
	actionUUID     ActionUUID
	createInstance ActionInstanceFactory
	instances      map[EventContext]ActionInstance

	pluginUUID PluginUUID
	publisher ActionPublisher
}

// ActionUUID implements the Action interface.
func (a *InstancedAction) ActionUUID() ActionUUID {
	return a.actionUUID
}

// InitializeAction implements the Action interface.
func (a *InstancedAction) InitializeAction(pluginUUID PluginUUID, publisher ActionPublisher) {
	a.pluginUUID = pluginUUID
	a.publisher = publisher
}

// HandleEvent implements the streamdeckcore.Handler interface.
func (a *InstancedAction) HandleEvent(ctx context.Context, raw json.RawMessage) error {
	var eventHeader struct {
		Event EventName `json:"event"`
		Action ActionUUID `json:"action"`
		Context EventContext `json:"context"`
	}
	if err := json.Unmarshal(raw, &eventHeader); err != nil {
		return fmt.Errorf("unmarshalling action and eventHeader: %w", err)
	}

	// If the action doesn't match, it means this event was sent here improperly.
	if eventHeader.Action != "" && eventHeader.Action != a.actionUUID {
		return fmt.Errorf("received mismatched action, %s != %s", a.actionUUID, eventHeader.Action)
	}

	// If the context is empty, the event is intended for all instances of this action.
	if eventHeader.Context == "" {
		for _, instance := range a.instances {
			if err := dispatchEvent(ctx, instance, eventHeader.Event, raw); err != nil {
				return fmt.Errorf("dispatching event %q to action instance %q: %w", eventHeader.Event, eventHeader.Context, err)
			}
		}

		return nil
	}

	// If the instance doesn't yet exist, create one and save it off.
	instance, ok := a.instances[eventHeader.Context]
	if !ok {
		publisher := newCoreActionInstancePublisher(eventHeader.Context, a.publisher)
		instance = a.createInstance(eventHeader.Context, publisher)
		a.instances[eventHeader.Context] = instance
	}

	if err := dispatchEvent(ctx, instance, eventHeader.Event, raw); err != nil {
		return fmt.Errorf("dispatching event %q to action instance %q: %w", eventHeader.Event, eventHeader.Context, err)
	}

	return nil
}
