package utils

import (
	"time"
)

var LOC, _ = time.LoadLocation("Asia/Shanghai")

const TIME_LAYOUT = "2006-01-02 15:04:05"

// 字符串转时间戳
func StrToTime(stime string) time.Time {
	return StrToLayoutTime(stime, TIME_LAYOUT)
}

// 根据模板获取时间戳
func StrToLayoutTime(stime, layout string) time.Time {
	t, _ := time.ParseInLocation(layout, stime, LOC)
	return t
}

// 时间戳转字符串时间
func TimestampToStr(timestamp int64) string {
	t := time.Unix(timestamp, 0)
	return t.Format(TIME_LAYOUT)
}

func TimestampToLayoutStr(timestamp int64, layout string) string {
	t := time.Unix(timestamp, 0)
	return t.Format(layout)
}
