package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/Galdoba/ffstuff/constant"
	"github.com/Galdoba/ffstuff/pkg/config"
	"github.com/Galdoba/ffstuff/pkg/glog"
	"github.com/Galdoba/ffstuff/pkg/scanner"
	"github.com/Galdoba/utils"
	"github.com/urfave/cli"
)

var configMap map[string]string
var logger glog.Logger

func init() {
	//err := errors.New("Initial obstract error")
	conf, err := config.ReadProgramConfig("ffstuff")
	if err != nil {
		fmt.Println(err)
	}
	configMap = conf.Field
	if err != nil {
		switch err.Error() {
		default:
			fmt.Print("Unknown error: ", err.Error())
			os.Exit(42)
		case "Config file not found":
			fmt.Print("Expecting config file in:\n", conf.Path)
			os.Exit(1)
		}
	}
}

func main() {
	muxRoot := configMap[constant.MuxPath] + "\\"
	muxFolder := muxRoot + "MUX_" + utils.DateStamp() + "\\"
	logPath := muxFolder + "logfile.txt"
	logger = glog.New(logPath, glog.LogLevelINFO)
	app := cli.NewApp()
	app.Version = "v 0.0.1"
	app.Name = "grabber"
	app.Usage = "dowloads files and sort it to working directories"
	app.Flags = []cli.Flag{}
	app.Commands = []cli.Command{
		//////////////////////////////////////
		{
			Name:  "today",
			Usage: "------",
			Action: func(c *cli.Context) error {
				fmt.Println(muxFolder)
				fileList, err := scanner.Scan(muxFolder, ".ac3")
				if err != nil {
					fmt.Println("ERROR:", err)
					return err
				}
				for _, v := range fileList {
					fmt.Println(v)
				}
				for _, audioFile := range fileList {
					if !strings.Contains(audioFile, "-ebur128") {
						continue
					}
					if strings.Contains(audioFile, "51-ebur128.ac3") {
						base := strings.TrimSuffix(audioFile, "51-ebur128.ac3")
						er := os.Rename(audioFile, base+"51.ac3")
						if er != nil {
							logger.ERROR("rename failed: " + er.Error())
						}
						logger.TRACE("renamed: " + audioFile + " -->" + muxFolder + base + "51.ac3")
						continue
					}
					if strings.Contains(audioFile, "20-ebur128.ac3") {
						base := strings.TrimSuffix(audioFile, "20-ebur128.ac3")
						er := os.Rename(audioFile, base+"20.ac3")
						if er != nil {
							logger.ERROR("rename failed: " + er.Error())
						}
						logger.TRACE("renamed: " + audioFile + " -->" + muxFolder + base + "20.ac3")
						continue
					}
					if strings.Contains(audioFile, "51-ebur128-stereo.ac3") {
						base := strings.TrimSuffix(audioFile, "51-ebur128-stereo.ac3")
						er := os.Rename(audioFile, base+"51.ac3")
						if er != nil {
							logger.ERROR("rename failed: " + er.Error())
						}
						logger.TRACE("renamed: " + audioFile + " -->" + muxFolder + base + "20.ac3")
						continue
					}
					if strings.Contains(audioFile, "20-ebur128-stereo.ac3") {
						base := strings.TrimSuffix(audioFile, "20-ebur128-stereo.ac3")
						er := os.Rename(audioFile, base+"51.ac3")
						if er != nil {
							logger.ERROR("rename failed: " + er.Error())
						}
						logger.TRACE("renamed: " + audioFile + " -->" + muxFolder + base + "20.ac3")
						continue
					}
					logger.ERROR("rename failed: " + audioFile)
				}
				return nil
			},
		},
		////////////////////////////////////
		// {
		// 	Name:  "todays",
		// 	Usage: "Create one or more new directories",
		// 	Action: func(c *cli.Context) error {
		// 		paths := c.Args().Tail() //	path := c.String("path") //*cli.Context.String(key) - вызывает флаг с именем key и возвращает значение Value
		// 		for _, path := range paths {
		// 			dir := fldr.New("",
		// 				fldr.Set(fldr.AddressFormula, path),
		// 			)
		// 			dir.Make()
		// 		}
		// 		return nil
		// 	},
		// },

	}
	args := os.Args
	if len(args) < 2 {
		args = append(args, "help") //Принудительно зовем помощь если нет других аргументов
	}
	fmt.Println("RUN")
	if err := app.Run(args); err != nil {

		fmt.Println(err.Error())
	}
	fmt.Println("END")
}
