package main

import (
	"fmt"
	"time"

	"os"
	"strconv"
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
	logger = glog.New(logLocation, glog.LogLevelINFO)
	app := cli.NewApp()
	app.Version = "v 0.0.3"
	app.Name = "searcher"
	app.Usage = "Scans root directory and all subdirectories to create list of files that matches queary"
	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:  "vocal",
			Usage: "If flag is active searcher will print ALL log entries (level INFO is set by default)",
		},
		&cli.IntFlag{
			Name:  "delay",
			Usage: "If flag is active searcher will delay start for N seconds",
		},
	}

	app.Commands = []cli.Command{
		//////////////////////////////////////
		{
			Name:  "probe",
			Usage: "Searches all files in the root which associated with marker",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:  "check, c",
					Usage: "If flag is active run incheker on every found file individualy",
				},
				&cli.BoolFlag{
					Name:        "grab, g",
					Usage:       "If flag is active grabber will try to download all new files",
					Required:    false,
					Hidden:      false,
					Destination: new(bool),
				},
				&cli.IntFlag{
					Name:        "repeat, r",
					Usage:       "repeat action every N seconds",
					EnvVar:      "",
					FilePath:    "",
					Required:    false,
					Hidden:      false,
					Value:       0,
					Destination: new(int),
				},
			},
			Action: func(c *cli.Context) error {
				if c.GlobalInt("delay") > 0 {
					for i := 0; i <= c.GlobalInt("delay"); i++ {
						fmt.Print("Searcher will start in ", stamp.Seconds(int64(c.GlobalInt("delay")-i)), "                       \r")
						time.Sleep(time.Second)
					}
					fmt.Print("\n")
					fcli.RunConsole("dirmaker", "daily")
				}
				restart := true
				for restart {
					restart = false

					if c.GlobalBool("vocal") {
						logger.ShoutWhen(glog.LogLevelALL)
					}
					takeFile, err := scanner.Scan(root, marker)
					if err != nil {
						fmt.Println(err)
						logger.ERROR(err.Error())
						return err
					}
					fileList := scanner.ListReady(takeFile)
					for _, fl := range fileList {
						logger.TRACE("detected " + fl)
					}
					if c.Bool("check") {
						fcli.RunConsole("inchecker", fileList...)
					}
					logger.INFO(strconv.Itoa(len(fileList)-len(takeFile)) + " new files found")

					if c.Bool("grab") {
						prog := "grabber"
						args := []string{}
						if c.Bool("vocal") {
							args = append(args, "--vocal")
						}
						args = append(args, "takenew")
						fcli.RunConsole(prog, args...)
					}
					if c.Int("repeat") > 0 {
						restart = true
						for i := 0; i < c.Int("repeat"); i++ {
							fmt.Print("Probe in ", stamp.Seconds(int64(c.Int("repeat")-i)), "                 \r")
							time.Sleep(time.Second)
						}
						//fmt.Print("\n")
					}
				}
				//fmt.Print("Flag is |", app.Flags[0].String())
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

/*

search -new
search -all

search -take

search -today
search -thisweek
search -lastweek
search -repeat=60 -incheck -grab -until:202127020900


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
