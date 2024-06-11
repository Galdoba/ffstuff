package main

import (
	"fmt"
	"os"

	"github.com/Galdoba/devtools/printer"
	"github.com/urfave/cli/v2"
)

/*
run


*/

const (
	programName = "bashgen"
)

// var cfg config.Config
var pm printer.Printer

func init() {

	// err := fmt.Errorf("config not loaded")
	// cfg, err = config.Load()
	// if err != nil {
	// 	cfg = config.New()
	// 	cfg.SetDefault()
	// 	if err := cfg.Save(); err != nil {
	// 		fmt.Printf("initialisation failed: %v", err.Error())
	// 		os.Exit(1)

	// 	}
	// 	fmt.Printf("config file generated at %v \n", cfg.Path())
	// 	fmt.Println("restart application")
	// 	os.Exit(0)

	// }

}

/*
функции:
run - работаем в режиме демона, если можно создать амедиевский скрипт - делаем его
config - выводим информацию о конфиге
*/

func main() {

	app := cli.NewApp()
	app.Version = "v 0.0.1"
	app.Name = programName
	app.Usage = "generate bash script for job"
	app.Flags = []cli.Flag{}

	app.Before = func(c *cli.Context) error {
		return nil
	}
	app.Commands = []*cli.Command{
		// cmd.Config(cfg),
		// cmd.Serial(cfg),
	}
	// app.DefaultCommand = "run"

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
