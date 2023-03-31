package events

import (
	"log"
	"strings"

	"github.com/yeslayla/birdbot/common"
	"github.com/yeslayla/birdbot/discord"
)

type RecurringEventsComponent struct {
	session *discord.Discord
}

func NewRecurringEventsComponent() *RecurringEventsComponent {
	return &RecurringEventsComponent{}
}

func (c *RecurringEventsComponent) Initialize(birdbot common.ComponentManager) error {
	_ = birdbot.OnEventComplete(c.OnEventComplete)

	return nil
}

func (c *RecurringEventsComponent) OnEventComplete(e common.Event) error {

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
