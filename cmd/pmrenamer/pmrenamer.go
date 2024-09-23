package main

// import (
// 	"fmt"
// 	"os"

// 	"github.com/Galdoba/ffstuff/constant"
// 	"github.com/Galdoba/ffstuff/pkg/config"
// 	"github.com/Galdoba/ffstuff/pkg/glog"
// 	"github.com/Galdoba/ffstuff/pkg/namedata"
// 	"github.com/Galdoba/ffstuff/pkg/scanner"
// 	"github.com/Galdoba/utils"
// 	"github.com/urfave/cli/v2"
// )

// var configMap map[string]string
// var logger glog.Logger

// func init() {
// 	//err := errors.New("Initial obstract error")
// 	conf, err := config.ReadProgramConfig("ffstuff")
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	configMap = conf.Field
// 	if err != nil {
// 		switch err.Error() {
// 		default:
// 			fmt.Print("Unknown error: ", err.Error())
// 			os.Exit(42)
// 		case "Config file not found":
// 			fmt.Print("Expecting config file in:\n", conf.Path)
// 			os.Exit(1)
// 		}
// 	}
// }

// func main() {
// 	muxRoot := configMap[constant.MuxPath] + "\\"
// 	inRoot := configMap[constant.InPath] + "\\"
// 	muxFolder := muxRoot + "MUX_" + utils.DateStamp() + "\\"
// 	inFolder := inRoot + "IN_" + utils.DateStamp() + "\\"

// 	logger := glog.New(glog.LogPathDEFAULT, glog.LogLevelINFO)
// 	app := cli.NewApp()
// 	app.Version = "v 0.0.2"
// 	app.Name = "grabber"
// 	app.Usage = "dowloads files and sort it to working directories"
// 	app.Flags = []cli.Flag{}
// 	app.Commands = []*cli.Command{
// 		//////////////////////////////////////
// 		{
// 			Name:  "today",
// 			Usage: "------",
// 			Action: func(c *cli.Context) error {
// 				fmt.Println(muxFolder)
// 				fileList, err := scanner.Scan(muxFolder, ".ac3")
// 				if err != nil {
// 					fmt.Println("ERROR:", err)
// 					return err
// 				}
// 				for _, oldName := range fileList {
// 					newName, err := namedata.TrimLoudnormPrefix(oldName)
// 					if err != nil {
// 						logger.ERROR(err.Error())
// 						continue
// 					}
// 					os.Rename(oldName, newName)
// 				}
// 				// for _, audioFile := range fileList {
// 				// 	if !strings.Contains(audioFile, "-ebur128") {
// 				// 		continue
// 				// 	}
// 				// 	if strings.Contains(audioFile, "51-ebur128.ac3") {
// 				// 		base := strings.TrimSuffix(audioFile, "51-ebur128.ac3")
// 				// 		er := os.Rename(audioFile, base+"51.ac3")
// 				// 		if er != nil {
// 				// 			logger.ERROR("rename failed: " + er.Error())
// 				// 		}
// 				// 		logger.TRACE("renamed: " + audioFile + " -->" + muxFolder + base + "51.ac3")
// 				// 		continue
// 				// 	}
// 				// 	if strings.Contains(audioFile, "20-ebur128.ac3") {
// 				// 		base := strings.TrimSuffix(audioFile, "20-ebur128.ac3")
// 				// 		er := os.Rename(audioFile, base+"20.ac3")
// 				// 		if er != nil {
// 				// 			logger.ERROR("rename failed: " + er.Error())
// 				// 		}
// 				// 		logger.TRACE("renamed: " + audioFile + " -->" + muxFolder + base + "20.ac3")
// 				// 		continue
// 				// 	}
// 				// 	if strings.Contains(audioFile, "51-ebur128-stereo.ac3") {
// 				// 		base := strings.TrimSuffix(audioFile, "51-ebur128-stereo.ac3")
// 				// 		er := os.Rename(audioFile, base+"51.ac3")
// 				// 		if er != nil {
// 				// 			logger.ERROR("rename failed: " + er.Error())
// 				// 		}
// 				// 		logger.TRACE("renamed: " + audioFile + " -->" + muxFolder + base + "20.ac3")
// 				// 		continue
// 				// 	}
// 				// 	if strings.Contains(audioFile, "20-ebur128-stereo.ac3") {
// 				// 		base := strings.TrimSuffix(audioFile, "20-ebur128-stereo.ac3")
// 				// 		er := os.Rename(audioFile, base+"51.ac3")
// 				// 		if er != nil {
// 				// 			logger.ERROR("rename failed: " + er.Error())
// 				// 		}
// 				// 		logger.TRACE("renamed: " + audioFile + " -->" + muxFolder + base + "20.ac3")
// 				// 		continue
// 				// 	}
// 				// 	logger.ERROR("rename failed: " + audioFile)
// 				// }
// 				return nil
// 			},
// 		},
// 		////////////////////////////////////
// 		{
// 			Name:        "preapare",
// 			Aliases:     []string{},
// 			Usage:       "------",
// 			UsageText:   "",
// 			Description: "",
// 			ArgsUsage:   "",
// 			Category:    "",
// 			Action: func(c *cli.Context) error {
// 				if inFolder == "000" {
// 					fmt.Println(inFolder)
// 				}
// 				file := c.String("input")
// 				nf := namedata.ParseName(file)
// 				newName, err := nf.ReconstructName()
// 				if err != nil {
// 					fmt.Println(err.Error())
// 					return fmt.Errorf("Cannot rename %v\nError: %v", file, err.Error())
// 				}
// 				logger.INFO(fmt.Sprintf("renaming: '%v' => '%v'", file, newName))
// 				renamingErr := os.Rename(file, newName)
// 				if renamingErr != nil {
// 					fmt.Println(renamingErr.Error())
// 					return renamingErr
// 				}
// 				return nil
// 			},
// 			Flags: []cli.Flag{
// 				&cli.StringFlag{
// 					Name:     "input",
// 					Usage:    "sets file to preapare name",
// 					Required: true,
// 				},
// 			},
// 		},
// 		//////////////////////////////
// 		// {
// 		// 	Name:  "todays",
// 		// 	Usage: "Create one or more new directories",
// 		// 	Action: func(c *cli.Context) error {
// 		// 		paths := c.Args().Tail() //	path := c.String("path") //*cli.Context.String(key) - вызывает флаг с именем key и возвращает значение Value
// 		// 		for _, path := range paths {
// 		// 			dir := fldr.New("",
// 		// 				fldr.Set(fldr.AddressFormula, path),
// 		// 			)
// 		// 			dir.Make()
// 		// 		}
// 		// 		return nil
// 		// 	},
// 		// },

// 	}
// 	args := os.Args
// 	if len(args) < 2 {
// 		args = append(args, "help") //Принудительно зовем помощь если нет других аргументов
// 	}
// 	fmt.Println("RUN")
// 	if err := app.Run(args); err != nil {
// 		logger.FATAL(err.Error())
// 		//fmt.Println(err.Error())
// 	}
// 	fmt.Println("END")
// }
