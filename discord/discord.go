package discord

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/mock"
	"github.com/yeslayla/birdbot/core"
)

type Discord struct {
	mock.Mock

	guildID string
	session *discordgo.Session

	// Signal for shutdown
	stop chan os.Signal
}

// New creates a new Discord session
func New(guildID string, token string) *Discord {

	// Create Discord Session
	session, err := discordgo.New(fmt.Sprint("Bot ", token))
	if err != nil {
		log.Fatalf("Failed to create Discord session: %v", err)
	}

	return &Discord{
		session: session,
		guildID: guildID,
	}
}

// Run opens the Discod session until exit
func (discord *Discord) Run() error {

	if err := discord.session.Open(); err != nil {
		return fmt.Errorf("failed to open Discord session: %v", err)
	}
	defer discord.session.Close()

	// Keep alive
	discord.stop = make(chan os.Signal, 1)
	signal.Notify(discord.stop, os.Interrupt)
	<-discord.stop
	return nil
}

// Stop tells the Discord session to exit
func (discord *Discord) Stop() {
	discord.stop <- os.Kill
}

// OnReady registers a handler for when the Discord session is ready
func (discord *Discord) OnReady(handler func(*Discord)) {
	discord.session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		handler(discord)
	})
}

// OnEventCreate registers a handler when a guild scheduled event is created
func (discord *Discord) OnEventCreate(handler func(*Discord, *core.Event)) {
	discord.session.AddHandler(func(s *discordgo.Session, r *discordgo.GuildScheduledEventCreate) {
		if r.GuildID != discord.guildID {
			return
		}
		event := NewEvent(r.GuildScheduledEvent)
		handler(discord, event)
	})
}

// OnEventDelete registers a handler when a guild scheduled event is deleted
func (discord *Discord) OnEventDelete(handler func(*Discord, *core.Event)) {
	discord.session.AddHandler(func(s *discordgo.Session, r *discordgo.GuildScheduledEventDelete) {
		if r.GuildID != discord.guildID {
			return
		}
		event := NewEvent(r.GuildScheduledEvent)
		handler(discord, event)
	})
}

// OnEventUpdate registers a handler when a guild scheduled event is updated
func (discord *Discord) OnEventUpdate(handler func(*Discord, *core.Event)) {
	discord.session.AddHandler(func(s *discordgo.Session, r *discordgo.GuildScheduledEventUpdate) {
		if r.GuildID != discord.guildID {
			return
		}
		event := NewEvent(r.GuildScheduledEvent)
		handler(discord, event)
	})
}

func (discord *Discord) SetStatus(status string) {
	if err := discord.session.UpdateGameStatus(0, status); err != nil {
		log.Fatal("Failed to update status: ", err)
	}
}
