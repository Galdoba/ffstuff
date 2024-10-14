package logman

import (
	"fmt"
	"testing"
)

func TestColors(t *testing.T) {
	formatterInfo("")
}

func formatterInfo(format string, args ...interface{}) string {

	// colMap := make(map[int]string)
	// this := 0
	// last = this - 1
	// totalArgs := len(args)
	// currentArg := -1

	for _, r := range `e` {
		fmt.Println("")
		fmt.Println("string", string(r), "=", int(r))

	}

	return ""
}

/*
sumrunes :
% = 37
v = 118
d = 100
s = 115
f = 102
t = 116
q = 113


*/

/*
% = 37
*/

/*
%s  - white / greenish?
%d  - green /green yellow
%f  - orange/red yellow
%t  - dark blue
%x  - hi white
%c  - white / redish?
%U
%q  - hi orange
%p  - gray / dark white
*/
