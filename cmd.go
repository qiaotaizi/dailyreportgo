package main

import (
	"flag"
	"reflect"
)

const commandName = "dailyreportgo"

type cmd struct {
	filler
	Help    bool //显示帮助
	Verbose bool //啰嗦模式

	JiraUserName string `must:"true" usage:"jira登录用户名"`
	JiraPwd      string `must:"true" usage:"jira登录用户密码"`

	History bool
}

func parseCmd() *cmd {
	result := new(cmd)
	flag.BoolVar(&result.Help, "h", false, "显示帮助")
	flag.StringVar(&result.Receiver, "r", "", "收件人email,若希望有多个收件人email,请使用';'分割")
	flag.BoolVar(&result.Verbose, "v", false, "开启啰嗦模式,输出日志. 啰嗦模式或静默模式都将输出警告")
	flag.StringVar(&result.Cc, "c", "", "抄送人email,若希望有多个抄送人email,请使用';'分割")
	flag.StringVar(&result.DepName, "d", "", "报告人部门名称")
	flag.StringVar(&result.ReporterName, "n", "", "报告人姓名")
	flag.StringVar(&result.JiraUserName, "un", "", "jira用户名")
	flag.StringVar(&result.JiraPwd, "up", "", "jira用户密码")
	flag.BoolVar(&result.History, "history", false, "显示历史命令,最多100条")
	flag.Parse()
	return result
}

//检查cmd中是否所有成员都是0值
//都是零值,返回true
//不完全是零值,返回false
func (c *cmd) empty() bool {
	v := reflect.ValueOf(c).Elem() //注意这里必须有.Elem(),否则将panic
	return v.IsZero()
}

//检查必填项
//返回false,检查失败,并返回必填参数的usage标签,若不存在标签,返回字段名
//返回true,检查通过
func (c *cmd) checkMust() (string, bool) {
	v := reflect.ValueOf(c).Elem()

	return structCheckMust(v)
}

//校验结构体所有字段是否标记了must,如果有的话,检查这个字段值是否为0值
//若true,表示通过检查,若false,表示未通过检查
func structCheckMust(structValue reflect.Value) (string, bool) {
	structType := structValue.Type()
	for i := 0; i < structValue.NumField(); i++ {
		ft := structType.Field(i)
		fv := structValue.Field(i)
		if fv.Kind() == reflect.Struct {
			if failField, ok := structCheckMust(fv); !ok {
				return failField, false
			}
			continue
		}

		//检查是否有must标记
		must := ft.Tag.Get("must")
		if must != "true" {
			continue
		}
		if !fv.IsZero() {
			continue
		}
		usage := ft.Tag.Get("usage")
		if usage != "" {
			return usage, false
		}
		//must为true且该字段是零值
		return ft.Name, false
	}
	return "", true
}
