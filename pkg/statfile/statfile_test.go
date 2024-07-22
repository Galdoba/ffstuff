package statfile

import (
	"fmt"
	"testing"
)

func TestCreate(t *testing.T) {
	sf, err := Create(Default_DIR, "test", "aaa")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	sf.AddOption("bbb", "ddd")
	sf.AddWeight("bbb", 1)
	fmt.Println(sf)
	fmt.Println(sf.List())
}
