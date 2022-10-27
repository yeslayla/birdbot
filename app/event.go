package app

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

const REMOTE_LOCATION string = "online"

type Event struct {
	Name     string
	Location string
	DateTime time.Time

	OrganizerID string
}

func (event *Event) GetChannelName() string {
	month := event.GetMonthPrefix()
	day := event.DateTime.Day() - 1
	city := event.GetCityFromLocation()

	channel := fmt.Sprint(month, "-", day, "-", city, "-", event.Name)
	channel = strings.ReplaceAll(channel, " ", "-")
	channel = strings.ToLower(channel)

	re, _ := regexp.Compile(`[^\w\-]`)
	channel = re.ReplaceAllString(channel, "")

	return channel
}

func (event *Event) GetCityFromLocation() string {

	if event.Location == REMOTE_LOCATION {
		return REMOTE_LOCATION
	}
	parts := strings.Split(event.Location, " ")
	index := -1
	loc := event.Location

	for i, v := range parts {
		part := strings.ToLower(v)
		if part == "mi" || part == "michigan" {
			index = i - 1
			if index < 0 {
				return "unknown"
			}
			if index > 0 && parts[index] == "," {
				index -= 1
			}

			if index > 1 && strings.Contains(parts[index-2], ",") {
				loc = fmt.Sprintf("%s-%s", parts[index-1], parts[index])
				break
			}

			loc = parts[index]
			break
		}
	}

	return loc
}

func (event *Event) GetMonthPrefix() string {
	month := event.DateTime.Month()
	data := map[time.Month]string{
		time.January:   "jan",
		time.February:  "feb",
		time.March:     "march",
		time.April:     "april",
		time.May:       "may",
		time.June:      "june",
		time.July:      "july",
		time.August:    "aug",
		time.September: "sept",
		time.October:   "oct",
		time.November:  "nov",
		time.December:  "dec",
	}

	return data[month]
}