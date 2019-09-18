package main

import (
	"fmt"
	"log"
	"os"
)

//模板控制

var templateFile *os.File

func init() {
	if testFlag{
		return
	}
	//检查userhome/Documents/dailyReport/dr.template文件是否存在
	//若不存在,创建之
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("用户目录获取失败,程序终止: %v\n", err);
	}
	const TemplateFileName = "dr.template"
	const docDir = "Documents"
	const drDir = "dailyReport"
	targetFile := fmt.Sprintf("%s%c%s%c%s%c%s",
		home, os.PathSeparator,
		docDir, os.PathSeparator,
		drDir, os.PathSeparator,
		TemplateFileName)
	templateFile, err = os.Open(targetFile)
	if templateFile != nil {
		lg("发现模板文件%s\n", targetFile)
		return //找到文件,直接结束函数
	}
	if !os.IsNotExist(err) {
		//如果返回的err不是文件存在,停止进程
		log.Fatalf("打开文件异常: %v\n", err)
	}
	//文件不存在,创建之,并写入默认模板
	templateFile, err = os.Create(targetFile)
	if err != nil {
		log.Fatalf("创建文件异常: %v\n", err)
	}

	const templateDefaultContent = `收件人：
{{.Receiver}}

抄送：
{{.Cc}}

主题：
{{.DepName}}日报 - {{.Date}} - {{.ReporterName}}

内容：

{{.DepName}}日报 - {{.Date}} - {{.ReporterName}}
----------------------------------------------------------------------------------------------------------------------
当天工作成果
{{.Achievements}}
----------------------------------------------------------------------------------------------------------------------
碰到问题&解决方案
{{.PAndS}}
----------------------------------------------------------------------------------------------------------------------
明天工作计划
{{.Targets}}
`
	_,err=templateFile.WriteString(templateDefaultContent)
	if err!=nil{
		log.Fatalf("模板文件写入异常: %v\n", err)
	}
	lg("已创建模板文件%s\n", targetFile)
	//写入后重新打开
	templateFile.Close()
	templateFile,err=os.Open(targetFile)
	if err!=nil{
		log.Fatalf("创建并重新打开模板文件异常: ",err)
	}
}

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
