package internal

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/craiggwilson/streamdeck-plugins/cmd/streamdeck-example/internal/counter"
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

	ctx, cancel := context.WithCancel(context.Background())
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	defer func() {
		signal.Stop(interrupt)
		cancel()
	}()

	go func() {
		select {
		case <-interrupt:
			cancel()
		case <-ctx.Done():
		}
	}()

	if err = streamdeckcore.Serve(ctx, cfg, plugin); err != nil {
		log.Fatalf("serving: %v", err)
	}
}
