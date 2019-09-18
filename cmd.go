package main

import (
	"flag"
	"fmt"
	"reflect"
)

const commandName="dailyreport"

type cmd struct {
	filler
	help    bool //显示帮助
	verbose bool //啰嗦模式

	jiraUserName string `must:true` //jira登录用户名
	jiraPwd      string `must:true`//jira登录用户密码
}

func parseCmd() *cmd {
	result := new(cmd)
	flag.BoolVar(&result.help, "h", false, "显示帮助")
	flag.StringVar(&result.receiver, "r", "", "收件人email,若希望有多个收件人email,请使用';'分割")
	flag.BoolVar(&result.verbose, "v", false, "开启罗嗦模式,输出日志")
	flag.StringVar(&result.cc, "c", "", "抄送人email,若希望有多个抄送人email,请使用';'分割")
	flag.StringVar(&result.depName, "d", "", "报告人部门名称")
	flag.StringVar(&result.reporterName, "n", "", "报告人姓名")
	flag.StringVar(&result.jiraUserName,"un","","jira用户名")
	flag.StringVar(&result.jiraPwd,"up","","jira用户密码")
	flag.Parse()
	return result
}

//检查cmd中是否所有成员都是0值
//都是零值,返回true
//不完全是零值,返回false
func (c *cmd) empty() bool {
	v := reflect.ValueOf(c).Elem() //注意这里必须有.Elem(),否则将panic
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		if !f.IsZero(){
			return false;
		}
		//fv := reflect.ValueOf(f)
		//switch fv.Kind() {
		//case reflect.Bool:
		//	if fv.Bool() { //bool类型参数,如果为true,则不是默认值
		//		return false
		//	}
		//case reflect.String:
		//	if fv.String() != "" {
		//		return false
		//	}
		//}
	}
	return true
}

//检查必填项
//返回false,检查失败,并返回必填的参数名
//返回true,检查通过
func (c *cmd) checkMust() (string,bool) {
	t:=reflect.TypeOf(c).Elem()
	v := reflect.ValueOf(c).Elem()
	fmt.Println(t.NumField())
	fmt.Println(v.NumField())
	for  i := 0; i < t.NumField(); i++{
		tf:=t.Field(i)
		must:=tf.Tag.Get("must")
		if must!="true"{
			continue
		}
		//must为true,校验当前字段的值是否为零值
		vf:=v.Field(i)
		fmt.Printf("%s is zero?%t\n",tf.Name,vf.IsZero())
		if vf.IsZero(){
			return tf.Name,false
		}
	}
	return "",true
}
