package app

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

func CreateChannelIfNotExists(discord *discordgo.Session, guildID string, channel_name string) (*discordgo.Channel, error) {

	// Grab channels to query
	channels, err := discord.GuildChannels(guildID)
	if err != nil {
		return nil, fmt.Errorf("failed to list channels when creating new channel: '%s': %v", channel_name, err)
	}
	for _, channel := range channels {

		// Found channel!
		if channel.Name == channel_name {
			log.Printf("Tried to create channel, but it already exists '%s'", channel_name)
			return channel, nil
		}
	}

	// Since a channel was not found, create one
	channel, err := discord.GuildChannelCreate(guildID, channel_name, discordgo.ChannelTypeGuildText)
	if err != nil {
		return nil, fmt.Errorf("failed to created channel '%s': %v", channel_name, err)
	}

	log.Printf("Created channel: '%s'", channel_name)
	return channel, nil
}

func DeleteChannel(discord *discordgo.Session, guildID string, channel_name string) (bool, error) {

	channels, err := discord.GuildChannels(guildID)
	if err != nil {
		return false, fmt.Errorf("failed to list channels when deleting channel: '%s': %v", channel_name, err)
	}
	for _, channel := range channels {
		if channel.Name == channel_name {
			_, err = discord.ChannelDelete(channel.ID)
			if err != nil {
				return false, fmt.Errorf("failed to delete channel: %v", err)
			}
			return true, nil
		}
	}

	log.Printf("Tried to delete channel, but it didn't exist '%s'", channel_name)
	return false, nil
}

func GetChannelID(discord *discordgo.Session, guildID string, channel_name string) (string, error) {
	channels, err := discord.GuildChannels(guildID)
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
