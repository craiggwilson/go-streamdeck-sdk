package streamdeckevent

import (
	"encoding/json"

	"github.com/craiggwilson/go-streamdeck-sdk/streamdeckcore"
)

const (
	ApplicationDidLaunchName          streamdeckcore.EventName = "applicationDidLaunch"
	ApplicationDidTerminateName       streamdeckcore.EventName = "applicationDidTerminate"
	DeviceDidConnectName              streamdeckcore.EventName = "deviceDidConnect"
	DeviceDidDisconnectName           streamdeckcore.EventName = "deviceDidDisconnect"
	DidReceiveGlobalSettingsName      streamdeckcore.EventName = "didReceiveGlobalSettings"
	DidReceiveSettingsName            streamdeckcore.EventName = "didReceiveSettings"
	KeyDownName                       streamdeckcore.EventName = "keyDown"
	KeyUpName                         streamdeckcore.EventName = "keyUp"
	PropertyInspectorDidAppearName    streamdeckcore.EventName = "propertyInspectorDidAppear"
	PropertyInspectorDidDisappearName streamdeckcore.EventName = "propertyInspectorDidDisappear"
	SendToPluginName                  streamdeckcore.EventName = "sendToPlugin"
	SystemDidWakeUpName               streamdeckcore.EventName = "systemDidWakeUp"
	TitleParametersDidChangeName      streamdeckcore.EventName = "titleParametersDidChange"
	WillAppearName                    streamdeckcore.EventName = "willAppear"
	WillDisappearName                 streamdeckcore.EventName = "willDisappear"
)

type ApplicationDidLaunch struct {
	Event   streamdeckcore.EventName    `json:"event"`
	Payload ApplicationDidLaunchPayload `json:"payload"`
}

type ApplicationDidLaunchPayload struct {
	Application string `json:"application"`
}

type ApplicationDidTerminate struct {
	Event   streamdeckcore.EventName       `json:"event"`
	Payload ApplicationDidTerminatePayload `json:"payload"`
}

type ApplicationDidTerminatePayload struct {
	Application string `json:"application"`
}

type DeviceDidConnect struct {
	Event      streamdeckcore.EventName  `json:"event"`
	Device     streamdeckcore.DeviceUUID `json:"device"`
	DeviceInfo DeviceInfo                `json:"deviceInfo"`
}

type DeviceDidDisconnect struct {
	Event  streamdeckcore.EventName  `json:"event"`
	Device streamdeckcore.DeviceUUID `json:"device"`
}

type DidReceiveGlobalSettings struct {
	Event   streamdeckcore.EventName        `json:"event"`
	Payload DidReceiveGlobalSettingsPayload `json:"settings"`
}

type DidReceiveGlobalSettingsPayload struct {
	Settings json.RawMessage `json:"settings"`
}

type DidReceiveSettings struct {
	Action  streamdeckcore.ActionUUID   `json:"action"`
	Event   streamdeckcore.EventName    `json:"event"`
	Context streamdeckcore.EventContext `json:"context"`
	Device  streamdeckcore.DeviceUUID   `json:"device"`
	Payload DidReceiveSettingsPayload   `json:"payload"`
}

type DidReceiveSettingsPayload struct {
	Settings        json.RawMessage `json:"settings"`
	Coordinates     Coordinates     `json:"coordinates"`
	IsInMultiAction bool            `json:"isInMultiAction"`
}

type KeyDown struct {
	Action  streamdeckcore.ActionUUID   `json:"action"`
	Event   streamdeckcore.EventName    `json:"event"`
	Context streamdeckcore.EventContext `json:"context"`
	Device  streamdeckcore.DeviceUUID   `json:"device"`
	Payload KeyDownPayload              `json:"payload"`
}

type KeyDownPayload struct {
	Settings         json.RawMessage `json:"settings"`
	Coordinates      Coordinates     `json:"coordinates"`
	State            int             `json:"state"`
	UserDesiredState int             `json:"userDesiredState"`
	IsInMultiAction  bool            `json:"isInMultiAction"`
}

type KeyUp struct {
	Action  streamdeckcore.ActionUUID   `json:"action"`
	Event   streamdeckcore.EventName    `json:"event"`
	Context streamdeckcore.EventContext `json:"context"`
	Device  streamdeckcore.DeviceUUID   `json:"device"`
	Payload KeyUpPayload                `json:"payload"`
}

type KeyUpPayload struct {
	Settings         json.RawMessage `json:"settings"`
	Coordinates      Coordinates     `json:"coordinates"`
	State            int             `json:"state"`
	UserDesiredState int             `json:"userDesiredState"`
	IsInMultiAction  bool            `json:"isInMultiAction"`
}

type PropertyInspectorDidAppear struct {
	Action  streamdeckcore.ActionUUID   `json:"action"`
	Event   streamdeckcore.EventName    `json:"event"`
	Context streamdeckcore.EventContext `json:"context"`
	Device  streamdeckcore.DeviceUUID   `json:"device"`
}

type PropertyInspectorDidDisappear struct {
	Action  streamdeckcore.ActionUUID   `json:"action"`
	Event   streamdeckcore.EventName    `json:"event"`
	Context streamdeckcore.EventContext `json:"context"`
	Device  streamdeckcore.DeviceUUID   `json:"device"`
}

type SendToPlugin struct {
	Action  streamdeckcore.ActionUUID   `json:"action"`
	Event   streamdeckcore.EventName    `json:"event"`
	Context streamdeckcore.EventContext `json:"context"`
	Payload json.RawMessage             `json:"payload"`
}

type SystemDidWakeUp struct {
	Event streamdeckcore.EventName `json:"event"`
}

type TitleParametersDidChange struct {
	Action  streamdeckcore.ActionUUID       `json:"action"`
	Event   streamdeckcore.EventName        `json:"event"`
	Context streamdeckcore.EventContext     `json:"context"`
	Device  streamdeckcore.DeviceUUID       `json:"device"`
	Payload TitleParametersDidChangePayload `json:"payload"`
}

type TitleParametersDidChangePayload struct {
	Coordinates Coordinates     `json:"coordinates"`
	Settings    json.RawMessage `json:"settings"`
	State       int             `json:"state"`
	Title       string          `json:"title"`
}

type TitleParameters struct {
	FontFamily     string            `json:"fontFamily"`
	FontSize       int               `json:"fontSize"`
	FontStyle      string            `json:"fontStyle"`
	FontUnderline  bool              `json:"fontUnderline"`
	ShowTitle      bool              `json:"showTitle"`
	TitleAlignment VerticalAlignment `json:"titleAlignment"`
	TitleColor     Color             `json:"titleColor"`
}

type WillAppear struct {
	Action  streamdeckcore.ActionUUID   `json:"action"`
	Event   streamdeckcore.EventName    `json:"event"`
	Context streamdeckcore.EventContext `json:"context"`
	Device  streamdeckcore.DeviceUUID   `json:"device"`
	Payload WillAppearPayload           `json:"payload"`
}

type WillAppearPayload struct {
	Settings        json.RawMessage `json:"settings"`
	Coordinates     Coordinates     `json:"coordinates"`
	State           int             `json:"state"`
	IsInMultiAction bool            `json:"isInMultiAction"`
}

type WillDisappear struct {
	Action  streamdeckcore.ActionUUID   `json:"action"`
	Event   streamdeckcore.EventName    `json:"event"`
	Context streamdeckcore.EventContext `json:"context"`
	Device  streamdeckcore.DeviceUUID   `json:"device"`
	Payload WillDisappearPayload        `json:"payload"`
}

type WillDisappearPayload struct {
	Settings        json.RawMessage `json:"settings"`
	Coordinates     Coordinates     `json:"coordinates"`
	State           int             `json:"state"`
	IsInMultiAction bool            `json:"isInMultiAction"`
}
