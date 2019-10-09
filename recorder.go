package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const latestFileName = "latest.json"

//命令记录

//将本次命令写入文件
func recordCmd() {

	commandLine2Record := func() []byte {
		var buf bytes.Buffer
		for _, arg := range os.Args[1:] {
			buf.WriteString(arg)
			buf.WriteRune(' ')
		}
		buf.WriteRune('\n')
		return buf.Bytes()
	}

	fileName := fmt.Sprintf(
		"%s%c%s%c%s%c%s",
		userHome, os.PathSeparator,
		docDir, os.PathSeparator,
		drDir, os.PathSeparator,
		latestFileName)
	//读取文件
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		if !os.IsNotExist(err) {
			//若有异常且异常非文件不存在,终止该线程
			warn("读取命令记录文件时出现问题,-history参数可用性将受到影响: %v", err)
			return
		}
		//文件不存在,给data赋初始值,长度为128字节
		data = make([]byte, 128)
	}
	//将历史命令以行为单位读出,检查行数
	historyContent := string(data)
	historyArr := strings.Split(historyContent, "\n")
	//当前命令
	currentCommand := commandLine2Record()

	if len(historyArr) >= 100 {
		fmt.Println("11111")
		//超过100行时,加入本次记录并使用最近100行覆盖文本
		historyArr = append(historyArr, string(currentCommand))
		var buf bytes.Buffer
		for _, history := range historyArr[1:] {
			buf.WriteString(history)
			buf.WriteRune('\n')
		}
		err := ioutil.WriteFile(fileName, buf.Bytes(), 0644)
		if err != nil {
			warn("历史命令超出100行,命令记录文件写入时出现问题,-history参数可用性将受到影响: %v", err)
			return
		}
	} else {
		//未超过100行时,直接将本次参数拼接进去
		f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			warn("命令记录文件打开写入时出现问题,-history参数可用性将受到影响: %v", err)
			return
		}
		defer f.Close()
		_, err = f.Write(currentCommand)
		if err != nil {
			warn("命令记录文件打开写入时出现问题,-history参数可用性将受到影响: %v", err)
		}
	}
}

//打印历史记录
func printHistory() {
	fileName := fmt.Sprintf(
		"%s%c%s%c%s%c%s",
		userHome, os.PathSeparator,
		docDir, os.PathSeparator,
		drDir, os.PathSeparator,
		latestFileName)
	//读取文件
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		if !os.IsNotExist(err) {
			//若有异常且异常非文件不存在,终止该线程
			warn("读取命令记录文件时出现问题: %v", err)
			return
		} else {
			fmt.Println("暂无记录")
			return
		}
	}
	fmt.Println(string(data))
}
