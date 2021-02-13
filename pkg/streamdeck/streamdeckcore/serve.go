package streamdeckcore

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"golang.org/x/net/websocket"
)

type EventHandler interface {
	Initialize(eventPublisher EventPublisher)
	HandleEvent(ctx context.Context, raw json.RawMessage) error
}

type EventPublisher interface {
	PublishEvent(raw json.RawMessage) error
}

// Serve wraps a websocket to handle receiving and publishing events. To shutdown, cancel the provided context.
func Serve(ctx context.Context, cfg *Config, eventHandler EventHandler) error {
	url := fmt.Sprintf("ws://127.0.0.1:%d", cfg.Port)
	log.Printf("[core] pluginUUID %q connecting to %s", cfg.PluginUUID, url)

	ws, err := websocket.Dial(url, "", "http://127.0.0.1")
	if err != nil {
		return fmt.Errorf("dialing %s: %w", url, err)
	}

	var publishLock sync.Mutex
	publishFunc := eventPublisherFunc(func(raw json.RawMessage) error {
		publishLock.Lock()
		defer publishLock.Unlock()
		log.Print("[core] sending message")
		if err := websocket.JSON.Send(ws, raw); err != nil {
			return fmt.Errorf("sending event: %w", err)
		}

		return nil
	})
	eventHandler.Initialize(publishFunc)

	go func() {
		for {
			var event json.RawMessage
			if err = websocket.JSON.Receive(ws, &event); err != nil {
				log.Printf("[core] ERROR receiving event: %v", err)
				continue
			}

			if err = eventHandler.HandleEvent(ctx, event); err != nil {
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

	_ = ws.Close()

	return ctx.Err()
}

type eventPublisherFunc func(raw json.RawMessage) error
func (f eventPublisherFunc) PublishEvent(raw json.RawMessage) error {
	return f(raw)
}