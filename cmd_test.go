package main

import "testing"

func TestCmdEmpty(t *testing.T){
	c:=new(cmd)
	got:=c.empty()
	if !got{
		t.Errorf("cmd应为空,实际上为%t\n",got)
	}
}
