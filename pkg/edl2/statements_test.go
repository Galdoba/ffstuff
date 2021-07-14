package edl2

import (
	"testing"
)

var sampleLines = []string{
	"001  AX       V     H        00:10:51:21 00:10:57:13 00:00:00:00 00:00:05:17",
	"001  AX       V     C        00:10:51:21 00:10:57:13 00:00:00:00 00:00:05:17",
	"* FROM CLIP NAME: Young_Rock_s01e02__HD.mp4",
}

func TestNewStandard(t *testing.T) {
	for _, line := range sampleLines {
		ss, err := newStandard(line)
		if ss != nil {
			if err == nil && ss.Type() != STATEMENT_STANDARD {
				t.Errorf("Error expected but not passed")
			}
			if ss.editType == "C" && ss.editDuration != 0 {
				t.Errorf("editDuration = %v, expected 0", ss.editDuration)
			}
			if ss.editType != "C" && ss.editDuration == 0 {
				t.Errorf("editDuration = 0, expected to be float [float64] > 0 (%v from line %v)", ss.editDuration, line)
			}
		}
	}
}
