package main

import (
	"bytes"
	"io/ioutil"
	"text/template"
)

//输出报告文本
func genReportString()(string,error){
	//从文件中读取模板
	bs,err:=ioutil.ReadAll(templateFile)
	if err!=nil{
		return "",err
	}

	templateText:=string(bs)
	lg("模板文件中的内容:\n%s\n",templateText)

	parsedTmpl,err:=template.New("dailyReport").Parse(templateText)
	if err!=nil{
		return "",err
	}

	buf:=new(bytes.Buffer)

	err=parsedTmpl.Execute(buf, params.filler)
	if err!=nil{
		return "",err
	}

	return buf.String(), nil
}
