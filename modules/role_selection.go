package modules

import (
	"fmt"
	"log"

	"github.com/yeslayla/birdbot-common/common"
	"github.com/yeslayla/birdbot/core"
	"github.com/yeslayla/birdbot/discord"
	"github.com/yeslayla/birdbot/persistence"
)

type roleSelectionModule struct {
	session *discord.Discord
	db      persistence.Database

	cfg       core.RoleSelectionConfig
	exclusive bool
}

// NewRoleSelectionComponent creates a new component
func NewRoleSelectionComponent(discord *discord.Discord, db persistence.Database, cfg core.RoleSelectionConfig) common.Module {
	return &roleSelectionModule{
		session:   discord,
		cfg:       cfg,
		db:        db,
		exclusive: true,
	}
}

// Initialize setups component on discord and registers handlers
func (c *roleSelectionModule) Initialize(birdbot common.ModuleManager) error {

	roles := []*discord.Role{}
	roleButtons := []discord.Component{}

	for _, roleConfig := range c.cfg.Roles {

		// Create & Validate Roles
		role := c.session.GetRoleAndCreate(roleConfig.RoleName)
		configColor, _ := core.HexToColor(roleConfig.Color)

		var emoji *discord.Emoji = nil
		if c.cfg.GenerateColorEmoji.IsEnabledByDefault() {
			emoji = c.session.GetEmoji(roleConfig.RoleName)
			if emoji == nil {
				emoji = c.session.CreateEmoji(roleConfig.RoleName, core.ColorToImage(configColor))
				if emoji == nil {
					log.Printf("Failed to create emoji for role: %s", roleConfig.RoleName)
					return nil
				}
			}

			// If role.ID in emoji.Roles
			var found bool
			for _, r := range emoji.Roles {
				if r == role.ID {
					found = true
					break
				}
			}
			if !found {
				emoji.Roles = append(emoji.Roles, role.ID)
				emoji.Save()
			}
		}

		if role.Color != configColor {
			role.Color = configColor
			role.Save()
		}

		// Create button
		var btn *discord.Button
		if emoji != nil {
			btn = c.session.NewButtonWithEmoji(fmt.Sprint(c.cfg.Title, roleConfig.RoleName), roleConfig.RoleName, emoji)
		} else {
			btn = c.session.NewButton(fmt.Sprint(c.cfg.Title, role.Name), role.Name)
		}
		btn.OnClick(func(user common.User) {

			// Assign the roles asynchronously to avoid Discord's response timeout
			go func() {
				// Remove other roles if exclusive
				if c.exclusive {
					for _, r := range roles {
						if r.ID == role.ID {
							continue
						}

						if c.session.HasRole(user, r) {
							c.session.UnassignRole(user, r)
						}
					}
				}

				// Toggle role
				if c.session.HasRole(user, role) {
					if err := c.session.UnassignRole(user, role); err != nil {
						log.Printf("Failed to unassign role: %s", err)
					}
				} else if err := c.session.AssignRole(user, role); err != nil {
					log.Printf("Failed to assign role: %s", err)
				}
			}()

		})

		roles = append(roles, role)
		roleButtons = append(roleButtons, btn)
	}

	components := []discord.Component{}
	var actionRow *discord.ActionRow
	for i, btn := range roleButtons {
		if i%5 == 0 {
			actionRow = c.session.NewActionRow()
			components = append(components, actionRow)
		}

		actionRow.AddComponent(btn)
	}

	localID := fmt.Sprint("role_selection:", c.cfg.Title)
	messageID, err := c.db.GetDiscordMessage(localID)
	if err != nil {
		return err
	}

	// Update message
	if messageID != "" {
		resultID := c.session.UpdateMessageComponent(messageID, c.cfg.SelectionChannel, fmt.Sprintf("**%s**\n%s", c.cfg.Title, c.cfg.Description), components)
		if resultID != "" {
			return nil
		}
	}

	// Create new message
	messageID = c.session.CreateMessageComponent(c.cfg.SelectionChannel, fmt.Sprintf("**%s**\n%s", c.cfg.Title, c.cfg.Description), components)
	return c.db.SetDiscordMessage(localID, messageID)

}
