package main

import (
	"time"
)

//手动维护一套节假日

const (
	rest = 1 + iota //法定节假日休息
	work            //法定节假日调休工作
)

type holidayBalanceFlag int //存余假期标志

const (
	exhausting = 1 + iota //假期将耗尽标志
	exhausted             //假期已耗尽标志
	enough
)

type holiday struct {
	m int //月份
	d int //日期
	t int //类型 rest/work
}

//年份-休假表map
//硬编码是不好的
//后续应改为由外置json文件来维护
var holidaysMap = map[int][]holiday{
	2019: {
		{1, 1, rest},
		{2, 2, work},
		{2, 3, work},
		{2, 4, rest},
		{2, 5, rest},
		{2, 6, rest},
		{2, 7, rest},
		{2, 8, rest},
		{2, 9, rest},
		{2, 10, rest},
		{4, 5, rest},
		{4, 6, rest},
		{4, 7, rest},
		{4, 28, work},
		{5, 1, rest},
		{5, 2, rest},
		{5, 3, rest},
		{5, 4, rest},
		{5, 5, work},
		{6, 7, rest},
		{6, 8, rest},
		{6, 9, rest},
		{9, 13, rest},
		{9, 14, rest},
		{9, 15, rest},
		{9, 29, work},
		{10, 1, rest},
		{10, 2, rest},
		{10, 3, rest},
		{10, 4, rest},
		{10, 5, rest},
		{10, 6, rest},
		{10, 7, rest},
		{10, 7, rest},
		{10, 12, work},
	},
	//2020年后续完善
	2020: {
		{1, 1, rest},
		{1, 19, work},
		{1, 24, rest},
		{1, 25, rest},
		{1, 26, rest},
		{1, 27, rest},
		{1, 28, rest},
		{1, 29, rest},
		{1, 30, rest},
		{2, 1, work},
		{4, 4, rest},
		{4, 5, rest},
		{4, 6, rest},
		{4, 26, work},
		{5, 1, rest},
		{5, 2, rest},
		{5, 3, rest},
		{5, 4, rest},
		{5, 5, rest},
		{5, 9, work},
		{6, 25, rest},
		{6, 26, rest},
		{6, 27, rest},
		{6, 28, work},
		{10, 1, rest},
		{10, 2, rest},
		{10, 3, rest},
		{10, 4, rest},
		{10, 5, rest},
		{10, 6, rest},
		{10, 7, rest},
		{10, 8, rest},
	},
}

//对holidayExhaustingInAMonth函数进行优化
//思路:
//如果休假map中的最后一个假期距给定日期小于30天,返回true
func holidayBalanceByGivenDate(date time.Time) holidayBalanceFlag {
	if holidaysMap == nil || len(holidaysMap) == 0 {
		//map尚未被维护
		return exhausted //直接返回已耗尽
	}
	//获取已维护的最大假期
	maxYear := 0
	for year := range holidaysMap {
		if year > maxYear {
			maxYear = year
		}
	}

	//获取当年最大假期
	hArray := holidaysMap[maxYear]
	if hArray == nil || len(hArray) == 0 {
		warn("年份表中每个被维护的年份应该至少有一个假期")
		return exhausted //年份应该至少维护一个假期,否则不应维护
	}
	var maxHoliday holiday
	for i := len(hArray) - 1; i > 0; i-- {
		h := hArray[i]
		if h.t == rest {
			maxHoliday = h
		}
	}
	//计算两天之间的日期数
	daysMinus := daysBetweenTimes(date, time.Date(maxYear, time.Month(maxHoliday.m), maxHoliday.d, 0, 0, 0, 0, time.Local))
	if daysMinus <= 0 {
		return exhausted
	}
	if daysMinus < 30 {
		return exhausting
	}
	return enough
}

//根据今天时间,判断是否即将耗尽
func holidayBalanceByNow() holidayBalanceFlag {
	return holidayBalanceByGivenDate(now)
}
