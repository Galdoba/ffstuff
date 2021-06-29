package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/Galdoba/ffstuff/pkg/config"
	"github.com/Galdoba/ffstuff/pkg/glog"
	"github.com/urfave/cli"
)

var configMap map[string]string

var logger glog.Logger

func init() {
	conf, err := config.ReadProgramConfig("ffstuff")
	if err != nil {
		fmt.Println(err)
	}
	configMap = conf.Field
	if err != nil {
		switch err.Error() {
		case "Config file not found":
			fmt.Print("Expecting config file in:\n", conf.Path)
			os.Exit(1)
		}
	}
}

func main() {
	app := cli.NewApp()
	app.Version = "v 0.0.0"
	app.Name = "cutter"
	app.Commands = []cli.Command{
		{
			Name:        "cut",
			ShortName:   "",
			Aliases:     []string{},
			Usage:       "",
			UsageText:   "",
			Description: "",
			ArgsUsage:   "",
			Category:    "",
			BashComplete: func(*cli.Context) {
				fmt.Println("Start bashcomplete action")
			},
			Before: func(c *cli.Context) error {
				fmt.Println("Start before action")
				edlFound := 0
				for _, filepath := range c.Args() {
					f, err := os.Stat(filepath)
					switch {
					default:
						return errors.New("unknown before action error: " + err.Error())
					case err == nil:
						edlFound++
						fmt.Println(f)
						fmt.Println(f.Name(), "is valid file")
					case cannotFindFile(err):
						fmt.Println("Error: Can't find file specified:", filepath)
						fmt.Println("Solution: skip file")
					}
				}
				fmt.Println(c.Args())
				if c.Bool("testflag") {
					fmt.Println("testflag is active")
				}
				if edlFound == 0 {
					fmt.Println("no valid edl-files detected")
					return errors.New("Before Action end error")
				}
				return nil
			},
			After: func(*cli.Context) error {
				fmt.Println("Start after action")
				return nil
			},
			Action: func(c *cli.Context) error {
				fmt.Println("Start main action")
				//MAINACTION:
				/*
					обязательные входящие данные: EDL-file
					пример: cutter cut filename.edl


				*/
				return nil
			},
			OnUsageError: func(context *cli.Context, err error, isSubcommand bool) error {
				fmt.Println("Start on error action")
				return nil
			},
			Subcommands: []cli.Command{},
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:  "testflag",
					Usage: "If flag is active run sometests",
				},
			},
			SkipFlagParsing:        false,
			SkipArgReorder:         false,
			HideHelp:               false,
			Hidden:                 false,
			UseShortOptionHandling: false,
			HelpName:               "",
			CustomHelpTemplate:     "",
		},
	}
	args := os.Args
	if err := app.Run(args); err != nil {
		fmt.Println(err.Error())
	}
}

func cannotFindFile(err error) bool {
	if strings.Contains(err.Error(), "The system cannot find the file specified.") {
		return true
	}
	return false
}
