package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/Galdoba/ffstuff/fldr"
	"github.com/Galdoba/ffstuff/pkg/config"
	"github.com/Galdoba/ffstuff/pkg/grabber"
	"github.com/Galdoba/ffstuff/pkg/logfile"
	"github.com/Galdoba/ffstuff/pkg/namedata"
	"github.com/Galdoba/ffstuff/pkg/scanner"
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

type mode struct {
	logging bool
	vocal   bool
}

var configMap map[string]string
var moduleMode mode
var logger logfile.Logger

func init() {
	err := errors.New("Initial obstract error")
	configMap, err = config.Read() //CHECK config file
	if err != nil {
		switch err.Error() {
		case "Config file not found":
			fmt.Println("Creating config...")
			_, err := config.Construct()
			if err != nil {
				panic(err)
			}
			//config.SetField("marker", ".ready") //marker files
			config.SetField("INPATH", "default")
			config.SetField("OUTPATH", "default")
			config.SetField("LOGLOCATION", "default")
		}
	}
}

func main() {
	if err := config.Verify(); err != nil {
		fmt.Println(err.Error())
		return
	}
	conf, _ := config.Read()
	dest := conf["OUTPATH"]
	if dest == "default" {
		dest = fldr.InPath()
	}
	app := cli.NewApp()
	app.Version = "v 0.0.1"
	app.Name = "grabber"
	app.Usage = "dowloads files and sort it to working directories"
	app.Commands = []*cli.Command{
		//////////////////////////////////////
		{
			Name:  "takeonly",
			Usage: "Download only those files, that was received as arguments",
			Action: func(c *cli.Context) error {
				paths := c.Args().Slice() //	path := c.String("path") //*cli.Context.String(key) - вызывает флаг с именем key и возвращает значение Value
				for _, path := range paths {
					fmt.Println("GRABBER DOWNLOADING FILE:", path)
					err := grabber.CopyFile(path, dest)
					fmt.Println(err)
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
		//args = append(args, "help") //Принудительно зовем помощь если нет других аргументов
	}
	if err := app.Run(args); err != nil {
		fmt.Println(err.Error())
	}
}

func argsReceived() []string {
	outArgs := []string{}
	for i, val := range os.Args {
		if len(os.Args) == 1 {
			fmt.Println("No аrguments received. Pass '-h' to get Help")
			os.Exit(0)
		}
		if i == 0 {
			continue
		}

		outArgs = append(outArgs, val)
	}
	return outArgs
}

func describeArg(arg string) int {
	actionCode := -999
	switch arg {
	default:
		if strings.Contains(arg, ".ready") {
			//fmt.Println(arg, "- Valid argument: grabber will call scanner and download all associated files to InFolder")
			//fmt.Println("In this case:", scanForAssociatedFiles(arg))
			return 0
		}
		fmt.Println(arg, "- Invalid argument: grabber will ignore it")
	case "-h", "--help":
		//fmt.Println("This flag prints help text and exits program")
		return 1
	}
	return actionCode
}

func scanForAssociatedFiles(readyFile string) []string {
	directory := namedata.RetrieveDirectory(readyFile)
	base := namedata.RetrieveBase(readyFile)
	found, err := scanner.Scan(directory, base)
	fmt.Println("---", found)
	if err != nil {
		fmt.Println(err)
	}
	return found
}

func printHelp() {
	fmt.Println("'grabber' отвечает за скачивание и распределение входящих файлов по папкам.")
	fmt.Println("")
	fmt.Println("Типовая команда для использования в консоли:")
	fmt.Println("grabber \\\\nas\\ROOT\\EDIT\\21_02_20\\Critical_Thinking.ready")
	fmt.Println("")
	fmt.Println("Принцип работы:")
	fmt.Println("Отталкиваясь от имени ready-файла (аргумент) grabber ищет все файлы в этой же папке с повторяющейся базой имени,")
	fmt.Println("после чего все закачивает в папку входящих (берется из конфига или модуля 'fldr')")
	fmt.Println("")
	fmt.Println("Keys:")
	fmt.Println("-a, --all   -  запускает модуль 'scanner' который ищет все ready файлы и использует их как аргументы")
	fmt.Println("-v, --valid -  запускает модуль 'inchecker' перед скачиванием и скачивает только если нет ошибок")
	fmt.Println("-p, --proxy -  скачивает только прокси и звук")

}
