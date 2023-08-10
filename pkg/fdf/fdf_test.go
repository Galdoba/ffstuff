package fdf

import (
	"fmt"
	"testing"
)

func TestFDF(t *testing.T) {
	fd := newFD(`\\nas\ROOT\EDIT\_amedia\Krovnye_s01\Krovnye_s01_05_PRT230809182151.srt`)
	fmt.Println(fd)
}
