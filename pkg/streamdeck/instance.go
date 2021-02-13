package streamdeck

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/craiggwilson/streamdeck-plugins/pkg/streamdeck/streamdeckcore"
)

// ActionInstanceFactory creates instances of actions.
type ActionInstanceFactory interface {
	CreateActionInstance(eventContext streamdeckcore.EventContext) Action
}

// ActionInstanceFactoryFunc is a functional implementation of ActionInstanceFactory.
type ActionInstanceFactoryFunc func(eventContext streamdeckcore.EventContext) Action

// CreateActionInstance implements the ActionInstanceFactory interface.
func (f ActionInstanceFactoryFunc) CreateActionInstance(eventContext streamdeckcore.EventContext) Action {
	return f(eventContext)
}

// NewInstancedAction makes an InstancedAction.
func NewInstancedAction(factory ActionInstanceFactory) *InstancedAction {
	return &InstancedAction{
		instanceFactory: factory,
		instances: make(map[streamdeckcore.EventContext]Action),
	}
}

// InstancedAction delegates received events to the correct instance of an Action.
type InstancedAction struct {
	instanceFactory ActionInstanceFactory
	instances map[streamdeckcore.EventContext]Action

	eventPublisher streamdeckcore.EventPublisher
}

// UUID implements the Action interface.
func (mux *InstancedAction) UUID() streamdeckcore.ActionUUID {
	return "com.craiggwilson.streamdeck.counter"
}

// HandleEvent implements the streamdeckcore.EventHandler interface.
func (mux *InstancedAction) Initialize(eventPublisher streamdeckcore.EventPublisher) {
	mux.eventPublisher = eventPublisher
}

// HandleEvent implements the streamdeckcore.EventHandler interface.
func (mux *InstancedAction) HandleEvent(ctx context.Context, raw json.RawMessage) error {
	var eventAndContext struct {
		Event streamdeckcore.EventName `json:"event"`
		Context streamdeckcore.EventContext `json:"context"`
	}
	if err := json.Unmarshal(raw, &eventAndContext); err != nil {
		return fmt.Errorf("unmarshalling action and event: %w", err)
	}

	if eventAndContext.Context == "" {
		for _, a := range mux.instances {
			if err := dispatchEvent(ctx, a, eventAndContext.Event, raw); err != nil {
				return fmt.Errorf("dispatching to instance %q: %w", eventAndContext.Context, err)
			}
		}

		return nil
	}

	a, ok := mux.instances[eventAndContext.Context]
	if !ok {
		a = mux.instanceFactory.CreateActionInstance(eventAndContext.Context)
		a.Initialize(mux.eventPublisher)
		mux.instances[eventAndContext.Context] = a
	}

	err := dispatchEvent(ctx, a, eventAndContext.Event, raw)
	if err != nil {
		return fmt.Errorf("dispatching to instance %q: %w", eventAndContext.Context, err)
	}

	return nil
}
