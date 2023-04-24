package util

import (
	"regexp"
	"time"
	"unicode"
)

var IsPhone = regexp.MustCompile(`^(\+66|0)?[-\s]?[1-9]\d{8}$`).MatchString

func ParseDateTime(dateTime string) (time.Time, error) {
	return time.ParseInLocation(DateTimeFormat, dateTime, time.Local)
}

func IsLetter(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsSpace(r) {
			return false
		}
	}
	return true
}
