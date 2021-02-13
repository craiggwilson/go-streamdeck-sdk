package streamdeck

import (
	"encoding/json"
)

type Publisher interface {
	SetSettings(settings json.RawMessage)
}