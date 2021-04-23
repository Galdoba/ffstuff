package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/Galdoba/devtools/cli/user"
	"github.com/Galdoba/ffstuff/constant"
	"github.com/Galdoba/ffstuff/pkg/config"
	"github.com/Galdoba/ffstuff/pkg/stamp"
	"github.com/Galdoba/utils"
	"github.com/urfave/cli"
)

const (
	FlagFilter = "filter"
	FlagDate   = "date"
)

var configMap map[string]string

func init() {
	//err := errors.New("Initial obstract error")
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
	//dateStamp := stamp.Date()
	logPath := configMap[constant.MuxPath] + "MUX_" + stamp.Date() + "\\logfile.txt"
	app := cli.NewApp()
	app.Version = "v 0.0.0"
	app.Name = "logviewer"
	app.Usage = "Controls ffstuff logs"
	app.Flags = []cli.Flag{}

	app.Commands = []cli.Command{

		//////////////////////////////////////
		{
			Name:        "show",
			Usage:       "Show todays logs",
			UsageText:   "TODO Usage",
			Description: "TODO Description",
			Action: func(c *cli.Context) error {
				if c.String(FlagDate) != "" {
					fmt.Println("Set date")
					panic(9)
				}
				logfile := logPath // + dateStamp + "\\logfile.txt"
				entries := utils.LinesFromTXT(logfile)
				filters := strings.Split(c.String(FlagFilter), ",")
				matches := []string{}
			lineProc:
				for _, line := range entries {
					lowLine := strings.ToLower(line)
					for _, f := range filters {
						lowFilter := strings.ToLower(f)
						if !strings.Contains(lowLine, lowFilter) {
							continue lineProc
						}
					}
					matches = append(matches, line)
				}
				switch len(matches) {
				default:
					for _, ln := range matches {
						fmt.Println(ln)
					}
				case 0:
					fmt.Println("No valid entries found in " + logfile)
				}
				fmt.Println("Press 'ENTER' to finish programm")
				user.InputStr()
				return nil
			},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "date, d",
					Usage: "Set day stamp for log that needs to be shown (if absent takes todays logs)",
					Value: "",
				},
				&cli.StringFlag{
					Name:  "filter, f",
					Usage: "Set filters separated by comma ','",
					Value: "",
				},
			},
		},
		////////////////////////////////////

	}
	args := os.Args
	if len(args) < 2 {
		args = append(args, "help") //Принудительно зовем помощь если нет других аргументов
	}
	if err := app.Run(args); err != nil {
		fmt.Println(err.Error())
	}
}
