package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"
	"time"
)

const (
	dailyReportTemplate = `收件人：
{{.receivers}}

抄送：
{{.ccReceivers}}

主题：
{{.departmentName}}日报 - {{.reportDate}} - {{.reporterName}}

内容：
{{.departmentName}}日报 - {{.reportDate}} - {{.reporterName}}
----------------------------------------------------------------------------------------------------------------------
当天工作成果
{{range .todayAchievements}}
{{.num}}.{{.item}}
{{end}}
----------------------------------------------------------------------------------------------------------------------
碰到问题&解决方案
1.无
----------------------------------------------------------------------------------------------------------------------
明天工作计划
{{range .todayAchievements}}
{{.num}}.{{.item}}
{{end}}
`

	codeReviewTemplate = `
review报告：
-------------------------------------------------------------------------------------------
日期：周二{{.date01}}
项目：{{.date01Projects}}
未发现问题/未找到code review
------------------------------------------------------------------------------------------
日期：周三{{.date02}}
项目：{{.date02Projects}}
未发现问题/未找到code review
------------------------------------------------------------------------------------------
日期：周四{{.date03}}
项目：{{.date03Projects}}
未发现问题/未找到code review
------------------------------------------------------------------------------------------
日期：周五{{.date04}}
项目：{{.date04Projects}}
未发现问题/未找到code review
------------------------------------------------------------------------------------------
日期：周一{{.date05}}
项目：{{.date05Projects}}
未发现问题/未找到code review
`
)

const reportTemplateName = "dailyReportGO.template"

const codeReviewTemplateName = "codeReviewGO.template"

var reportTemplatePosition string

var codeReviewTemplatePosision string

var reportGeneratePosition string

func init() {
	reportTemplatePosition = strings.Join([]string{userHome, "Documents", "dailyReport", reportTemplateName}, string(os.PathSeparator))
	codeReviewTemplatePosision = strings.Join([]string{userHome, "Documents", "dailyReport", codeReviewTemplateName}, string(os.PathSeparator))
	reportGeneratePosition = strings.Join([]string{userHome, "Documents", "dailyReport", "%s日报.dr"}, string(os.PathSeparator))
}

//如果模板不存在则生成一套默认的
func createTemplatesIfNotExist() {
	//日报模板
	if _, err := os.Stat(reportTemplatePosition); os.IsNotExist(err) {
		f, err := os.Create(reportTemplatePosition)
		if err != nil {
			log.Fatalf("creating template file %s failed", reportTemplatePosition)
		}
		_, err = fmt.Fprintf(f, dailyReportTemplate)
		if err != nil {
			log.Fatalf("writing config file %s failed", reportTemplatePosition)
		}
	}

	//code review模板
	if _, err := os.Stat(codeReviewTemplatePosision); os.IsNotExist(err) {
		f, err := os.Create(codeReviewTemplatePosision)
		if err != nil {
			log.Fatalf("creating template file %s failed", codeReviewTemplatePosision)
		}

		_, err = fmt.Fprintf(f, codeReviewTemplate)
		if err != nil {
			log.Fatalf("writing config file %s failed", codeReviewTemplatePosision)
		}
	}
}

//将内容填充至模板并输出文件
func fillTemplate(missions []jiraMissionVo) {
	rDate := time.Now().Format("2006年1月2日")
	reportContent := reportContent{
		receivers:      reportConfig.EmailReceivers,
		ccReceivers:    reportConfig.EmailCcReceivers,
		departmentName: reportConfig.DepartmentName,
		reportDate:     rDate,
		reporterName:   reportConfig.ReporterName,
	}
	//筛选任务,填入todayAchievements和tomorrowTargets
	selectMissionsIntoReportContent(&reportContent, missions)
	t, _ := template.ParseFiles(dailyReportTemplate)
	reportFile, err := os.Create(fmt.Sprintf(reportGeneratePosition, rDate))
	if err != nil {
		log.Fatalf("creating report file %s failed", reportGeneratePosition)
	}
	t.Execute(reportFile, reportContent)
}

func selectMissionsIntoReportContent(rc *reportContent, missions []jiraMissionVo) {
	c1, c2 := 0, 0
	for _, m := range missions {
		if m.inProgress() {
			c1++
			rc.todayAchievements = append(rc.todayAchievements, reportListVo{
				c1,
				fmt.Sprintf(
					"%s http://jira.ttpai.cn/browse/%s %d%%",
					m.Title,
					m.Id,

				)})
		}
		if m.inProgressNextWorkDay() {
			c2++
			rc.tomorrowTargets = append(rc.tomorrowTargets, reportListVo{c2, m.Title})
		}
	}
	if len(rc.todayAchievements) == 0 {
		rc.todayAchievements = append(rc.todayAchievements, reportListVo{1, "无"})
	}
	if len(rc.tomorrowTargets) == 0 {
		rc.tomorrowTargets = append(rc.tomorrowTargets, reportListVo{1, "无"})
	}
}

type reportContent struct {
	receivers         string
	ccReceivers       string
	departmentName    string
	reportDate        string
	reporterName      string
	todayAchievements []reportListVo
	tomorrowTargets   []reportListVo
}

type reportListVo struct {
	num  int
	item string
}
