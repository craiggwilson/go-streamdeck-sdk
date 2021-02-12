package internal

import (
	"context"
	"log"

	"github.com/samwho/streamdeck"

	"streamdeckpihole/cmd/streamdeck-pihole/internal/statusplugin"
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
	statusplugin.Register(client)
	return client
}
