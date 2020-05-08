package model

import "time"

func parseDate(value, layout string) (time.Time, error) {
	return time.Parse(layout, value)
}
