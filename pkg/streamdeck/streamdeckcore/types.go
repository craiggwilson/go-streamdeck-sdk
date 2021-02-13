package streamdeckcore

import (
	"encoding/json"
)

type ActionUUID string

type Coordinates struct {
	Column int `json:"column"`
	Row int `json:"row"`
}

type DeviceUUID string

type DeviceName string

type DidReceiveGlobalSettingsEvent struct {
	Event EventName `json:"event"`
	Payload DidReceiveGlobalSettingsPayload `json:"settings"`
}

type DidReceiveGlobalSettingsPayload struct {
	Settings json.RawMessage `json:"settings"`
}

type DidReceiveSettingsEvent struct {
	Action     ActionUUID                `json:"action"`
	Event      EventName                 `json:"event"`
	Context    EventContext              `json:"context"`
	DeviceUUID DeviceUUID                `json:"device"`
	Payload    DidReceiveSettingsPayload `json:"payload"`
}

type DidReceiveSettingsPayload struct {
	Settings json.RawMessage `json:"settings"`
	Coordinates Coordinates `json:"coordinates"`
	IsInMultiAction bool `json:"isInMultiAction"`
}

type EventContext string

type EventName string
const (
	DidReceiveSettings EventName = "didReceiveSettings"
	DidReceiveGlobalSettings EventName = "didReceiveGlobalSettings"
	GetGlobalSettings EventName = "getGlobalSettings"
	GetSettings EventName = "getSettings"
	KeyDown EventName = "keyDown"
	KeyUp EventName = "keyUp"
	SetGlobalSettings EventName = "setGlobalSettings"
	SetImage EventName = "setImage"
	SetSettings EventName = "setSettings"
	SetState EventName = "setState"
	SetTitle EventName = "setTitle"
	ShowAlert EventName = "showAlert"
	ShowOK EventName = "showOk"
	WillAppear EventName = "willAppear"
	WillDisappear EventName = "willDisappear"
)

type GetGlobalSettingsEvent struct {
	Event EventName `json:"event"`
	Context PluginUUID `json:"context"`
}

type GetSettingsEvent struct {
	Event EventName `json:"event"`
	Context EventContext `json:"context"`
}

type Base64String string

type KeyDownEvent struct {
	Action     ActionUUID     `json:"action"`
	Event      EventName      `json:"event"`
	Context    EventContext   `json:"context"`
	DeviceUUID DeviceUUID     `json:"device"`
	Payload    KeyDownPayload `json:"payload"`
}

type KeyDownPayload struct {
	Settings json.RawMessage `json:"settings"`
	Coordinates Coordinates `json:"coordinates"`
	State int `json:"state"`
	UserDesiredState int `json:"userDesiredState"`
	IsInMultiAction bool `json:"isInMultiAction"`
}

type KeyUpEvent struct {
	Action     ActionUUID   `json:"action"`
	Event      EventName    `json:"event"`
	Context    EventContext `json:"context"`
	DeviceUUID DeviceUUID   `json:"device"`
	Payload    KeyUpPayload `json:"payload"`
}

type KeyUpPayload struct {
	Settings json.RawMessage `json:"settings"`
	Coordinates Coordinates `json:"coordinates"`
	State int `json:"state"`
	UserDesiredState int `json:"userDesiredState"`
	IsInMultiAction bool `json:"isInMultiAction"`
}

type PluginUUID string

type SetGlobalSettingsEvent struct {
	Event EventName `json:"event"`
	Context PluginUUID `json:"context"`
	Payload json.RawMessage `json:"payload"`
}

type SetImageEvent struct {
	Event EventName `json:"event"`
	Context EventContext `json:"context"`
	Payload SetImagePayload `json:"payload"`
}

type SetImagePayload struct {
	Image Base64String `json:"image"`
	Target Target `json:"target"`
	State *int `json:"state,omitempty"`
}

type SetSettingsEvent struct {
	Event EventName `json:"event"`
	Context EventContext `json:"context"`
	Payload json.RawMessage `json:"payload"`
}

type SetStateEvent struct {
	Event EventName `json:"event"`
	Context EventContext `json:"context"`
	Payload SetStatePayload `json:"payload"`
}

type SetStatePayload struct {
	State int `json:"state"`
}

type SetTitleEvent struct {
	Event EventName `json:"event"`
	Context EventContext `json:"context"`
	Payload SetTitlePayload `json:"payload"`
}

type SetTitlePayload struct {
	Title string `json:"title"`
	Target Target `json:"target"`
	State *int `json:"state,omitempty"`
}

type ShowAlertEvent struct {
	Event EventName `json:"event"`
	Context EventContext `json:"context"`
}

type ShowOKEvent struct {
	Event EventName `json:"event"`
	Context EventContext `json:"context"`
}

type Target int
const (
	HardwareAndSoftware Target = 0
	OnlyHardware        Target = 1
	OnlySoftware        Target = 2
)

type WillAppearEvent struct {
	Action     ActionUUID        `json:"action"`
	Event      EventName         `json:"event"`
	Context    EventContext      `json:"context"`
	DeviceUUID DeviceUUID        `json:"device"`
	Payload    WillAppearPayload `json:"payload"`
}

type WillAppearPayload struct {
	Settings json.RawMessage `json:"settings"`
	Coordinates Coordinates `json:"coordinates"`
	State int `json:"state"`
	IsInMultiAction bool `json:"isInMultiAction"`
}

type WillDisappearEvent struct {
	Action     ActionUUID           `json:"action"`
	Event      EventName            `json:"event"`
	Context    EventContext         `json:"context"`
	DeviceUUID DeviceUUID           `json:"device"`
	Payload    WillDisappearPayload `json:"payload"`
}

type WillDisappearPayload struct {
	Settings json.RawMessage `json:"settings"`
	Coordinates Coordinates `json:"coordinates"`
	State int `json:"state"`
	IsInMultiAction bool `json:"isInMultiAction"`
}