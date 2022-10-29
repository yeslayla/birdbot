package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/yeslayla/birdbot/core"
)

// NewEvent converts a discordgo.GuildScheduledEvent to birdbot event
func NewEvent(guildEvent *discordgo.GuildScheduledEvent) *core.Event {
	event := &core.Event{
		Name:      guildEvent.Name,
		ID:        guildEvent.ID,
		Organizer: NewUser(guildEvent.Creator),
		DateTime:  guildEvent.ScheduledStartTime,
	}

	event.Completed = guildEvent.Status == discordgo.GuildScheduledEventStatusCompleted

	if guildEvent.EntityType != discordgo.GuildScheduledEventEntityTypeExternal {
		event.Location = core.REMOTE_LOCATION
	} else {
		event.Location = guildEvent.EntityMetadata.Location
	}

	return event
}
