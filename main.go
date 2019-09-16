package main

import (
	"flag"
	"fmt"
	"log"
)

//重构这个项目
//使用命令行参数进行各项参数的输入

var lg func(message string,args ...interface{})

func lgVerbose(message string,args ...interface{}){
	if len(args)==0{
		log.Println(message)
	}else{
		log.Printf(message,args...)
	}
}

func lgSilence(message string,args ...interface{}){
	//什么也不做
}

func main(){
	c:=parseCmd()

	if c.help || c.empty(){
		flag.Usage()
		return
	}
	if c.verbose{
		lg=lgVerbose
	}else{
		lg=lgSilence
	}
	fmt.Println(c.receiver)
}



var reportConfig *drCfg

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

