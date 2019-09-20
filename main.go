package main

import (
	"flag"
	"fmt"
)

//重构这个项目
//使用命令行参数进行各项参数的输入

func warn(message string, args ...interface{}) {
	fmt.Printf("%c[0;31m警告: %s%c[0m\n", 0x1B, fmt.Sprintf(message, args...), 0x1B)
}

func main() {
	defer func() func() {
		now_ := now
		nwd_ := nextWorkDay
		return func() {
			now__ := now
			nwd__ := nextWorkDay
			if now_ != now__ {
				warn("全局变量now在程序运行期间发生了变化")
			}
			if nwd_ != nwd__ {
				warn("全局变量nwd在程序运行期间发生了变化")
			}
		}
	}()()

	defer releaseResources() //释放资源

	//打印帮助
	if params.Help || params.empty() {
		fmt.Printf("%s命令的参数及其含义: \n", commandName)
		flag.PrintDefaults()
		fmt.Println("你可以通过修改下面命令的参数, 快速生成日报:")
		fmt.Println(`go run . -d="你所在的部门" -n="你的名字" -c="抄送人邮箱" -r="收件人邮箱" -un="你的jira登录用户名" -up="你的jira登录密码"`)
		return
	}

	//执行生成命令前,校验必填项
	if failField, ok := params.checkMust(); !ok {
		warn("字段%s必填,请使用%s -h查看帮助", failField, commandName)
		return
	}

	reportContent, err := genReportString()
	if err != nil {
		warn("生成日报文本时发生异常: %v", err)
		return
	}
	//将结果写入文件,并且调用notepad打开文件
	filePath, err := writeReportIntoFile(reportContent)
	if err != nil {
		warn("写入日报文本时异常: %v", err)
	}
	if err = openReportFileWithNotepad(filePath); err != nil {
		warn("使用记事本打开日志文件时异常%v", err)
	}

}

func releaseResources() {
	err := templateFile.Close()
	if err != nil {
		lg("释放模板文件时异常: %v\n", err)
		return
	}
	lg("已成功释放模板文件资源")
}
