package internal

import (
	"context"
	"log"

	"github.com/craiggwilson/streamdeck-plugins/cmd/streamdeck-counter/internal/counter"
	"github.com/craiggwilson/streamdeck-plugins/pkg/streamdeck"
	"github.com/craiggwilson/streamdeck-plugins/pkg/streamdeck/streamdeckcore"
)

func Execute(args []string) {
	cfg, err := streamdeckcore.ParseConfig(args)
	if err != nil {
		log.Fatalf("parsing config args: %v", err)
	}

	plugin := streamdeck.NewPlugin(
		counter.New(),
	)

	if err = streamdeckcore.Serve(context.Background(), cfg, plugin); err != nil {
		log.Fatalf("serving: %v", err)
	}
}
