package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

var jiraHttpClient http.Client

const (
	sessionIdName = "JSESSIONID"
	jiraLoginParam="登录"
	jiraLoginUrl="http://jira.ttpai.cn/login.jsp"
	jiraCalendarForUserUrl="http://jira.ttpai.cn/rest/mailrucalendar/1.0/calendar/forUser"
	jiraCalendarJsonRequestUrl ="http://jira.ttpai.cn/rest/mailrucalendar/1.0/calendar/events/%d?start=%s&end=%s"
)
var sessionIdCookie *http.Cookie

//初始化httpClient,并设置cookie管理
func init() {
	jiraCookieJar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalf("jira http client init failed: %v", err)
	}
	jiraHttpClient = http.Client{Jar: jiraCookieJar}
}

//jira登录
func jiraLogin() {
	lg("登录jira,用户名: %s\n", c.jiraUserName)

	form := url.Values{
		"os_username": {c.jiraUserName},
		"os_password": {c.jiraUserName},
		"login":       {jiraLoginParam},
	}
	formString := form.Encode()
	req, err := http.NewRequest(http.MethodPost, jiraLoginUrl, strings.NewReader(formString))
	if err != nil {
		log.Fatalf("构建jira登录请求失败: %v\n", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := jiraHttpClient.Do(req)
	defer resp.Body.Close()
	if err != nil {
		log.Fatalf("jira登录失败: %v\n", err)
	}
	cks := jiraHttpClient.Jar.Cookies(req.URL)
	for _, ck := range cks {
		if ck.Name == sessionIdName {
			sessionIdCookie = ck
			break
		}
	}
}

//获取jira用户数据
//返回用户id
func jiraForUser() int {
	req, err := http.NewRequest(http.MethodGet, jiraCalendarForUserUrl, nil)
	if err != nil {
		log.Fatalf("构建jira用户数据请求失败: %v\n", err)
	}
	req.AddCookie(sessionIdCookie)
	resp, err := jiraHttpClient.Do(req)
	defer resp.Body.Close()
	if err != nil {
		log.Fatalf("获取jira用户数据失败: %v\n", err)
	}
	var users []jiraUserVo
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		log.Fatalf("用户数据json数据转对象失败: %v", err)
	}
	for _, user := range users {
		if strings.HasPrefix(user.Name, c.reporterName) {
			return user.Id
		}
	}
	return users[0].Id
}

//获取任务列表
func jiraCalendarMission(userId int) []jiraMissionVo {
	start, end := monthStartAndEnd()
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(jiraCalendarJsonRequestUrl, userId, start, end), nil)
	if err != nil {
		log.Fatalf("构建jira任务数据请求失败: %v", err)
	}
	req.AddCookie(sessionIdCookie)
	resp, err := jiraHttpClient.Do(req)
	defer resp.Body.Close()
	if err != nil {
		log.Fatalf("get jira missions request failed: %v", err)
	}
	var missions []jiraMissionVo
	if err := json.NewDecoder(resp.Body).Decode(&missions); err != nil {
		log.Fatalf("parse jira missions json failed: %v", err)
	}
	//时间转换,完成度计算
	for _, m := range missions {
		m.startTime, err = time.Parse(format, m.Start)
		if err != nil {
			log.Fatalf("parsing time %s error", m.Start)
		}
		m.endTime, err = time.Parse(format, m.End)
		if err != nil {
			log.Fatalf("parsing time %s error", m.End)
		}
	}
	return missions
}


//接收jira用户json数据
type jiraUserVo struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Color      string `json:"color"`
	Source     string `json:"source"`
	Changable  bool   `json:"changable"`
	Viewable   bool   `json:"viewable"`
	Visible    bool   `json:"visible"`
	Favorite   bool   `json:"favorite"`
	HasError   bool   `json:"hasError"`
	UsersCount int    `json:"usersCount"`
}

//接收jira任务json数据
type jiraMissionVo struct {
	Id               string `json:"id"`
	Status           string `json:"status"`
	StatusColor      string `json:"statusColor"`
	CalendarId       int    `json:"calendarId"`
	AllDay           bool   `json:"allDay"`
	Color            string `json:"color"`
	End              string `json:"end"`
	Start            string `json:"start"`
	Title            string `json:"title"`
	DurationEditable bool   `json:"durationEditable"`
	StartEditable    bool   `json:"startEditable"`
	DatesError       bool   `json:"datesError"`
	startTime        time.Time
	endTime          time.Time
}

//任务当前日期执行中
func (mission jiraMissionVo) inProgress() bool {
	now := time.Now()
	return !mission.startTime.After(now) && !mission.endTime.Before(now)
}

//任务在下个工作日执行中
func (mission jiraMissionVo) inProgressNextWorkDay() bool {
	return !mission.startTime.After(nextWorkDay) && !mission.startTime.Before(nextWorkDay)
}
