package status

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/craiggwilson/go-streamdeck-sdk"
	"github.com/craiggwilson/go-streamdeck-sdk/cmd/streamdeck-pihole/internal/pihole"
	"github.com/craiggwilson/go-streamdeck-sdk/streamdeckevent"
)

const actionUUID = "com.craiggwilson.streamdeck.pihole.status"

func New() *streamdeck.InstancedAction {
	return streamdeck.NewInstancedAction(
		actionUUID,
		func(eventContext streamdeck.EventContext, publisher streamdeck.ActionInstancePublisher) streamdeck.ActionInstance {
			instance := &ActionInstance{
				eventContext: eventContext,
				publisher:    publisher,
			}

			return instance
		},
	)
}

type ActionInstance struct {
	eventContext streamdeck.EventContext
	publisher    streamdeck.ActionInstancePublisher

	ph      *pihole.PiHole
	monitor *pihole.Monitor
}

type instanceSettings struct {
	AdminURL               string `json:"adminURL,omitempty"`
	APIKey                 string `json:"apiKey,omitempty"`
	RefreshIntervalSeconds int    `json:"refreshIntervalSeconds,omitempty"`
	DisableDurationSeconds int    `json:"disableDurationSeconds,omitempty"`
}

func (a *ActionInstance) ActionUUID() streamdeck.ActionUUID {
	return actionUUID
}

func (a *ActionInstance) EventContext() streamdeck.EventContext {
	return a.eventContext
}

func (a *ActionInstance) HandleDidReceiveSettings(_ context.Context, event streamdeckevent.DidReceiveSettings) error {
	err := a.publisher.SetSettings(event.Payload.Settings)
	if err != nil {
		_ = a.publisher.ShowAlert()
		return fmt.Errorf("applying settings: %w", err)
	}

	if _, err := a.applySettings(event.Payload.Settings); err != nil {
		_ = a.publisher.ShowAlert()
		return fmt.Errorf("receiving settings: %w", err)
	}

	return nil
}

func (a *ActionInstance) HandleKeyDown(_ context.Context, event streamdeckevent.KeyDown) error {
	settings, err := a.applySettings(event.Payload.Settings)
	if err != nil {
		_ = a.publisher.ShowAlert()
		return fmt.Errorf("receiving settings: %w", err)
	}

	if a.monitor == nil {
		return nil
	}

	status := a.monitor.Status()
	switch status.Status { //nolint:exhaustive
	case pihole.Enabled:
		a.monitor.Disable(settings.DisableDurationSeconds)
		duration := time.Duration(settings.DisableDurationSeconds) * time.Second
		a.monitor.RefreshIn(duration)
	default:
		a.monitor.Enable()
	}

	return nil
}

func (a *ActionInstance) HandleWillAppear(_ context.Context, event streamdeckevent.WillAppear) error {
	if _, err := a.applySettings(event.Payload.Settings); err != nil {
		_ = a.publisher.ShowAlert()
		return fmt.Errorf("receiving settings: %w", err)
	}

	return nil
}

func (a *ActionInstance) applySettings(settingsRaw json.RawMessage) (instanceSettings, error) {
	var settings instanceSettings
	if err := json.Unmarshal(settingsRaw, &settings); err != nil {
		return instanceSettings{}, fmt.Errorf("unmarshalling settings: %w", err)
	}

	if a.monitor != nil {
		if a.ph.AdminURL() == settings.AdminURL &&
			a.ph.APIKey() == settings.APIKey &&
			a.monitor.RefreshInterval() == time.Duration(settings.RefreshIntervalSeconds)*time.Second {
			return settings, nil
		}

		a.monitor.Stop()
	}

	a.ph = pihole.New(settings.AdminURL, settings.APIKey)
	a.monitor = a.ph.Monitor(time.Duration(settings.RefreshIntervalSeconds) * time.Second)
	currentStatus := a.monitor.Status()
	ch, _ := a.monitor.Subscribe()

	go func() {
		ticker := time.NewTicker(1 * time.Second)
		for {
			a.updateStatus(currentStatus)
			select {
			case <-ticker.C:
			case su, ok := <-ch:
				if !ok {
					ticker.Stop()
					return
				}
				currentStatus = su
			}
		}
	}()

	return settings, nil
}

func (a *ActionInstance) updateStatus(status pihole.StatusUpdate) {
	var title string
	var state int

	switch status.Status { //nolint:exhaustive
	case pihole.Enabled:
		state = 0
		title = ""
	case pihole.Disabled:
		state = 1
		if !status.DisabledUntil.IsZero() {
			remaining := time.Duration(int(time.Until(status.DisabledUntil).Seconds())) * time.Second
			if remaining > 0 {
				title = remaining.String()
			}
		}
	default:
		state = 1
		title = "(unknown)"
	}

	_ = a.publisher.SetTitle(streamdeckevent.SetTitlePayload{Title: title})
	_ = a.publisher.SetState(streamdeckevent.SetStatePayload{State: state})
}
