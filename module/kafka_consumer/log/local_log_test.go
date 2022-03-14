package log

import (
	"os"
	"regexp"
	"testing"
)

func TestLocalLogger(t *testing.T) {
	pwd, _ := os.Getwd()
	logPath := pwd + string(os.PathSeparator) + "main.log"
	Info("test")
	Infof("%s %s %s %s", "test1", "test2", "test3", "test4")
	file, err := os.ReadFile(logPath)
	if err != nil {
		t.Errorf("os.ReadFile %s", err.Error())
		return
	}

	match, err := regexp.Match(`\"M\":\"\[test\]\"`, file)
	if err != nil {
		t.Errorf("regexp.Match %s", err.Error())
		return
	}

	if !match {
		t.Errorf("match faild")
	}
}
