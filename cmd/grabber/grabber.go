package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Galdoba/ffstuff/constant"
	"github.com/Galdoba/ffstuff/fldr"
	"github.com/Galdoba/ffstuff/pkg/config"
	"github.com/Galdoba/ffstuff/pkg/glog"
	"github.com/Galdoba/ffstuff/pkg/grabber"
	"github.com/Galdoba/ffstuff/pkg/scanner"
	"github.com/Galdoba/utils"
	"github.com/urfave/cli"
)

/*
TZ:
>> grab only [path]				-- забрать только указанные пути
>> grab filename.ready			-- забрать все связанное с ready файлом
>> grab help (-h)				-- вывести на экран помогалку 							--help
>> grab new (-n)				-- забрать все новое (предварительное сканирование)		--new
>> grab -v						-- забрать только если одобряет инчекер					--valid
>> grab -p						-- забрать только звук и прокси							--proxy
>> grab -fc						-- забрать только если одобряет fflite @check0			--fflitecheck0

пред проверки:
-папка куда копировать
-отсуствие файла с таким же именем и размером
-наличие свободного места для копии

пост проверки:
-копия равна по имени и размеру с источником

*/

var configMap map[string]string
var logger glog.Logger

func init() {
	err := errors.New("Initial obstract error")
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
	searchRoot := configMap[constant.SearchRoot]
	searchMarker := configMap[constant.SearchMarker]
	//dest := configMap[constant.InPath] + "IN_" + utils.DateStamp() + "\\"
	logPath := configMap[constant.MuxPath] + "MUX_" + utils.DateStamp() + "\\logfile.txt"
	logger = glog.New(logPath, glog.LogLevelINFO)

	app := cli.NewApp()
	app.Version = "v 0.0.2"
	app.Name = "grabber"
	app.Usage = "dowloads files and sort it to working directories"
	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:  "vocal",
			Usage: "If flag is active grabber set logLevel to TRACE (level INFO is set by default)",
		},
	}

	app.Commands = []cli.Command{
		//////////////////////////////////////
		{
			Name:  "takeonly",
			Usage: "Download only those files, that was received as arguments",
			Action: func(c *cli.Context) error {
				//paths := c.Args().Slice() //	path := c.String("path") //*cli.Context.String(key) - вызывает флаг с именем key и возвращает значение Value
				paths := c.Args().Tail()
				dest := configMap[constant.InPath] + "IN_" + utils.DateStamp() + "\\"
				for _, path := range paths {
					fmt.Println("GRABBER DOWNLOADING FILE:", path)
					if strings.Contains(path, "_Proxy_") {
						dest = dest + "proxy\\"
					}
					if strings.Contains(path, ".srt") {
						dest = fldr.MuxPath()
					}
					err := grabber.CopyFile(path, dest)
					if err != nil {
						fmt.Println(err.Error())
					}
				}
				return nil
			},
		},
		////////////////////////////////////
		{
			Name:  "takenew",
			Usage: "Call Scanner to get list of new and ready files",
			Action: func(c *cli.Context) error {

				if c.Bool("vocal") {
					logger.ShoutWhen(glog.LogLevelALL)
				}
				takeFile, err := scanner.Scan(searchRoot, searchMarker)
				if err != nil {
					fmt.Println(err)
					return err
				}
				fileList := scanner.ListReady(takeFile)
				for _, path := range fileList {
					dest := configMap[constant.InPath] + "IN_" + utils.DateStamp() + "\\"
					if strings.Contains(path, "_Proxy_") {
						dest = dest + "proxy\\"
					}
					grabber.CopyFile(path, dest)
					logger.TRACE("downloaded from:" + path)
				}
				logger.INFO(strconv.Itoa(len(fileList)) + " files downloaded")
				//paths := c.Args().Slice() //	path := c.String("path") //*cli.Context.String(key) - вызывает флаг с именем key и возвращает значение Value
				// for _, path := range paths {
				// 	fmt.Println("GRABBER DOWNLOADING FILE:", path)
				// 	err := grabber.CopyFile(path, dest)
				// 	fmt.Println(err)
				// 	if err != nil {
				// 		fmt.Println(err.Error())
				// 	}
				// }
				return nil
			},
		},
	}
	args := os.Args
	if len(args) < 2 {
		args = append(args, "help") //Принудительно зовем помощь если нет других аргументов
	}
	if err := app.Run(args); err != nil {
		fmt.Println(err.Error())
	}
}
