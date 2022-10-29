package discord

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/yeslayla/birdbot/core"
)

// NewChannelFromID creates a channel from a given ID
func (discord *Discord) NewChannelFromID(ID string) *core.Channel {
	channel, err := discord.session.Channel(ID)
	if err != nil {
		return nil
	}

	return &core.Channel{
		ID:       ID,
		Name:     channel.Name,
		Verified: true,
	}
}

// NewChannelFromName creates a channel object with its name
func (discord *Discord) NewChannelFromName(channel_name string) (*core.Channel, error) {

	// Grab channels to query
	channels, err := discord.session.GuildChannels(discord.guildID)
	if err != nil {
		return nil, fmt.Errorf("failed to list channels when creating new channel: '%s': %v", channel_name, err)
	}
	for _, channel := range channels {

		// Found channel!
		if channel.Name == channel_name {
			log.Printf("Tried to create channel, but it already exists '%s'", channel_name)
			return &core.Channel{
				Name:     channel.Name,
				ID:       channel.ID,
				Verified: true,
			}, nil
		}
	}

	// Since a channel was not found, create one
	channel, err := discord.session.GuildChannelCreate(discord.guildID, channel_name, discordgo.ChannelTypeGuildText)
	if err != nil {
		return nil, fmt.Errorf("failed to created channel '%s': %v", channel_name, err)
	}

	log.Printf("Created channel: '%s'", channel_name)
	return &core.Channel{
		Name: channel.Name,
		ID:   channel.ID,
	}, nil
}

// DeleteChannel deletes a channel
func (discord *Discord) DeleteChannel(channel *core.Channel) (bool, error) {
	if channel.Verified == false {
		return false, fmt.Errorf("failed to delete channel: given channel object is not verified")
	}

	_, err := discord.session.ChannelDelete(channel.ID)
	if err != nil {
		return false, fmt.Errorf("failed to delete channel: %v", err)
	}
	return true, nil
}

// getChannelID returns a channel ID from its name
func (discord *Discord) getChannelID(channel_name string) (string, error) {

	// Get list of all channels
	channels, err := discord.session.GuildChannels(discord.guildID)
	if err != nil {
		return "", fmt.Errorf("failed to list channels when getting channel id: '%s': %v", channel_name, err)
	}

	// Loop through to find channel
	for _, ch := range channels {

		// Find and return ID!
		if ch.Name == channel_name {
			return ch.ID, nil
		}
	}

	return "", fmt.Errorf("failed to get channel id for '%s': channel not found", channel_name)
}

// SendMessage sends a message to a given channel
func (discord *Discord) SendMessage(channel *core.Channel, message string) error {
	if channel.Verified == false {
		return fmt.Errorf("failed to delete channel: given channel object is not verified")
	}

	_, err := discord.session.ChannelMessageSend(channel.ID, message)
	return err
}

// MoveChannelToCategory places a channel in a given category
func (discord *Discord) MoveChannelToCategory(channel *core.Channel, categoryID string) error {
	if channel.Verified == false {
		return fmt.Errorf("failed to delete channel: given channel object is not verified")
	}

	// Move to archive category
	if _, err := discord.session.ChannelEdit(channel.ID, &discordgo.ChannelEdit{
		ParentID: categoryID,
	}); err != nil {
		return fmt.Errorf("failed to move channel to archive category: %v", err)
	}

	return nil
}
