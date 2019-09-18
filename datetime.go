package main

import (
	"log"
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
	nextWorkDay = time.Now()
	for true {
		nextWorkDay=nextWorkDay.AddDate(0, 0, 1)
		if isWorkDay(nextWorkDay) {
			break
		}
	}
}

//判断日期是否是工作日
//逻辑:判断是否是法定节假日,是,返回false,不是,判断是否是周末,是,返回false,不是返回true
func isWorkDay(date time.Time) bool {

	isWeekend := func(d time.Time) bool {
		wkd := date.Weekday()
		return wkd == time.Saturday || wkd == time.Sunday
	}

	isHoliday_ := func(d time.Time) bool {
		y := date.Year()
		holidaysOfYear, ok := holidaysMap[y]
		if !ok {
			log.Fatalf("please maintain holidays of year %d before generating the daily report", y)
		}
		for _, h := range holidaysOfYear {
			if time.Month(h.m) == date.Month() && h.d == date.Day() {
				//找到当前日期
				return h.t == rest
			}
		}
		return false
	}

	//判断是否是调休
	isTX_ := func(d time.Time) bool {
		y := date.Year()
		holidaysOfYear, ok := holidaysMap[y]
		if !ok {
			log.Fatalf("please maintain holidays of year %d before generating the daily report", y)
		}
		for _, h := range holidaysOfYear {
			if time.Month(h.m) == date.Month() && h.d == date.Day() {
				//找到当前日期
				return h.t == work
			}
		}
		return false
	}

	var isHoliday,isTX func(d time.Time) bool


	balanceFlag:=holidayBalanceByNow()

	switch balanceFlag {
	case enough://维护的假期充足,应用正常的假期/调休判断方法
		isHoliday=isHoliday_
		isTX=isTX_
	case exhausting://维护的假期即将耗尽,应用正常的假期/调休判断方法,但给出警告
		isHoliday=isHoliday_
		isTX=isTX_
		warn("假期库即将耗尽, 请尽快维护")
	default://exhausted 维护的假期库已经耗尽,假期判断和调休判断总是返回false,并且给出警告
		isHoliday= func(time.Time) bool {
			return false
		}
		isTX= func(time.Time) bool {
			return false
		}
		warn("假期库已经耗尽, 法定节假日的推断逻辑将被禁用, 下个工作日的推断可能会不准确,请尽快维护假期库")
	}
	if isHoliday(date) {
		return false
	}
	if !isWeekend(date) {
		return true
	}
	return isTX(date)
}

//计算两个时间之间的日期数
//不足一天的向下取整
func daysBetweenTimes(startTime time.Time,endTime time.Time) int{
	hourCount:=int(endTime.Sub(startTime).Hours())
	a,b:=hourCount/24,hourCount%24
	if b==0{
		return a
	}
	return a+1
}
