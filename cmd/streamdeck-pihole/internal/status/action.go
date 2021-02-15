package status

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/samwho/streamdeck"

	"github.com/craiggwilson/go-streamdeck-sdk/cmd/streamdeck-pihole/internal/pihole"
)

// Register adds the CurrentStatus action to the client.
func Register(client *streamdeck.Client) {
	action := client.Action("com.craiggwilson.streamdeck.pihole.status")

	var stateLock sync.Mutex
	states := make(map[string]*state)

	loadState := func(eventContext string) *state {
		stateLock.Lock()
		defer stateLock.Unlock()
		s, ok := states[eventContext]
		if !ok {
			s = newState()
			states[eventContext] = s
		}

		return s
	}

	deleteState := func(eventContext string) {
		stateLock.Lock()
		defer stateLock.Unlock()
		s, ok := states[eventContext]
		if ok {
			s.Close()
			delete(states, eventContext)
		}
	}

	action.RegisterHandler(streamdeck.WillAppear, func(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
		log.Printf("received WillAppear event (%s): %s", event.Context, event.Payload)
		var p streamdeck.WillAppearPayload
		if err := json.Unmarshal(event.Payload, &p); err != nil {
			return handleError("unmarshaling WillAppear payload: %w", err)
		}

		s := loadState(event.Context)
		if err := s.UpdateSettingsFromJSON(p.Settings); err != nil {
			return handleError("unmarshaling WillAppear payload states: %w", err)
		}

		m := s.Monitor()
		go func() {
			for su := range m {
				var title string
				if su.RemainingDisabledSeconds > 0 {
					d := time.Duration(su.RemainingDisabledSeconds) * time.Second
					title = d.String()
				}

				if err := client.SetTitle(ctx, title, streamdeck.HardwareAndSoftware); err != nil {
					_ = handleError("setting title: %w", err)
				}

				switch su.Status {
				case pihole.Enabled:
					if err := client.SetState(ctx, 0); err != nil {
						_ = handleError("setting state: %w", err)
					}
				case pihole.Disabled:
					if err := client.SetState(ctx, 1); err != nil {
						_ = handleError("setting state: %w", err)
					}
				}
			}
		}()

		return nil
	})

	action.RegisterHandler(streamdeck.WillDisappear, func(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
		log.Printf("received Will Disappear event (%s): %s", event.Context, event.Payload)
		deleteState(event.Context)
		return nil
	})

	action.RegisterHandler(streamdeck.DidReceiveSettings, func(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
		log.Printf("received DidReceiveSettings event (%s): %s", event.Context, event.Payload)
		var p streamdeck.DidReceiveSettingsPayload
		if err := json.Unmarshal(event.Payload, &p); err != nil {
			return handleError("unmarshalling DidReceiveSettings payload: %w", err)
		}

		s := loadState(event.Context)
		if err := s.UpdateSettingsFromJSON(p.Settings); err != nil {
			return handleError("unmarshalling DidReceiveSettings payload settings: %w", err)
		}

		if err := client.SetSettings(ctx, s.Settings); err != nil {
			return handleError("setting states: %w", err)
		}

		return nil
	})

	action.RegisterHandler(streamdeck.KeyUp, func(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
		log.Printf("received KeyUp event (%s): %s", event.Context, event.Payload)
		s, ok := states[event.Context]
		if !ok {
			return handleError("no states for context %q", event.Context)
		}

		currentStatus, err := s.CurrentStatus()
		if err != nil {
			return client.ShowAlert(ctx)
		}

		switch currentStatus {
		case pihole.Enabled:
			if err := s.SetStatus(pihole.Disabled); err != nil {
				return client.ShowAlert(ctx)
			}
		case pihole.Disabled:
			if err := s.SetStatus(pihole.Enabled); err != nil {
				return client.ShowAlert(ctx)
			}
		}

		return nil
	})
}
