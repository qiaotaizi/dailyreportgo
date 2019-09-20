package main

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
)

func TestOpenReportFileWithNotepad(t *testing.T) {
	h, _ := os.UserHomeDir()
	sep := os.PathSeparator
	err := openReportFileWithNotepad(fmt.Sprintf("%s%c%s%c%s%c%s", h, sep, docDir, sep, drDir, sep, "2019年9月5日日报.dr"))
	if err != nil {
		fmt.Printf("%v\n", err)
	}
}

func TestCmd(t *testing.T) {
	h, _ := os.UserHomeDir()
	sep := os.PathSeparator
	c := exec.Command("notepad", fmt.Sprintf("%s%c%s%c%s%c%s", h, sep, docDir, sep, drDir, sep, "2019年9月5日日报.dr"))
	c.Start()
}

func TestGenReportFileName(t *testing.T) {
	fmt.Println(genReportFileName())
}

func TestFileMode(t *testing.T) {
	fmt.Printf("%x\n", uint32(os.ModeType))
}
