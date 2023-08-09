package fdf

import (
	"fmt"
	"testing"
)

func TestFDF(t *testing.T) {
	fd := newFD(`\\192.168.31.4\buffer\IN\_DONE\bash\amedia_I_prosto_tak.sh`)
	fmt.Println(fd)
}
