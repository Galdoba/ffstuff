package edl2

import (
	"testing"
)

func TestSum(t *testing.T) {
	a, b, s := 5, 2, 8
	c := Sum(a, b)
	// if c != s {
	// 	t.Errorf("unexpested c=%d", c)
	// }
	if c == s {
		t.Errorf("even more unexpested c=%d", c)
	}
}
