package main

import (
	"fmt"
	"os"

	"github.com/Galdoba/ffstuff/app/demux/commands"
	"github.com/urfave/cli/v2"
)

const (
	CONFIG = "cfg"
)

func main() {
	app := cli.NewApp()

	app.Version = "v 0.0.0"
	app.Usage = "demux media with prompts"
	app.Flags = []cli.Flag{}

	//ДО НАЧАЛА ДЕЙСТВИЯ
	app.Before = func(c *cli.Context) error {

		return nil
	}
	app.Commands = []*cli.Command{
		commands.Prompt(),
	}

	//ПО ОКОНЧАНИЮ ДЕЙСТВИЯ
	app.After = func(c *cli.Context) error {
		return nil
	}
	args := os.Args
	if err := app.Run(args); err != nil {
		errOut := fmt.Sprintf("%v error: %v", app.Name, err.Error())
		println(errOut)
		os.Exit(1)
	}

}
