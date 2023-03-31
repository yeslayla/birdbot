package events

import (
	"log"

	"github.com/yeslayla/birdbot/common"
	"github.com/yeslayla/birdbot/core"
	"github.com/yeslayla/birdbot/discord"
)

type ManageEventChannelsComponent struct {
	session    *discord.Discord
	categoryID string
}

func NewManageEventChannelsComponent(categoryID string, session *discord.Discord) *ManageEventChannelsComponent {
	return &ManageEventChannelsComponent{
		session:    session,
		categoryID: categoryID,
	}
}

func (c *ManageEventChannelsComponent) Initialize(birdbot common.ComponentManager) error {
	_ = birdbot.OnEventCreate(c.OnEventCreate)
	_ = birdbot.OnEventComplete(c.OnEventComplete)
	_ = birdbot.OnEventDelete(c.OnEventDelete)

	return nil
}

func (c *ManageEventChannelsComponent) OnEventCreate(e common.Event) error {
	channel, err := c.session.NewChannelFromName(core.GenerateChannel(e).Name)
	if err != nil {
		log.Print("Failed to create channel for event: ", err)
	}

	if c.categoryID != "" {
		err = c.session.MoveChannelToCategory(channel, c.categoryID)
		if err != nil {
			log.Printf("Failed to move channel to events category '%s': %v", channel.Name, err)
		}
	}
	return nil
}

func (c *ManageEventChannelsComponent) OnEventDelete(e common.Event) error {
	_, err := c.session.DeleteChannel(core.GenerateChannel(e))
	if err != nil {
		log.Print("Failed to create channel for event: ", err)
	}
	return nil
}

func (c *ManageEventChannelsComponent) OnEventComplete(e common.Event) error {
	channel := core.GenerateChannel(e)

	if c.categoryID != "" {

		if err := c.session.MoveChannelToCategory(channel, c.categoryID); err != nil {
			log.Print("Failed to move channel to archive category: ", err)
		}

		if err := c.session.ArchiveChannel(channel); err != nil {
			log.Print("Failed to archive channel: ", err)
		}

		log.Printf("Archived channel: '%s'", channel.Name)

	} else {

		// Delete Channel
		_, err := c.session.DeleteChannel(channel)
		if err != nil {
			log.Print("Failed to delete channel: ", err)
		}

		log.Printf("Deleted channel: '%s'", channel.Name)
	}

	return nil
}
