package common

import "github.com/bwmarrin/discordgo"

type CommandOptionType uint64

const (
	CommandTypeString CommandOptionType = CommandOptionType(discordgo.ApplicationCommandOptionString)
	CommandTypeInt    CommandOptionType = CommandOptionType(discordgo.ApplicationCommandOptionInteger)
	CommandTypeBool   CommandOptionType = CommandOptionType(discordgo.ApplicationCommandOptionBoolean)
	CommandTypeFloat  CommandOptionType = CommandOptionType(discordgo.ApplicationCommandOptionNumber)
)

type ChatCommandConfiguration struct {
	Description       string
	EphemeralResponse bool
	Options           map[string]ChatCommandOption
}

type ChatCommandOption struct {
	Description string
	Type        CommandOptionType
	Required    bool
}
