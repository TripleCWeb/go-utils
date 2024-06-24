package main

import "time"

// IsToday按照utc0时区来判断自然天
func IsTodayUTC(ts int32) bool {
	now := time.Now().UTC()
	t := time.Unix(int64(ts), 0).UTC()
	return now.Year() == t.Year() && now.Month() == t.Month() && now.Day() == t.Day()
}
