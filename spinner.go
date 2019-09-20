package main

import (
	"fmt"
	"time"
)

func spinner(delay time.Duration) {
	symbols := `-\|/`

	//为了运用锁
	//把打印单独摘出来作为一个方法
	prt := func(r rune) {
		outputLock.Lock()
		defer outputLock.Unlock()
		fmt.Printf("\r日报生成中 %c", r)
	}

	for {
		for _, r := range symbols {
			prt(r)
			time.Sleep(delay)
		}
	}
}
