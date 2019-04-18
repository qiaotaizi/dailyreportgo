package holiday

import "time"

//手动维护一套节假日

const (
	Rest = 1 + iota //法定节假日休息
	Work            //法定节假日调休工作
)

type Holiday struct {
	M int  //月份
	D int  //日期
	T int  //类型 Rest/Work
}

var HolidayMap = map[int][]Holiday{
	2019:
	{
		{1, 1, Rest},
		{2, 2, Work},
		{2, 3, Work},
		{2, 4, Rest},
		{2, 5, Rest},
		{2, 6, Rest},
		{2, 7, Rest},
		{2, 8, Rest},
		{2, 9, Rest},
		{2, 10, Rest},
		{4, 5, Rest},
		{4, 6, Rest},
		{4, 7, Rest},
		{4, 28, Work},
		{5, 1, Rest},
		{5, 2, Rest},
		{5, 3, Rest},
		{5, 4, Rest},
		{5, 5, Work},
		{6, 7, Rest},
		{6, 8, Rest},
		{6, 9, Rest},
		{9, 13, Rest},
		{9, 14, Rest},
		{9, 15, Rest},
		{10, 1, Rest},
		{10, 2, Rest},
		{10, 3, Rest},
		{10, 4, Rest},
		{10, 5, Rest},
		{10, 6, Rest},
		{10, 7, Rest},
		{10, 7, Rest},
		{10, 12, Work},
	},
	//2020年后续完善
}

//检查当前时间是不是即将耗尽已维护的假期
//未定义年份返回true
//当前日期之后没有假期返回true
func HolidayExhaustingInAMonth(date time.Time) bool {
	hs,ok:=HolidayMap[date.Year()]
	if !ok{
		return true
	}
	//nd:= yearLastDay(date)
	//dur:=nd.Sub(date)
	//if dur<24*time.Hour*30{
	//	return true
	//}
	for _,h:=range hs{
		if h.T==Rest{
			//发现后续假期
			if time.Month(h.M)>date.Month(){
				return false
			}
			if time.Month(h.M)==date.Month() && h.D>date.Day(){
				return false
			}
		}
	}
	return true
}

//获取年底
//func yearLastDay(date time.Time) time.Time {
//	return date.AddDate(1,int(1-date.Month()),-date.Day())
//}

