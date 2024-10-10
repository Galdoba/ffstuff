package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/Galdoba/ffstuff/app/grabber/commands"
	"github.com/Galdoba/ffstuff/app/grabber/config"
	"github.com/Galdoba/ffstuff/pkg/stdpath"
	"github.com/urfave/cli/v2"
)

func init() {
	stdpath.SetAppName("grabber")
}

func main() {

	app := cli.NewApp()
	app.Version = "0.1.0:" + dateTag()
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
		commands.Setup(),
		commands.Health(),
		commands.Search(),
		commands.Queue(),
		commands.Run(),
	}

	//ПО ОКОНЧАНИЮ ДЕЙСТВИЯ
	app.After = func(c *cli.Context) error {
		return nil
	}
	args := os.Args
	if err := app.Run(args); err != nil {
		errOut := fmt.Sprintf("%v error: %v", app.Name, err.Error())
		println(errOut)
		if errors.Is(err, config.ErrNoConfig) {
			fmt.Printf("suggestion: run 'grabber setup'")
		}
		os.Exit(1)
	}

}

func dateTag() string {
	tm := time.Now()
	dur := time.Since(tm)
	dur.Hours()
	return time.Now().Format("060102")
}
