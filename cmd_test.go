package main

import (
	"fmt"
	"testing"
)

func TestCmdEmpty(t *testing.T){
	c:=new(cmd)
	//c.reporterName="姜志恒"
	got:=c.empty()
	if !got{
		t.Errorf("cmd应为空,实际上为%t\n",got)
	}
}

func TestCheckMust(t *testing.T){
	c:=new(cmd)
	field,ok:=c.checkMust()
	fmt.Println(field,ok)
}
