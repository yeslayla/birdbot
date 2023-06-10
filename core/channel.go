package core

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/yeslayla/birdbot/common"
)

type Channel struct {
	Name     string
	ID       string
	Verified bool
}

// GenerateEventChannelName deciphers a channel name from a given set of event data
func GenerateEventChannelName(eventName string, location string, dateTime time.Time) string {
	month := GetMonthPrefix(dateTime)
	day := dateTime.Day()
	city := GetCityFromLocation(location)
	year := dateTime.Year()

	// Remove special characters
	eventName = regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(eventName, "")
	eventName = strings.Trim(eventName, " ")

	channel := fmt.Sprint(month, "-", day, city, "-", eventName, "-", year)
	channel = strings.ReplaceAll(channel, " ", "-")
	channel = strings.ToLower(channel)

	re, _ := regexp.Compile(`[^\w\-]`)
	channel = re.ReplaceAllString(channel, "")

	return channel
}

// GenerateChannelFromEvent returns a channel object associated with an event
func GenerateChannelFromEvent(event common.Event) *Channel {

	channelName := GenerateEventChannelName(event.Name, event.Location, event.DateTime)

	return &Channel{
		Name:     channelName,
		Verified: false,
	}
}
