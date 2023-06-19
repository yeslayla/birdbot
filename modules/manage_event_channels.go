package modules

import (
	"log"

	"github.com/yeslayla/birdbot-common/common"
	"github.com/yeslayla/birdbot/core"
	"github.com/yeslayla/birdbot/discord"
)

type manageEventChannelsModule struct {
	session           *discord.Discord
	categoryID        string
	archiveCategoryID string
}

// NewManageEventChannelsComponent creates a new component
func NewManageEventChannelsComponent(categoryID string, archiveCategoryID string, session *discord.Discord) common.Module {
	return &manageEventChannelsModule{
		session:           session,
		categoryID:        categoryID,
		archiveCategoryID: archiveCategoryID,
	}
}

// Initialize registers event listeners
func (c *manageEventChannelsModule) Initialize(birdbot common.ModuleManager) error {
	_ = birdbot.OnEventCreate(c.OnEventCreate)
	_ = birdbot.OnEventComplete(c.OnEventComplete)
	_ = birdbot.OnEventDelete(c.OnEventDelete)

	return nil
}

// OnEventCreate creates a new channel for an event and moves it to a given category
func (c *manageEventChannelsModule) OnEventCreate(e common.Event) error {
	channel, err := c.session.NewChannelFromName(core.GenerateChannelFromEvent(e).Name)
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

// OnEventDelete deletes the channel associated with the given event
func (c *manageEventChannelsModule) OnEventDelete(e common.Event) error {
	_, err := c.session.DeleteChannel(core.GenerateChannelFromEvent(e))
	if err != nil {
		log.Print("Failed to create channel for event: ", err)
	}
	return nil
}

// OnEventComplete archives a given event channel if not given
// an archive category will delete the channel instead
func (c *manageEventChannelsModule) OnEventComplete(e common.Event) error {
	channel := core.GenerateChannelFromEvent(e)

	if c.archiveCategoryID != "" {

		if err := c.session.MoveChannelToCategory(channel, c.archiveCategoryID); err != nil {
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
