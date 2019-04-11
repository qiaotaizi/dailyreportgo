package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"os/user"
	"reflect"
	"strconv"
	"strings"
)

type drCfg struct {
	JiraLoginName              string `xml:"jira_login_name" nilable:"false"`
	JiraLoginPwd               string `xml:"jira_login_pwd" nilable:"false"`
	JiraLoginParam             string `xml:"jira_login_param" nilable:"false"`
	JiraLoginUrl               string `xml:"jira_login_url" nilable:"true" default:"http://jira.ttpai.cn/login.jsp"`
	JiraCalendarForUserUrl     string `xml:"jira_calendar_for_user_url" nilable:"true" default:"http://jira.ttpai.cn/rest/mailrucalendar/1.0/calendar/forUser"`
	JiraCalendarJsonRequestUrl string `xml:"jira_calendar_json_request_url" nilable:"true" default:"http://jira.ttpai.cn/rest/mailrucalendar/1.0/calendar/events/%d?start=%s&end=%s"`
	ReporterName               string `xml:"reporter_name" nilable:"false"`
	DepartmentName             string `xml:"department_name" nilable:"false"`
	EmailReceivers             string `xml:"email_receivers" nilable:"false"`
	EmailCcReceivers           string `xml:"email_cc_receivers" nilable:"false"`
	CodeReviewReceiver         string `xml:"code_review_receiver" nilable:"true"`
}

const configFileName = "dailyReportGO.xml"

var configPositions []string

func init() {
	//获取当前用户
	sysUser, err := user.Current();
	if err != nil {
		log.Fatalf("config init: %v", err)
	}
	//初始化文件位置
	//user_home/Documents/DailyReport
	configPositions = append(configPositions, strings.Join([]string{sysUser.HomeDir, "Documents", "dailyReport", configFileName}, string(os.PathSeparator)))
}

//如果配置不存在则创建配置
func createConfigIfNotExist() {
	configFileExist := false
	for i := 0; i < len(configPositions); i++ {
		_, err := os.Stat(configPositions[i])
		if err == nil {
			//文件存在
			configFileExist = true
			log.Printf("%s found, use it", configPositions[i])
			break
		}
	}
	if !configFileExist {
		//在默认位置创建文件
		f, err := os.Create(configPositions[0])
		defer f.Close()
		if err != nil {
			log.Fatalf("creating config file %s failed", configPositions[0])
		}
		_, err = fmt.Fprintf(f, dailyReportConfig)
		if err != nil {
			log.Fatalf("writing config file %s failed", configPositions[0])
		}
		log.Printf("%s created", configPositions[0])
	}
}

//读取配置
func parseConfig() *drCfg {
	cfg := new(drCfg)
	dailyReportGOXml, err := os.Open(configPositions[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "parseConfig: %v", err)
		os.Exit(1)
	}
	defer dailyReportGOXml.Close()
	if err := xml.NewDecoder(dailyReportGOXml).Decode(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "parseConfig: %v", err)
		os.Exit(1)
	}
	return cfg
}

//非空配置项校验
func antiInvalidConfig(cfg *drCfg) {
	typeOfCfg := reflect.TypeOf(cfg)
	valueOfCfg := reflect.ValueOf(cfg)
	//入参是指针,这里获取对象,否则下面获取字段值的时候会报错
	typeOfCfg = typeOfCfg.Elem()
	valueOfCfg = valueOfCfg.Elem()
	for i := 0; i < typeOfCfg.NumField(); i++ {
		ft := typeOfCfg.Field(i)
		fv := valueOfCfg.Field(i)
		nilable, err := strconv.ParseBool(ft.Tag.Get("nilable"))
		if err != nil {
			log.Fatalf("config field parse bool fail: %v", err)
		}
		if !nilable && fv.String() == "" {
			log.Fatalf("config %s can't be nil", ft.Name)
		}
		defaultValue, ok := ft.Tag.Lookup("default")
		if ok {
			//有默认值
			fv.SetString(defaultValue)
		}
	}
}
