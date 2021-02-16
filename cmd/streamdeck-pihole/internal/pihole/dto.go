package pihole

// StatusResponse contains the status of the Pi-Hole.
type StatusResponse struct {
	Status Status `json:"status,omitempty"`
}

// Status indicates whether a Pi-Hole is enabled or disabled.
type Status string
const (
	Unknown Status = ""
	Enabled  Status = "enabled"
	Disabled Status = "disabled"
)
