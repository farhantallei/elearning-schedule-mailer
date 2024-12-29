package main

import (
	"fmt"
	"time"
)

func defaultIfNil(value *string, defaultValues ...string) string {
	if value == nil {
		if len(defaultValues) > 0 {
			return defaultValues[0]
		}
		return "-"
	}
	return *value
}

func getEndTime(endTime string) time.Time {
	now := time.Now()
	layout := "15:04" // Time format "HH:mm"
	endTimeParsed, err := time.ParseInLocation(layout, endTime, now.Location())
	if err != nil {
		fmt.Println("Error parsing time:", err)
		return now
	}

	return time.Date(now.Year(), now.Month(), now.Day(), endTimeParsed.Hour(), endTimeParsed.Minute(), 0, 0, now.Location())
}
