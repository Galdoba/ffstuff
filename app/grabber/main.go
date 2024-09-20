package main

import (
	"fmt"
	"os"
	"time"

	"github.com/Galdoba/ffstuff/app/grabber/commands"
	"github.com/Galdoba/ffstuff/pkg/stdpath"
	"github.com/urfave/cli/v2"
)

func init() {
	stdpath.SetAppName("grabber")
}

func main() {

	app := cli.NewApp()
	app.Version = "0.0.1:" + dateTag()
	app.Usage = "utility for copy/move operations management"
	app.UsageText = "grabber [global options] command [command options] [arguments]..."
	app.Description = "grabber allows copy/move files by direct command or in delayed manner via crone-like shedule \n" +
		"or loop with dormancy periods."
	app.Flags = []cli.Flag{}
	stdpath.SetAppName(app.Name)

	//ДО НАЧАЛА ДЕЙСТВИЯ
	app.Before = func(c *cli.Context) error {

		return nil
	}
	app.Commands = []*cli.Command{
		commands.Grab(),
		commands.Search(),
		commands.Queue(),
		commands.Run(),
		commands.Setup(),
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

func dateTag() string {
	tm := time.Now()
	dur := time.Since(tm)
	dur.Hours()
	return time.Now().Format("060102")
}
