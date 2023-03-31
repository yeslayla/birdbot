package components

import (
	"fmt"

	"github.com/yeslayla/birdbot/common"
	"github.com/yeslayla/birdbot/mastodon"
)

type announceEventsComponent struct {
	bot      common.ComponentManager
	mastodon *mastodon.Mastodon
	guildID  string
}

// NewAnnounceEventsComponent creates a new component
func NewAnnounceEventsComponent(mastodon *mastodon.Mastodon, guildID string) common.Component {
	return &announceEventsComponent{
		mastodon: mastodon,
		guildID:  guildID,
	}
}

// Initialize registers event listeners
func (c *announceEventsComponent) Initialize(birdbot common.ComponentManager) error {
	c.bot = birdbot

	_ = birdbot.OnEventCreate(c.OnEventCreate)
	_ = birdbot.OnEventDelete(c.OnEventDelete)

	return nil
}

// OnEventCreate notifies about the event creation to given providers
func (c *announceEventsComponent) OnEventCreate(e common.Event) error {
	eventURL := fmt.Sprintf("https://discordapp.com/events/%s/%s", c.guildID, e.ID)
	c.bot.Notify(fmt.Sprintf("%s is organizing an event '%s': %s", e.Organizer.DiscordMention(), e.Name, eventURL))

	// Toot an announcement if Mastodon is configured
	if c.mastodon != nil {
		err := c.mastodon.Toot(fmt.Sprintf("A new event has been organized '%s': %s", e.Name, eventURL))
		if err != nil {
			fmt.Println("Failed to send Mastodon Toot:", err)
		}
	}

	return nil
}

func (c *announceEventsComponent) OnEventDelete(e common.Event) error {
	_ = c.bot.Notify(fmt.Sprintf("%s cancelled '%s' on %s, %d!", e.Organizer.DiscordMention(), e.Name, e.DateTime.Month().String(), e.DateTime.Day()))

	if c.mastodon != nil {
		err := c.mastodon.Toot(fmt.Sprintf("'%s' cancelled on %s, %d!", e.Name, e.DateTime.Month().String(), e.DateTime.Day()))
		if err != nil {
			fmt.Println("Failed to send Mastodon Toot:", err)
		}
	}

	return nil
}
