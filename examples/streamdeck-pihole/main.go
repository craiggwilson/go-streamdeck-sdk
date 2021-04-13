package main

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/craiggwilson/go-streamdeck-sdk/examples/streamdeck-pihole/internal/status"
	"github.com/craiggwilson/go-streamdeck-sdk/streamdeckutil"
)

func main() {
	lf, err := os.OpenFile(filepath.Join(os.TempDir(), "streamdeck-pihole.log"), os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Panicf("opening log file: %v", err)
	}
	defer func() {
		_ = lf.Close()
	}()

	log.SetOutput(lf)
	if err = streamdeckutil.Serve(context.Background(), os.Args, status.New()); err != nil {
		log.Fatal(err.Error())
	}
}