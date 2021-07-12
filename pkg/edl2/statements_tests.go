package edl2

import (
	"fmt"
	"testing"
)

func TestNewStandard(t *testing.T) {
	line := "001  AX       V     C        00:10:51:21 00:10:57:13 00:00:00:00 00:00:05:17"
	ss, err := NewStandard(line)
	switch {
	default:
		fmt.Errorf("%v\n  result is not '*standard' type %v", line, ss)
	case ss.Type() == "STANDARD":
	}
	if err != nil {
		fmt.Errorf("%v\n  Error: %v", line, err)
	}
}
