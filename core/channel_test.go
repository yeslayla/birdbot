package core

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGenerateEventChannelName(t *testing.T) {
	assert := assert.New(t)

	// Test Valid Address
	channelName := GenerateEventChannelName("Hello World", "1234 Place Rd, Ann Arbor, MI 00000", time.Date(2022, time.January, 5, 0, 0, 0, 0, time.UTC))
	assert.Equal("jan-5-ann-arbor-hello-world-2022", channelName)

	// Test Unparsable
	// lmanley: Note it'd be nice to expand support for this
	channelName = GenerateEventChannelName("Hello World", "Michigan Theater, Ann Arbor", time.Date(2022, time.January, 5, 0, 0, 0, 0, time.UTC))
	assert.Equal("jan-5-hello-world-2022", channelName)

	// Test Short Location
	channelName = GenerateEventChannelName("Hello World", "Monroe, MI", time.Date(2022, time.January, 5, 0, 0, 0, 0, time.UTC))
	assert.Equal("jan-5-monroe-hello-world-2022", channelName)

	// Test Remote Event
	channelName = GenerateEventChannelName("Hello World", RemoteLocation, time.Date(2022, time.January, 5, 0, 0, 0, 0, time.UTC))
	assert.Equal("jan-5-online-hello-world-2022", channelName)

	channelName = GenerateEventChannelName("Hangout :)", "Quickly Livonia", time.Date(2022, time.January, 5, 0, 0, 0, 0, time.UTC))
	assert.Equal("jan-5-quickly-livonia-hangout-2022", channelName)
}
