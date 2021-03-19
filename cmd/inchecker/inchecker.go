package main

import (
	"fmt"
	"os"

	"github.com/Galdoba/ffstuff/fldr"
	"github.com/Galdoba/ffstuff/pkg/inchecker"
	"github.com/Galdoba/ffstuff/pkg/logfile"
)

func main() {
	fmt.Println("Create at:", fldr.MuxPath()+"logfile.txt")
	logger := logfile.New(fldr.MuxPath()+"logfile.txt", logfile.LogLevelINFO)
	fmt.Println("GO INCHECKER")
	checker := inchecker.NewChecker()
	for _, path := range pathsReceived() {
		checker.AddTask(path)
		logger.TRACE("Checking: " + path)
	}
	fmt.Println("GO INCHECKER CHECK")
	allErrors := checker.Check()
	fmt.Println("GO INCHECKER REPORT")
	checker.Report()
	if len(allErrors) == 0 {
		//logger.INFO("All files valid")
	}
	for _, err := range allErrors {
		logger.WARN(err.Error())
	}
	fmt.Println("END")
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
