package core

import (
	"fmt"
	"strings"
)

const REMOTE_LOCATION string = "online"

// GetCityFromLocation returns the city name of an event's location
func GetCityFromLocation(location string) string {

	if location == REMOTE_LOCATION {
		return fmt.Sprint("-", REMOTE_LOCATION)
	}
	parts := strings.Split(location, " ")
	index := -1
	loc := location

	for i, v := range parts {
		part := strings.ToLower(v)
		if part == "mi" || part == "michigan" {
			index = i - 1
			if index < 0 {
				return ""
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

	return fmt.Sprint("-", loc)
}
