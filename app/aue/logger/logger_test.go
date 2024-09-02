package logger

import "testing"

func TestLogger(t *testing.T) {
	Setup()

}

/*
LOG MESSAGES

type	stop app	print		color
FATAL	yes			stdout		red/dark
ERROR	no			stdout		red/hi
WARN	no			stdout		yellow
INFO	no			stderr		white
DEBUG	no			file/stderr	gray
TRACE	no			file/stderr	gray


*/
