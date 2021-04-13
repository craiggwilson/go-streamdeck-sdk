package streamdeck

import (
	"github.com/craiggwilson/go-streamdeck-sdk/streamdeckcore"
)

// ActionUUID is an action's unique identifier. It is an alias for streamdeckcore.ActionUUID.
type ActionUUID = streamdeckcore.ActionUUID

// DeviceUUID is a device's unique identifier. It is an alias for streamdeckcore.DeviceUUID.
type DeviceUUID = streamdeckcore.DeviceUUID

// EventContext is the unique identifier for an instance of an action. It is an alias for streamdeckcore.EventContext.
type EventContext = streamdeckcore.EventContext

// EventName is the name of an event. It is an alias for streamdeckcore.EventName.
type EventName = streamdeckcore.EventName

// PluginUUID is the unique identifier assigned to a plugin by a device. It is an alias for streamdeckcore.PluginUUID.
type PluginUUID = streamdeckcore.PluginUUID
