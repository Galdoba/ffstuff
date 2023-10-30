package main

import (
	"fmt"
	"os"

	"github.com/Galdoba/ffstuff/pkg/fdf"
	"github.com/Galdoba/ffstuff/pkg/mdm/inputinfo"
)

func main() {
	//
	args := os.Args

	if len(args) > 2 {
		fmt.Println("profiler failed: to many arguments")
		return
	}
	if len(args) < 2 {
		fmt.Println("profiler failed: not enough arguments")
		return
	}
	file := args[1]
	f, err := os.Stat(file)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if f.IsDir() {
		fmt.Println("profiler failed: can't make profile of a directory")
		return
	}
	pi, err := inputinfo.ParseFile(file)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(pi.String())
	prf := fdf.FMP(pi)

	fmt.Println(prf)

}
