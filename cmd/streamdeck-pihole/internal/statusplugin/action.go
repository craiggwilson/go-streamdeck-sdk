package statusplugin

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/samwho/streamdeck"
)

// Register adds the Status action to the client.
func Register(client *streamdeck.Client) {
	action := client.Action("com.craiggwilson.streamdeck.pihole.status")

	settings := make(map[string]*StatusSettings)

	loadSettings := func(eventContext string) *StatusSettings{
		s, ok := settings[eventContext]
		if !ok {
			s = &StatusSettings{}
			settings[eventContext] = s
		}

		return s
	}

	action.RegisterHandler(streamdeck.WillAppear, func(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
		log.Printf("received WillAppear event (%s): %s", event.Context, event.Payload)
		var p streamdeck.WillAppearPayload
		if err := json.Unmarshal(event.Payload, &p); err != nil {
			return fmt.Errorf("unmarshaling WillAppear payload: %w", err)
		}

		s := loadSettings(event.Context)
		if err := json.Unmarshal(p.Settings, s); err != nil {
			return fmt.Errorf("unmarshaling WillAppear payload settings: %w", err)
		}

		if err := client.SendToPropertyInspector(ctx, p.Settings); err != nil {
			return fmt.Errorf("sending to property inspector: %w", err)
		}

		if err := updateDisplay(ctx, client, s); err != nil {
			return fmt.Errorf("setting title: %w", err)
		}

		return nil
	})

	action.RegisterHandler(streamdeck.SendToPlugin, func(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
		log.Printf("received SendToPlugin event (%s): %s", event.Context, event.Payload)
		s := loadSettings(event.Context)
		if err := json.Unmarshal(event.Payload, s); err != nil {
			return fmt.Errorf("unmarshaling SendToPlugin payload: %w", err)
		}

		if err := client.SetSettings(ctx, s); err != nil {
			return fmt.Errorf("setting settings: %w", err)
		}

		return nil
	})

	action.RegisterHandler(streamdeck.WillDisappear, func(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
		log.Printf("received Will Disappear event (%s): %s", event.Context, event.Payload)
		delete(settings, event.Context)
		return nil
	})

	action.RegisterHandler(streamdeck.KeyUp, func(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
		log.Printf("received KeyUp event (%s): %s", event.Context, event.Payload)
		s, ok := settings[event.Context]
		if !ok {
			return fmt.Errorf("no settings for context %q", event.Context)
		}

		s.Enabled = !s.Enabled
		if err := client.SetSettings(ctx, s); err != nil {
			return fmt.Errorf("setting settings: %w", err)
		}

		if err := updateDisplay(ctx, client, s); err != nil {
			return fmt.Errorf("setting title: %w", err)
		}

		return nil
	})
}

func updateDisplay(ctx context.Context, client *streamdeck.Client, s *StatusSettings) error {
	title := "Enabled"
	if !s.Enabled {
		title = "Disabled"
	}

	return client.SetTitle(ctx, title, streamdeck.HardwareAndSoftware)
}
