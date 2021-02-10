package main

import (
	"fmt"
	"os"

	"github.com/Galdoba/ffstuff/pkg/inchecker"
)

func main() {
	for i, val := range os.Args {
		if len(os.Args) == 1 {
			fmt.Println("No Arguments received")
		}
		if i == 0 {
			continue
		}
		fmt.Print("File:	", val)
		checker := inchecker.NewChecker()
		if err := checker.CheckValidity(val); err != nil {
			fmt.Println("\n------------------------------------------------------------")
			fmt.Println(err.Error())
			fmt.Println("------------------------------------------------------------")
			continue
		}
		fmt.Print(" . . . ok\n")
	}
}
