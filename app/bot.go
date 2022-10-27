package app

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/ilyakaznacheev/cleanenv"
)

type Bot struct {
	discord *discordgo.Session

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
	cfg := &Config{}

	_, err := os.Stat(config_path)
	if errors.Is(err, os.ErrNotExist) {
		log.Printf("Config file not found: '%s'", config_path)
		err := cleanenv.ReadEnv(&cfg)
		if err != nil {
			return nil
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
	app.discord, err = discordgo.New(fmt.Sprint("Bot ", cfg.Discord.Token))
	if err != nil {
		return fmt.Errorf("failed to create Discord session: %v", err)
	}

	// Register Event Handlers
	app.discord.AddHandler(app.onReady)
	app.discord.AddHandler(app.onEventCreate)
	app.discord.AddHandler(app.onEventDelete)
	app.discord.AddHandler(app.onEventUpdate)

	return nil
}

// Run opens the session with Discord until exit
func (app *Bot) Run() error {

	if err := app.discord.Open(); err != nil {
		return fmt.Errorf("failed to open Discord session: %v", err)
	}
	defer app.discord.Close()

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
		return
	}

	_, err := app.discord.ChannelMessageSend(app.notificationChannelID, message)

	log.Println("Notification: ", message)
	if err != nil {
		log.Println("Failed notification: ", err)
	}
}

func (app *Bot) onReady(s *discordgo.Session, r *discordgo.Ready) {
	app.Notify("BirdBot is ready!")
	log.Print("BirdBot is ready!")
}

func (app *Bot) onEventCreate(s *discordgo.Session, r *discordgo.GuildScheduledEventCreate) {
	event := &Event{}
	event.Name = r.Name
	event.OrganizerID = r.CreatorID
	event.DateTime = r.ScheduledStartTime
	if r.EntityType != discordgo.GuildScheduledEventEntityTypeExternal {
		event.Location = REMOTE_LOCATION
	} else {
		event.Location = r.EntityMetadata.Location
	}
	log.Print("Event Created: '", event.Name, "':'", event.Location, "'")

	channel, err := CreateChannelIfNotExists(s, app.guildID, event.GetChannelName())
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
	app.Notify(fmt.Sprintf("<@%s> is organizing an event '%s': %s", event.OrganizerID, event.Name, eventURL))
}

func (app *Bot) onEventDelete(s *discordgo.Session, r *discordgo.GuildScheduledEventDelete) {

	// Create Event Object
	event := &Event{}
	event.Name = r.Name
	event.OrganizerID = r.CreatorID
	event.DateTime = r.ScheduledStartTime
	if r.EntityType != discordgo.GuildScheduledEventEntityTypeExternal {
		event.Location = REMOTE_LOCATION
	} else {
		event.Location = r.EntityMetadata.Location
	}

	_, err := DeleteChannel(app.discord, app.guildID, event.GetChannelName())
	if err != nil {
		log.Print("Failed to create channel for event: ", err)
	}

	app.Notify(fmt.Sprintf("<@%s> cancelled '%s' on %s, %d!", event.OrganizerID, event.Name, event.DateTime.Month().String(), event.DateTime.Day()))
}

func (app *Bot) onEventUpdate(s *discordgo.Session, r *discordgo.GuildScheduledEventUpdate) {

	// Create Event Object
	event := &Event{}
	event.Name = r.Name
	event.OrganizerID = r.CreatorID
	event.DateTime = r.ScheduledStartTime
	if r.EntityType != discordgo.GuildScheduledEventEntityTypeExternal {
		event.Location = REMOTE_LOCATION
	} else {
		event.Location = r.EntityMetadata.Location
	}

	// Pass event onwards
	switch r.Status {
	case discordgo.GuildScheduledEventStatusCompleted:
		app.onEventComplete(s, event)
	}
}

func (app *Bot) onEventComplete(s *discordgo.Session, event *Event) {
	channel_name := event.GetChannelName()

	if app.archiveCategoryID != "" {

		// Get Channel ID
		id, err := GetChannelID(s, app.guildID, channel_name)
		if err != nil {
			log.Printf("Failed to archive channel: %v", err)
			return
		}

		// Move to archive category
		if _, err := s.ChannelEdit(id, &discordgo.ChannelEdit{
			ParentID: app.archiveCategoryID,
		}); err != nil {
			log.Printf("Failed to move channel to archive category: %v", err)
			return
		}

		log.Printf("Archived channel: '%s'", channel_name)
	} else {

		// Delete Channel
		_, err := DeleteChannel(s, app.guildID, channel_name)
		if err != nil {
			log.Print("Failed to delete channel: ", err)
		}

		log.Printf("Deleted channel: '%s'", channel_name)
	}
}

func NewBot() *Bot {
	return &Bot{}
}
