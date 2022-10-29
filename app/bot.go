package app

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/yeslayla/birdbot/core"
	"github.com/yeslayla/birdbot/discord"
)

var Version string
var Build string

type Bot struct {
	session *discord.Discord

	// Discord Objects
	guildID               string
	eventCategoryID       string
	archiveCategoryID     string
	notificationChannelID string
}

// Initalize creates the discord session and registers handlers
func (app *Bot) Initialize(config_path string) error {
	log.Printf("Using config: %s", config_path)
	cfg := &core.Config{}

	_, err := os.Stat(config_path)
	if errors.Is(err, os.ErrNotExist) {
		log.Printf("Config file not found: '%s'", config_path)
		err := cleanenv.ReadEnv(cfg)
		if err != nil {
			return err
		}
	} else {
		err := cleanenv.ReadConfig(config_path, cfg)
		if err != nil {
			return err
		}
	}
	// Load directly from config
	app.guildID = cfg.Discord.GuildID
	app.eventCategoryID = cfg.Discord.EventCategory
	app.archiveCategoryID = cfg.Discord.ArchiveCategory
	app.notificationChannelID = cfg.Discord.NotificationChannel

	if app.guildID == "" {
		return fmt.Errorf("discord Guild ID is not set")
	}

	app.session = discord.New(app.guildID, cfg.Discord.Token)

	// Register Event Handlers
	app.session.OnReady(app.onReady)
	app.session.OnEventCreate(app.onEventCreate)
	app.session.OnEventDelete(app.onEventDelete)
	app.session.OnEventUpdate(app.onEventUpdate)

	return nil
}

// Run opens the session with Discord until exit
func (app *Bot) Run() error {
	return app.session.Run()
}

// Stop triggers a graceful shutdown of the app
func (app *Bot) Stop() {
	log.Print("Shuting down...")
	app.session.Stop()
}

// Notify sends a message to the notification channe;
func (app *Bot) Notify(message string) {
	if app.notificationChannelID == "" {
		log.Println(message)
		return
	}

	log.Print("Notification: ", message)

	channel := app.session.NewChannelFromID(app.notificationChannelID)
	if channel == nil {
		log.Printf("Failed notification: channel was not found with ID '%v'", app.notificationChannelID)
	}

	err := app.session.SendMessage(channel, message)
	if err != nil {
		log.Print("Failed notification: ", err)
	}
}

func (app *Bot) onReady(d *discord.Discord) {
	app.Notify(fmt.Sprintf("BirdBot %s is ready!", Version))
}

func (app *Bot) onEventCreate(d *discord.Discord, event *core.Event) {

	log.Print("Event Created: '", event.Name, "':'", event.Location, "'")

	channel, err := app.session.NewChannelFromName(event.Channel().Name)
	if err != nil {
		log.Print("Failed to create channel for event: ", err)
	}

	if app.eventCategoryID != "" {
		err = app.session.MoveChannelToCategory(channel, app.eventCategoryID)
		if err != nil {
			log.Printf("Failed to move channel to events category '%s': %v", channel.Name, err)
		}
	}

	eventURL := fmt.Sprintf("https://discordapp.com/events/%s/%s", app.guildID, event.ID)
	app.Notify(fmt.Sprintf("%s is organizing an event '%s': %s", event.Organizer.Mention(), event.Name, eventURL))
}

func (app *Bot) onEventDelete(d *discord.Discord, event *core.Event) {

	_, err := app.session.DeleteChannel(event.Channel())
	if err != nil {
		log.Print("Failed to create channel for event: ", err)
	}

	app.Notify(fmt.Sprintf("%s cancelled '%s' on %s, %d!", event.Organizer.Mention(), event.Name, event.DateTime.Month().String(), event.DateTime.Day()))
}

func (app *Bot) onEventUpdate(d *discord.Discord, event *core.Event) {
	// Pass event onwards
	if event.Completed {
		app.onEventComplete(d, event)
	}
}

func (app *Bot) onEventComplete(d *discord.Discord, event *core.Event) {

	channel := event.Channel()

	if app.archiveCategoryID != "" {

		if err := app.session.MoveChannelToCategory(channel, app.archiveCategoryID); err != nil {
			log.Print("Failed to move channel to archive category: ", err)
		}

		if err := app.session.ArchiveChannel(channel); err != nil {
			log.Print("Failed to archive channel: ", err)
		}

		log.Printf("Archived channel: '%s'", channel.Name)

	} else {

		// Delete Channel
		_, err := app.session.DeleteChannel(channel)
		if err != nil {
			log.Print("Failed to delete channel: ", err)
		}

		log.Printf("Deleted channel: '%s'", channel.Name)
	}
}

func NewBot() *Bot {
	return &Bot{}
}
