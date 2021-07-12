package edl2

import (
	"testing"
)

var sampleLines = []string{
	"001  AX       V     U        00:10:51:21 00:10:57:13 00:00:00:00 00:00:05:17",
	//"* FROM CLIP NAME: Young_Rock_s01e02__HD.mp4",
}

func TestNewStandard(t *testing.T) {
	for _, line := range sampleLines {
		ss, err := NewStandard(line)
		if ss == nil {
			continue
		}
		//expected := &standard{}
		switch {
		case err != nil:
			t.Errorf("%v\n  result is not '*standard' type %v", line, ss)

		}
		// switch {
		// default:
		// 	t.Errorf("%v\n  result is not '*standard' type %v", line, ss)
		// case ss.Type() == "STANDARD" && isStandard(line):

		// }
		// if err != nil {
		// 	t.Errorf("%v\nError: %v", line, err)
		// }
	}
}
