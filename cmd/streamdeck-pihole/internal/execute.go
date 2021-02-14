package internal

import (
	"context"
	"log"

	"github.com/samwho/streamdeck"

	"github.com/craiggwilson/streamdeck-plugins/cmd/streamdeck-pihole/internal/status"
)

func Execute(args []string) {
	params, err := streamdeck.ParseRegistrationParams(args)
	if err != nil {
		log.Fatalf("parsing registration params: %v", err)
	}

	client := makeClient(params)
	defer func() {
		_ = client.Close()
	}()

	if err = client.Run(); err != nil {
		log.Fatalf("running: %v", err)
	}
}

func makeClient(params streamdeck.RegistrationParams) *streamdeck.Client {
	client := streamdeck.NewClient(context.Background(), params)
	status.Register(client)
	return client
}
