package pihole

import (
	"time"
)

// Status indicates whether a Pi-Hole is enabled or disabled.
type Status string
const (
	Unknown Status = ""
	Enabled  Status = "enabled"
	Disabled Status = "disabled"
)

// StatusResponse contains the status of the Pi-Hole.
type StatusResponse struct {
	Status Status `json:"status,omitempty"`
}

// StatusUpdate is used to indicate when a status change has occurred while watching a Pi-Hole.
type StatusUpdate struct {
	Status Status
	DisabledUntil time.Time
	Err error
}
