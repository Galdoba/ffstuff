package main

import (
	"fmt"
	"os"
	"os/user"
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
var username string

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
	currentUser, userErr := user.Current()
	if err != nil {
		fmt.Printf("Initialisation failed: %v", userErr.Error())
	}
	username = currentUser.Name
}

func main() {
	searchRoot := configMap[constant.SearchRoot]
	searchMarker := configMap[constant.SearchMarker]
	//dest := configMap[constant.InPath] + "IN_" + utils.DateStamp() + "\\"
	//logPath := configMap[constant.MuxPath] + "MUX_" + utils.DateStamp() + "\\logfile.txt"
	//logger = glog.New(logPath, glog.LogLevelINFO)
	logger = glog.New(glog.LogPathDEFAULT, glog.LogLevelINFO)
	destination := configMap[constant.InPath] + "IN_" + utils.DateStamp() + "\\"
	app := cli.NewApp()
	app.Version = "v 0.0.3"
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
				dest := destination
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
				invalidErrors := 0
				if c.GlobalBool("vocal") {
					logger.ShoutWhen(glog.LogLevelALL)
				}
				takeFile, err := scanner.Scan(searchRoot, searchMarker)
				if err != nil {
					logger.ERROR(err.Error())
					return err
				}
				fileList := scanner.ListReady(takeFile)
				logger.TRACE(strconv.Itoa(len(fileList)) + " files detected")
				for _, path := range fileList {
					dest := destination
					if strings.Contains(path, "_Proxy_") {
						dest = dest + "proxy\\"
					}
					//grabber.CopyFile(path, dest, c.GlobalBool("vocal"))
					logger.TRACE("Start downloading:")
					if err := grabber.Download(logger, path, dest); err != nil {
						switch err.Error() {
						default:
							invalidErrors++
						case "valid copy exists":
						}
					}

				}
				if invalidErrors == 0 {
					for _, val := range takeFile {
						if !strings.Contains(val, ".ready") {
							continue
						}
						body := strings.TrimSuffix(val, ".ready")
						logger.TRACE("rename: " + val + " >> " + body + "." + username)
						if err := os.Rename(val, body+"."+username); err != nil {
							logger.ERROR(err.Error())
						}
					}
				}
				logger.INFO(strconv.Itoa(len(fileList)) + " files downloaded")
				return nil
			},
		},
		////////////////////////////////////
		{
			Name:        "takeready",
			ShortName:   "",
			Aliases:     []string{},
			Usage:       "Call Scanner to get list of new and ready files",
			UsageText:   "TODO:Usage",
			Description: "TODO:Descr",
			ArgsUsage:   "TODO:ArgsUsage",
			Category:    "TODO:Category",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "connectedwith",
					Usage:    "Setups exact .ready file to grab associated files with",
					Required: true,
				},
			},
			Action: func(c *cli.Context) error {
				readyfile := c.String("connectedwith")
				//	fmt.Println(readyfile)
				if err := checkMediaFile(readyfile, ".ready"); err != nil {
					logger.ERROR(err.Error())
					return err
				}
				allFiles := scanner.ListReady([]string{readyfile})
				allFiles = ensureValidOrder(allFiles)
				if err := downloadAssociatedWith(logger, allFiles, destination); err != nil {
					logger.ERROR(err.Error())
				}
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

func checkMediaFile(path string, keys ...string) error {
	f, err := os.Stat(path)
	if err != nil {
		return err
	}
	for _, key := range keys {
		switch key {
		default:
		case ".ready", ".mp4", ".m4a":
			if !strings.Contains(f.Name(), key) {
				return fmt.Errorf("%v is not a '%v' file", path, key)
			}
		}
	}
	return nil
}

func downloadAssociatedWith(l glog.Logger, paths []string, destination string) error {
	marker := ""
	for _, path := range paths {
		if strings.Contains(path, ".ready") {
			marker = strings.TrimSuffix(path, ".ready") + "." + username
			err := os.Rename(path, marker)
			//fmt.Println("Rename", path)
			if err != nil {
				return err
			}
			continue
		}
		if err := grabber.Download(logger, path, destination); err != nil {
			return err
		}
	}

	return nil
}

func ensureValidOrder(sl []string) []string {
	valid := []string{}
	for _, val := range sl {
		if strings.Contains(val, ".ready") {
			valid = append(valid, val)
		}
	}
	for _, val := range sl {
		valid = utils.AppendUniqueStr(valid, val)
	}
	return valid
}

/*

1 2 3 4 5 6
6 5 4 3 2 1 = 7 * (6/2) = 21

1 2 3 4 5 6 7
7 6 5 4 3 2 1 = 28

X1 + Xn = 8
X2 + Xn-1 = 8 + y - y
X3 + Xn-2 = 8 + 2y - 2y
X4 + Xn-3 = 8 + 3y - 3y
28/(8/2) = 7





*/
