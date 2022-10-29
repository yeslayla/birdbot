package app

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/yeslayla/birdbot/core"
	"github.com/yeslayla/birdbot/discord"
)

type Bot struct {
	session *discordgo.Session
	d       *discord.Discord

	// Discord Objects
	guildID               string
	eventCategoryID       string
	archiveCategoryID     string
	notificationChannelID string

	// Signal for shutdown
	stop chan os.Signal
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

	// Create Discord Session
	app.session, err = discordgo.New(fmt.Sprint("Bot ", cfg.Discord.Token))
	if err != nil {
		return fmt.Errorf("failed to create Discord session: %v", err)
	}
	app.d = discord.NewDiscord(app.session)

	// Register Event Handlers
	app.session.AddHandler(app.onReady)
	app.session.AddHandler(app.onEventCreate)
	app.session.AddHandler(app.onEventDelete)
	app.session.AddHandler(app.onEventUpdate)

	return nil
}

// Run opens the session with Discord until exit
func (app *Bot) Run() error {

	if err := app.session.Open(); err != nil {
		return fmt.Errorf("failed to open Discord session: %v", err)
	}
	defer app.session.Close()

	// Keep alive
	app.stop = make(chan os.Signal, 1)
	signal.Notify(app.stop, os.Interrupt)
	<-app.stop
	return nil
}

// Stop triggers a graceful shutdown of the app
func (app *Bot) Stop() {
	log.Print("Shuting down...")
	app.stop <- os.Kill
}

// Notify sends a message to the notification channe;
func (app *Bot) Notify(message string) {
	if app.notificationChannelID == "" {
		log.Println(message)
		return
	}

	log.Println("Notification: ", message)

	channel := app.d.NewChannelFromID(app.notificationChannelID)
	if channel == nil {
		log.Printf("Failed notification: channel was not found with ID '%v'", app.notificationChannelID)
	}

	err := app.d.SendMessage(channel, message)
	if err != nil {
		log.Print("Failed notification: ", err)
	}
}

func (app *Bot) onReady(s *discordgo.Session, r *discordgo.Ready) {
	app.Notify("BirdBot is ready!")
}

func (app *Bot) onEventCreate(s *discordgo.Session, r *discordgo.GuildScheduledEventCreate) {
	if r.GuildID != app.guildID {
		return
	}

	event := discord.NewEvent(r.GuildScheduledEvent)
	log.Print("Event Created: '", event.Name, "':'", event.Location, "'")

	channel, err := app.d.CreateChannelIfNotExists(app.guildID, event.Channel().Name)
	if err != nil {
		log.Print("Failed to create channel for event: ", err)
	}

	if app.eventCategoryID != "" {
		if _, err = s.ChannelEdit(channel.ID, &discordgo.ChannelEdit{
			ParentID: app.eventCategoryID,
		}); err != nil {
			log.Printf("Failed to move channel to events category '%s': %v", channel.Name, err)
		}
	}

	eventURL := fmt.Sprintf("https://discordapp.com/events/%s/%s", app.guildID, r.ID)
	app.Notify(fmt.Sprintf("%s is organizing an event '%s': %s", event.Organizer.Mention(), event.Name, eventURL))
}

func (app *Bot) onEventDelete(s *discordgo.Session, r *discordgo.GuildScheduledEventDelete) {
	if r.GuildID != app.guildID {
		return
	}

	// Create Event Object
	event := discord.NewEvent(r.GuildScheduledEvent)

	_, err := app.d.DeleteChannel(app.guildID, event.Channel().Name)
	if err != nil {
		log.Print("Failed to create channel for event: ", err)
	}

	app.Notify(fmt.Sprintf("%s cancelled '%s' on %s, %d!", event.Organizer.Mention(), event.Name, event.DateTime.Month().String(), event.DateTime.Day()))
}

func (app *Bot) onEventUpdate(s *discordgo.Session, r *discordgo.GuildScheduledEventUpdate) {
	if r.GuildID != app.guildID {
		return
	}

	// Create Event Object
	event := discord.NewEvent(r.GuildScheduledEvent)

	// Pass event onwards
	switch r.Status {
	case discordgo.GuildScheduledEventStatusCompleted:
		app.onEventComplete(s, event)
	}
}

func (app *Bot) onEventComplete(s *discordgo.Session, event *core.Event) {

	channel := event.Channel()

	if app.archiveCategoryID != "" {

		if err := app.d.MoveChannelToCategory(app.guildID, app.archiveCategoryID, channel); err != nil {
			log.Print("Failed to archive channel: ", err)
		}
		log.Printf("Archived channel: '%s'", channel.Name)

	} else {

		// Delete Channel
		_, err := app.d.DeleteChannel(app.guildID, channel.Name)
		if err != nil {
			log.Print("Failed to delete channel: ", err)
		}

		log.Printf("Deleted channel: '%s'", channel.Name)
	}
}

func NewBot() *Bot {
	return &Bot{}
}
