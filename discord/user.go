package discord

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/yeslayla/birdbot-common/common"
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
		DisplayName: user.Username,
		ID:          user.ID,
	}
}

// AssignRole adds a role to a user
func (discord *Discord) AssignRole(user common.User, role *Role) error {
	return discord.session.GuildMemberRoleAdd(discord.guildID, user.ID, role.ID)
}

// UnassignRole removes a role from a user
func (discord *Discord) UnassignRole(user common.User, role *Role) error {
	return discord.session.GuildMemberRoleRemove(discord.guildID, user.ID, role.ID)
}

// HasRole returns true when a user has a given role
func (discord *Discord) HasRole(user common.User, role *Role) bool {
	return discord.HasAtLeastOneRole(user, []*Role{role})
}

// HasAtLeastOneRole returns true when a user has at one role from a given array
func (discord *Discord) HasAtLeastOneRole(user common.User, roles []*Role) bool {

	member, err := discord.session.GuildMember(discord.guildID, user.ID)
	if err != nil {
		log.Printf("Failed to get member: %s", err)
		return false
	}

	for _, v := range member.Roles {
		for _, targetRole := range roles {
			if v == targetRole.ID {
				return true
			}
		}
	}

	return false
}
