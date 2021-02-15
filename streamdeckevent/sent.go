package streamdeckevent

import (
	"encoding/json"

	"github.com/craiggwilson/go-streamdeck-sdk/streamdeckcore"
)

const (
	GetGlobalSettingsName       streamdeckcore.EventName = "getGlobalSettings"
	GetSettingsName             streamdeckcore.EventName = "getSettings"
	LogMessageName              streamdeckcore.EventName = "logMessage"
	OpenURLName                 streamdeckcore.EventName = "openUrl"
	SendToPropertyInspectorName streamdeckcore.EventName = "sendToPropertyInspector"
	SetGlobalSettingsName       streamdeckcore.EventName = "setGlobalSettings"
	SetSettingsName             streamdeckcore.EventName = "setSettings"
	SetImageName                streamdeckcore.EventName = "setImage"
	SetStateName                streamdeckcore.EventName = "setState"
	SetTitleName                streamdeckcore.EventName = "setTitle"
	ShowAlertName               streamdeckcore.EventName = "showAlert"
	ShowOKName                  streamdeckcore.EventName = "showOk"
	SwitchToProfileName         streamdeckcore.EventName = "switchToProfile"
)

type GetGlobalSettings struct {
	Event   streamdeckcore.EventName  `json:"event"`
	Context streamdeckcore.PluginUUID `json:"context"`
}

type GetSettings struct {
	Event   streamdeckcore.EventName    `json:"event"`
	Context streamdeckcore.EventContext `json:"context"`
}

type LogMessage struct {
	Event   streamdeckcore.EventName `json:"event"`
	Payload LogMessagePayload        `json:"payload"`
}

type LogMessagePayload struct {
	Message string `json:"message"`
}

type OpenURL struct {
	Event   streamdeckcore.EventName `json:"event"`
	Payload OpenURLPayload           `json:"payload"`
}

type OpenURLPayload struct {
	URL string `json:"url"`
}

type SendToPropertyInspector struct {
	Action  streamdeckcore.ActionUUID   `json:"action"`
	Event   streamdeckcore.EventName    `json:"event"`
	Context streamdeckcore.EventContext `json:"context"`
	Payload json.RawMessage             `json:"payload"`
}

type SetGlobalSettings struct {
	Event   streamdeckcore.EventName  `json:"event"`
	Context streamdeckcore.PluginUUID `json:"context"`
	Payload json.RawMessage           `json:"payload"`
}

type SetImage struct {
	Event   streamdeckcore.EventName    `json:"event"`
	Context streamdeckcore.EventContext `json:"context"`
	Payload SetImagePayload             `json:"payload"`
}

type SetImagePayload struct {
	Image  Base64String `json:"image"`
	Target Target       `json:"target"`
	State  *int         `json:"state,omitempty"`
}

type SetSettings struct {
	Event   streamdeckcore.EventName    `json:"event"`
	Context streamdeckcore.EventContext `json:"context"`
	Payload json.RawMessage             `json:"payload"`
}

type SetState struct {
	Event   streamdeckcore.EventName    `json:"event"`
	Context streamdeckcore.EventContext `json:"context"`
	Payload SetStatePayload             `json:"payload"`
}

type SetStatePayload struct {
	State int `json:"state"`
}

type SetTitle struct {
	Event   streamdeckcore.EventName    `json:"event"`
	Context streamdeckcore.EventContext `json:"context"`
	Payload SetTitlePayload             `json:"payload"`
}

type SetTitlePayload struct {
	Title  string `json:"title"`
	Target Target `json:"target"`
	State  *int   `json:"state,omitempty"`
}

type ShowAlert struct {
	Event   streamdeckcore.EventName    `json:"event"`
	Context streamdeckcore.EventContext `json:"context"`
}

type ShowOK struct {
	Event   streamdeckcore.EventName    `json:"event"`
	Context streamdeckcore.EventContext `json:"context"`
}

type SwitchToProfile struct {
	Event   streamdeckcore.EventName    `json:"event"`
	Context streamdeckcore.EventContext `json:"context"`
	Device  streamdeckcore.DeviceUUID   `json:"device"`
	Payload SwitchToProfilePayload      `json:"payload"`
}

type SwitchToProfilePayload struct {
	Profile streamdeckcore.ProfileName `json:"profile"`
}
