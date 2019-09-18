package main

import (
	"flag"
	"fmt"
	"log"
)

//重构这个项目
//使用命令行参数进行各项参数的输入

//日志函数定义
//啰嗦模式下
//若args长度为0,message后会主动追加换行
//若args长度非零,message只会进行格式化,不会进行换行符追加
var lg func(message string,args ...interface{})

//啰嗦日志输出
func lgVerbose(message string,args ...interface{}){
	if len(args)==0{
		log.Println(message)
	}else{
		log.Printf(message,args...)
	}
}

//静默日志输出
func lgSilence(message string,args ...interface{}){
	//什么也不做
}

func warn(message string, args ...interface{}){
	fmt.Printf("%c[0;31m",0x1B)
	fmt.Printf("警告: "+message,args...)
	fmt.Printf("%c[0m",0x1B)
}

//命令行参数保存
var params *cmd

//使用这个标志在单元测试阶段关闭出发init函数
var testFlag=false

func init(){
	if testFlag{
		return
	}

	params =parseCmd()
	//控制日志输出模式
	if params.Verbose{
		lg=lgVerbose
	}else{
		lg=lgSilence
	}
}

func main(){
	defer releaseResources()//释放资源

	//打印帮助
	if params.Help || params.empty(){
		fmt.Printf("%s命令的参数及其含义: \n",commandName)
		flag.PrintDefaults()
		return
	}

	//执行生成命令前,校验必填项
	if failField,ok:= params.checkMust();!ok{
		warn("字段%s必填,请使用%s -h查看帮助\n",failField,commandName)
		return
	}

	reportContent,err:=genReportString()
	if err!=nil{
		warn("生成日报文本时发生异常: %v\n",err)
		return
	}
	lg("填充模板后的内容:\n%s\n",reportContent)
}

func releaseResources(){
	err:=templateFile.Close()
	if err!=nil{
		lg("释放模板文件时异常: %v\n",err)
	}
}



//var reportConfig *drCfg

//func main() {
//	if !isWorkDay(time.Now()){
//		log.Fatal("today is not work day, generating aborted")
//	}
//	log.Println("welcome!")
//	//解析配置(配置不存在时生成一套默认配置,和模板)
//	//使用xml进行配置
//	//如果不存在配置,生成一套默认的
//	createConfigIfNotExist()
//	//xml转对象
//	reportConfig = parseConfig()
//	//空值校验
//	antiInvalidConfig(reportConfig)
//	log.Println("config validate ok")
//	log.Println("generating daily report according to jira")
//	//登录jira获取cookie
//	jiraLogin()
//	//获取用户信息
//	userId := jiraForUser()
//	//获取calendar数据
//	missions:=jiraCalendarMission(userId)
//	//周二获取gitlab结对人提交数据
//	if needGitLabData() {
//		//登录gitlab
//
//		//循环分页获取gitlab提交信息,直到时间到达上周二
//	}
//
//	//检查模板是否存在,如果不存在,生成一套默认的
//	createTemplatesIfNotExist()
//	//获取模板 填充模板 生成日报
//	fillTemplate(missions)
//
//	//使用指定应用打开
//}

