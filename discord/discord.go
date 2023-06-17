package discord

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/mock"
	"github.com/yeslayla/birdbot/common"
	"github.com/yeslayla/birdbot/persistence"
)

type Discord struct {
	mock.Mock

	guildID       string
	applicationID string
	session       *discordgo.Session

	commands        map[string]*discordgo.ApplicationCommand
	commandHandlers map[string]func(session *discordgo.Session, i *discordgo.InteractionCreate)

	db persistence.Database

	// Signal for shutdown
	stop chan os.Signal
}

// New creates a new Discord session
func New(applicationID string, guildID string, token string, db persistence.Database) *Discord {

	// Create Discord Session
	session, err := discordgo.New(fmt.Sprint("Bot ", token))
	if err != nil {
		log.Fatalf("Failed to create Discord session: %v", err)
	}
	session.ShouldReconnectOnError = true
	return &Discord{
		db: db,

		session:         session,
		applicationID:   applicationID,
		guildID:         guildID,
		commands:        make(map[string]*discordgo.ApplicationCommand),
		commandHandlers: make(map[string]func(*discordgo.Session, *discordgo.InteractionCreate)),
	}
}

// Run opens the Discod session until exit
func (discord *Discord) Run() error {

	if err := discord.session.Open(); err != nil {
		return fmt.Errorf("failed to open Discord session: %v", err)
	}
	defer discord.session.Close()

	// Register command handler
	discord.session.AddHandler(func(session *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.GuildID != discord.guildID {
			return
		}

		if handler, ok := discord.commandHandlers[i.ApplicationCommandData().Name]; ok {
			handler(session, i)
		}
	})

	// Keep alive
	discord.stop = make(chan os.Signal, 1)
	signal.Notify(discord.stop, os.Interrupt)
	<-discord.stop

	discord.ClearCommands()

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
func (discord *Discord) OnEventCreate(handler func(*Discord, common.Event)) {
	discord.session.AddHandler(func(s *discordgo.Session, r *discordgo.GuildScheduledEventCreate) {
		if r.GuildID != discord.guildID {
			return
		}
		event := NewEvent(r.GuildScheduledEvent)
		handler(discord, event)
	})
}

// OnEventDelete registers a handler when a guild scheduled event is deleted
func (discord *Discord) OnEventDelete(handler func(*Discord, common.Event)) {
	discord.session.AddHandler(func(s *discordgo.Session, r *discordgo.GuildScheduledEventDelete) {
		if r.GuildID != discord.guildID {
			return
		}
		event := NewEvent(r.GuildScheduledEvent)
		handler(discord, event)
	})
}

// OnEventUpdate registers a handler when a guild scheduled event is updated
func (discord *Discord) OnEventUpdate(handler func(*Discord, common.Event)) {
	discord.session.AddHandler(func(s *discordgo.Session, r *discordgo.GuildScheduledEventUpdate) {
		if r.GuildID != discord.guildID {
			return
		}
		event := NewEvent(r.GuildScheduledEvent)
		handler(discord, event)
	})
}

// OnMessageRecieved registers a handler when a message is recieved
func (discord *Discord) OnMessageRecieved(handler func(*Discord, string, common.User, string)) {
	discord.session.AddHandler(func(s *discordgo.Session, r *discordgo.MessageCreate) {
		if r.GuildID != discord.guildID {
			return
		}

		handler(discord, r.ChannelID, NewUser(r.Author), r.Content)
	})
}

func (discord *Discord) SetStatus(status string) {
	if err := discord.session.UpdateGameStatus(0, status); err != nil {
		log.Fatal("Failed to update status: ", err)
	}
}
