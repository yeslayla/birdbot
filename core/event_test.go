package core

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetChannelName(t *testing.T) {
	assert := assert.New(t)

	// Test Valid Address
	event := Event{
		Name:     "Hello World",
		Location: "1234 Place Rd, Ann Arbor, MI 00000",
		DateTime: time.Date(2022, time.January, 5, 0, 0, 0, 0, time.UTC),
	}
	assert.Equal("jan-5-ann-arbor-hello-world-2022", event.Channel().Name)

	// Test Unparsable Location
	// lmanley: Note it'd be nice to expand support for this
	event = Event{
		Name:     "Hello World",
		Location: "Michigan Theater, Ann Arbor",
		DateTime: time.Date(2022, time.January, 5, 0, 0, 0, 0, time.UTC),
	}
	assert.Equal("jan-5-hello-world-2022", event.Channel().Name)

	// Test Short Location
	event = Event{
		Name:     "Hello World",
		Location: "Monroe, MI",
		DateTime: time.Date(2022, time.January, 5, 0, 0, 0, 0, time.UTC),
	}
	assert.Equal("jan-5-monroe-hello-world-2022", event.Channel().Name)

	// Test Short Location
	event = Event{
		Name:     "Hello World",
		Location: "Monroe St, Monroe , MI",
		DateTime: time.Date(2022, time.January, 5, 0, 0, 0, 0, time.UTC),
	}
	assert.Equal("jan-5-monroe-hello-world-2022", event.Channel().Name)

	// Test Remote Event
	event = Event{
		Name:     "Hello World",
		Location: REMOTE_LOCATION,
		DateTime: time.Date(2022, time.January, 5, 0, 0, 0, 0, time.UTC),
	}
	assert.Equal("jan-5-online-hello-world-2022", event.Channel().Name)
}

func TestMonthPrefix(t *testing.T) {
	assert := assert.New(t)

	event := Event{
		DateTime: time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC),
	}
	assert.Equal("jan", event.GetMonthPrefix())
}
