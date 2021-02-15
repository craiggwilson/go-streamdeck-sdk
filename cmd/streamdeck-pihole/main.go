package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/craiggwilson/go-streamdeck-sdk/cmd/streamdeck-pihole/internal"
)

func main() {
	lf, err := os.OpenFile(filepath.Join(os.TempDir(), "streamdeck-pihole.log"), os.O_CREATE | os.O_TRUNC, 0644)
	if err != nil {
		log.Panicf("opening log file: %v", err)
	}
	defer func() {
		_ = lf.Close()
	}()

	log.SetOutput(lf)
	internal.Execute(os.Args)
}
