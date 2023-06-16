package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/yeslayla/birdbot/common"
)

type Button struct {
	Label string
	ID    string

	discord *Discord
}

// NewButton creates a new button component
func (discord *Discord) NewButton(id string, label string) *Button {
	return &Button{
		discord: discord,
		ID:      id,
		Label:   label,
	}
}

// OnClick registers an event when the button is clicked
func (button *Button) OnClick(action func(user common.User)) {
	button.discord.session.AddHandler(func(s *discordgo.Session, r *discordgo.InteractionCreate) {
		if r.Interaction.Type != discordgo.InteractionMessageComponent {
			return
		}

		if r.MessageComponentData().CustomID == button.ID {

			action(NewUser(r.Member.User))

			s.InteractionRespond(r.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
			})
		}
	})

}

func (button *Button) toMessageComponent() discordgo.MessageComponent {
	return discordgo.Button{
		Label:    button.Label,
		CustomID: button.ID,
		Style:    discordgo.PrimaryButton,
	}
}
