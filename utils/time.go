package utils

import (
	"time"
)

func ParseTime(s string) (time.Time, error) {
	layoutSimple := "02-01-2006"
	return time.Parse(layoutSimple, s)
}
