package main

import (
	"fmt"
	"testing"

	"github.com/Galdoba/ffstuff/pkg/stdpath"
)

func TestBuildTag(t *testing.T) {
	tmn := dateTag()
	fmt.Println("v 0.0.1:" + tmn)
	fmt.Println(stdpath.ConfigFile())
}
