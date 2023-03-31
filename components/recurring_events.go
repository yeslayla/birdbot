package components

import (
	"log"
	"strings"

	"github.com/yeslayla/birdbot/common"
	"github.com/yeslayla/birdbot/discord"
)

type recurringEventsComponent struct {
	session *discord.Discord
}

// NewRecurringEventsComponent creates a new component instance
func NewRecurringEventsComponent() common.Component {
	return &recurringEventsComponent{}
}

// Initialize registers event listeners
func (c *recurringEventsComponent) Initialize(birdbot common.ComponentManager) error {
	_ = birdbot.OnEventComplete(c.OnEventComplete)

	return nil
}

// OnEventComplete checks for keywords before creating a new event
func (c *recurringEventsComponent) OnEventComplete(e common.Event) error {

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
