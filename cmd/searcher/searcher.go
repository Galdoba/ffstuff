package main

import (
	"fmt"
	"time"

	"os"
	"strings"

	"github.com/Galdoba/ffstuff/constant"
	"github.com/Galdoba/ffstuff/fldr"
	"github.com/Galdoba/utils"
	"github.com/urfave/cli"

	fcli "github.com/Galdoba/ffstuff/pkg/cli"
	"github.com/Galdoba/ffstuff/pkg/scanner"
	"github.com/Galdoba/ffstuff/pkg/stamp"

	"github.com/Galdoba/ffstuff/pkg/config"
	"github.com/Galdoba/ffstuff/pkg/glog"
)

var configMap map[string]string

var logger glog.Logger
var logLocation string

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
	root := configMap[constant.SearchRoot]
	marker := configMap[constant.SearchMarker]
	if configMap[constant.LogDirectory] == "default" {
		logLocation = fldr.MuxPath() + "logfile.txt"
	}
	//logger = glog.New(logLocation, glog.LogLevelINFO)
	logger = glog.New(glog.LogPathDEFAULT, glog.LogLevelINFO)
	app := cli.NewApp()
	app.Version = "v 0.0.4"
	app.Name = "searcher"
	app.Usage = "Scans root directory and all subdirectories to create list of files that matches queary"
	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:  "vocal",
			Usage: "If flag is active searcher will print ALL log entries (level INFO is set by default)",
		},
		&cli.BoolFlag{
			Name:  "clear, cl",
			Usage: "If flag is active searcher Clear Terminal before every Search",
		},
		&cli.StringFlag{
			Name:  "delay",
			Usage: "If flag is active searcher will delay start for N seconds",
		},
	}

	app.Commands = []cli.Command{
		//////////////////////////////////////
		// {
		// 	Name:  "probe",
		// 	Usage: "Searches all files in the root which associated with marker",
		// 	Flags: []cli.Flag{
		// 		&cli.BoolFlag{
		// 			Name:  "check, c",
		// 			Usage: "If flag is active run incheker on every found file individualy",
		// 		},
		// 		&cli.BoolFlag{
		// 			Name:     "grab, g",
		// 			Usage:    "If flag is active grabber will try to download all new files",
		// 			Required: false,
		// 			Hidden:   false,
		// 		},
		// 		&cli.StringFlag{
		// 			Name:     "repeat, r",
		// 			Usage:    "repeat action every N seconds",
		// 			EnvVar:   "",
		// 			FilePath: "",
		// 			Required: false,
		// 			Hidden:   false,
		// 			Value:    "",
		// 		},
		// 	},
		// 	Action: func(c *cli.Context) error {
		// 		if c.GlobalString("delay") != "" {
		// 			sec := utils.TimeStampToSeconds(c.GlobalString("delay"))
		// 			for i := 0; i <= sec; i++ {
		// 				fmt.Print("Searcher will start in ", stamp.Seconds(int64(sec-i)), "                       \r")
		// 				time.Sleep(time.Second)
		// 			}
		// 			fmt.Print("\n")
		// 			fcli.RunConsole("dirmaker", "daily")
		// 		}
		// 		restart := false
		// 		if c.String("repeat") != "" {
		// 			restart = true
		// 		}
		// 	maincycle:
		// 		for {
		// 			if c.GlobalBool("vocal") {
		// 				logger.ShoutWhen(glog.LogLevelALL) //вещаем в терминал все сообщения логгера
		// 			}
		// 			if c.GlobalBool("clear") {
		// 				utils.ClearScreen() //обновляем экран
		// 			}
		// 			takeFile, err := scanner.Scan(root, marker) //сканируем
		// 			if err != nil {
		// 				//fmt.Println(err)
		// 				logger.ERROR("scan failed: " + err.Error())
		// 				//return err
		// 			}
		// 			fileList := scanner.ListReady(takeFile)
		// 			if len(fileList) == 0 { //если найдено 0 новых файлов - то дальже либо ждем либо прекращаем работу
		// 				logger.INFO("No new files found")
		// 				switch restart {
		// 				case true:
		// 					repeatIfNeeded(c)
		// 					continue
		// 				case false:
		// 					break maincycle
		// 				}
		// 			}
		// 			for _, fl := range fileList {
		// 				logger.TRACE("detected " + fl)
		// 			}
		// 			if c.Bool("check") {
		// 				arguments := append([]string{"check"}, fileList...)
		// 				fcli.RunConsole("inchecker", arguments...) //проверяем инчекером
		// 			}
		// 			logger.INFO(strconv.Itoa(len(fileList)-len(takeFile)) + " new files found")
		// 			if c.Bool("grab") {
		// 				prog := "grabber"
		// 				args := []string{}
		// 				if c.Bool("vocal") {
		// 					args = append(args, "--vocal")
		// 				}
		// 				args = append(args, "takenew")
		// 				fcli.RunConsole(prog, args...) //хватем найденое
		// 			}
		// 			repeatIfNeeded(c)
		// 			// if c.String("repeat") != "" {
		// 			// 	break
		// 			// }
		// 		}
		// 		return nil
		// 	},
		// },
		//////////////////////////////////////
		{
			Name:  "list",
			Usage: "List all files in NAS",
			Flags: []cli.Flag{},
			Action: func(c *cli.Context) error {
				takeFile, err := scanner.Scan(root, "") //сканируем
				if err != nil {
					//fmt.Println(err)
					logger.ERROR("scan failed: " + err.Error())
					//return err
				}
				lenFiles := len(takeFile)
				for i, val := range takeFile {
					fmt.Println("Scan position:", i+1, "/", lenFiles)
					stat, _ := os.Stat(val)
					if stat.IsDir() {
						fmt.Println("Skip")
						continue
					}
					fcli.RunConsole("inchecker", val) //проверяем инчекером
				}
				return nil
			},
		},
		//////////////////////////////////////
		{
			Name:  "probe",
			Usage: "Searches all files in the root which associated with marker and downloading by order",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:  "check, c",
					Usage: "If flag is active run incheker on every found file individualy",
				},
				&cli.BoolFlag{
					Name:  "grab, g",
					Usage: "If flag is active grabber will try to download all new files",
				},
				&cli.StringFlag{
					Name:  "repeat, r",
					Usage: "repeat action every N seconds",
				},
			},
			Action: func(c *cli.Context) error {
				logger.INFO("Run Searcher")
				if c.GlobalString("delay") != "" {
					sec := utils.TimeStampToSeconds(c.GlobalString("delay"))
					for i := 0; i <= sec; i++ {
						fmt.Print("Searcher will start in ", stamp.Seconds(int64(sec-i)), "                       \r")
						time.Sleep(time.Second)
					}
					fmt.Print("\n")
					fcli.RunConsole("dirmaker", "daily", "-clean")
				}
				restart := true
				if c.String("repeat") != "" {
					restart = true
				}
				for restart {
					if c.GlobalBool("vocal") {
						logger.ShoutWhen(glog.LogLevelALL) //вещаем в терминал все сообщения логгера
					}
					if c.GlobalBool("clear") {
						utils.ClearScreen() //обновляем экран
					}
					takeFile, err := scanner.Scan(root, marker) //сканируем
					if err != nil {
						logger.ERROR("scan failed: " + err.Error())
					}
					if len(takeFile) == 0 { //если ничего не найдено смотрим надо ли ждать
						err = repeatIfNeeded(c)
						switch err { //                    если ничего не найдено смотрим надо ли ждать
						default: //                        и ждем поле чего начинаем поиск с начала
							logger.TRACE(err.Error())
							fmt.Printf("%v: New files not found\n", time.Now().Format("2006-01-02 15:04:05.000"))
							continue
						case nil: //                       или выходим
							logger.TRACE("End Program")
							fmt.Print("Searcher: End Program")
							restart = false
							os.Exit(0)
						}
					}
					takeFile = scanner.SortPriority(takeFile) //сортируем список readyFile
					if c.Bool("check") {
						fileList := scanner.ListReady(takeFile)
						arguments := append([]string{"check"}, fileList...)
						fcli.RunConsole("inchecker", arguments...) //проверяем инчекером
					}
					if c.Bool("grab") {
						fileReady := takeFile[0] // берем первый
						prog := "grabber"
						args := []string{}
						if c.Bool("vocal") {
							args = append(args, "--vocal")
						}
						args = append(args, "takeready")
						args = append(args, "-connectedwith")
						args = append(args, fileReady)
						fcli.RunConsole(prog, args...) //хватем найденое
						//fmt.Println("RUN:", prog, args)
						restart = true
					}
				}
				repeatIfNeeded(c)
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

func repeatIfNeeded(c *cli.Context) error {
	if c.String("repeat") != "" {
		wait := utils.TimeStampToSeconds(c.String("repeat"))
		for i := 0; i < wait; i++ {
			fmt.Print("Next scan in ", stamp.Seconds(int64(wait-i)), "                 \r")
			time.Sleep(time.Second)
		}
		return fmt.Errorf("Waited %v\r", stamp.Seconds(int64(wait)))
	}
	return nil
}

/*

скачать video1, audio1, audio2
если ошибок == 0 {
	loudnorm audio1
	ren audio1-ebur128.ac3 audio1.ac3

}



*/

func sortResults(list []string) []string {
	sorted := []string{}
	for _, val := range list {
		if strings.Contains(val, ".srt") {
			sorted = append(sorted, val)
		}
	}
	for _, val := range list {
		if strings.Contains(val, ".ready") {
			sorted = append(sorted, val)
		}
	}
	for _, val := range list {
		if strings.Contains(val, "_Proxy_") {
			sorted = append(sorted, val)
		}
	}
	for _, val := range list {
		if strings.Contains(val, ".m4a") {
			sorted = append(sorted, val)
		}
	}
	for _, val := range list {
		sorted = utils.AppendUniqueStr(sorted, val)
	}
	return sorted
}

/*

:process
set f_fullname=%~1
set f_path=%~p1
set f_name=%~n1
set f_ext=%~x1


rem ffmpeg -i %f_name%.ac3 -map 0:0 -acodec ac3 -ab 640k %f_name%.ac3
rem ffmpeg -i eng_%f_name%.ac3 -map 0:0 -acodec ac3 -ab 640k eng_%f_name%.ac3

rem mkvmerge -o "\\192.168.32.3\root\#KIRILL\%f_name%_ar6e2.mkv" -d 0 --language 0:rus --default-track 0:1 -A ="%f_name%.mp4" -a 0 --language 0:rus --default-track 0:1 ="%f_name%.ac3" -a 0 --language 0:eng --default-track 0:0 ="eng_%f_name%.ac3"
rem -s 0 --language 0:rus --default-track 0:0 ="%f_name%.srt"

rem ffmpeg -i "%f_name%.mp4" -i "%f_name%.ac3" -i "eng_%f_name%.ac3" -i "%f_name%.srt" -map 0:v -map 1:a -map 2:a -map 3:s -codec copy -codec:s mov_text -metadata:s:a:0 language=rus -metadata:s:a:1 language=eng -metadata:s:s:0 language=rus "1232_%f_name%_ar6e2_sr.mp4"

ffmpeg ^
-i "%f_name%.mp4" ^
-i "%f_name%_rus20.ac3" ^
-i "%f_name%_eng51.ac3" ^
-i "%f_name%.srt" ^
-codec copy -codec:s mov_text ^
    -map 0:v ^
    -map 1:a -metadata:s:a:0 language=rus ^
    -map 2:a -metadata:s:a:1 language=eng ^
    -map 3:s -metadata:s:s:0 language=rus ^
"\\192.168.32.3\ROOT\#PETR\toCheck\%f_name%_ar2e6.mp4"



exit /b 0

rem ffmpeg ^
rem -i "%f_name%.mp4" ^
rem -i "%f_name%_rus20.ac3" ^
rem -i "%f_name%_eng51.ac3" ^
rem -i "%f_name%.srt" ^
rem -codec copy -codec:s mov_text ^
rem     -map 0:v ^
rem     -map 1:a -metadata:s:a:0 language=rus ^
rem     -map 2:a -metadata:s:a:1 language=eng ^
rem     -map 3:s -metadata:s:s:0 language=rus ^
rem "\\192.168.32.3\ROOT\#PETR\toCheck\%f_name%_ar2e6.mp4"

*/
