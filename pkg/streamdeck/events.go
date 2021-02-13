package streamdeck

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/craiggwilson/streamdeck-plugins/pkg/streamdeck/streamdeckcore"
)

// DidReceiveSettingsHandler is implemented by Actions that wish to receive the streamdeckcore.DidReceiveSettingsEvent.
type DidReceiveSettingsHandler interface {
	HandleDidReceiveSettings(ctx context.Context, event streamdeckcore.DidReceiveSettingsEvent) error
}

// DidReceiveGlobalSettingsHandler is implemented by Actions that wish to receive the streamdeckcore.DidReceiveGlobalSettingsEvent.
type DidReceiveGlobalSettingsHandler interface {
	HandleDidReceiveGlobalSettings(ctx context.Context, event streamdeckcore.DidReceiveGlobalSettingsEvent) error
}

// KeyDownHandler is implemented by Actions that wish to receive the streamdeckcore.KeyDownEvent.
type KeyDownHandler interface {
	HandleKeyDown(ctx context.Context, event streamdeckcore.KeyDownEvent) error
}

// KeyUpHandler is implemented by Actions that wish to receive the streamdeckcore.KeyUpEvent.
type KeyUpHandler interface {
	HandleKeyUp(ctx context.Context, event streamdeckcore.KeyUpEvent) error
}

// WillAppearHandler is implemented by Actions that wish to receive the streamdeckcore.WillAppearEvent.
type WillAppearHandler interface {
	HandleWillAppear(ctx context.Context, event streamdeckcore.WillAppearEvent) error
}

// WillDisappearHandler is implemented by Actions that wish to receive the streamdeckcore.WillDisappearEvent.
type WillDisappearHandler interface {
	HandleWillDisappear(ctx context.Context, event streamdeckcore.WillDisappearEvent) error
}

func dispatchEvent(ctx context.Context, action ActionInstance, eventName streamdeckcore.EventName, raw json.RawMessage) error {
	switch eventName {
	case streamdeckcore.DidReceiveSettings:
		if h, ok := action.(DidReceiveSettingsHandler); ok {
			var event streamdeckcore.DidReceiveSettingsEvent
			if err := json.Unmarshal(raw, &event); err != nil {
				return fmt.Errorf("unmarshalling %s: %w", streamdeckcore.DidReceiveSettings, err)
			}
			return h.HandleDidReceiveSettings(ctx, event)
		}
	case streamdeckcore.DidReceiveGlobalSettings:
		if h, ok := action.(DidReceiveGlobalSettingsHandler); ok {
			var event streamdeckcore.DidReceiveGlobalSettingsEvent
			if err := json.Unmarshal(raw, &event); err != nil {
				return fmt.Errorf("unmarshalling %s: %w", streamdeckcore.DidReceiveGlobalSettings, err)
			}
			return h.HandleDidReceiveGlobalSettings(ctx, event)
		}
	case streamdeckcore.KeyDown:
		if h, ok := action.(KeyDownHandler); ok {
			var event streamdeckcore.KeyDownEvent
			if err := json.Unmarshal(raw, &event); err != nil {
				return fmt.Errorf("unmarshalling %s: %w", streamdeckcore.KeyDown, err)
			}
			return h.HandleKeyDown(ctx, event)
		}
	case streamdeckcore.KeyUp:
		if h, ok := action.(KeyUpHandler); ok {
			var event streamdeckcore.KeyUpEvent
			if err := json.Unmarshal(raw, &event); err != nil {
				return fmt.Errorf("unmarshalling %s: %w", streamdeckcore.KeyUp, err)
			}
			return h.HandleKeyUp(ctx, event)
		}
	case streamdeckcore.WillAppear:
		if h, ok := action.(WillAppearHandler); ok {
			var event streamdeckcore.WillAppearEvent
			if err := json.Unmarshal(raw, &event); err != nil {
				return fmt.Errorf("unmarshalling %s: %w", streamdeckcore.WillAppear, err)
			}
			return h.HandleWillAppear(ctx, event)
		}
	case streamdeckcore.WillDisappear:
		if h, ok := action.(WillDisappearHandler); ok {
			var event streamdeckcore.WillDisappearEvent
			if err := json.Unmarshal(raw, &event); err != nil {
				return fmt.Errorf("unmarshalling %s: %w", streamdeckcore.WillDisappear, err)
			}
			return h.HandleWillDisappear(ctx, event)
		}
	}

	if h, ok := action.(streamdeckcore.Plugin); ok {
		return h.HandleEvent(ctx, raw)
	}

	return nil
}