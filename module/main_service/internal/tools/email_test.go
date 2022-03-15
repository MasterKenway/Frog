package tools

import "testing"

func TestSendEmail(t *testing.T) {
	err := SendEmail("northcountrys@163.com", "TEST")
	if err != nil {
		t.Errorf("failed to send email, %s", err.Error())
	}
}
