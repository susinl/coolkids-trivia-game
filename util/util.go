package util

import (
	"time"
)

func ParseDateTime(dateTime string) (time.Time, error) {
	return time.ParseInLocation(DateTimeFormat, dateTime, time.Local)
}
