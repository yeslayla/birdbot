package main

import (
	"flag"
	"log"

	"github.com/yeslayla/bird-bot/app"
)

func main() {
	var config_file string
	flag.StringVar(&config_file, "c", "birdbot.yaml", "Path to config file")
	flag.Parse()

	bot := app.NewBot()
	if err := bot.Initialize(config_file); err != nil {
		log.Fatal("Failed to initialize: ", err)
	}

	if err := bot.Run(); err != nil {
		log.Fatal(err)
	}
}
