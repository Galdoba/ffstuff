package main

import (
	"fmt"
	"os"

	"github.com/Galdoba/ffstuff/fldr"
	"github.com/Galdoba/ffstuff/pkg/glog"
	"github.com/Galdoba/ffstuff/pkg/inchecker"
)

func main() {
	logger := glog.New(fldr.MuxPath()+"logfile.txt", glog.LogLevelINFO)
	checker := inchecker.NewChecker()
	for _, path := range pathsReceived() {
		checker.AddTask(path)
		logger.TRACE("Checking: " + path)
	}
	allErrors := checker.Check()
	checker.Report()
	if len(allErrors) == 0 {
		if len(pathsReceived()) > 1 {
			logger.INFO("All files valid")
		}
		os.Exit(0)
	}
	for _, err := range allErrors {
		logger.WARN(err.Error())
	}
	os.Exit(1)
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
