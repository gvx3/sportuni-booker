package main

import (
	"fmt"
	"os"

	"github.com/gvx3/sportuni-book/pkg/app"
	"github.com/gvx3/sportuni-book/pkg/config"
	flag "github.com/spf13/pflag"
)

func main() {
	configPath := flag.StringP("file", "f", "", "Path to config YAML file")
	flag.Parse()

	var err error
	var cfg *config.Config

	if *configPath != "" {
		cfg, err = config.LoadConfigFromFile(*configPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
			os.Exit(1)
		}
	} else {
		cfg, err = config.NewConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
			os.Exit(1)
		}
	}

	if err := app.RunAppWithConfig(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Application error: %v\n", err)
		os.Exit(1)
	}
}
