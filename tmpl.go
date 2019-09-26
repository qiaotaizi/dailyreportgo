package main

//模板控制

//日报填充器
type filler struct {
	Receiver     string `cmd:"r" default:"" name:"收件人email列表" must:"true" usage:"收件人email列表"`
	Cc           string `cmd:"c" default:"" name:"抄送人email列表" must:"false" usage:"抄送人email列表"`
	DepName      string `cmd:"d" default:"" name:"部门名称" must:"true" usage:"部门名称"`
	Date         string `name:"日报日期"`
	ReporterName string `cmd:"n" default:"" name:"报告人姓名" must:"true" usage:"报告人姓名"`
	Achievements string `name:"当天工作成果"`
	PAndS        string `name:"碰到问题&解决方案"`
	Targets      string `name:"明天工作计划"`
}
