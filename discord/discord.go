package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/mock"
)

type Discord struct {
	mock.Mock

	session *discordgo.Session
}

func NewDiscord(session *discordgo.Session) *Discord {
	return &Discord{
		session: session,
	}
}
