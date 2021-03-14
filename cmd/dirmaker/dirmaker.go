package main

import (
	"fmt"
	"os"

	"github.com/Galdoba/ffstuff/fldr"
	"github.com/Galdoba/ffstuff/pkg/config"
	"github.com/Galdoba/utils"
	"github.com/urfave/cli"
)

//$ app [global options] command [command options] arguments
//dirmaker new [directory]
//dirmaker daily

func init() {
	file, err := os.Stat(config.StandardPath())
	if err == nil {
		return
	}
	config.Construct()
	config.SetField("INROOT", config.FieldUndefined)
	config.SetField("MUXROOT", config.FieldUndefined)
	config.SetField("OUTROOT", config.FieldUndefined)
	fmt.Println("Please set root folders in", file.Name())
}

func main() {
	if err := config.Verify(); err != nil {
		fmt.Println(err.Error())
		return
	}
	conf, _ := config.Read()
	app := cli.NewApp()
	app.Version = "v 0.0.1"
	app.Name = "dirmaker"
	app.Usage = "checks, creates and (TODO: Deletes) directories"
	app.Commands = []*cli.Command{
		//////////////////////////////////////
		{
			Name:  "new",
			Usage: "Create one or more new directories",
			Action: func(c *cli.Context) error {
				paths := c.Args().Slice() //	path := c.String("path") //*cli.Context.String(key) - вызывает флаг с именем key и возвращает значение Value
				for _, path := range paths {
					dir := fldr.New("",
						fldr.Set(fldr.AddressFormula, path),
					)
					dir.Make()
				}
				return nil
			},
		},
		//////////////////////////////////////
		{
			Name:  "daily",
			Usage: "Create today's work directories",
			Action: func(c *cli.Context) error {
				paths := []string{
					conf["INROOT"] + "\\IN_" + utils.DateStamp(),
					conf["INROOT"] + "\\IN_" + utils.DateStamp() + "\\proxy",
					conf["MUXROOT"] + "\\MUX_" + utils.DateStamp(),
					conf["OUTROOT"] + "\\OUT_" + utils.DateStamp(),
				}
				for _, path := range paths {
					dir := fldr.New("",
						fldr.Set(fldr.AddressFormula, path),
					)
					dir.Make()
				}
				return nil
			},
		},
		//////////////////////////////////////
	}
	args := os.Args
	if len(args) < 2 {
		args = append(args, "help") //Принудительно зовем помощь если нет других аргументов
	}
	if err := app.Run(args); err != nil {
		fmt.Println(err.Error())
	}
}
