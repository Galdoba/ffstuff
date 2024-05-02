package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

/*
run


*/

var configPath string

const (
	programName = "aamed"
)

/*
работаем в режиме демона, если можно создать амедиевский скрипт - делаем его
*/

func main() {

	app := cli.NewApp()
	app.Version = "v 0.0.1"
	app.Name = programName
	app.Usage = "работаем в режиме демона, если можно создать амедиевский скрипт - делаем его"
	app.Flags = []cli.Flag{}

	app.Before = func(c *cli.Context) error {
		return nil
	}
	app.Commands = []*cli.Command{

		// cmd.NewCommand(),
	}
	app.DefaultCommand = "run"

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
