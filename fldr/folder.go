package fldr

// import (
// 	"fmt"
// 	"io/ioutil"
// 	"log"
// 	"os"
// 	"path/filepath"
// 	"runtime"
// 	"strconv"
// 	"strings"

// 	"github.com/Galdoba/devtools/cli/user"
// 	"github.com/Galdoba/ffstuff/constant"
// 	"github.com/Galdoba/ffstuff/pkg/config"
// 	"github.com/Galdoba/utils"
// )

// // const (
// // 	inFolder  = "d:\\IN\\IN_"   //TODO: запихнуть эти адресса в конфиг
// // 	muxFolder = "d:\\MUX\\MUX_" //
// // 	outFolder = "d:\\OUT\\OUT_" //
// // )

// var workdate string
// var inFolder string
// var muxFolder string
// var outFolder string

// func Init() {
// 	//fmt.Print("Create directory: '", InPath(), "' \n")
// }

// func init() {

// 	//fmt.Print("Initiate 'fldr' module...\n")
// 	workdateTemp := utils.DateStamp()
// 	workdate = workdateTemp
// 	conf, err := config.ReadProgramConfig("ffstuff")
// 	if err != nil {
// 		fmt.Println("fldr init()")
// 		fmt.Println(err)
// 		panic(err)
// 	}
// 	inFolder = conf.Field[constant.InPath]
// 	muxFolder = conf.Field[constant.MuxPath]
// 	outFolder = conf.Field[constant.OutPath]

// 	// if _, err := os.Stat(InPath()); os.IsNotExist(err) {
// 	// 	os.Mkdir(InPath(), 0700)
// 	// 	fmt.Print("Create directory: '", InPath(), "' \n")
// 	// }
// 	// if _, err := os.Stat(MuxPath()); os.IsNotExist(err) {
// 	// 	os.Mkdir(MuxPath(), 0700)
// 	// 	fmt.Print("Create directory: '", MuxPath(), "'\n")
// 	// }
// 	// if _, err := os.Stat(OutPath()); os.IsNotExist(err) {
// 	// 	os.Mkdir(OutPath(), 0700)
// 	// 	fmt.Print("Create directory: '", OutPath(), "'\n")
// 	// }
// 	//fmt.Print("'fldr'...ok\n")
// }

// func Test() {
// 	//fmt.Println("Test")
// }

// //InPath - Возвращает сегодняшнюю папку для скачивания
// func InPath() string {
// 	return inFolder + "IN_" + utils.DateStamp() + "\\"
// }

// //MuxPath - Возвращает сегодняшнюю папку для мукса
// func MuxPath() string {
// 	return muxFolder + "MUX_" + utils.DateStamp() + "\\"
// }

// //OutPath - Возвращает сегодняшнюю папку для проверки/отправки
// func OutPath(dynamic bool) string {
// 	if dynamic {
// 		return outFolder + "OUT_" + utils.DateStamp() + "\\"
// 	}
// 	return outFolder + "OUT_unchecked" + "\\"
// }

// // func SelectEDL() string {
// // 	return ""
// // }

// func SelectEDL() string {
// 	files := []string{}
// 	files = append(files, filesByExtention(".edl")...)
// 	files = append(files, "Exit")
// 	_, edlFile := menu("Select EDL file:", files...)
// 	if edlFile == "Exit" {
// 		os.Exit(1)
// 	}
// 	return InPath() + edlFile
// }

// func filesByExtention(extention string) []string {
// 	var names []string
// 	files, err := ioutil.ReadDir(InPath()) //InPath() +
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	for _, f := range files {
// 		if strings.Contains(f.Name(), extention) {
// 			names = append(names, f.Name())
// 		}
// 	}
// 	return names
// }

// func menu(question string, options ...string) (int, string) {
// 	fmt.Println(question)
// 	for i := range options {
// 		prefix := " [" + strconv.Itoa(i) + "] - "
// 		fmt.Println(prefix + options[i])
// 	}
// 	answerGl := 0
// 	gotIt := false
// 	for !gotIt {
// 		answer, err := user.InputInt()
// 		if err != nil {
// 			fmt.Println("\033[FError: " + err.Error())
// 			fmt.Println(question)
// 			continue
// 		}
// 		if answer >= len(options) || answer < 0 {
// 			fmt.Println("\033[FError: Option", answer, "is invalid")
// 			fmt.Println(question)
// 			continue
// 		}

// 		if answer < len(options) {
// 			gotIt = true
// 			answerGl = answer
// 		}
// 	}
// 	//fmt.Println(answerGl, options[answerGl])
// 	return answerGl, options[answerGl]
// 	//return a, text
// }

// /*

//  */
// func ExecutableBase() string {
// 	processInit, err := os.Executable()
// 	if err != nil {
// 		fmt.Println(err)
// 		panic(err)
// 	}
// 	processInit = filepath.Base(processInit)
// 	processInit = strings.TrimSuffix(processInit, ".exe")
// 	return processInit

// }

// func AutoPath(file string) string {
// 	home, err := os.UserHomeDir()
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	exe := ExecutableBase()
// 	autoPath := home + "\\galdoba" + "\\" + exe + "\\" + file
// 	return autoPath
// }

// func configPath() (string, string) {
// 	home, err := os.UserHomeDir()
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	exe, err := os.Executable()
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	exe = filepath.Base(exe)
// 	configDir := ""
// 	switch runtime.GOOS {
// 	case "windows":
// 		exe = strings.TrimSuffix(exe, ".exe")
// 		configDir = home + "\\config\\" + exe // + exe + ".config"

// 	}
// 	return configDir, exe + ".config"
// }

// func LogPathDefault() string {
// 	home, err := os.UserHomeDir()
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	logDir := ""
// 	switch runtime.GOOS {
// 	case "windows":
// 		logDir = home + "\\.logs\\" + utils.DateStamp() + "_ffstuff.log"
// 	}
// 	return logDir
// }
