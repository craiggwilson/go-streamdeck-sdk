package pihole

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// New creates an instance of a PiHole.
func New(adminURL string, apiKey string) *PiHole {
	return &PiHole{
		AdminURL: adminURL,
		APIKey: apiKey,
	}
}

// PiHole is a connection to a Pi-Hole.
type PiHole struct {
	AdminURL string
	APIKey string
}

// Disable disables the Pi-Hole.
func (ph *PiHole) Disable(duration time.Duration) error {
	endpoint := "disable"
	if duration != 0 {
		endpoint = fmt.Sprintf("disable=%v", int(duration.Seconds()))
	}

	return ph.setStatus(endpoint, Disabled)
}

// Enable enables the Pi-Hole.
func (ph *PiHole) Enable() error {
	return ph.setStatus("enable", Enabled)
}

// Status returns the current status of the Pi-Hole.
func (ph *PiHole) Status() (Status, error) {
	raw, err := ph.get("status")
	if err != nil {
		return "", err
	}

	var resp StatusResponse
	if err := json.Unmarshal(raw, &resp); err != nil {
		return "", fmt.Errorf("unmarshaling summaryRaw response: %w", err)
	}

	return resp.Status, nil
}

func (ph *PiHole) get(endpoint string) ([]byte, error) {
	url := fmt.Sprintf("%s?%s", ph.AdminURL, endpoint)
	if ph.APIKey != "" {
		url += fmt.Sprintf("&auth=%s", ph.APIKey)
	}

	resp, err := http.Get(url) //nolint:gosec
	if err != nil {
		return nil, fmt.Errorf("hitting %q: %w", url, err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	return body, nil
}

func (ph *PiHole) setStatus(endpoint string, expectedStatus Status) error {
	raw, err := ph.get(endpoint)
	if err != nil {
		return err
	}

	var resp StatusResponse
	if err := json.Unmarshal(raw, &resp); err != nil {
		return fmt.Errorf("unmarshaling enable response: %w", err)
	}

	if resp.Status != expectedStatus {
		log.Printf("Response status %q, expected status %q: %s", resp.Status, expectedStatus, raw)
		if expectedStatus == Enabled {
			return errors.New("failed to enable")
		} else {
			return errors.New("failed to disable")
		}
	}

	return nil
}
