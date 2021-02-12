package statusplugin

import (
	"time"
)

// StatusSettings holds the settings for the pi-hole status plugin.
type StatusSettings struct {
	AdminURL string `json:"adminURL"`
	APIKey string `json:"apiKey"`
	Duration time.Duration `json:"duration"`
	Enabled bool `json:"enabled"`
}

