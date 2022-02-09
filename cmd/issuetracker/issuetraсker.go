package main

import (
	"fmt"

	"os"

	"github.com/urfave/cli"

	"github.com/Galdoba/ffstuff/fldr"
	"github.com/Galdoba/ffstuff/pkg/config"
	"github.com/Galdoba/ffstuff/pkg/glog"
	"github.com/Galdoba/ffstuff/pkg/scanner"
)

var configMap map[string]string

var logger glog.Logger
var logLocation string
var issueFilePath string

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
	issueFilePath = fldr.InPath() + "issues.txt"
	if _, err := os.Stat(issueFilePath); err != nil {
		os.Create(issueFilePath)
		fmt.Println("File created:", issueFilePath)
	}
}

func main() {
	//root := configMap[constant.SearchRoot]
	//marker := configMap[constant.SearchMarker]
	// if configMap[constant.LogDirectory] == "default" {
	// 	logLocation = fldr.MuxPath() + "logfile.txt"
	// }
	//logger = glog.New(logLocation, glog.LogLevelINFO)
	logger = glog.New(glog.LogPathDEFAULT, glog.LogLevelINFO)

	app := cli.NewApp()
	app.Version = "v 0.0.1"
	app.Name = "issuetracker"
	app.Usage = "Scans audio with Loudnorm and Soundscan and creates report files for analisys"
	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:  "vocal",
			Usage: "If flag is active soundscan will print data on terminal",
		},
	}
	app.Commands = []cli.Command{
		//////////////////////////////////////
		{
			Name:  "track",
			Usage: "Listens all files and checks silence in audio stream",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "repeat, r ",
					Usage: "If used creates loop to run issuetraker every x seconds",
					Value: "180",
				},
			},
			Action: func(c *cli.Context) error {
				found, err := scanner.Scan(fldr.InPath(), ".m4a")
				fmt.Println(err)
				for _, f := range found {
					fmt.Println(f)
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

/*
issuetracker run --repeat 42
init()
1. определяем папку IN
2. Создаем issues.txt
Run()
1. Создаем список Аудио
	1а. Ищем все видео в папке IN
	1б. Для каждого найденного: Добавляем в список все аудио ассациированные с ним
2. Удоставеряемся что для каждого аудио есть Report
	2a. FALSE: создаем Report
	2б. TRUE : Добавляем Report в список к анализу
3. Для каждого Report
	3a. анализируем Loudnorm
	3б. анализируем Soundscan
	3в. Принтуем аномалии в issues.txt
4.Публикуем issues.txt

track - запускает треккер
--repeat x  - повторяет процесс каждые x секунд

*/

func createIssueFile() {
	os.Create(issueFilePath)
	f, err := os.Open(issueFilePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()
}

type issuesReport struct {
	files []fileReport
}

type fileReport struct {
	data      []string
	loudnorm  bool
	soundscan bool
	issues    []string
}

/*
issueFile Exapmle:
START
File [file1.m4a]:
 Loudnorm warning:
 warning 1
 warning 2
 ...
 warning n
 Loudnorm report end.
 --------------------
 Soundscan report start:
 [report text line 1]
 [report text line 2]
 ...
 [report text line n]
 Soundscan report end.
 --------------------
[blank line]
File [file2.m4a]:
...
[blank line]
END
*/
