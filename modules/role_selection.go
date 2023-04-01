package modules

import (
	"fmt"
	"log"

	"github.com/yeslayla/birdbot/common"
	"github.com/yeslayla/birdbot/core"
	"github.com/yeslayla/birdbot/discord"
)

type roleSelectionModule struct {
	session  *discord.Discord
	cfg      core.RoleSelectionConfig
	exlusive bool
}

// NewRoleSelectionComponent creates a new component
func NewRoleSelectionComponent(discord *discord.Discord, cfg core.RoleSelectionConfig) common.Module {
	return &roleSelectionModule{
		session:  discord,
		cfg:      cfg,
		exlusive: true,
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

		if role.Color != configColor {
			role.Color = configColor
			role.Save()
		}

		// Create button
		btn := c.session.NewButton(fmt.Sprint(c.cfg.Title, role.Name), role.Name)
		btn.OnClick(func(user common.User) {

			// Remove other roles if exclusive
			if c.exlusive {
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

	c.session.CreateMessageComponent(c.cfg.SelectionChannel, fmt.Sprintf("**%s**\n%s", c.cfg.Title, c.cfg.Description), components)

	return nil
}
