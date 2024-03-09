package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/yeslayla/birdbot-common/common"
)

type Button struct {
	Label string
	ID    string

	Emoji   *Emoji
	discord *Discord
}

// NewButton creates a new button component
func (discord *Discord) NewButton(id string, label string) *Button {
	return &Button{
		discord: discord,
		ID:      id,
		Label:   label,
		Emoji:   nil,
	}
}

// NewButtonWithEmoji creates a new button component with a emoji
func (discord *Discord) NewButtonWithEmoji(id string, label string, emoji *Emoji) *Button {
	return &Button{
		discord: discord,
		ID:      id,
		Label:   label,
		Emoji:   emoji,
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
	cmp := discordgo.Button{
		Label:    button.Label,
		CustomID: button.ID,
		Style:    discordgo.PrimaryButton,
	}
	if button.Emoji != nil {
		cmp.Emoji = &discordgo.ComponentEmoji{
			Name: button.Emoji.Name,
			ID:   button.Emoji.ID,
		}
	}
	return cmp
}
