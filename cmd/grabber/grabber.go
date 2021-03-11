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
)

/*
TZ:
>> grab filename.ready			-- забрать все связанное с ready файлом
>> grab -h						-- вывести на экран помогалку 							--help
>> grab -n						-- забрать все новое (предварительное сканирование)		--new
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
	copyList := []string{}
	mainDestDir := configMap["INPATH"]
	if mainDestDir == "default" {
		mainDestDir = fldr.InPath()
	}
	secondaryDestDir := configMap["OUTPATH"]
	if secondaryDestDir == "default" {
		secondaryDestDir = fldr.OutPath()
	}
	//fmt.Println("START GRABBER")
	//fmt.Println("GRABBER ARGUMENTS:")
	for _, arg := range argsReceived() {
		//	fmt.Println(i, arg)
		switch describeArg(arg) {
		case 0:
			copyList = append(copyList, scanForAssociatedFiles(arg)...)
		case 1:
			printHelp()

		}
	}

	for _, file := range copyList {
		fmt.Println("GRABBING", file)

		if err := grabber.CopyFile(file, mainDestDir); err != nil {
			fmt.Println("Grabber Error:", err.Error())
		}
	}

	// logger := logfile.New(fldr.MuxPath()+"logfile.txt", logfile.LogLevelINFO)
	// args := os.Args
	//fmt.Println("END GRABBER")
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
