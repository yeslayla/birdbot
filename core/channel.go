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

func GenerateEventChannelName(eventName string, location string, dateTime time.Time) string {
	month := GetMonthPrefix(dateTime)
	day := dateTime.Day()
	city := GetCityFromLocation(location)
	year := dateTime.Year()

	channel := fmt.Sprint(month, "-", day, city, "-", eventName, "-", year)
	channel = strings.ReplaceAll(channel, " ", "-")
	channel = strings.ToLower(channel)

	re, _ := regexp.Compile(`[^\w\-]`)
	channel = re.ReplaceAllString(channel, "")

	return channel
}

// GenerateChannel returns a channel object associated with an event
func GenerateChannel(event common.Event) *Channel {

	channelName := GenerateEventChannelName(event.Name, event.Location, event.DateTime)

	return &Channel{
		Name:     channelName,
		Verified: false,
	}
}
