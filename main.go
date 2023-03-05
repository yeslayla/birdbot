package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/yeslayla/birdbot/app"
)

func main() {

	configDir, _ := os.UserConfigDir()

	defaultConfigPath := path.Join(configDir, "birdbot", "config.yaml")

	var config_file string
	var version bool
	flag.StringVar(&config_file, "c", defaultConfigPath, "Path to config file")
	flag.BoolVar(&version, "v", false, "List version")
	flag.Parse()

	if version {
		fmt.Printf("BirdBot %s (%s)\n", app.Version, app.Build)
		return
	}

	bot := app.NewBot()

	if err := bot.Initialize(config_file); err != nil {
		log.Fatal("Failed to initialize: ", err)
	}

	if err := bot.Run(); err != nil {
		log.Fatal(err)
	}
}
