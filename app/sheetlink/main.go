package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Galdoba/ffstuff/app/sheetlink/dataconnection"
	"github.com/Galdoba/ffstuff/pkg/spreadsheet"
	"github.com/urfave/cli/v2"
)

const (
	programName = "sheetlink"
	/*
		gmscan

	*/
)

func main() {
	app := cli.NewApp()
	app.Version = "v 0.0.1"
	app.Name = programName
	app.Usage = "информация о файле из рабочей таблицы"
	app.Flags = []cli.Flag{}
	//ДО НАЧАЛА ДЕЙСТВИЯ
	app.Before = func(c *cli.Context) error {

		return nil
	}
	app.Commands = []*cli.Command{
		//ДЕЙСТВИЕ 1
		{
			Name:      "draw",
			Usage:     "draw [infoflag]... [argument]...",
			ArgsUsage: "",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:    "tablename",
					Aliases: []string{"tn"},
				},
				//ФЛАГИ ДЕЙСТВИЯ
			},
			Action: func(c *cli.Context) error {
				args := c.Args().Slice()
				if fail := argumentsConnectedCheck(args); fail != nil {
					return fail
				}
				names := fileNamesOnly(args)
				fmt.Println(names)
				sheet, err := spreadsheet.New()
				if err != nil {
					return err
				}

				dc := dataconnection.New(names...)
				if err := dc.Link(sheet.Data()); err != nil {
					return fmt.Errorf("link failed: %v", err.Error())
				}
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

const (
	Tag_TRL   = "--TRL--"
	Tag_FILM  = "--FLM--"
	Tag_SER   = "--SER--"
	Tag_FILM2 = "--FILM--"
)

func typeTags() []string {
	return []string{
		Tag_TRL,
		Tag_FILM,
		Tag_FILM2,
		Tag_SER,
	}
}

func argumentsConnectedCheck(args []string) error {
	for _, arg := range args {
		haveTypeTag := false
		for _, typeTag := range typeTags() {
			if strings.Contains(arg, typeTag) {
				haveTypeTag = true
			}
		}
		if !haveTypeTag {
			return fmt.Errorf("argument unconnected: %v", arg)
		}
	}
	return nil
}

func fileNamesOnly(fullNames []string) []string {
	short := []string{}
	for _, path := range fullNames {
		short = append(short, filepath.Base(path))
	}
	return short
}
