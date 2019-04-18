package main

import (
	"fmt"
	"testing"
	"time"
)

func TestNeedGitLabData(t *testing.T) {
	b := needGitLabData()
	fmt.Println(b)
}

func TestCreateTemplatesIfNotExist(t *testing.T) {
	createTemplatesIfNotExist()
}

func TestHoliday(t *testing.T){
	ds:=[]string{
		"2019-05-02",
		"2019-05-03",
		"2019-10-12",
		"2019-04-28",
		"2019-04-04",
		"2019-04-21",
	}
	for _,dstr:=range ds{
		d,err:=time.Parse(format,dstr)
		if err!=nil{
			fmt.Println(err)
		}
		b:=isWorkDay(d)
		fmt.Printf("%s %v\n",dstr,b)
	}
}
