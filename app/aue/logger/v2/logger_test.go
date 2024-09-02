package v2

import "testing"

func TestLogMan(t *testing.T) {
	Setup()
	if logMan.appLoglevel == -1 {
		t.Errorf("app log level was not set")
	}
}
