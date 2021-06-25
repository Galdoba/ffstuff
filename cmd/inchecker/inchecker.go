package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/Galdoba/ffstuff/fldr"
	"github.com/Galdoba/ffstuff/pkg/glog"
	"github.com/Galdoba/ffstuff/pkg/inchecker"
	"github.com/urfave/cli"
)

func main() {
	logger := glog.New(fldr.MuxPath()+"logfile.txt", glog.LogLevelINFO)
	checker := inchecker.NewChecker()
	pathsReceived := pathsReceived()
	app := cli.NewApp()
	app.Version = "v 0.0.2"
	app.Name = "inchecker"
	app.Usage = "Checks media files for standard format"
	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:  "vocal",
			Usage: "If flag is active searcher will print ALL log entries (level INFO is set by default)",
		},
		&cli.BoolFlag{
			Name:  "clear, cl",
			Usage: "If flag is active searcher Clear Terminal before every Search",
		},
		&cli.StringFlag{
			Name:  "delay",
			Usage: "If flag is active searcher will delay start for N seconds",
		},
	}
	args := os.Args
	if len(args) < 2 {
		args = append(args, "help") //Принудительно зовем помощь если нет других аргументов
	}
	if err := app.Run(args); err != nil {
		fmt.Println(err.Error())
	}

	//////////////
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
