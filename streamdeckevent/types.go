package streamdeckevent

import (
	"github.com/craiggwilson/go-streamdeck-sdk/streamdeckcore"
)

// Base64String is a base64 string.
type Base64String string

// Color is a color.
type Color string

// Coordinates is the column and row of a button.
type Coordinates struct {
	Column int `json:"column"`
	Row int `json:"row"`
}

// DeviceInfo is the device information provided by the device.
type DeviceInfo struct {
	Type       DeviceType `json:"type,omitempty"`
	Size       DeviceSize `json:"size,omitempty"`
	DeviceName streamdeckcore.DeviceName     `json:"deviceName,omitempty"`
}

// DeviceSize is the size of a Streamdeck.
type DeviceSize struct {
	Columns int `json:"columns,omitempty"`
	Rows    int `json:"rows,omitempty"`
}

// DeviceType is the type of device.
type DeviceType int
const (
	StreamDeck       DeviceType = 0
	StreamDeckMini   DeviceType = 1
	StreamDeckXL     DeviceType = 2
	StreamDeckMobile DeviceType = 3
	CorsairGKeys DeviceType = 4
)

// Target indicates where to apply an event.
type Target int
const (
	HardwareAndSoftware Target = 0
	OnlyHardware        Target = 1
	OnlySoftware        Target = 2
)

// VerticalAlignment is a vertical alignment.
type VerticalAlignment string
const (
	Bottom VerticalAlignment = "bottom"
	Middle VerticalAlignment = "middle"
	Top    VerticalAlignment = "top"
)
