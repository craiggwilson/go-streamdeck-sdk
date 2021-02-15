package streamdeckutil

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/craiggwilson/go-streamdeck-sdk"
	"github.com/craiggwilson/go-streamdeck-sdk/streamdeckcore"
)

// Serve is a helper method for launching a plugin. It parse the arguments and listens for the os.Interrupt event
// to shutdown.
func Serve(ctx context.Context, args []string, actions ...streamdeck.Action) error {
	cfg, err := streamdeckcore.ParseConfig(args)
	if err != nil {
		return fmt.Errorf("parsing config args: %w", err)
	}

	plugin := streamdeck.NewDefaultPlugin(actions...)

	ctx, cancel := context.WithCancel(ctx)
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

	return streamdeckcore.Serve(ctx, cfg, plugin)
}
