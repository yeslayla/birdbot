package app

import (
	"fmt"
	"log"

	"github.com/yeslayla/birdbot/common"
	"github.com/yeslayla/birdbot/core"
	"github.com/yeslayla/birdbot/discord"
	"github.com/yeslayla/birdbot/mastodon"
	"github.com/yeslayla/birdbot/persistence"
)

var Version string
var Build string

type Bot struct {
	Session  *discord.Discord
	Mastodon *mastodon.Mastodon

	Database persistence.Database

	// Discord Objects
	guildID               string
	eventCategoryID       string
	archiveCategoryID     string
	notificationChannelID string

	onReadyHandlers  [](func() error)
	onNotifyHandlers [](func(string) error)

	onEventCreatedHandlers   [](func(common.Event) error)
	onEventDeletedHandlers   [](func(common.Event) error)
	onEventUpdatedHandlers   [](func(common.Event) error)
	onEventCompletedHandlers [](func(common.Event) error)

	channelChats map[string][]common.ExternalChatModule
}

// Initalize creates the discord session and registers handlers
func (app *Bot) Initialize(cfg *core.Config) error {

	// Load directly from config
	app.guildID = cfg.Discord.GuildID
	app.eventCategoryID = cfg.Discord.EventCategory
	app.archiveCategoryID = cfg.Discord.ArchiveCategory
	app.notificationChannelID = cfg.Discord.NotificationChannel

	if app.guildID == "" {
		return fmt.Errorf("discord Guild ID is not set")
	}
	if cfg.Discord.ApplicationID == "" {
		return fmt.Errorf("discord Application ID is not set")
	}

	if cfg.Mastodon.ClientID != "" && cfg.Mastodon.ClientSecret != "" &&
		cfg.Mastodon.Username != "" && cfg.Mastodon.Password != "" &&
		cfg.Mastodon.Server != "" {
		app.Mastodon = mastodon.NewMastodon(cfg.Mastodon.Server, cfg.Mastodon.ClientID, cfg.Mastodon.ClientSecret,
			cfg.Mastodon.Username, cfg.Mastodon.Password)
	}

	app.Session = discord.New(cfg.Discord.ApplicationID, app.guildID, cfg.Discord.Token, app.Database)

	// Intialize submodules
	for channelID, chats := range app.channelChats {
		channel := app.Session.NewChannelFromID(channelID)
		for _, chat := range chats {
			app.InitalizeExternalChat(channel, chat)
		}
	}

	// Register Event Handlers
	app.Session.OnReady(app.onReady)
	app.Session.OnEventCreate(app.onEventCreate)
	app.Session.OnEventDelete(app.onEventDelete)
	app.Session.OnEventUpdate(app.onEventUpdate)

	if len(app.channelChats) > 0 {
		app.Session.OnMessageRecieved(app.onMessageRecieved)
	}

	return nil
}

// Run opens the session with Discord until exit
func (app *Bot) Run() error {
	return app.Session.Run()
}

// Stop triggers a graceful shutdown of the app
func (app *Bot) Stop() {
	log.Print("Shuting down...")
	app.Session.Stop()
}

// Notify sends a message to the notification channe;
func (app *Bot) Notify(message string) {
	if app.notificationChannelID == "" {
		log.Println(message)
		return
	}

	log.Print("Notification: ", message)

	channel := app.Session.NewChannelFromID(app.notificationChannelID)
	if channel == nil {
		log.Printf("Failed notification: channel was not found with ID '%v'", app.notificationChannelID)
	}

	err := app.Session.SendMessage(channel, message)
	if err != nil {
		log.Print("Failed notification: ", err)
	}

	for _, handler := range app.onNotifyHandlers {
		if err := handler(message); err != nil {
			log.Println(err)
		}
	}
}

func (app *Bot) onReady(d *discord.Discord) {
	app.Session.SetStatus(fmt.Sprintf("with fire! (%s)", Version))

	for _, handler := range app.onReadyHandlers {
		if err := handler(); err != nil {
			log.Println(err)
		}
	}
}

func (app *Bot) onEventCreate(d *discord.Discord, event common.Event) {

	log.Print("Event Created: '", event.Name, "':'", event.Location, "'")
	for _, handler := range app.onEventCreatedHandlers {
		if err := handler(event); err != nil {
			log.Println(err)
		}
	}

}

func (app *Bot) onEventDelete(d *discord.Discord, event common.Event) {

	for _, handler := range app.onEventDeletedHandlers {
		if err := handler(event); err != nil {
			log.Println(err)
		}
	}

}

func (app *Bot) onEventUpdate(d *discord.Discord, event common.Event) {

	for _, handler := range app.onEventUpdatedHandlers {
		if err := handler(event); err != nil {
			log.Println(err)
		}
	}

	// Pass event onwards
	if event.Completed {
		app.onEventComplete(d, event)
	}
}

func (app *Bot) onEventComplete(d *discord.Discord, event common.Event) {

	for _, handler := range app.onEventCompletedHandlers {
		if err := handler(event); err != nil {
			log.Println(err)
		}
	}

}

func (app *Bot) onMessageRecieved(d *discord.Discord, channelID string, user common.User, message string) {
	chats, ok := app.channelChats[channelID]
	if !ok {
		return
	}

	for _, chat := range chats {
		chat.RecieveMessage(user, message)
	}
}

// NewBot creates a new bot instance
func NewBot(db persistence.Database) *Bot {
	return &Bot{
		Database:     db,
		channelChats: make(map[string][]common.ExternalChatModule),
	}
}
