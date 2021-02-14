package streamdeckcore

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"nhooyr.io/websocket"
)

// Plugin is implemented by a plugin in order to interact with a Streamdeck.
type Plugin interface {
	Initialize(pluginUUID PluginUUID, publisher Publisher)
	HandleEvent(ctx context.Context, raw json.RawMessage) error
}

// Publisher is provided to Plugins so they can communicate with a Streamdeck.
type Publisher interface {
	PublishEvent(raw json.RawMessage) error
}

// Serve wraps a websocket to handle receiving and publishing events. To shutdown, cancel the provided context.
func Serve(ctx context.Context, cfg *Config, plugin Plugin) error {
	url := fmt.Sprintf("ws://127.0.0.1:%d", cfg.Port)
	log.Printf("[core] pluginUUID %q connecting to %s", cfg.PluginUUID, url)

	c, _, err := websocket.Dial(ctx, url, nil)
	if err != nil {
		return fmt.Errorf("dialing %s: %w", url, err)
	}

	var publishLock sync.Mutex
	publishFunc := eventPublisherFunc(func(raw json.RawMessage) error {
		publishLock.Lock()
		defer publishLock.Unlock()
		log.Printf("[core] sending message %v", string(raw))
		if err := c.Write(ctx, websocket.MessageText, raw); err != nil {
			return fmt.Errorf("sending event: %w", err)
		}
		log.Printf("[core] message sent")

		return nil
	})
	plugin.Initialize(cfg.PluginUUID, publishFunc)

	go func() {
		for {
			var msg json.RawMessage
			if _, msg, err = c.Read(ctx); err != nil {
				log.Printf("[core] ERROR receiving event: %v", err)
				continue
			}

			log.Printf("[core] received message: %s", string(msg))

			if err = plugin.HandleEvent(ctx, msg); err != nil {
				log.Printf("[core] ERROR handling event: %v", err)
			}
		}
	}()

	type registerEvent struct {
		PluginUUID PluginUUID `json:"pluginUUID,omitempty"`
		Event EventName `json:"event,omitempty"`
	}

	raw, _ := json.Marshal(registerEvent{
		PluginUUID: cfg.PluginUUID,
		Event: cfg.RegisterEvent,
	})

	if err := publishFunc(raw); err != nil {
		log.Printf("[core] ERROR registering plugin: %v", err)
		return fmt.Errorf("registering plugin: %w", err)
	}

	<-ctx.Done()

	_ = c.Close(websocket.StatusNormalClosure, "shutdown")

	return ctx.Err()
}

type eventPublisherFunc func(raw json.RawMessage) error
func (f eventPublisherFunc) PublishEvent(raw json.RawMessage) error {
	return f(raw)
}