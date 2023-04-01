package discord

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

// Component is an object that can be formatted as a discord component
type Component interface {
	toMessageComponent() discordgo.MessageComponent
}

// CreateMessageComponent creates a discord component
func (discord *Discord) CreateMessageComponent(channelID string, content string, components []Component) string {

	dComponents := make([]discordgo.MessageComponent, len(components))
	for i, v := range components {
		dComponents[i] = v.toMessageComponent()
	}

	result, err := discord.session.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
		Components: dComponents,
		Content:    content,
	})
	if err != nil {
		log.Print(err)
		return ""
	}

	return result.ID
}