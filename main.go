package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/yeslayla/birdbot/app"
)

func main() {
	var config_file string
	var version bool
	flag.StringVar(&config_file, "c", "birdbot.yaml", "Path to config file")
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
