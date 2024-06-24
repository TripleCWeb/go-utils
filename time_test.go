package main

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestTimeStrToTimeStamp(t *testing.T) {
	layout := "Jan-02-2006 03:04:05 PM -0700"
	dateTimeStr := "Nov-28-2021 11:29:26 PM +UTC"

	// 将 "+UTC" 替换为 "+0000"
	dateTimeStr = strings.Replace(dateTimeStr, "+UTC", "+0000", 1)

	dateTime, err := time.Parse(layout, dateTimeStr)
	if err != nil {
		fmt.Println("解析错误:", err)
		return
	}

	timestamp := dateTime.Unix()
	fmt.Println("时间戳:", timestamp)
}
