package statusplugin

import (
	"encoding/json"
	"time"

	"streamdeckpihole/pkg/pihole"
)

func defaultSettings() *settings {
	return &settings{}
}

type settings struct {
	AdminURL string `json:"adminURL,omitempty"`
	APIKey string `json:"apiKey,omitempty"`
	DisableDurationSeconds int `json:"disableDurationSeconds,omitempty"`
	RefreshIntervalSeconds int `json:"refreshIntervalSeconds,omitempty"`
}

func newState() *state {
	return &state{
		Settings:               defaultSettings(),
		forcedStatusUpdateChan: make(chan forcedStatusUpdate, 10),
		done:                   make(chan struct{}, 1),
	}
}

type state struct {
	Settings *settings

	forcedStatusUpdateChan chan forcedStatusUpdate
	done                   chan struct{}
}

type forcedStatusUpdate struct {
	Status pihole.Status
	RemainingDisabledSeconds int
}

type statusUpdate struct {
	Status pihole.Status
	RemainingDisabledSeconds int
}

func (s *state) Monitor() <-chan statusUpdate {
	monitor := make(chan statusUpdate, 10)
	go func() {
		var remainingDisabledSeconds int
		var lastReportedStatus pihole.Status
		var lastReportedStatusTime time.Time

		ticker := time.NewTicker(time.Second)
		for {
			select {
			case ch := <-s.forcedStatusUpdateChan:
				remainingDisabledSeconds = ch.RemainingDisabledSeconds
				lastReportedStatus = ch.Status
				lastReportedStatusTime = time.Now()
				monitor <- statusUpdate{RemainingDisabledSeconds: remainingDisabledSeconds, Status: lastReportedStatus}
			case <-ticker.C:
				if lastReportedStatus == "" || (s.Settings.RefreshIntervalSeconds > 0 && time.Since(lastReportedStatusTime) > (time.Duration(s.Settings.RefreshIntervalSeconds) * time.Second)) {
					status, err := s.CurrentStatus()
					if err == nil {
						lastReportedStatus = status
						lastReportedStatusTime = time.Now()
						if lastReportedStatus == pihole.Enabled {
							remainingDisabledSeconds = 0
						}
						monitor <- statusUpdate{RemainingDisabledSeconds: remainingDisabledSeconds, Status: lastReportedStatus}
					}
				}

				if remainingDisabledSeconds > 0 {
					remainingDisabledSeconds--
					if remainingDisabledSeconds == 0 {
						lastReportedStatusTime = time.Unix(0, 0)
					} else {
						monitor <- statusUpdate{RemainingDisabledSeconds: remainingDisabledSeconds, Status: lastReportedStatus}
					}
				}
			case <-s.done:
				close(monitor)
				ticker.Stop()
				return
			}
		}
	}()

	return monitor
}

func (s *state) Close() {
	close(s.done)
}

func (s *state) CurrentStatus() (pihole.Status, error) {
	ph := pihole.New(s.Settings.AdminURL, s.Settings.APIKey)
	status, err := ph.Status()
	if err != nil {
		return "", handleError("getting status: %w", err)
	}

	return status, nil
}

func (s *state) SetStatus(status pihole.Status) error {
	ph := pihole.New(s.Settings.AdminURL, s.Settings.APIKey)
	switch status {
	case pihole.Enabled:
		err := ph.Enable()
		if err != nil {
			return handleError("enabling: %w", err)
		}

		s.forcedStatusUpdateChan <- forcedStatusUpdate{Status: pihole.Enabled}
		return nil
	case pihole.Disabled:
		err := ph.Disable(time.Duration(s.Settings.DisableDurationSeconds) * time.Second)
		if err != nil {
			return handleError("disabling: %w", err)
		}

		s.forcedStatusUpdateChan <- forcedStatusUpdate{Status: pihole.Disabled, RemainingDisabledSeconds: s.Settings.DisableDurationSeconds}
		return nil
	}

	return handleError("unknown status %q", status)
}

func (s *state) UpdateSettingsFromJSON(settings json.RawMessage) error {
	if err := json.Unmarshal(settings, s.Settings); err != nil {
		return handleError("unmarshalling settings: %w", err)
	}

	return nil
}