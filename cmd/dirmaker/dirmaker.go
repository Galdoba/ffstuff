package main

import (
	"fmt"
	"os"

	"github.com/Galdoba/ffstuff/constant"
	"github.com/Galdoba/ffstuff/fldr"
	"github.com/Galdoba/ffstuff/pkg/config"
	"github.com/Galdoba/utils"
	"github.com/urfave/cli"
)

//$ app [global options] command [command options] arguments
//dirmaker new [directory]
//dirmaker daily

var configMap map[string]string

func init() {
	// file, err := os.Stat(config.StandardPath())
	// if err == nil {
	// 	return
	// }
	// config.Construct()
	// config.SetField("INROOT", config.FieldUndefined)
	// config.SetField("MUXROOT", config.FieldUndefined)
	// config.SetField("OUTROOT", config.FieldUndefined)
	// fmt.Println("Please set root folders in", file.Name())
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
					configMap[constant.InPath] + "IN_" + utils.DateStamp() + "\\",
					configMap[constant.InPath] + "IN_" + utils.DateStamp() + "\\proxy\\",
					configMap[constant.MuxPath] + "MUX_" + utils.DateStamp() + "\\",
					configMap[constant.OutPath] + "OUT_" + utils.DateStamp() + "\\",
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
