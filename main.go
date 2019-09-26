package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/go-toast/toast"
)

//重构这个项目
//使用命令行参数进行各项参数的输入

var outputLock sync.Mutex
var warns []string

func warn(message string, args ...interface{}) {
	outputLock.Lock()
	defer outputLock.Unlock()
	content := fmt.Sprintf(message, args...)
	fmt.Printf("警告: %s\n", content)
	warns = append(warns, content)
}

func warnNotify() {
	//如果关闭了系统通知,不执行下面的程序
	if !params.OsWarn {
		return
	}
	if warns == nil || len(warns) == 0 {
		return
	}
	var message string
	if len(warns) == 1 {
		message = warns[0]
	} else {
		var buf bytes.Buffer
		for i, w := range warns {
			buf.WriteString(fmt.Sprintf("%d.%s\n", i+1, w))
		}
		message = buf.String()
	}
	notify := toast.Notification{
		AppID:   "Daily.Report.Generator",
		Title:   "日报警告",
		Message: message,
	}
	_ = notify.Push()
}

func main() {

	//defer recordCmd()

	defer warnNotify() //收集所有警告,并调用系统通知

	defer releaseResources() //释放资源

	//打印帮助,或者命令行参数长度为1
	//也就是用户命令为dailyreportgo或者dailyreportgo -h
	if len(os.Args) == 1 || params.Help {
		fmt.Printf("%s命令的参数及其含义: \n", commandName)
		flag.PrintDefaults()
		fmt.Println("你可以通过修改下面命令的参数, 快速生成日报:")
		fmt.Printf(`%s -d="你所在的部门" -n="你的名字" -c="抄送人邮箱" -r="收件人邮箱" -un="你的jira登录用户名" -up="你的jira登录密码"
`, commandName)
		return
	}

	fmt.Println("-----",params)

	//执行生成命令前,校验必填项
	failField, err := params.checkMust()

	if failField != "" {
		warn("字段%s必填,请使用%s -h查看帮助", failField, commandName)
		return
	}
	if err != nil {
		warn("参数校验时发生错误: %v\n", err)
		return
	}

	//所有校验完成,开始生成日报
	go spinner(100 * time.Millisecond)

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
