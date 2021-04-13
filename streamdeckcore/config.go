package streamdeckcore

import (
	"flag"
	"fmt"
)

// Config holds the launch configuration for a plugin.
type Config struct {
	Port          int
	PluginUUID    PluginUUID
	RegisterEvent EventName
	Info          string
}

// ParseConfig parses the configuration from the provide arguments.
func ParseConfig(args []string) (*Config, error) {
	f := flag.NewFlagSet("config", flag.ContinueOnError)

	port := f.Int("port", -1, "")
	pluginUUID := f.String("pluginUUID", "", "")
	registerEvent := f.String("registerEvent", "", "")
	info := f.String("info", "", "")

	if err := f.Parse(args[1:]); err != nil {
		return nil, err
	}

	if *port == -1 {
		return nil, fmt.Errorf("missing -port flag")
	}
	if *pluginUUID == "" {
		return nil, fmt.Errorf("missing -pluginUUID flag")
	}
	if *registerEvent == "" {
		return nil, fmt.Errorf("missing -registerEvent flag")
	}
	if *info == "" {
		return nil, fmt.Errorf("missing -info flag")
	}

	return &Config{
		Port:          *port,
		PluginUUID:    PluginUUID(*pluginUUID),
		RegisterEvent: EventName(*registerEvent),
		Info:          *info,
	}, nil
}
