package streamdeck

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/craiggwilson/go-streamdeck-sdk/streamdeckcore"
	"github.com/craiggwilson/go-streamdeck-sdk/streamdeckevent"
)

// ApplicationDidLaunchHandler is implemented by ActionInstances that wish to receive the streamdeckevent.ApplicationDidLaunch event.
type ApplicationDidLaunchHandler interface {
	HandleApplicationDidLaunch(ctx context.Context, event streamdeckevent.ApplicationDidLaunch) error
}

// ApplicationDidTerminateHandler is implemented by ActionInstances that wish to receive the streamdeckevent.ApplicationDidTerminate event.
type ApplicationDidTerminateHandler interface {
	HandleApplicationDidTerminate(ctx context.Context, event streamdeckevent.ApplicationDidTerminate) error
}

// DeviceDidConnectHandler is implemented by ActionInstances that wish to receive the streamdeckevent.DeviceDidConnect event.
type DeviceDidConnectHandler interface {
	HandleDeviceDidConnect(ctx context.Context, event streamdeckevent.DeviceDidConnect) error
}

// DeviceDidDisconnectHandler is implemented by ActionInstances that wish to receive the streamdeckevent.DeviceDidDisconnect event.
type DeviceDidDisconnectHandler interface {
	HandleDeviceDidDisconnect(ctx context.Context, event streamdeckevent.DeviceDidDisconnect) error
}

// DidReceiveSettingsHandler is implemented by ActionInstances that wish to receive the streamdeckevent.DidReceiveSettings event.
type DidReceiveSettingsHandler interface {
	HandleDidReceiveSettings(ctx context.Context, event streamdeckevent.DidReceiveSettings) error
}

// DidReceiveGlobalSettingsHandler is implemented by ActionInstances that wish to receive the streamdeckevent.DidReceiveGlobalSettings event.
type DidReceiveGlobalSettingsHandler interface {
	HandleDidReceiveGlobalSettings(ctx context.Context, event streamdeckevent.DidReceiveGlobalSettings) error
}

// KeyDownHandler is implemented by ActionInstances that wish to receive the streamdeckevent.KeyDown event.
type KeyDownHandler interface {
	HandleKeyDown(ctx context.Context, event streamdeckevent.KeyDown) error
}

// KeyUpHandler is implemented by ActionInstances that wish to receive the streamdeckevent.KeyUp event.
type KeyUpHandler interface {
	HandleKeyUp(ctx context.Context, event streamdeckevent.KeyUp) error
}

// PropertyInspectorDidAppearHandler is implemented by ActionInstances that wish to receive the streamdeckevent.PropertyInspectorDidAppear event.
type PropertyInspectorDidAppearHandler interface {
	HandlePropertyInspectorDidAppear(ctx context.Context, event streamdeckevent.PropertyInspectorDidAppear) error
}

// PropertyInspectorDidDisappearHandler is implemented by ActionInstances that wish to receive the streamdeckevent.PropertyInspectorDidDisappear event.
type PropertyInspectorDidDisappearHandler interface {
	HandlePropertyInspectorDidDisappear(ctx context.Context, event streamdeckevent.PropertyInspectorDidDisappear) error
}

// SendToPluginHandler is implemented by ActionInstances that wish to receive the streamdeckevent.SendToPlugin event.
type SendToPluginHandler interface {
	HandleSendToPlugin(ctx context.Context, event streamdeckevent.SendToPlugin) error
}

// SystemDidWakeUpHandler is implemented by ActionInstances that wish to receive the streamdeckevent.SystemDidWakeUp event.
type SystemDidWakeUpHandler interface {
	HandleSystemDidWakeUp(ctx context.Context, event streamdeckevent.SystemDidWakeUp) error
}

// TitleParametersDidChangeHandler is implemented by ActionInstances that wish to receive the streamdeckevent.TitleParametersDidChange event.
type TitleParametersDidChangeHandler interface {
	HandleTitleParametersDidChange(ctx context.Context, event streamdeckevent.TitleParametersDidChange) error
}

// WillAppearHandler is implemented by ActionInstances that wish to receive the streamdeckevent.WillAppear event.
type WillAppearHandler interface {
	HandleWillAppear(ctx context.Context, event streamdeckevent.WillAppear) error
}

// WillDisappearHandler is implemented by ActionInstances that wish to receive the streamdeckevent.WillDisappear event.
type WillDisappearHandler interface {
	HandleWillDisappear(ctx context.Context, event streamdeckevent.WillDisappear) error
}

func dispatchEvent(ctx context.Context, target interface{}, eventName streamdeckcore.EventName, raw json.RawMessage) error {
	switch eventName {
	case streamdeckevent.ApplicationDidLaunchName:
		if h, ok := target.(ApplicationDidLaunchHandler); ok {
			var event streamdeckevent.ApplicationDidLaunch
			if err := json.Unmarshal(raw, &event); err != nil {
				return fmt.Errorf("unmarshalling %s: %w", streamdeckevent.ApplicationDidLaunchName, err)
			}
			return h.HandleApplicationDidLaunch(ctx, event)
		}
	case streamdeckevent.ApplicationDidTerminateName:
		if h, ok := target.(ApplicationDidTerminateHandler); ok {
			var event streamdeckevent.ApplicationDidTerminate
			if err := json.Unmarshal(raw, &event); err != nil {
				return fmt.Errorf("unmarshalling %s: %w", streamdeckevent.ApplicationDidTerminateName, err)
			}
			return h.HandleApplicationDidTerminate(ctx, event)
		}
	case streamdeckevent.DeviceDidConnectName:
		if h, ok := target.(DeviceDidConnectHandler); ok {
			var event streamdeckevent.DeviceDidConnect
			if err := json.Unmarshal(raw, &event); err != nil {
				return fmt.Errorf("unmarshalling %s: %w", streamdeckevent.DeviceDidConnectName, err)
			}
			return h.HandleDeviceDidConnect(ctx, event)
		}
	case streamdeckevent.DeviceDidDisconnectName:
		if h, ok := target.(DeviceDidDisconnectHandler); ok {
			var event streamdeckevent.DeviceDidDisconnect
			if err := json.Unmarshal(raw, &event); err != nil {
				return fmt.Errorf("unmarshalling %s: %w", streamdeckevent.DeviceDidDisconnectName, err)
			}
			return h.HandleDeviceDidDisconnect(ctx, event)
		}
	case streamdeckevent.DidReceiveSettingsName:
		if h, ok := target.(DidReceiveSettingsHandler); ok {
			var event streamdeckevent.DidReceiveSettings
			if err := json.Unmarshal(raw, &event); err != nil {
				return fmt.Errorf("unmarshalling %s: %w", streamdeckevent.DidReceiveSettingsName, err)
			}
			return h.HandleDidReceiveSettings(ctx, event)
		}
	case streamdeckevent.DidReceiveGlobalSettingsName:
		if h, ok := target.(DidReceiveGlobalSettingsHandler); ok {
			var event streamdeckevent.DidReceiveGlobalSettings
			if err := json.Unmarshal(raw, &event); err != nil {
				return fmt.Errorf("unmarshalling %s: %w", streamdeckevent.DidReceiveGlobalSettingsName, err)
			}
			return h.HandleDidReceiveGlobalSettings(ctx, event)
		}
	case streamdeckevent.KeyDownName:
		if h, ok := target.(KeyDownHandler); ok {
			var event streamdeckevent.KeyDown
			if err := json.Unmarshal(raw, &event); err != nil {
				return fmt.Errorf("unmarshalling %s: %w", streamdeckevent.KeyDownName, err)
			}
			return h.HandleKeyDown(ctx, event)
		}
	case streamdeckevent.KeyUpName:
		if h, ok := target.(KeyUpHandler); ok {
			var event streamdeckevent.KeyUp
			if err := json.Unmarshal(raw, &event); err != nil {
				return fmt.Errorf("unmarshalling %s: %w", streamdeckevent.KeyUpName, err)
			}
			return h.HandleKeyUp(ctx, event)
		}
	case streamdeckevent.PropertyInspectorDidAppearName:
		if h, ok := target.(PropertyInspectorDidAppearHandler); ok {
			var event streamdeckevent.PropertyInspectorDidAppear
			if err := json.Unmarshal(raw, &event); err != nil {
				return fmt.Errorf("unmarshalling %s: %w", streamdeckevent.PropertyInspectorDidAppearName, err)
			}
			return h.HandlePropertyInspectorDidAppear(ctx, event)
		}
	case streamdeckevent.PropertyInspectorDidDisappearName:
		if h, ok := target.(PropertyInspectorDidDisappearHandler); ok {
			var event streamdeckevent.PropertyInspectorDidDisappear
			if err := json.Unmarshal(raw, &event); err != nil {
				return fmt.Errorf("unmarshalling %s: %w", streamdeckevent.PropertyInspectorDidDisappearName, err)
			}
			return h.HandlePropertyInspectorDidDisappear(ctx, event)
		}
	case streamdeckevent.SendToPluginName:
		if h, ok := target.(SendToPluginHandler); ok {
			var event streamdeckevent.SendToPlugin
			if err := json.Unmarshal(raw, &event); err != nil {
				return fmt.Errorf("unmarshalling %s: %w", streamdeckevent.SendToPluginName, err)
			}
			return h.HandleSendToPlugin(ctx, event)
		}
	case streamdeckevent.SystemDidWakeUpName:
		if h, ok := target.(SystemDidWakeUpHandler); ok {
			var event streamdeckevent.SystemDidWakeUp
			if err := json.Unmarshal(raw, &event); err != nil {
				return fmt.Errorf("unmarshalling %s: %w", streamdeckevent.SystemDidWakeUpName, err)
			}
			return h.HandleSystemDidWakeUp(ctx, event)
		}
	case streamdeckevent.TitleParametersDidChangeName:
		if h, ok := target.(TitleParametersDidChangeHandler); ok {
			var event streamdeckevent.TitleParametersDidChange
			if err := json.Unmarshal(raw, &event); err != nil {
				return fmt.Errorf("unmarshalling %s: %w", streamdeckevent.TitleParametersDidChangeName, err)
			}
			return h.HandleTitleParametersDidChange(ctx, event)
		}
	case streamdeckevent.WillAppearName:
		if h, ok := target.(WillAppearHandler); ok {
			var event streamdeckevent.WillAppear
			if err := json.Unmarshal(raw, &event); err != nil {
				return fmt.Errorf("unmarshalling %s: %w", streamdeckevent.WillAppearName, err)
			}
			return h.HandleWillAppear(ctx, event)
		}
	case streamdeckevent.WillDisappearName:
		if h, ok := target.(WillDisappearHandler); ok {
			var event streamdeckevent.WillDisappear
			if err := json.Unmarshal(raw, &event); err != nil {
				return fmt.Errorf("unmarshalling %s: %w", streamdeckevent.WillDisappearName, err)
			}
			return h.HandleWillDisappear(ctx, event)
		}
	}

	if h, ok := target.(streamdeckcore.Handler); ok {
		return h.HandleEvent(ctx, raw)
	}

	return nil
}
