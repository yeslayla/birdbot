package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/yeslayla/birdbot/core"
)

func NewUser(user *discordgo.User) *core.User {
	return &core.User{
		Name: user.Username,
		ID:   user.ID,
	}
}
