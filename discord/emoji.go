package discord

import (
	"image"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/yeslayla/birdbot/core"
)

type Emoji struct {
	discord *Discord
	ID      string

	Name  string
	Roles []string
}

// GetEmoji returns a emoji that exists on Discord
func (discord *Discord) GetEmoji(name string) *Emoji {

	emojis, err := discord.session.GuildEmojis(discord.guildID)
	if err != nil {
		log.Printf("Error occured listing roles: %s", err)
		return nil
	}

	for _, emoji := range emojis {
		if emoji.Managed {
			continue
		}
		if emoji.Name == name {

			return &Emoji{
				ID:      emoji.ID,
				discord: discord,
				Name:    emoji.Name,
				Roles:   emoji.Roles,
			}
		}
	}

	return nil
}

// CreateEmoji creates a new emoji on Discord
func (discord *Discord) CreateEmoji(name string, image image.Image) *Emoji {
	result, err := discord.session.GuildEmojiCreate(discord.guildID, &discordgo.EmojiParams{
		Name:  name,
		Image: core.ImageToBase64(image),
	})
	if err != nil {
		log.Printf("Failed to create emoji: %s", err)
		return nil
	}

	if result == nil {
		log.Print("Failed to create emoji: result is nil")
		return nil
	}

	return &Emoji{
		ID:      result.ID,
		Name:    result.Name,
		Roles:   result.Roles,
		discord: discord,
	}
}

// Save updates the emoji on Discord
func (emoji *Emoji) Save() {
	if _, err := emoji.discord.session.GuildEmojiEdit(emoji.discord.guildID, emoji.ID, &discordgo.EmojiParams{
		Name:  emoji.Name,
		Roles: emoji.Roles,
	}); err != nil {
		log.Printf("Failed to save role: %s", err)
	}
}
