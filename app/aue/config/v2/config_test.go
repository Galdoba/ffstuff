package v2

import (
	"fmt"
	"testing"
)

func TestMain(t *testing.T) {
	ex := NewConfig("0.1.2")
	bt, err := marshal(ex)
	if err != nil {
		t.Errorf(err.Error())
	}
	fmt.Println("--")
	fmt.Println(string(bt))
	fmt.Println("--")
}
