package main

import (
	"log"
	"time"
)

var reportConfig *drCfg

func main() {
	if !isWorkDay(time.Now()){

	}
	log.Println("welcome!")
	//解析配置(配置不存在时生成一套默认配置,和模板)
	//使用xml进行配置
	//如果不存在配置,生成一套默认的
	createConfigIfNotExist()
	//xml转对象
	reportConfig = parseConfig()
	//空值校验
	antiInvalidConfig(reportConfig)
	log.Println("config validate ok")
	log.Println("generating daily report according to jira")
	//登录jira获取cookie
	jiraLogin()
	//获取用户信息
	userId := jiraForUser()
	//获取calendar数据
	missions:=jiraCalendarMission(userId)
	//周二获取gitlab结对人提交数据
	if needGitLabData() {
		//登录gitlab

		//循环分页获取gitlab提交信息,直到时间到达上周二
	}

	//检查模板是否存在,如果不存在,生成一套默认的
	createTemplatesIfNotExist()
	//获取模板 填充模板 生成日报
	fillTemplate(missions)

	//使用指定应用打开
}

