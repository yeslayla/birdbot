package modules

import (
	"log"
	"strings"

	"github.com/yeslayla/birdbot-common/common"
	"github.com/yeslayla/birdbot/discord"
)

type recurringEventsModule struct {
	session *discord.Discord
}

// NewRecurringEventsComponent creates a new component instance
//
// Deprecated: This recurring events are now a native feature in Discord,
// so this will be removed in the next version.
func NewRecurringEventsComponent(session *discord.Discord) common.Module {
	return &recurringEventsModule{
		session: session,
	}
}

// Initialize registers event listeners
func (c *recurringEventsModule) Initialize(birdbot common.ModuleManager) error {
	_ = birdbot.OnEventComplete(c.OnEventComplete)

	return nil
}

// OnEventComplete checks for keywords before creating a new event
func (c *recurringEventsModule) OnEventComplete(e common.Event) error {

	if strings.Contains(strings.ToLower(e.Description), "recurring weekly") {
		startTime := e.DateTime.AddDate(0, 0, 7)
		finishTime := e.CompleteDateTime.AddDate(0, 0, 7)
		nextEvent := e
		nextEvent.DateTime = startTime
		nextEvent.CompleteDateTime = finishTime

		if err := c.session.CreateEvent(nextEvent); err != nil {
			log.Print("Failed to create recurring event: ", err)
		}
	}

	return nil
}
