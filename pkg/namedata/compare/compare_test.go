package compare

import (
	"fmt"
	"testing"
)

func TestCompare(t *testing.T) {
	inputs := []string{
		`\\192.168.31.4\buffer\IN\Test_serial_s01e01.mp4`,
	}
	for _, inp := range inputs {
		fmt.Println("INPUT:", inp)
		name := SuggestNameFromTable(inp)
		fmt.Println(inp, "=>", name)

	}
}
