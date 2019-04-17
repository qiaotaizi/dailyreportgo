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

var nextWorkDay time.Time

//初始化下个工作日
func init() {
	nextWorkDay=time.Now()
	for true{
		nextWorkDay.AddDate(0,0,1)
		if isWorkDay(nextWorkDay){
			break
		}
	}
}

//判断日期是否是工作日
func isWorkDay(date time.Time) bool {
	wkd:=date.Weekday()
	return wkd!=time.Saturday && wkd!=time.Sunday && isNotHoliday(date)
}

//判断日期不是法定节日
//后续完善
func isNotHoliday(weekday time.Time) bool {
	return true
}


