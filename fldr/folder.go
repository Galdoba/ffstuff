package fldr

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/Galdoba/devtools/cli/user"
	"github.com/Galdoba/utils"
)

const (
	inFolder  = "f:\\Work\\petr_proj\\___IN\\IN_" //TODO: запихнуть эти адресса в конфиг
	muxFolder = "e:\\_OUT\\MUX_"                  //
	outFolder = "d:\\SENDER\\DONE_"               //
)

var workdate string

func Init() {
	//fmt.Print("Create directory: '", InPath(), "' \n")
}

func init() {
	//fmt.Print("Initiate 'fldr' module...\n")
	workdateTemp := utils.DateStamp()
	workdate = workdateTemp
	if _, err := os.Stat(InPath()); os.IsNotExist(err) {
		os.Mkdir(InPath(), 0700)
		fmt.Print("Create directory: '", InPath(), "' \n")
	}
	if _, err := os.Stat(MuxPath()); os.IsNotExist(err) {
		os.Mkdir(MuxPath(), 0700)
		fmt.Print("Create directory: '", MuxPath(), "'\n")
	}
	if _, err := os.Stat(OutPath()); os.IsNotExist(err) {
		os.Mkdir(OutPath(), 0700)
		fmt.Print("Create directory: '", OutPath(), "'\n")
	}

	//fmt.Print("'fldr'...ok\n")
}

func Test() {
	//fmt.Println("Test")
}

//InPath - Возвращает сегодняшнюю папку для скачивания
func InPath() string {
	return inFolder + workdate + "\\"
}

//MuxPath - Возвращает сегодняшнюю папку для мукса
func MuxPath() string {
	return muxFolder + workdate + "\\"
}

//OutPath - Возвращает сегодняшнюю папку для проверки/отправки
func OutPath() string {
	return outFolder + workdate + "\\"
}

// func SelectEDL() string {
// 	return ""
// }

func SelectEDL() string {
	files := []string{}
	files = append(files, filesByExtention(".edl")...)
	files = append(files, "Exit")
	_, edlFile := menu("Select EDL file:", files...)
	if edlFile == "Exit" {
		os.Exit(1)
	}
	return InPath() + edlFile
}

func filesByExtention(extention string) []string {
	var names []string
	files, err := ioutil.ReadDir(InPath() + ".")
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		if strings.Contains(f.Name(), extention) {
			names = append(names, f.Name())
		}
	}
	return names
}

func menu(question string, options ...string) (int, string) {
	fmt.Println(question)
	for i := range options {
		prefix := " [" + strconv.Itoa(i) + "] - "
		fmt.Println(prefix + options[i])
	}
	answerGl := 0
	gotIt := false
	for !gotIt {
		answer, err := user.InputInt()
		if err != nil {
			fmt.Println("\033[FError: " + err.Error())
			fmt.Println(question)
			continue
		}
		if answer >= len(options) || answer < 0 {
			fmt.Println("\033[FError: Option", answer, "is invalid")
			fmt.Println(question)
			continue
		}

		if answer < len(options) {
			gotIt = true
			answerGl = answer
		}
	}
	//fmt.Println(answerGl, options[answerGl])
	return answerGl, options[answerGl]
	//return a, text
}
