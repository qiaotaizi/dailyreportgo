package main

import (
	"time"
)

const dateFormat = "2006-01-02"

//获取当前月份的第一天和最后一天
func monthStartAndEnd() (string, string) {
	start := now.AddDate(0, 0, -now.Day()+1)
	end := start.AddDate(0, 1, -1)
	startStr := start.Format(dateFormat)
	endStr := end.Format(dateFormat)
	return startStr, endStr
}

//判断日期是否是工作日
//逻辑:判断是否是法定节假日,是,返回false,不是,判断是否是周末,是,返回false,不是返回true
func isWorkDay(date time.Time) bool {
	if isHoliday(date) {
		return false
	}
	if !isWeekend(date) {
		return true
	}
	return isTX(date)
}

//判断给定日期是否是周末
func isWeekend(date time.Time) bool {
	wkd := date.Weekday()
	return wkd == time.Saturday || wkd == time.Sunday
}

//在包初始化阶段(init函数中)决定isHoliday和isTX函数
//判断标准是假期库是否已经耗尽
var isHoliday, isTX func(d time.Time) bool

//计算两个时间之间的日期数
//不足一天的向下取整
func daysBetweenTimes(startTime time.Time, endTime time.Time) int {
	hourCount := int(endTime.Sub(startTime).Hours())
	a, b := hourCount/24, hourCount%24
	if b == 0 {
		return a
	}
	return a + 1
}
