package util

import (
	"net/mail"
	"time"
)

func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func ParseDateTime(dateTime string) (time.Time, error) {
	return time.ParseInLocation(DateTimeFormat, dateTime, time.Local)
}
