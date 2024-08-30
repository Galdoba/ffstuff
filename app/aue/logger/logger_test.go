package logger

import "testing"

func TestLogger(t *testing.T) {
	lgr := New()
	logg = lgr
	Func1()
}
