package main

import (
	"fmt"
	"os"

	"github.com/Galdoba/ffstuff/pkg/mdm/inputinfo"
)

// func main() {
// 	inputinfo.CleanScanData()
// }

func main() {

	for _, arg := range os.Args {
		fmt.Println("=========")
		pi, err := inputinfo.ParseFile(arg)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		for _, w := range pi.Warnings() {
			fmt.Println("!!!   ", w)
		}
		fmt.Println("----------")
	}
}
