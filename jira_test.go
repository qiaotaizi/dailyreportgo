package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestJsonDecode(t *testing.T) {
	str := `[{"id":"TECH-12636","status":"To Do","statusColor":"blue-gray","calendarId":158,"allDay":true,"color":"#5dab3e","end":"2019-09-20","start":"2019-09-19","title":"【开发自测任务】【二组】短链接微服务调用迁移至 MS-COMMUNAL --二期","durationEditable":true,"startEditable":true,"datesError":false},{"id":"TECH-12635","status":"In Progress","statusColor":"yellow","calendarId":158,"allDay":true,"color":"#5dab3e","end":"2019-09-19","start":"2019-09-17","title":"【开发任务】【二组】短链接微服务调用迁移至 MS-COMMUNAL --二期","durationEditable":true,"startEditable":true,"datesError":false},{"id":"TECH-12479","status":"To Do","statusColor":"blue-gray","calendarId":158,"allDay":true,"color":"#5dab3e","end":"2019-09-28","start":"2019-09-26","title":"【联调任务】成交模块代看模块重构","durationEditable":true,"startEditable":true,"datesError":false},{"id":"TECH-12478","status":"To Do","statusColor":"blue-gray","calendarId":158,"allDay":true,"color":"#5dab3e","end":"2019-09-26","start":"2019-09-20","title":"【开发任务】成交模块代看模块重构","durationEditable":true,"startEditable":true,"datesError":false},{"id":"INVITATION-2474","status":"Done","statusColor":"green","calendarId":158,"allDay":true,"color":"#5dab3e","end":"2019-09-04","start":"2019-08-30","title":"【PRD阅读+业务文档梳理整理】[车主意向录入]邀约环节车主互动1.0","durationEditable":true,"startEditable":true,"datesError":false},{"id":"INVITATION-2473","status":"Done","statusColor":"green","calendarId":158,"allDay":true,"color":"#5dab3e","end":"2019-09-17","start":"2019-09-11","title":"【前后端联调任务】[车主意向录入]邀约环节车主互动1.0","durationEditable":true,"startEditable":true,"datesError":false},{"id":"INVITATION-2468","status":"Done","statusColor":"green","calendarId":158,"allDay":true,"color":"#5dab3e","end":"2019-09-11","start":"2019-09-04","title":"【开发任务】[车主意向录入]邀约环节车主互动1.0","durationEditable":true,"startEditable":true,"datesError":false}]`
	buf := new(bytes.Buffer)
	buf.WriteString(str)
	var missions []jiraMissionVo
	err := json.NewDecoder(buf).Decode(&missions)
	if err != nil {
		t.Errorf("%v\n", err)
		return
	}
	fmt.Println("ok")
}

func TestJsonDecodeSingle(t *testing.T) {
	str := `{"id":"TECH-12636","status":"To Do","statusColor":"blue-gray","calendarId":158,"allDay":true,"color":"#5dab3e","end":"2019-09-20","start":"2019-09-19","title":"【开发自测任务】【二组】短链接微服务调用迁移至 MS-COMMUNAL --二期","durationEditable":true,"startEditable":true,"datesError":false}`
	buf := new(bytes.Buffer)
	buf.WriteString(str)
	var missions jiraMissionVo
	err := json.NewDecoder(buf).Decode(&missions)
	if err != nil {
		t.Errorf("%v\n", err)
		return
	}
	fmt.Println(missions)
	fmt.Println("ok")
}

func TestProgress(t *testing.T) {
	start, _ := time.Parse(dateFormat, "2019-09-20")
	end, _ := time.Parse(dateFormat, "2019-09-26")
	var mission jiraMissionVo
	mission.Start = jsonTime(start)
	mission.End = jsonTime(end)
	progress := mission.progress(now)
	fmt.Printf("%d%%\n", progress)
}
