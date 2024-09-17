package main

import (
	"fmt"
	"os"

	"github.com/Galdoba/ffstuff/app/archivator/commands"
	"github.com/urfave/cli/v2"
)

func main() {

	app := cli.NewApp()

	app.Version = "0.0.1"
	app.Usage = "track/send/retrive files from archve"
	app.Description = "TODO: Description"
	app.Flags = []cli.Flag{}

	//ДО НАЧАЛА ДЕЙСТВИЯ
	app.Before = func(c *cli.Context) error {

		return nil
	}
	app.Commands = []*cli.Command{
		commands.Health(),
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
