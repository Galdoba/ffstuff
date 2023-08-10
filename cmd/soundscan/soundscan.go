package main

import (
	"fmt"
	"strconv"

	"os"

	"github.com/urfave/cli/v2"

	"github.com/Galdoba/ffstuff/pkg/silence"
)

var configMap map[string]string

//var logger glog.Logger
var logLocation string

// func init() {
// 	conf, err := config.ReadProgramConfig("ffstuff")
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	configMap = conf.Field
// 	if err != nil {
// 		switch err.Error() {
// 		case "Config file not found":
// 			fmt.Print("Expecting config file in:\n", conf.Path)
// 			os.Exit(1)
// 		}
// 	}
// }

func main() {
	//root := configMap[constant.SearchRoot]
	//marker := configMap[constant.SearchMarker]
	// if configMap[constant.LogDirectory] == "default" {
	// 	logLocation = fldr.MuxPath() + "logfile.txt"
	// }
	//logger = glog.New(logLocation, glog.LogLevelINFO)
	//logger = glog.New(glog.LogPathDEFAULT, glog.LogLevelINFO)
	app := cli.NewApp()
	app.Version = "v 0.0.1"
	app.Name = "soundscan"
	app.Usage = "Scans audio stream for it's loudness data"
	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:  "vocal",
			Usage: "If flag is active soundcan will print data on terminal",
		},
	}
	app.Commands = []*cli.Command{
		//////////////////////////////////////
		{
			Name:  "listen",
			Usage: "Listens all files and checks silence in audio stream",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "loudnessborder, lb",
					Usage: "Sets loudness border (everything lower than ln treated as silence)",
					Value: "-72.0",
				},
				&cli.StringFlag{
					Name:  "duration, d",
					Usage: "Sets minimal duration for silence to record",
					Value: "2",
				},
			},
			Action: func(c *cli.Context) error {
				args := c.Args().Tail()
				for _, v := range args {
					lb, err := strconv.ParseFloat(c.String("loudnessborder"), 64)
					if lb > 0 {
						lb = lb * -1
					}
					if err != nil {
						fmt.Println(err)
						return err
					}
					d, err := strconv.ParseFloat(c.String("duration"), 64)
					if err != nil {
						fmt.Println(err)
						return err
					}
					si, err := silence.Detect(v, lb, d, c.Bool("vocal"))
					if err != nil {
						fmt.Println(err)
					} else {
						fmt.Println(si)
						fmt.Println(si.Timings())
					}
				}
				return nil
			},
		},
		//////////////////////////////////////
		// {
		// Name:  "list",
		// Usage: "List all files in NAS",
		// Flags: []cli.Flag{},
		// Action: func(c *cli.Context) error {
		// takeFile, err := scanner.Scan(root, "") //сканируем
		// if err != nil {
		//fmt.Println(err)
		// logger.ERROR("scan failed: " + err.Error())
		//return err
		// }
		// lenFiles := len(takeFile)
		// for i, val := range takeFile {
		// fmt.Println("Scan position:", i+1, "/", lenFiles)
		// stat, _ := os.Stat(val)
		// if stat.IsDir() {
		// fmt.Println("Skip")
		// continue
		// }
		// fcli.RunConsole("inchecker", val) //проверяем инчекером
		// }
		// return nil
		// },
		// },
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

скачать video1, audio1, audio2
если ошибок == 0 {
	loudnorm audio1
	ren audio1-ebur128.ac3 audio1.ac3

}



*/

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
