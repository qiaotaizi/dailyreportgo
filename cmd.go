package main

import (
	"flag"
	"fmt"
	"reflect"
)

const commandName = "dailyreportgo"

type cmd struct {
	//tag说明:
	//cmd:命令行参数名,必须填写
	//name:警告时显示的命令行参数解释
	//must:是否必填
	//usage:打印命令帮助时显示的参数解释
	//default:默认值,请注意如果must是true,default不应设为对应类型的零值,否则在做参数校验时会有意想不到的问题
	filler
	Help    bool `cmd:"h" default:"false" name:"显示帮助控制符" must:"false" usage:"显示帮助"`
	Verbose bool `cmd:"v" default:"false" name:"啰嗦模式控制符" must:"false" usage:"开启啰嗦模式,输出日志. 啰嗦模式或静默模式都将输出警告."`

	JiraUserName string `cmd:"un" default:"" name:"jira登录用户名" must:"true" usage:"jira登录用户名"`
	JiraPwd      string `cmd:"up" default:"" name:"jira登录用户密码" must:"true" usage:"jira登录用户密码"`

	Latest bool `cmd:"latest" default:"false" name:"显示最近命令控制符" must:"false" usage:"打印最近执行的命令列表"`

	OsWarn bool `cmd:"warn" default:"true" name:"操作系统通知控制符" must:"false" usage:"操作系统通知开关,默认开启"`
}

func parseCmd() (*cmd, error) {
	result := new(cmd)

	//使用用反射初始化参数表
	//v:=reflect.ValueOf(result).Elem()//获取反射对象
	//for i:=0;i<v.NumField() ;i++  {
	//	f:=v.Field(i) //反射属性
	//	kind:=f.Kind() //属性类型
	//	switch kind {
	//	case reflect.String:
	//		flag.BoolVar(&f.Bool(),)
	//
	//	}
	//}

	hook := func(fieldType reflect.StructField, fieldValue *reflect.Value) (string, error) {
		tag := fieldType.Tag
		cmdTag := tag.Get("cmd")
		defaultTag := tag.Get("default")
		usageTag := tag.Get("usage")
		//该成员没有cmd标签,跳过设置
		if cmdTag == "" {
			return "", nil
		}
		switch fieldValue.Kind() {
		case reflect.Bool:
			boolElem := fieldValue.Bool()
			flag.BoolVar(&boolElem, cmdTag, defaultTag == "true", usageTag)
			fmt.Println("aaa---",boolElem,fieldValue.Bool())
		case reflect.String:
			stringElem := fieldValue.String()
			flag.StringVar(&stringElem, cmdTag, defaultTag, usageTag)
			fmt.Println("aaa---",stringElem,fieldValue.String())
		default:
			return fieldType.Name, fmt.Errorf("不识别的属性类型: %s.请在parseCmd函数中进行完善", fieldValue.Kind())
		}
		return "", nil
	}

	if failField, err := iterateStructField(reflect.ValueOf(result).Elem(), hook); err != nil {
		return nil, fmt.Errorf("parseCmd操作%s时失败: %v", failField, err)
	}
	//flag.BoolVar(&result.Help, "h", false, "显示帮助")
	//flag.StringVar(&result.Receiver, "r", "", "收件人email,若希望有多个收件人email,请使用';'分割")
	//flag.BoolVar(&result.Verbose, "v", false, "开启啰嗦模式,输出日志. 啰嗦模式或静默模式都将输出警告")
	//flag.StringVar(&result.Cc, "c", "", "抄送人email,若希望有多个抄送人email,请使用';'分割")
	//flag.StringVar(&result.DepName, "d", "", "报告人部门名称")
	//flag.StringVar(&result.ReporterName, "n", "", "报告人姓名")
	//flag.StringVar(&result.JiraUserName, "un", "", "jira用户名")
	//flag.StringVar(&result.JiraPwd, "up", "", "jira用户密码")
	//flag.BoolVar(&result.Latest, "latest", false, "打印最近执行的10条命令记录")
	//flag.BoolVar(&result.OsWarn, "warn", true, "是否在程序运行完毕时将收集的警告信息发送给操作系统通知,默认开启")
	flag.Parse()
	return result, nil
}

//检查必填项
//若返回"",nil 表明结构体通过检查
//若返回"xx",err 表明结构体未通过检查,xx字段的name标签或者字段名
//若返回"",err 表明检查结构体时出现了其他异常
func (c *cmd) checkMust() (string, error) {
	v := reflect.ValueOf(c).Elem()

	return structCheckMust(v)
}

//校验结构体所有字段是否标记了must,如果有的话,检查这个字段值是否为0值
//若返回"",nil 表明结构体通过检查
//若返回"xx",err 表明结构体未通过检查,xx字段的name标签或者字段名
//若返回"",err 表明检查结构体时出现了其他异常
func structCheckMust(structValue reflect.Value) (string, error) {

	//定义钩子
	//若返回"",nil 表明字段通过检查
	//若返回"xx",err 表明字段未通过检查,xx字段的name标签或者字段名
	//若返回"",err 表明检查字段中出现了其他异常
	hook := func(ft reflect.StructField, fv *reflect.Value) (string, error) {
		//检查是否有must标记
		must := ft.Tag.Get("must")
		if must != "true" {
			return "", nil //非must通过检查
		}
		if !fv.IsZero() {
			return "", nil //非零值通过检查
		}
		name := ft.Tag.Get("name")
		if name != "" {
			//usage非空,返回usage
			return name, fmt.Errorf("%s未通过检查", name)
		}
		//usage为空,返回参数名
		return ft.Name, fmt.Errorf("%s未通过检查", ft.Name)
	}

	if failField, err := iterateStructField(structValue, hook); err != nil {
		return failField, err
	}
	return "", nil
}
