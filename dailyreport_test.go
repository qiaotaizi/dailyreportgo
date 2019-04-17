package main

import (
	"fmt"
	"testing"
)

func TestNeedGitLabData(t *testing.T) {
	b := needGitLabData()
	fmt.Println(b)
}

func TestCreateTemplatesIfNotExist(t *testing.T) {
	createTemplatesIfNotExist()
}
