package main

import (
	"fmt"
	"os"

	"github.com/Galdoba/ffstuff/constant"
	"github.com/Galdoba/ffstuff/fldr"
	"github.com/Galdoba/ffstuff/pkg/config"
	"github.com/Galdoba/ffstuff/pkg/glog"
	"github.com/Galdoba/ffstuff/pkg/muxer"
	"github.com/Galdoba/utils"
	"github.com/urfave/cli"
)

var configMap map[string]string

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
	logger := glog.New(glog.LogPathDEFAULT, glog.LogLevelINFO)
	app := cli.NewApp()
	app.Version = "v 0.0.1"
	app.Name = "muxer2"
	app.Usage = "Muxes media files using 'muxlist.txt' as a directions"
	app.Commands = []cli.Command{
		//////////////////////////////////////
		{
			Name:  "today",
			Usage: "Muxes files using current day muxer only",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:  "unsafe, us",
					Usage: "If flag is active muxer will try",
				},
			},
			Action: func(c *cli.Context) error {
				///
				tasks, err := muxer.MuxList()
				if err != nil {
					logger.ERROR(err.Error())
					if !c.Bool("unsafe") {
						fmt.Printf("end program")
						os.Exit(2)
					}
				}
				for i, task := range tasks {
					fmt.Print("Task ", i+1, "/", len(tasks), ":\n")
					files, muxTask, err := muxer.ChooseMuxer(task)
					if err != nil {
						logger.ERROR(err.Error())
						fmt.Println(err)
						continue
					}
					err = muxer.Run(muxTask, files)
					if err != nil {
						logger.ERROR(err.Error())
						fmt.Println(err)
					}
					logger.TRACE("Task Complete: " + task)
				}
				logger.INFO("Muxig Complete")
				///
				return nil
			},
		},
		//////////////////////////////////////
		{
			Name:  "daily",
			Usage: "Create today's work directories and daily files",
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
				ensureFileExistiense(configMap[constant.MuxPath] + "MUX_" + utils.DateStamp() + "\\" + "muxlist.txt")

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

func ensureFileExistiense(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0666)
	}
	_, err = os.Stat(path)
	return err
}
