package streamdeck

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/craiggwilson/streamdeck-plugins/pkg/streamdeck/streamdeckcore"
)

// NewPlugin makes a Plugin.
func NewPlugin(actionFactories ...ActionInstanceFactory) *Plugin {
	m := make(map[streamdeckcore.ActionUUID]actionHolder, len(actionFactories))
	for _, actionFactory := range actionFactories {
		m[actionFactory.UUID()] = actionHolder{
			factory: actionFactory,
			instances: make(map[streamdeckcore.EventContext]ActionInstance),
		}
	}

	return &Plugin{
		actions: m,
	}
}

// Plugin is the core implementation of a streamdeckcore.Plugin. It handles the raw events
// and dispatches them to the appropriate actions.
type Plugin struct {
	pluginUUID streamdeckcore.PluginUUID
	corePublisher streamdeckcore.Publisher

	actions map[streamdeckcore.ActionUUID]actionHolder
}

func (p *Plugin) Initialize(pluginUUID streamdeckcore.PluginUUID, publisher streamdeckcore.Publisher) {
	p.pluginUUID = pluginUUID
	p.corePublisher = publisher
}

func (p *Plugin) HandleEvent(ctx context.Context, raw json.RawMessage) error {
	var eventHeader struct {
		Event streamdeckcore.EventName `json:"event"`
		Action streamdeckcore.ActionUUID `json:"action"`
		Context streamdeckcore.EventContext `json:"context"`
	}
	if err := json.Unmarshal(raw, &eventHeader); err != nil {
		return fmt.Errorf("unmarshalling action and eventHeader: %w", err)
	}

	action, ok := p.actions[eventHeader.Action]
	if !ok {
		return fmt.Errorf("unknown action %q", eventHeader.Action)
	}

	// If there is no context, then dispatch the event to all action. This is true, for example, with the global
	// settings events.
	if eventHeader.Context == "" {
		for _, action := range action.instances {
			if err := dispatchEvent(ctx, action, eventHeader.Event, raw); err != nil {
				return fmt.Errorf("dispatching to action %q, instance %q: %w", eventHeader.Action, eventHeader.Context, err)
			}
		}

		return nil
	}

	instance, ok := action.instances[eventHeader.Context]
	if !ok {
		publisher := newCorePublisherAdapter(p.pluginUUID, eventHeader.Action, eventHeader.Context, p.corePublisher)
		instance = action.factory.CreateActionInstance(eventHeader.Context, publisher)
		action.instances[eventHeader.Context] = instance
	}

	return dispatchEvent(ctx, instance, eventHeader.Event, raw)
}

type actionHolder struct {
	factory ActionInstanceFactory
	instances map[streamdeckcore.EventContext]ActionInstance
}
