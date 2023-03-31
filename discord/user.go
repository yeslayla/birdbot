package discord

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/yeslayla/birdbot/common"
)

// NewUser creates a new user object from a discordgo.User object
func NewUser(user *discordgo.User) common.User {
	if user == nil {
		log.Print("Cannot user object, user is nil!")
		return common.User{
			ID: "-1",
		}
	}

	return common.User{
		ID: user.ID,
	}
}
