package streamdeck

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/craiggwilson/go-streamdeck-sdk/streamdeckcore"
)

// NewDefaultPlugin makes a DefaultPlugin.
func NewDefaultPlugin(actions ...Action) *DefaultPlugin {
	actionMap := make(map[streamdeckcore.ActionUUID]Action, len(actions))
	for _, action := range actions {
		actionMap[action.UUID()] = action
	}

	return &DefaultPlugin{
		actions: actionMap,
	}
}

// DefaultPlugin is the default implementation of a streamdeckcore.Plugin. It handles the raw events
// and dispatches them to the appropriate actions.
type DefaultPlugin struct {
	actions map[streamdeckcore.ActionUUID]Action
}

// Initialize implements the streamdeckcore.Plugin interface.
func (p *DefaultPlugin) Initialize(pluginUUID streamdeckcore.PluginUUID, publisher streamdeckcore.Publisher) {
	for _, action := range p.actions {
		action.Initialize(pluginUUID, publisher)
	}
}

// Initialize implements the streamdeckcore.Handler interface.
func (p *DefaultPlugin) HandleEvent(ctx context.Context, raw json.RawMessage) error {
	var eventHeader struct {
		Event streamdeckcore.EventName `json:"event"`
		Action streamdeckcore.ActionUUID `json:"action"`
	}
	if err := json.Unmarshal(raw, &eventHeader); err != nil {
		return fmt.Errorf("unmarshalling action and eventHeader: %w", err)
	}

	if eventHeader.Action == "" {
		for _, action := range p.actions {
			if err := dispatchEvent(ctx, action, eventHeader.Event, raw); err != nil {
				return fmt.Errorf("dispatching event %q to action %q: %w", eventHeader.Event, eventHeader.Action, err)
			}
		}

		return nil
	}

	action, ok := p.actions[eventHeader.Action]
	if !ok {
		return fmt.Errorf("unknown action %q", eventHeader.Action)
	}

	if err := dispatchEvent(ctx, action, eventHeader.Event, raw); err != nil {
		return fmt.Errorf("dispatching event %q to action %q: %w", eventHeader.Event, eventHeader.Action, err)
	}

	return nil
}
