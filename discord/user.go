package discord

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/yeslayla/birdbot/core"
)

// NewUser creates a new user object from a discordgo.User object
func NewUser(user *discordgo.User) *core.User {
	if user == nil {
		log.Print("Cannot user object, user is nil!")
		return nil
	}

	return &core.User{
		ID: user.ID,
	}
}
