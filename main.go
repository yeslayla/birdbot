package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/yeslayla/birdbot/app"
	"github.com/yeslayla/birdbot/core"
	"github.com/yeslayla/birdbot/modules"
)

const PluginsDirectory = "./plugins"

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

	log.Printf("Using config: %s", config_file)
	cfg := &core.Config{}

	_, err := os.Stat(config_file)
	if errors.Is(err, os.ErrNotExist) {
		log.Printf("Config file not found: '%s'", config_file)
		err := cleanenv.ReadEnv(cfg)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err := cleanenv.ReadConfig(config_file, cfg)
		if err != nil {
			log.Fatal(err)
		}
	}

	bot := app.NewBot()

	if err := bot.Initialize(cfg); err != nil {
		log.Fatal("Failed to initialize: ", err)
	}

	loader := app.NewComponentLoader(bot)

	if cfg.Features.AnnounceEvents.IsEnabledByDefault() {
		loader.LoadComponent(modules.NewAnnounceEventsComponent(bot.Mastodon, cfg.Discord.NotificationChannel))
	}
	if cfg.Features.ManageEventChannels.IsEnabledByDefault() {
		loader.LoadComponent(modules.NewManageEventChannelsComponent(cfg.Discord.EventCategory, cfg.Discord.ArchiveCategory, bot.Session))
	}
	if cfg.Features.ReccurringEvents.IsEnabledByDefault() {
		loader.LoadComponent(modules.NewRecurringEventsComponent())
	}

	if cfg.Features.RoleSelection.IsEnabledByDefault() {
		for _, v := range cfg.Discord.RoleSelections {
			loader.LoadComponent(modules.NewRoleSelectionComponent(bot.Session, v))
		}
	}

	if _, err := os.Stat(PluginsDirectory); !os.IsNotExist(err) {
		components := app.LoadPlugins(PluginsDirectory)
		for _, comp := range components {
			loader.LoadComponent(comp)
		}
	}

	if err := bot.Run(); err != nil {
		log.Fatal(err)
	}
}
