package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/Galdoba/ffstuff/app/audlite/commands"
	"github.com/Galdoba/ffstuff/app/audlite/config"
	"github.com/Galdoba/ffstuff/pkg/stdpath"
	"github.com/urfave/cli/v2"
)

func init() {
	stdpath.SetAppName("audlite")
}

func main() {
	app := cli.NewApp()
	app.Version = "0.1.0:" + dateTag()
	app.Usage = "utility for audio encoding management"
	app.UsageText = "audlite [global options] command [command options] [arguments]..."
	app.Description = "audlite used for creating ffmpeg task with short single line"
	app.Flags = []cli.Flag{}
	stdpath.SetAppName(app.Name)

	//ДО НАЧАЛА ДЕЙСТВИЯ
	app.Before = func(c *cli.Context) error {

		return nil
	}

	app.Commands = []*cli.Command{
		commands.Merge(),
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
			fmt.Printf("suggestion: run 'audlite setup'")
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
