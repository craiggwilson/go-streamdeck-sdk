package status

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/craiggwilson/go-streamdeck-sdk/cmd/streamdeck-pihole/internal/pihole"
	"github.com/craiggwilson/go-streamdeck-sdk/streamdeckevent"

	"github.com/craiggwilson/go-streamdeck-sdk"
)

const uuid = "com.craiggwilson.streamdeck.pihole.status"

func New() *streamdeck.DefaultAction {
	return streamdeck.NewDefaultAction(
		uuid,
		func(eventContext streamdeck.EventContext, publisher streamdeck.ActionInstancePublisher) streamdeck.ActionInstance {
			instance := &ActionInstance{
				eventContext: eventContext,
				publisher:    publisher,
				forcedStatusChange: make(chan forcedStatusChange, 1),
			}

			return instance
		},
	)
}

type ActionInstance struct {
	eventContext streamdeck.EventContext
	publisher    streamdeck.ActionInstancePublisher

	ph            *pihole.PiHole
	monitor       *pihole.Monitor
	forcedStatusChange chan forcedStatusChange
}

type forcedStatusChange struct {
	status pihole.Status
	disabledUntil time.Time
}

type instanceSettings struct {
	AdminURL               string `json:"adminURL,omitempty"`
	APIKey                 string `json:"apiKey,omitempty"`
	RefreshIntervalSeconds int    `json:"refreshIntervalSeconds,omitempty"`
	DisableDurationSeconds int    `json:"disableDurationSeconds,omitempty"`
}

func (a *ActionInstance) UUID() streamdeck.ActionUUID {
	return uuid
}

func (a *ActionInstance) Context() streamdeck.EventContext {
	return a.eventContext
}

func (a *ActionInstance) HandleDidReceiveSettings(_ context.Context, event streamdeckevent.DidReceiveSettings) error {
	log.Printf("handling did receive")
	err := a.publisher.SetSettings(event.Payload.Settings)
	if err != nil {
		_ = a.publisher.ShowAlert()
		return fmt.Errorf("applying settings: %w", err)
	}

	if err := a.receivedSettings(event.Payload.Settings); err != nil {
		_ = a.publisher.ShowAlert()
		return fmt.Errorf("receiving settings: %w", err)
	}

	return nil
}

func (a *ActionInstance) HandleKeyDown(_ context.Context, event streamdeckevent.KeyDown) error {
	log.Printf("handling key down")
	var settings instanceSettings
	if err := json.Unmarshal(event.Payload.Settings, &settings); err != nil {
		_ = a.publisher.ShowAlert()
		return fmt.Errorf("unmarshalling settings: %w", err)
	}

	if a.ph == nil {
		return nil
	}

	status, err := a.ph.Status()
	if err != nil {
		a.forcedStatusChange <- forcedStatusChange{status: status}
		return nil
	}

	switch status {
	case pihole.Enabled:
		if err := a.ph.Disable(settings.DisableDurationSeconds); err != nil {
			_ = a.publisher.ShowAlert()
			return fmt.Errorf("disabling pi-hole: %w", err)
		}

		duration := time.Duration(settings.DisableDurationSeconds) * time.Second
		a.forcedStatusChange <- forcedStatusChange{pihole.Disabled, time.Now().Add(duration)}
		a.monitor.RefreshIn(duration)
	case pihole.Disabled:
		if err := a.ph.Enable(); err != nil {
			_ = a.publisher.ShowAlert()
			return fmt.Errorf("enabling pi-hole: %w", err)
		}
		a.forcedStatusChange <- forcedStatusChange{pihole.Enabled, time.Time{}}
		a.monitor.RefreshIn(0)
	}

	return nil
}

func (a *ActionInstance) HandleWillAppear(_ context.Context, event streamdeckevent.WillAppear) error {
	log.Printf("handling will appear")
	if err := a.receivedSettings(event.Payload.Settings); err != nil {
		_ = a.publisher.ShowAlert()
		return fmt.Errorf("receiving settings: %w", err)
	}

	status, _ := a.ph.Status()
	a.forcedStatusChange <- forcedStatusChange{status: status, disabledUntil: time.Time{}}

	return nil
}

func (a *ActionInstance) receivedSettings(settingsRaw json.RawMessage) error {
	var settings instanceSettings
	if err := json.Unmarshal(settingsRaw, &settings); err != nil {
		_ = a.publisher.ShowAlert()
		return fmt.Errorf("unmarshalling settings: %w", err)
	}

	if a.monitor != nil &&
		a.ph.AdminURL() == settings.AdminURL &&
		a.ph.APIKey() == settings.APIKey &&
		a.monitor.RefreshInterval() == time.Duration(settings.RefreshIntervalSeconds) * time.Second {
		return nil
	}

	if a.monitor != nil {
		a.monitor.Stop()
	}
	a.ph = pihole.New(settings.AdminURL, settings.APIKey)
	a.monitor = a.ph.Monitor(time.Duration(settings.RefreshIntervalSeconds) * time.Second)
	ch, _ := a.monitor.Subscribe()
	a.monitor.RefreshIn(0)

	go func() {
		ticker := time.NewTicker(1 * time.Second)
		var status pihole.Status
		var disabledUntil time.Time
		var err error
		for {
			select {
			case <-ticker.C:
			case fsc := <-a.forcedStatusChange:
				status = fsc.status
				disabledUntil = fsc.disabledUntil
				err = nil
			case su, ok := <-ch:
				if !ok {
					ticker.Stop()
					return
				}

				status, err = su.Status, su.Err
			}

			if status == pihole.Enabled {
				disabledUntil = time.Time{}
			}

			a.updateStatus(status, disabledUntil, err)
		}
	}()

	return nil
}

func (a *ActionInstance) updateStatus(status pihole.Status, disabledUntil time.Time, err error) {
	var title string
	var state int

	if err != nil {
		state = 1
		title = "err"
	} else {
		switch status {
		case pihole.Enabled:
			state = 0
			title = ""
		case pihole.Disabled:
			state = 1
			if !disabledUntil.IsZero() {
				remaining := time.Duration(int(time.Until(disabledUntil).Seconds())) * time.Second
				if remaining > 0 {
					title = remaining.String()
				}
			}
		default:
			state = 1
			title = "(unknown)"
		}
	}

	_ = a.publisher.SetTitle(streamdeckevent.SetTitlePayload{Title: title})
	_ = a.publisher.SetState(streamdeckevent.SetStatePayload{State: state})
}
