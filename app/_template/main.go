package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

const (
	programName = "mediaprofiler"
	/*
		gmscan

	*/
)

func main() {
	app := cli.NewApp()
	app.Version = "v 0.0.1"
	app.Name = programName
	app.Usage = "send message to telegram channel"
	app.Flags = []cli.Flag{}
	//ДО НАЧАЛА ДЕЙСТВИЯ
	app.Before = func(c *cli.Context) error {

		return nil
	}
	app.Commands = []*cli.Command{
		//ДЕЙСТВИЕ 1
		{
			Name:      "ДЕЙСТВИЕ 1",
			Usage:     "",
			ArgsUsage: "",
			Flags:     []cli.Flag{
				//ФЛАГИ ДЕЙСТВИЯ
			},
			Action: func(c *cli.Context) error {
				//ТЕЛО ДЕЙСТВИЯ
				return nil
			},
		},
	}

	//ПО ОКОНЧАНИЮ ДЕЙСТВИЯ
	app.After = func(c *cli.Context) error {
		return nil
	}
	args := os.Args
	if err := app.Run(args); err != nil {
		errOut := fmt.Sprintf("%v error: %v", programName, err.Error())
		println(errOut)
		os.Exit(1)
	}

}
