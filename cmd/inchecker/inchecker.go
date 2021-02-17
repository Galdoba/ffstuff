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

}

func pathsReceived() []string {
	outArgs := []string{}
	for i, val := range os.Args {
		if len(os.Args) == 1 {
			fmt.Println("No аrguments received")
		}
		if i == 0 {
			continue
		}
		outArgs = append(outArgs, val)
	}
	return outArgs
}
