package utils

import "time"

// IsToday按照utc0时区来判断自然天
func IsTodayUTC(ts int32) bool {
	now := time.Now().UTC()
	t := time.Unix(int64(ts), 0).UTC()
	return now.Year() == t.Year() && now.Month() == t.Month() && now.Day() == t.Day()
}

// NextMonthsFirstTimeStamp获取下N个月1号00:00:00的时间戳
func NextMonthsFirstTimeStamp(ts int64, months int) int64 {
	// 转换为时间
	t := time.Unix(ts, 0)

	// 获取下一月
	nextMonth := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location()).AddDate(0, months, 0)

	// 获取时间戳
	return nextMonth.Unix()
}

// YesterdayLastTS前一天的23:59:59
func YesterdayLastTS() int64 {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Unix() - 1
}
