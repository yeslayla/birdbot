package core

import "time"

// GetMonthPrefix returns a month in short form
func GetMonthPrefix(dateTime time.Time) string {
	month := dateTime.Month()
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
