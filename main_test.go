package main

import (
	"fmt"
	"testing"
)

func TestWarnColor(t *testing.T) {
	fmt.Println("start")
	warn("hello\n")
	fmt.Println("end")
}
