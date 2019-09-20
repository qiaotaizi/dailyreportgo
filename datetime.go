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

//在包初始化阶段决定isHoliday和isTX函数
//判断标准是假期库是否已经耗尽
var isHoliday, isTX = func() (func(d time.Time) bool, func(d time.Time) bool) {
	isHoliday_ := func(date time.Time) bool {
		y := date.Year()
		holidaysOfYear, ok := holidaysMap[y]
		if !ok {
			//实际上是不会走到这里的
			warn("请在假期表中维护%d年的假期及调休数据", y)
			return false
		}
		for _, h := range holidaysOfYear {
			if time.Month(h.m) == date.Month() && h.d == date.Day() {
				//找到当前日期
				return h.t == rest
			}
		}
		return false
	}
	isHoliday__ := func(time.Time) bool { return false }
	//判断是否是调休
	isTX_ := func(date time.Time) bool {
		y := date.Year()
		holidaysOfYear, ok := holidaysMap[y]
		if !ok {
			//实际上是不会走到这里的
			warn("请在假期表中维护%d年的假期及调休数据 ", y)
			return false
		}
		for _, h := range holidaysOfYear {
			if time.Month(h.m) == date.Month() && h.d == date.Day() {
				//找到当前日期
				return h.t == work
			}
		}
		return false
	}
	isTX__ := func(time.Time) bool { return false }

	balanceFlag := holidayBalanceByNow()

	switch balanceFlag {
	case enough: //维护的假期充足,应用正常的假期/调休判断方法
		lg("假期库充足")
		return isHoliday_, isTX_
	case exhausting: //维护的假期即将耗尽,应用正常的假期/调休判断方法,但给出警告
		warn("假期库即将耗尽, 请尽快维护")
		return isHoliday_, isTX_
	default: //exhausted 维护的假期库已经耗尽,假期判断和调休判断总是返回false,并且给出警告
		warn("假期库已经耗尽, 法定节假日的推断逻辑将被禁用, 下个工作日的推断可能会不准确,请尽快维护假期库")
		return isHoliday__, isTX__
	}
}()

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
