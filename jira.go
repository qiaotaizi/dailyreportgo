package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	sessionIdName              = "JSESSIONID"
	jiraLoginParam             = "登录"
	jiraLoginUrl               = "http://jira.ttpai.cn/login.jsp"
	jiraCalendarForUserUrl     = "http://jira.ttpai.cn/rest/mailrucalendar/1.0/calendar/forUser"
	jiraCalendarJsonRequestUrl = "http://jira.ttpai.cn/rest/mailrucalendar/1.0/calendar/events/%d?start=%s&end=%s"
)

var sessionIdCookie *http.Cookie

//自定义jira登录异常
type jiraLoginFailErr struct{}

func (e jiraLoginFailErr) Error() string {
	return "jira登录失败异常, 请检查用户名和密码"
}

//jira登录,让sessionIdCookie变量管理登录信息
func jiraLogin() error {
	lg("登录jira,用户名: %s\n", params.JiraUserName)

	form := url.Values{
		"os_username": {params.JiraUserName},
		"os_password": {params.JiraPwd},
		"login":       {jiraLoginParam},
	}
	formString := form.Encode()
	req, err := http.NewRequest(http.MethodPost, jiraLoginUrl, strings.NewReader(formString))
	if err != nil {
		warn("构建jira登录请求失败: %v", err)
		return jiraLoginFailErr{}
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := jiraHttpClient.Do(req)
	defer resp.Body.Close()
	if err != nil {
		log.Fatalf("jira登录失败: %v\n", err)
	}

	contentBytes, _ := ioutil.ReadAll(resp.Body)
	content := string(contentBytes)
	if strings.Contains(content, "对不起,你的用户名或者密码不正确") {
		return jiraLoginFailErr{}
	}

	cks := jiraHttpClient.Jar.Cookies(req.URL)
	sessionIdFound := false
	for _, ck := range cks {
		if ck.Name == sessionIdName {
			sessionIdCookie = ck
			sessionIdFound = true
			break
		}
	}
	if !sessionIdFound {
		return jiraLoginFailErr{}
	}
	lg("jira登录成功")
	return nil
}

//获取jira用户数据
//返回用户id
func jiraForUser() (int, bool) {
	req, err := http.NewRequest(http.MethodGet, jiraCalendarForUserUrl, nil)
	if err != nil {
		warn("构建jira用户数据请求失败: %v", err)
		return 0, false
	}
	req.AddCookie(sessionIdCookie)
	resp, err := jiraHttpClient.Do(req)
	defer resp.Body.Close()
	if err != nil {
		warn("获取jira用户数据失败: %v", err)
		return 0, false
	}
	var users []jiraUserVo
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		warn("用户数据json数据转对象失败: %v", err)
		return 0, false
	}
	for _, user := range users {
		if strings.HasPrefix(user.Name, params.ReporterName) {
			return user.Id, true
		}
	}
	warn("未找到报告人姓名, 将以jira返回的首个用户id为准")
	return users[0].Id, true
}

//获取任务列表
func jiraCalendarMission(userId int) ([]jiraMissionVo, bool) {
	start, end := monthStartAndEnd()
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(jiraCalendarJsonRequestUrl, userId, start, end), nil)
	if err != nil {
		warn("构建jira任务数据请求失败: %v", err)
		return nil, false
	}
	req.AddCookie(sessionIdCookie)
	resp, err := jiraHttpClient.Do(req)
	defer resp.Body.Close()
	if err != nil {
		warn("jira任务数据请求失败: %v", err)
		return nil, false
	}
	var missions []jiraMissionVo
	if err := json.NewDecoder(resp.Body).Decode(&missions); err != nil {
		warn("jira任务列表json转换失败: %v", err)
		return nil, false
	}
	return missions, true
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
	Id               string   `json:"id"`
	Status           string   `json:"status"`
	StatusColor      string   `json:"statusColor"`
	CalendarId       int      `json:"calendarId"`
	AllDay           bool     `json:"allDay"`
	Color            string   `json:"color"`
	End              jsonTime `json:"end"`
	Start            jsonTime `json:"start"`
	Title            string   `json:"title"`
	DurationEditable bool     `json:"durationEditable"`
	StartEditable    bool     `json:"startEditable"`
	DatesError       bool     `json:"datesError"`
}

type jsonTime time.Time

func (t jsonTime) MarshalJSON() ([]byte, error) {
	stamp := time.Time(t).Format(dateFormat)
	return []byte(stamp), nil
}

func (t jsonTime) String() string {
	return time.Time(t).String()
}

func (t *jsonTime) UnmarshalJSON(data []byte) error {
	//返回一个零值时间
	if string(data) == "null" {
		return nil
	}
	tm, err := time.ParseInLocation(fmt.Sprintf("%q", dateFormat), string(data), time.Local)
	if err != nil {
		return err
	}
	*t = jsonTime(tm)
	return nil
}

//任务当前日期执行中
func (mission jiraMissionVo) inProgress() bool {
	return !time.Time(mission.Start).After(now) && !time.Time(mission.End).Before(now)
}

//任务在下个工作日执行中
func (mission jiraMissionVo) inProgressNextWorkDay() bool {
	return !time.Time(mission.Start).After(nextWorkDay) && !time.Time(mission.End).Before(nextWorkDay)
}

//计算任务进度
//返回整数类型的进度百分比点数
//比如进度为100%,返回100

func (mission jiraMissionVo) progress(targetDate time.Time) int {
	start := time.Time(mission.Start)
	end := time.Time(mission.End)
	//计算start到end有几个工作日,计算start到now有几个工作日
	step := 24 * time.Hour
	//初始化分子和分母
	var numerator, denominator int
	for date := start; date.Before(end); date = date.Add(step) {
		if !isWorkDay(date) {
			continue
		}
		denominator++
		if date.Before(targetDate) {
			numerator++
		}
	}
	return numerator * 100 / denominator
}
