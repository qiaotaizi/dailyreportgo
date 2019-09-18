package main

import (
	"fmt"
	"reflect"
	"testing"
)

type myStruct struct {
	a string
	b int
	c bool
	subStruct
}

type subStruct struct {
	d float64
	e rune `must:"true"`
	subSubStruct
}

type subSubStruct struct {
	f int `must:"true"`
}

func TestStuctCheckMust(t *testing.T){
	var s myStruct
	s.f=1
	v:=reflect.ValueOf(s)
	fmt.Println(structCheckMust(v))
}

func TestCmdEmpty(t *testing.T){
	c:=new(cmd)
	got:=c.empty()
	if !got{
		t.Errorf("cmd应为空,实际上为%t\n",got)
	}

	c2:=new(cmd)
	c2.reporterName="姜志恒"
	got=c2.empty()
	if got{
		t.Errorf("cmd应非空,实际上为%t\n",got)
	}
}

func TestCheckMust(t *testing.T){
	c:=new(cmd)
	field,ok:=c.checkMust()
	fmt.Println(field,ok)
}
