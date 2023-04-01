package discord

import (
	"image/color"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/yeslayla/birdbot/core"
)

type Role struct {
	discord *Discord
	ID      string

	Name  string
	Color color.Color
}

// GetRole returns a role that exists on Discord
func (discord *Discord) GetRole(name string) *Role {

	roles, err := discord.session.GuildRoles(discord.guildID)
	if err != nil {
		log.Printf("Error occured listing roles: %s", err)
		return nil
	}

	for _, role := range roles {
		if role.Managed {
			continue
		}
		if role.Name == name {

			return &Role{
				Name:    role.Name,
				Color:   core.IntToColor(role.Color),
				discord: discord,
				ID:      role.ID,
			}
		}
	}

	return nil
}

// GetRoleAndCreate gets a role and creates it if it doesn't exist
func (discord *Discord) GetRoleAndCreate(name string) *Role {
	role := discord.GetRole(name)
	if role != nil {
		return role
	}

	if _, err := discord.session.GuildRoleCreate(discord.guildID, &discordgo.RoleParams{
		Name:  name,
		Color: core.Int(0),
	}); err != nil {
		log.Printf("Failed to create role: %s", err)
		return nil
	}

	return discord.GetRole(name)
}

// Save updates the role on Discord
func (role *Role) Save() {
	if _, err := role.discord.session.GuildRoleEdit(role.discord.guildID, role.ID, &discordgo.RoleParams{
		Name:  role.Name,
		Color: core.Int(core.ColorToInt(role.Color)),
	}); err != nil {
		log.Printf("Failed to save role: %s", err)
	}
}
