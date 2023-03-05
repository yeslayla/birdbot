package discord

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/yeslayla/birdbot/core"
)

// NewEvent converts a discordgo.GuildScheduledEvent to birdbot event
func NewEvent(guildEvent *discordgo.GuildScheduledEvent) *core.Event {
	event := &core.Event{
		Name:        guildEvent.Name,
		Description: guildEvent.Description,
		ID:          guildEvent.ID,
		Organizer: &core.User{
			ID: guildEvent.CreatorID,
		},
		DateTime: guildEvent.ScheduledStartTime,
		Image:    guildEvent.Image,
	}

	if guildEvent.ScheduledEndTime != nil {
		event.CompleteTime = *guildEvent.ScheduledEndTime
	} else {
		year, month, day := guildEvent.ScheduledStartTime.Date()
		event.CompleteTime = time.Date(year, month, day, 0, 0, 0, 0, guildEvent.ScheduledStartTime.Location())
	}

	event.Completed = guildEvent.Status == discordgo.GuildScheduledEventStatusCompleted

	if guildEvent.EntityType != discordgo.GuildScheduledEventEntityTypeExternal {
		event.Location = core.REMOTE_LOCATION
	} else {
		event.Location = guildEvent.EntityMetadata.Location
	}

	return event
}

func (discord *Discord) CreateEvent(event *core.Event) error {

	params := &discordgo.GuildScheduledEventParams{
		Name:               event.Name,
		Description:        event.Description,
		ScheduledStartTime: &event.DateTime,
		ScheduledEndTime:   &event.CompleteTime,
		Image:              event.Image,
		EntityType:         discordgo.GuildScheduledEventEntityTypeExternal,
		PrivacyLevel:       discordgo.GuildScheduledEventPrivacyLevelGuildOnly,
	}

	if event.Location != "" {
		params.EntityMetadata = &discordgo.GuildScheduledEventEntityMetadata{
			Location: event.Location,
		}
	}

	_, err := discord.session.GuildScheduledEventCreate(discord.guildID, params)

	return err
}
