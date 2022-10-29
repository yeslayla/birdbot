package core

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetChannelName(t *testing.T) {
	assert := assert.New(t)

	event := Event{
		Name:     "Hello World",
		Location: "1234 Place Rd, Ann Arbor, MI 00000",
		DateTime: time.Date(2022, time.January, 5, 0, 0, 0, 0, time.UTC),
	}
	assert.Equal("jan-5-ann-arbor-hello-world", event.Channel().Name)

	event = Event{
		Name:     "Hello World",
		Location: "Michigan Theater, Ann Arbor",
		DateTime: time.Date(2022, time.January, 5, 0, 0, 0, 0, time.UTC),
	}
	assert.Equal("jan-5-hello-world", event.Channel().Name)

}
