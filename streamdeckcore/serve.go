package streamdeckcore

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

// Handler is implemented to handle events.
type Handler interface {
	HandleEvent(ctx context.Context, raw json.RawMessage) error
}

// Plugin is implemented by a plugin in order to interact with a device.
type Plugin interface {
	Handler
	Initialize(pluginUUID PluginUUID, publisher Publisher)
}

// Publisher is provided to Plugins so they can communicate with a device.
type Publisher interface {
	PublishEvent(raw json.RawMessage) error
}

// Serve wraps a websocket to handle receiving and publishing events. To shutdown, cancel the provided context.
func Serve(ctx context.Context, cfg *Config, plugin Plugin) error {
	url := fmt.Sprintf("ws://127.0.0.1:%d", cfg.Port)
	log.Printf("[core] pluginUUID %q connecting to %s", cfg.PluginUUID, url)

	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return fmt.Errorf("dialing %s: %w", url, err)
	}

	var publishLock sync.Mutex
	publishFunc := eventPublisherFunc(func(raw json.RawMessage) error {
		publishLock.Lock()
		defer publishLock.Unlock()
		log.Printf("[core] sending message %v", string(raw))
		if err := c.WriteMessage(websocket.TextMessage, raw); err != nil {
			return fmt.Errorf("sending event: %w", err)
		}

		return nil
	})
	plugin.Initialize(cfg.PluginUUID, publishFunc)

	go func() {
		for {
			mt, msg, err := c.ReadMessage()
			if err != nil {
				log.Printf("[core] ERROR receiving event: %v", err)
				return
			}

			if mt == websocket.PingMessage {
				if err := c.WriteMessage(websocket.PongMessage, []byte{}); err != nil {
					log.Printf("error while ponging: %v\n", err)
				}
				continue
			}

			log.Printf("[core] received message: %s", string(msg))

			if err = plugin.HandleEvent(ctx, msg); err != nil {
				log.Printf("[core] ERROR handling event: %v", err)
			}
		}
	}()

	type registerEvent struct {
		PluginUUID PluginUUID `json:"uuid,omitempty"`
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
	_ = c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	return ctx.Err()
}

type eventPublisherFunc func(raw json.RawMessage) error

func (f eventPublisherFunc) PublishEvent(raw json.RawMessage) error {
	return f(raw)
}