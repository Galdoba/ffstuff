package main

import (
	"os"
	"strings"

	"github.com/Galdoba/ffstuff/fldr"
	"github.com/Galdoba/ffstuff/pkg/glog"
	"github.com/Galdoba/ffstuff/pkg/inchecker"
)

func main() {
	logger := glog.New(fldr.MuxPath()+"logfile.txt", glog.LogLevelINFO)
	checker := inchecker.NewChecker()
	pathsReceived := pathsReceived()
	if len(pathsReceived) == 0 {
		logger.TRACE("No arguments received")
		return
	}
	for _, path := range pathsReceived {
		if strings.Contains(path, ".ready") {
			continue
		}
		checker.AddTask(path)
		logger.TRACE("Checking: " + path)
	}
	allErrors := checker.Check()
	checker.Report(allErrors)
	if len(allErrors) == 0 {
		if len(pathsReceived) > 1 {
			logger.INFO("All files valid")
		}
		os.Exit(0)
	}
	for _, err := range allErrors {
		logger.WARN(err.Error())
	}
}

func pathsReceived() []string {
	outArgs := []string{}
	for i, val := range os.Args {
		if i == 0 {
			continue
		}
		outArgs = append(outArgs, val)
	}
	return outArgs
}
