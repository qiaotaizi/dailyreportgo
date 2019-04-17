package main

import (
	"time"
)

//判断当前日期是否需要GitLab信息
func needGitLabData() bool {
	return time.Now().Weekday() == time.Tuesday
}
