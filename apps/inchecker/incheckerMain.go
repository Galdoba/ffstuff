package main

import (
	"fmt"
	"os"

	"github.com/Galdoba/ffstuff/inchecker"
	"github.com/Galdoba/utils"
)

func main() {
	for i, val := range os.Args {
		if i == 0 {
			continue
		}
		fmt.Print("File:	", val)
		checker := inchecker.NewChecker()
		if err := checker.CheckValidity(val); err != nil {
			fmt.Println(utils.ASCIIColor("red", err.Error()))
			continue
		}
		fmt.Print("	ok\n")
	}
}
