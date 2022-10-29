package discord

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/yeslayla/birdbot/core"
)

func (discord *Discord) MoveChannelToCategory(guildID string, categoryID string, channel *core.Channel) error {

	// Move to archive category
	if _, err := discord.session.ChannelEdit(channel.ID, &discordgo.ChannelEdit{
		ParentID: categoryID,
	}); err != nil {
		return fmt.Errorf("failed to move channel to archive category: %v", err)
	}

	return nil
}

func (discord *Discord) CreateChannelIfNotExists(guildID string, channel_name string) (*core.Channel, error) {

	// Grab channels to query
	channels, err := discord.session.GuildChannels(guildID)
	if err != nil {
		return nil, fmt.Errorf("failed to list channels when creating new channel: '%s': %v", channel_name, err)
	}
	for _, channel := range channels {

		// Found channel!
		if channel.Name == channel_name {
			log.Printf("Tried to create channel, but it already exists '%s'", channel_name)
			return &core.Channel{
				Name: channel.Name,
				ID:   channel.ID,
			}, nil
		}
	}

	// Since a channel was not found, create one
	channel, err := discord.session.GuildChannelCreate(guildID, channel_name, discordgo.ChannelTypeGuildText)
	if err != nil {
		return nil, fmt.Errorf("failed to created channel '%s': %v", channel_name, err)
	}

	log.Printf("Created channel: '%s'", channel_name)
	return &core.Channel{
		Name: channel.Name,
		ID:   channel.ID,
	}, nil
}

func (discord *Discord) DeleteChannel(guildID string, channel_name string) (bool, error) {

	channels, err := discord.session.GuildChannels(guildID)
	if err != nil {
		return false, fmt.Errorf("failed to list channels when deleting channel: '%s': %v", channel_name, err)
	}
	for _, channel := range channels {
		if channel.Name == channel_name {
			_, err = discord.session.ChannelDelete(channel.ID)
			if err != nil {
				return false, fmt.Errorf("failed to delete channel: %v", err)
			}
			return true, nil
		}
	}

	log.Printf("Tried to delete channel, but it didn't exist '%s'", channel_name)
	return false, nil
}

func (discord *Discord) GetChannelID(guildID string, channel_name string) (string, error) {
	channels, err := discord.session.GuildChannels(guildID)
	if err != nil {
		return "", fmt.Errorf("failed to list channels when getting channel id: '%s': %v", channel_name, err)
	}
	for _, channel := range channels {
		if channel.Name == channel_name {
			return channel.ID, nil
		}
	}

	return "", fmt.Errorf("failed to get channel id for '%s': channel not found", channel_name)
}

func (discord *Discord) SendMessage(channel *core.Channel, message string) error {
	_, err := discord.session.ChannelMessageSend(channel.ID, message)
	return err
}

func (discord *Discord) NewChannelFromID(ID string) *core.Channel {
	channel, err := discord.session.Channel(ID)
	if err != nil {
		return nil
	}

	return &core.Channel{
		ID:   ID,
		Name: channel.Name,
	}
}
