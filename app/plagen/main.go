package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/Galdoba/ffstuff/app/plagen/cmd"
	"github.com/Galdoba/ffstuff/app/plagen/config"
	"github.com/urfave/cli/v2"
)

const (
	CONFIG = "cfg"
)

func main() {
	app := cli.NewApp()

	app.Version = "v 0.0.1"
	app.Usage = "Генерит плашки для использования в связке с agelogo"
	app.Flags = []cli.Flag{}

	//ДО НАЧАЛА ДЕЙСТВИЯ
	app.Before = func(c *cli.Context) error {
		cfg := &config.Config{}
		err := fmt.Errorf("config not loaded")
		cfg, err = config.Load(app.Name)
		if err != nil {
			fmt.Println("config err:", err.Error())

			switch {
			case strings.Contains(err.Error(), "The system cannot find "), strings.Contains(err.Error(), "no such file or directory"):
				cfg, err = config.NewConfig(c.App.Name)
				if err = cfg.Save(); err != nil {
					return fmt.Errorf("can't setup config: %v", err.Error())
				}
				fmt.Printf("default config created at %v: restart %v\n", cfg.Location, cfg.AppName)
				os.Exit(0)
			default:
				panic("wut?")
				return nil
			}
		}
		fmt.Println("Before ENDED")
		return nil
	}
	app.Commands = []*cli.Command{
		// cmd.Standard(),
		cmd.Custom(),
	}

	//ПО ОКОНЧАНИЮ ДЕЙСТВИЯ
	app.After = func(c *cli.Context) error {
		fmt.Println("After ENDED")
		return nil
	}
	args := os.Args
	if err := app.Run(args); err != nil {
		errOut := fmt.Sprintf("%v error: %v", app.Name, err.Error())
		println(errOut)
		os.Exit(1)
	}

}
