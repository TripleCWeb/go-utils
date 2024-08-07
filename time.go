package utils

import (
	"strconv"
	"time"
)

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

func StringToTime(str string) (*time.Time, error) {
	if str == "" {
		return nil, nil
	}
	// 将字符串解析为 int64 类型的时间戳
	timestamp, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return nil, err
	}

	// 使用 time.Unix() 将时间戳转换为 time.Time 类型
	t := time.Unix(timestamp, 0)

	return &t, nil
}

// StringToTimestamp将时间字符串解析为时间戳
// 输入参数：
// timeStr: 时间字符串:"2024-06-06T1:38:32.541+08:00"
// format: 时间格式:"2006-01-02T15:04:05.999-07:00"
func StringToTimestamp(timeStr, format string) (int64, error) {
	// 解析时间字符串为 time.Time 类型
	t, err := time.Parse(format, timeStr)
	if err != nil {
		return 0, err
	}

	// 将 time.Time 类型的时间转换为时间戳
	timestamp := t.Unix()

	return timestamp, nil
}

// 获取指定日期的前/后N天的开始和结束时间戳
func NDaysTimestamps(t time.Time, n int) (int64, int64) {
	targetDay := t.AddDate(0, 0, n)

	// 获取当天的开始时间
	dayStart := time.Date(targetDay.Year(), targetDay.Month(), targetDay.Day(), 0, 0, 0, 0, targetDay.Location())

	// 获取下一天的开始时间
	nextDayStart := dayStart.AddDate(0, 0, 1)

	// 获取时间戳
	return dayStart.Unix(), nextDayStart.Unix() - 1
}

// 获取指定日期的前/后N月的开始和结束时间戳
func NMonthsTimestamps(t time.Time, n int) (int64, int64) {
	targetMonth := t.AddDate(0, n, 0)

	// 获取指定月份的第一天
	monthStart := time.Date(targetMonth.Year(), targetMonth.Month(), 1, 0, 0, 0, 0, targetMonth.Location())

	// 获取下一个月的第一天
	nextMonthStart := monthStart.AddDate(0, 1, 0)

	// 获取时间戳
	return monthStart.Unix(), nextMonthStart.Unix() - 1
}
