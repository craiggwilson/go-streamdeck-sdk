package streamdeck

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/craiggwilson/streamdeck-plugins/pkg/streamdeck/streamdeckcore"
)

// NewPlugin makes a Plugin.
func NewPlugin(actions ...Action) *Plugin {
	actionMap := make(map[streamdeckcore.ActionUUID]Action, len(actions))
	for _, action := range actions {
		actionMap[action.UUID()] = action
	}

	return &Plugin{
		actions: actionMap,
	}
}

// Plugin is the core implementation of a streamdeckcore.EventHandler. It handles the raw events
// and dispatches them to the appropriate actions.
type Plugin struct {
	actions map[streamdeckcore.ActionUUID]Action
}

func (mux *Plugin) Initialize(eventPublisher streamdeckcore.EventPublisher) {
	for _, action := range mux.actions {
		action.Initialize(eventPublisher)
	}
}

func (mux *Plugin) HandleEvent(ctx context.Context, raw json.RawMessage) error {
	var actionAndEvent struct {
		Action streamdeckcore.ActionUUID `json:"action"`
		Event streamdeckcore.EventName `json:"event"`
	}
	if err := json.Unmarshal(raw, &actionAndEvent); err != nil {
		return fmt.Errorf("unmarshalling action and event: %w", err)
	}

	a, ok := mux.actions[actionAndEvent.Action]
	if !ok {
		return fmt.Errorf("unknown action %q", a.UUID())
	}

	return dispatchEvent(ctx, a, actionAndEvent.Event, raw)
}
