package utils

import (
	"errors"
	"time"
)

func GetCurrentTime() (string, error) {
	currentTime := time.Now().UTC()
	return currentTime.Format(time.RFC3339), nil
}

func ParseTime(timeStr string) (time.Time, error) {
	parsedTime, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return time.Time{}, errors.New("invalid time format")
	}
	return parsedTime, nil
}
