package main

//模板控制

//日报填充器
type filler struct {
	Receiver     string `must:"true" usage:"收件人email列表"`
	Cc           string `usage:"抄送人email列表"`
	DepName      string `must:"true" usage:"部门名称"`
	Date         string `usage:"日报日期"`
	ReporterName string `must:"true" usage:"报告人姓名"`
	Achievements string `usage:"当天工作成果"`
	PAndS        string `usage:"碰到问题&解决方案"`
	Targets      string `usage:"明天工作计划"`
}
