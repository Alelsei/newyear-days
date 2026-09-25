package daysuntil

import "time"

func DaysUntilNewYear(from time.Time) int {
	year, month, day := from.Date()

	today := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	nextNewYear := time.Date(year+1, time.January, 1, 0, 0, 0, 0, time.UTC)

	return int(nextNewYear.Sub(today).Hours() / 24)
}
