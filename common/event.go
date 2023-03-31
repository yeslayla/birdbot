package common

import (
	"time"
)

type Event struct {
	Name             string
	ID               string
	Location         string
	Completed        bool
	DateTime         time.Time
	CompleteDateTime time.Time
	Description      string
	ImageURL         string

	Organizer User
}
