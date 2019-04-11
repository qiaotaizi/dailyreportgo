package main

import (
	"log"
)

var reportConfig *drCfg

func main() {
	log.Println("welcome!")
	//解析配置(配置不存在时生成一套默认配置,和模板)
	//使用xml进行配置
	createConfigIfNotExist()
	reportConfig = parseConfig()
	antiInvalidConfig(reportConfig)

	log.Println("config validate ok")

	log.Println("generating daily report according to jira")
	//登录jira获取cookie
	jiraLogin()
	//获取用户信息
	userId := jiraForUser()
	//获取calendar数据
	jiraCalendarMission(userId)
	//获取模板

	//填充模板

	//生成日报

	//使用指定应用打开
}
