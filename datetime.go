package main

import (
	"time"
)

const format = "2006-01-02"

//获取当前月份的第一天和最后一天
func monthStartAndEnd() (string, string) {
	nowTime := time.Now()
	start := nowTime.AddDate(0, 0, -nowTime.Day()+1)
	end := start.AddDate(0, 1, -1)
	startStr := start.Format(format)
	endStr := end.Format(format)
	return startStr, endStr
}
