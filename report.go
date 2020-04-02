package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"sync"
	"text/template"
	"time"
)

//输出报告文本
func genReportString() (string, error) {
	//从文件中读取模板
	bs, err := ioutil.ReadAll(templateFile)
	if err != nil {
		return "", err
	}

	templateText := string(bs)

	parsedTmpl, err := template.New("dailyReport").Parse(templateText)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)

	//登录jira拿cookie
	if err = jiraLogin(); err != nil {
		return "", err
	}

	//获取用户id
	userId, ok := jiraForUser()
	if !ok {
		return "", errors.New("获取jira用户id失败")
	}
	//获取calendar数据
	missions, ok := jiraCalendarMission(userId)
	if !ok {
		return "", errors.New("获取jira任务表失败")
	}

	fillerComplete := func(f *filler) {
		f.Date = now.Format("2006年1月2日")
		f.PAndS = "1.无"
		//并发拼接任务字符串
		missionOutPutFormat := "%d. %s %s %d%%\n"
		missionNextWorkdayOutPutFormat:="%d. %s %s\n"

		var wg sync.WaitGroup
		collector := func(missions []jiraMissionVo, stringField *string, boolMethod func(jiraMissionVo) bool, targetDate time.Time,showPercent bool) {
			defer wg.Done()
			var missionToday bytes.Buffer
			counter := 0
			for _, mission := range missions {
				if boolMethod(mission) {
					counter++
					progress := mission.progress(targetDate)

					if showPercent {
						missionToday.WriteString(fmt.Sprintf(missionOutPutFormat, counter, mission.Title, mission.Id, progress))
						lg("任务名: %s,任务开始日期: %s,任务结束日期: %s,目标日期: %s,计算进度: %d%%", mission.Title, mission.Start, mission.End, targetDate, progress)
					}else{
						missionToday.WriteString(fmt.Sprintf(missionNextWorkdayOutPutFormat, counter, mission.Title, mission.Id))
						lg("任务名: %s,任务开始日期: %s,任务结束日期: %s,目标日期: %s", mission.Title, mission.Start, mission.End, targetDate)

					}

				}
			}
			result := missionToday.String()
			*stringField = result
		}
		//负责收集当天任务,并写入相应字段
		wg.Add(1)
		go collector(missions, &(f.Achievements), jiraMissionVo.inProgress, now,true)
		//负责收集下个工作日任务,并写入相应字段
		wg.Add(1)
		go collector(missions, &(f.Targets), jiraMissionVo.inProgressNextWorkDay, nextWorkDay,false)
		wg.Wait()
	}

	fillerComplete(&params.filler)

	err = parsedTmpl.Execute(buf, params.filler)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

//将文本写入文件,并返回文件路径
func writeReportIntoFile(content string) (string, error) {

	fileName := genReportFileName()
	fullPath := fmt.Sprintf("%s%c%s%c%s%c%s",
		userHome,
		filepath.Separator,
		docDir,
		filepath.Separator,
		drDir,
		filepath.Separator,
		fileName)
	err := ioutil.WriteFile(fullPath, []byte(content), 0644)
	if err != nil {
		return "", err
	}
	return fullPath, nil
}

//生成日报文件名
func genReportFileName() string {
	return fmt.Sprintf("%s日报.gdr", now.Format("2006年1月2日"))
}

//windows使用记事本打开文件
func openReportFileWithNotepad(filePath string) error {
	c := exec.Command("notepad", filePath)
	return c.Start()
}

//mac使用sublime打开文件
func openReportFileWithSublimeText(filepath string) error {
	c := exec.Command("/Applications/Sublime Text.app/Contents/MacOS/Sublime Text", filepath)
	return c.Start()
}
