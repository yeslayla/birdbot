package discord

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/yeslayla/birdbot/common"
)

// RegisterCommand creates an new command that can be used to interact with bird bot
func (discord *Discord) RegisterCommand(name string, config common.ChatCommandConfiguration, handler func(common.User, map[string]any) string) {
	command := &discordgo.ApplicationCommand{
		Name:        name,
		Description: config.Description,
	}

	// Convert options to discordgo objects
	command.Options = make([]*discordgo.ApplicationCommandOption, len(config.Options))
	index := 0
	for name, option := range config.Options {
		command.Options[index] = &discordgo.ApplicationCommandOption{
			Name:        name,
			Description: option.Description,
			Required:    option.Required,
			Type:        discordgo.ApplicationCommandOptionType(option.Type),
		}
		index++
	}

	// Register handler
	discord.commandHandlers[name] = func(session *discordgo.Session, r *discordgo.InteractionCreate) {
		if r.Interaction.Type != discordgo.InteractionApplicationCommand {
			return
		}

		cmdOptions := r.ApplicationCommandData().Options

		// Parse option types
		optionsMap := make(map[string]any, len(cmdOptions))
		for _, opt := range cmdOptions {
			switch config.Options[opt.Name].Type {
			case common.CommandTypeString:
				optionsMap[opt.Name] = opt.StringValue()
			case common.CommandTypeInt:
				optionsMap[opt.Name] = opt.IntValue()
			case common.CommandTypeBool:
				optionsMap[opt.Name] = opt.BoolValue()
			case common.CommandTypeFloat:
				optionsMap[opt.Name] = opt.FloatValue()
			default:
				optionsMap[opt.Name] = opt.Value
			}
		}

		session.InteractionRespond(r.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags: discordgo.MessageFlagsEphemeral,
			},
		})

		result := handler(NewUser(r.Member.User), optionsMap)

		if result != "" {
			// Handle response
			responseData := &discordgo.WebhookParams{
				Content: result,
			}

			if config.EphemeralResponse {
				responseData.Flags = discordgo.MessageFlagsEphemeral
			}

			session.FollowupMessageCreate(r.Interaction, false, responseData)

		} else {
			session.FollowupMessageCreate(r.Interaction, false, &discordgo.WebhookParams{
				Content: "Command did not return a response!",
			})
			log.Printf("Command '%s' did not return a response: %v", name, optionsMap)
		}
	}

	cmd, err := discord.session.ApplicationCommandCreate(discord.applicationID, discord.guildID, command)
	if err != nil {
		log.Fatalf("Cannot create command '%s': %v", name, err)
	}
	discord.commands[name] = cmd
}

// ClearCommands deregisters all commands from the discord API
func (discord *Discord) ClearCommands() {
	for _, v := range discord.commands {
		err := discord.session.ApplicationCommandDelete(discord.session.State.User.ID, discord.guildID, v.ID)
		if err != nil {
			log.Fatalf("Cannot delete command '%s': %v", v.Name, err)
		}
	}
}
