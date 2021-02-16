package main

import (
	"fmt"
	"os"

	"github.com/Galdoba/ffstuff/pkg/inchecker"
)

func main() {
	checker := inchecker.NewChecker()
	for _, path := range pathsReceived() {
		checker.AddTask(path)
	}
	checker.Check()
	checker.Report()
	// for i, val := range os.Args {
	// 	if len(os.Args) == 1 {
	// 		fmt.Println("No Arguments received")
	// 	}
	// 	if i == 0 {
	// 		continue
	// 	}
	// 	fmt.Print("File:	", val)
	// 	checker.AddTask(val)
	// 	if err := checker.CheckValidity(val); err != nil {
	// 		fmt.Println("\n------------------------------------------------------------")
	// 		fmt.Println(err.Error())
	// 		fmt.Println("------------------------------------------------------------")
	// 		continue
	// 	}
	// 	fmt.Print(" . . . ok\n")
	// }
}

func pathsReceived() []string {
	outArgs := []string{}
	for i, val := range os.Args {
		if len(os.Args) == 1 {
			fmt.Println("No Arguments received")
		}
		if i == 0 {
			continue
		}
		outArgs = append(outArgs, val)
	}
	return outArgs
}
