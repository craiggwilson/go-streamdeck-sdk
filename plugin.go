package streamdeck

import (
	"context"
	"encoding/json"
	"fmt"
)

// NewPlugin makes a Plugin.
func NewPlugin(actions ...Action) *Plugin {
	actionMap := make(map[ActionUUID]Action, len(actions))
	for _, action := range actions {
		actionMap[action.ActionUUID()] = action
	}

	return &Plugin{
		actions: actionMap,
	}
}

// Plugin is the default implementation of a streamdeckcore.Plugin. It handles the raw events
// and dispatches them to the appropriate actions.
type Plugin struct {
	actions map[ActionUUID]Action
}

// Initialize implements the streamdeckcore.Plugin interface.
func (p *Plugin) Initialize(pluginUUID PluginUUID, publisher Publisher) {
	for _, action := range p.actions {
		ap := newCoreActionPublisher(pluginUUID, action.ActionUUID(), publisher)
		action.InitializeAction(pluginUUID, ap)
	}
}

// HandleEvent implements the streamdeckcore.Handler interface.
func (p *Plugin) HandleEvent(ctx context.Context, raw json.RawMessage) error {
	var eventHeader struct {
		Event EventName `json:"event"`
		Action ActionUUID `json:"action"`
	}
	if err := json.Unmarshal(raw, &eventHeader); err != nil {
		return fmt.Errorf("unmarshalling event header: %w", err)
	}

	// If the action is empty, the event is intended for all actions.
	if eventHeader.Action == "" {
		for _, action := range p.actions {
			if err := dispatchEvent(ctx, action, eventHeader.Event, raw); err != nil {
				return fmt.Errorf("dispatching event %q to action %q: %w", eventHeader.Event, eventHeader.Action, err)
			}
		}

		return nil
	}

	// If the action doesn't exist, it wasn't registered and this plugin should not have received this event.
	action, ok := p.actions[eventHeader.Action]
	if !ok {
		return fmt.Errorf("unknown action %q", eventHeader.Action)
	}

	if err := dispatchEvent(ctx, action, eventHeader.Event, raw); err != nil {
		return fmt.Errorf("dispatching event %q to action %q: %w", eventHeader.Event, eventHeader.Action, err)
	}

	return nil
}
