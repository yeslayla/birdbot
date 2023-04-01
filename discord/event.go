package discord

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/yeslayla/birdbot/common"
	"github.com/yeslayla/birdbot/core"
)

// NewEvent converts a discordgo.GuildScheduledEvent to birdbot event
func NewEvent(guildEvent *discordgo.GuildScheduledEvent) common.Event {
	event := common.Event{
		Name:        guildEvent.Name,
		Description: guildEvent.Description,
		ID:          guildEvent.ID,
		Organizer: common.User{
			ID: guildEvent.CreatorID,
		},
		DateTime: guildEvent.ScheduledStartTime,
		ImageURL: guildEvent.Image,
	}

	if guildEvent.ScheduledEndTime != nil {
		event.CompleteDateTime = *guildEvent.ScheduledEndTime
	} else {
		year, month, day := guildEvent.ScheduledStartTime.Date()
		event.CompleteDateTime = time.Date(year, month, day, 0, 0, 0, 0, guildEvent.ScheduledStartTime.Location())
	}

	event.Completed = guildEvent.Status == discordgo.GuildScheduledEventStatusCompleted

	if guildEvent.EntityType != discordgo.GuildScheduledEventEntityTypeExternal {
		event.Location = core.RemoteLocation
	} else {
		event.Location = guildEvent.EntityMetadata.Location
	}

	return event
}

// CreateEvent creates a new discord event
func (discord *Discord) CreateEvent(event common.Event) error {

	params := &discordgo.GuildScheduledEventParams{
		Name:               event.Name,
		Description:        event.Description,
		ScheduledStartTime: &event.DateTime,
		ScheduledEndTime:   &event.CompleteDateTime,
		Image:              event.ImageURL,
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
