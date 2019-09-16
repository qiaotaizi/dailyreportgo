package main

import (
	"flag"
	"reflect"
)

type cmd struct {
	receiver string //收件人email列表
	cc       string //抄送人email列表
	help     bool   //显示帮助
	verbose  bool   //啰嗦模式
}

func parseCmd() *cmd {
	result := new(cmd)
	flag.BoolVar(&result.help, "h", false, "显示帮助")
	flag.StringVar(&result.receiver, "r", "", "收件人email,若希望有多个收件人email,请使用';'分割")
	flag.BoolVar(&result.verbose, "v", false, "开启罗嗦模式,输出日志")
	flag.StringVar(&result.cc,"c","","抄送人email,若希望有多个抄送人email,请使用';'分割")
	flag.Parse()
	return result
}

//检查cmd中是否所有成员都是0值
//都是零值,返回true
//不完全是零值,返回false
func (c *cmd) empty() bool{
	v:=reflect.ValueOf(c).Elem()//注意这里必须有.Elem(),否则将panic
	for i:=0;i<v.NumField();i++{
		f:=v.Field(i)
		fv:=reflect.ValueOf(f)
		switch fv.Kind() {
		case reflect.Bool:
			if fv.Bool(){//bool类型参数,如果为true,则不是默认值
				return false
			}
		case reflect.String:
			if fv.String()!=""{
				return false
			}
		}
	}
	return true
}
