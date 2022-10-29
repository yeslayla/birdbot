package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/yeslayla/birdbot/core"
)

func NewEvent(guildEvent *discordgo.GuildScheduledEvent) *core.Event {
	event := &core.Event{
		Name:      guildEvent.Name,
		Organizer: NewUser(guildEvent.Creator),
		DateTime:  guildEvent.ScheduledStartTime,
	}
	if guildEvent.EntityType != discordgo.GuildScheduledEventEntityTypeExternal {
		event.Location = core.REMOTE_LOCATION
	} else {
		event.Location = guildEvent.EntityMetadata.Location
	}

	return event
}
