package shortcode_test

import "time"

type clock struct{}

func (clock) Now() time.Time {
	str := "2026-02-21T11:04:57.497"
	date, _ := time.Parse("2006-01-02T15:04:05.000", str)
	return date
}
