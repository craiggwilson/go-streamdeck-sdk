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
		adminURL: adminURL,
		apiKey:   apiKey,
	}
}

// PiHole is a connection to a Pi-Hole.
type PiHole struct {
	adminURL string
	apiKey   string
}

// AdminURL returns the admin URL used to communicate with the Pi-Hole.
func (ph *PiHole) AdminURL() string {
	return ph.adminURL
}

// APIKey returns the api key used for updating the Pi-Hole.
func (ph *PiHole) APIKey() string {
	return ph.apiKey
}

// Disable disables the Pi-Hole.
func (ph *PiHole) Disable(durationSeconds int) error {
	endpoint := "disable"
	if durationSeconds != 0 {
		endpoint = fmt.Sprintf("disable=%v", durationSeconds)
	}

	return ph.setStatus(endpoint, Disabled)
}

// Enable enables the Pi-Hole.
func (ph *PiHole) Enable() error {
	return ph.setStatus("enable", Enabled)
}

// Monitor monitors a pi-hole, refreshing at the specified interval. It returns a channel as well as function
// to be used to unsubscribe.
func (ph *PiHole) Monitor(refreshInterval time.Duration) *Monitor {
	m := newMonitor(ph, refreshInterval)
	m.start()
	return m
}

// Status returns the current status of the Pi-Hole.
func (ph *PiHole) Status() (Status, error) {
	raw, err := ph.get("status")
	if err != nil {
		return "", err
	}

	var resp StatusResponse
	if err := json.Unmarshal(raw, &resp); err != nil {
		return Unknown, fmt.Errorf("unmarshaling summaryRaw response: %w", err)
	}

	return resp.Status, nil
}

func (ph *PiHole) get(endpoint string) ([]byte, error) {
	url := fmt.Sprintf("%s?%s", ph.adminURL, endpoint)
	if ph.apiKey != "" {
		url += fmt.Sprintf("&auth=%s", ph.apiKey)
	}

	resp, err := http.Get(url)
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
		}

		return errors.New("failed to disable")
	}

	return nil
}
