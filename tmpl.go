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
{{.receiver}}

抄送：
{{.cc}}

主题：
{{.depName}}日报 - {{.date}} - {{.reporterName}}

内容：

{{.depName}}日报 - {{.date}} - {{.reporterName}}
----------------------------------------------------------------------------------------------------------------------
当天工作成果
{{.achievements}}
----------------------------------------------------------------------------------------------------------------------
碰到问题&解决方案
{{.pAndS}}
----------------------------------------------------------------------------------------------------------------------
明天工作计划
{{.targets}}
`
	templateFile.WriteString(templateDefaultContent)
	lg("已创建模板文件%s\n", targetFile)
}

//日报填充器
type filler struct {
	receiver     string `must:true` //收件人email列表
	cc           string //抄送人email列表
	depName      string `must:true`//部门名称
	date         string //日期
	reporterName string `must:true`//报告人姓名
	achievements string //当天工作成果
	pAndS        string //碰到问题&解决方案
	targets      string //明天工作计划
}
