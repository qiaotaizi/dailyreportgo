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
//逻辑:判断是否是法定节假日,是,返回false,不是,判断是否是周末,是,返回false,不是返回true
func isWorkDay(date time.Time) bool {
	if isHoliday(date){
		return false
	}
	if !isWeekend(date){
		return true
	}
	return isTX(date)
}

//判断日期是调休
func isTX(date time.Time) bool {
	return false
}

//判断日期是周末
func isWeekend(date time.Time) bool {
	wkd:=date.Weekday()
	return wkd==time.Saturday || wkd==time.Sunday
}

//判断日期不是法定节日
//后续完善
func isHoliday(date time.Time) bool {
	return false
}
