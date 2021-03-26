package main

import (
	"fmt"

	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Galdoba/devtools/cli/user"
	"github.com/Galdoba/ffstuff/constant"
	"github.com/Galdoba/ffstuff/fldr"
	"github.com/Galdoba/utils"
	"github.com/urfave/cli"

	fcli "github.com/Galdoba/ffstuff/pkg/cli"
	"github.com/Galdoba/ffstuff/pkg/scanner"

	"github.com/Galdoba/ffstuff/pkg/config"
	"github.com/Galdoba/ffstuff/pkg/logfile"
)

var configMap map[string]string

var logger logfile.Logger
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
	// fldr.Init()
	// argsReceived()

	// root := configMap[constant.SearchRoot]
	// f, err := os.Stat(root)
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }
	// if !f.IsDir() {
	// 	fmt.Println(f.Name(), "is not directory")
	// 	os.Exit(2)
	// }
	// marker = configMap[constant.SearchMarker]

	// if configMap[constant.LogDirectory] == "default" {
	// 	logLocation = fldr.MuxPath() + "logfile.txt"
	// }
	// logger = logfile.New(logLocation, logfile.LogLevelWARN)

	// takeFile, err := scanner.Scan(root, marker)
	// fileList := scanner.ListReady(takeFile)
	// // if err := filepath.Walk(root, visit); err != nil {
	// if err != nil {
	// 	logger.ERROR(err.Error())
	// }
	// fmt.Println("")

	// /////////NEXT STAGE TEST
	// if len(fileList) == 0 {
	// 	fmt.Println("\rNothing new")
	// 	logger.INFO("No new files found")
	// 	return
	// }

	// logger.INFO(strconv.Itoa(len(fileList)) + " new files found")
	// fileList = sortResults(fileList)
	// //runInchecker(takeFile)
	// for _, val := range fileList {
	// 	fmt.Println("Can take", val)
	// }

	//os.Exit(0)
	//autoGrab := false
	root := configMap[constant.SearchRoot]
	marker := configMap[constant.SearchMarker]
	if configMap[constant.LogDirectory] == "default" {
		logLocation = fldr.MuxPath() + "logfile.txt"
	}
	logger = logfile.New(logLocation, logfile.LogLevelWARN)
	app := cli.NewApp()
	app.Version = "v 0.0.1"
	app.Name = "searcher"
	app.Usage = "Scans root directory and all subdirectories to create list of files that matches queary"
	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:  "grab",
			Usage: "If flag is active grabber will try to download all results",
			Value: false,
		},
	}
	app.Commands = []*cli.Command{
		//////////////////////////////////////
		{
			Name:  "new",
			Usage: "Searches all files in the root which associated with marker",
			Action: func(c *cli.Context) error {
				takeFile, err := scanner.Scan(root, marker)
				if err != nil {
					fmt.Println(err)

					return err
				}
				fileList := scanner.ListReady(takeFile)
				logger.INFO(strconv.Itoa(len(fileList)) + " new files found")
				fileList = sortResults(fileList)
				for _, val := range fileList {
					fmt.Println(val)
					fcli.RunConsole("grabber", "only", val)
				}
				fmt.Print("Flag is |", app.Flags[0].String())
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

func runInchecker(takeFile []string) []string {
	validFiles := []string{}
	logger.INFO("Run: " + "inchecker " + strings.Join(takeFile, " "))
	for _, file := range takeFile {
		_, _, err := fcli.RunConsole("inchecker", file)
		if err != nil {
			logger.ERROR(err.Error())

			continue
		}
		logger.TRACE("valid: " + file)
		validFiles = append(validFiles, file)
	}
	return validFiles
	// _, _, err := cli.RunConsole("inchecker", takeFile...)
	// if err != nil {
	// 	logger.ERROR(err.Error())
	// }
	// return validFiles
}

func defineRoot() string {
	fmt.Println("Enter path to root folder:")
	fmt.Print("Root=")
	str, err := user.InputStr()
	if err != nil {
		logger.WARN(err.Error())
	}
	config.SetField("ROOT", str)
	return str
}

// func argsReceived() {
// 	for _, val := range os.Args {
// 		val = strings.ToLower(val)
// 		switch val {
// 		case "--incheck", "-c":
// 			afterCheck = true
// 		case "--help", "-h":
// 			printHelp()
// 		}
// 	}

// }

func printHelp() {
	fmt.Print("Searcher walk all directories under the ROOT, and search any '[base].ready' files.\n")
	fmt.Print("After that it constructs result list of paths containing '[base]' in their names.\n")
	fmt.Print("This list can be used as arguments for other ffstuff aplications.\n")
	fmt.Print("\n")
	fmt.Print("ROOT=", configMap["ROOT"], "\n")
	fmt.Print("\n")
	fmt.Print("Keys:\n")
	fmt.Print(" -h, --help      -   show this message\n")
	fmt.Print(" -c, --incheck   -   run inchecker module on all files in result list\n")
	fmt.Print(" -g, --grab      -   run grabber module on all files in result list\n")
	os.Exit(0)
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

// func visit(path string, f os.FileInfo, err error) error {
// 	if f.IsDir() {
// 		clearLine()
// 		fmt.Print("\rSearch: ", path)
// 	}
// 	if !strings.Contains(f.Name(), marker) {
// 		return nil
// 	}
// 	dir, base := filepath.Split(path)
// 	base = strings.TrimSuffix(base, marker)
// 	files, err := ioutil.ReadDir(dir)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	for _, fl := range files {
// 		if strings.Contains(fl.Name(), base) && !strings.Contains(fl.Name(), marker) {
// 			takeFile = append(takeFile, dir+fl.Name())
// 		}
// 	}
// 	return nil
// }

func clearLine() {
	clr := ""
	for len(clr) < 196 {
		clr += " "
	}
	fmt.Print("\r" + clr)
}

func timeValue(t time.Time) int64 {
	y, m, d := t.Date()
	hh, mm, ss := t.Clock()
	valClock := int64(ss) + int64(mm*100) + int64(hh*10000)
	valDate := int64(d*1000000) + int64(m*100000000) + int64(y*10000000000)
	return valDate + valClock
}

func timeStr(tVal int64) string {
	tStr := ""
	// sec := int(tVal % 100)
	// tStr = strconv.Itoa(sec)
	// if sec < 10 {
	// 	tStr = "0" + tStr
	// }
	min := int(tVal%10000) / 100
	tStr = strconv.Itoa(min) /*+ ":"*/ + tStr
	if min < 10 {
		tStr = "0" + tStr
	}
	hr := int(tVal%1000000) / 10000
	tStr = strconv.Itoa(hr) + ":" + tStr
	if hr < 10 {
		tStr = "0" + tStr
	}

	day := int(tVal%100000000) / 1000000
	tStr = strconv.Itoa(day) + " " + tStr
	if day < 10 {
		tStr = "0" + tStr
	}
	mon := int(tVal%10000000000) / 100000000
	tStr = strconv.Itoa(mon) + "." + tStr
	if mon < 10 {
		tStr = "0" + tStr
	}
	yr := int(tVal%100000000000000) / 10000000000
	tStr = strconv.Itoa(yr) + "." + tStr
	if yr < 10 {
		tStr = "0" + tStr
	}
	if yr < 100 {
		tStr = "0" + tStr
	}
	if yr < 1000 {
		tStr = "0" + tStr
	}

	return tStr
}

func sortResults(list []string) []string {
	sorted := []string{}
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
