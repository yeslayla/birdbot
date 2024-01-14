package modules

import (
	"fmt"

	"github.com/yeslayla/birdbot-common/common"
)

type statusModule struct {
	portalURL string
}

// NewStatusComponent creates a new component
func NewStatusComponent(portalURL string) common.Module {
	m := &statusModule{
		portalURL: portalURL,
	}

	return m
}

func (c *statusModule) Initialize(birdbot common.ModuleManager) error {
	birdbot.RegisterCommand("status", common.ChatCommandConfiguration{
		Description:       "Gets the current status of the bot",
		EphemeralResponse: false,
	}, func(user common.User, args map[string]any) string {

		return fmt.Sprintf("The bot is currently OK.\nSee Status Portal for more information: %s", c.portalURL)
	})
	return nil
}
