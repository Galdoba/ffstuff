package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/Galdoba/ffstuff/pkg/scanner"

	"github.com/Galdoba/ffstuff/fldr"
	"github.com/Galdoba/ffstuff/pkg/cli"
	"github.com/Galdoba/ffstuff/pkg/grabber"

	"github.com/Galdoba/devtools/cli/user"
	"github.com/Galdoba/ffstuff/pkg/config"
	"github.com/Galdoba/ffstuff/pkg/logfile"
)

var configMap map[string]string
var takeFile []string
var marker string
var logger logfile.Logger
var logLocation string
var afterCheck bool

func init() {
	err := errors.New("Initial obstract error")

	configMap, err = config.Read() //CHECK config file
	if err != nil {
		switch err.Error() {
		case "Config file not found":
			fmt.Print("Creating config...")
			_, err := config.Construct()
			if err != nil {
				panic(err)
			}
			config.SetField("marker", ".ready") //marker files
			config.SetField("root", "UNDEFINED")
			config.SetField("logLocation", "default")
			fmt.Print("	done\n")
		}
	}
}

func main() {
	fldr.Init()
	argsReceived()
	takeFile = []string{}

	results, err := scanner.Scan("\\\\192.168.31.4\\root\\EDIT\\", ".ready")
	if err != nil {
		fmt.Println(err.Error())
	}
	for i, val := range results {
		fmt.Println(i, val)
	}

	//os.Exit(5)

	root := configMap["ROOT"]
	if root == "UNDEFINED" {
		fmt.Println("Search root undefined:")
		fmt.Println("Set root in: " + config.ConfigFile())
		fmt.Println("End Program")
		os.Exit(3)
	}

	marker = configMap["MARKER"]

	if configMap["LOGLOCATION"] == "default" {
		logLocation = fldr.MuxPath() + "logfile.txt"
	}
	logger = logfile.New(logLocation, logfile.LogLevelWARN)

	//takeFile, err = scan.ScanReady(root, marker)

	if err := filepath.Walk(root, visit); err != nil {
		logger.ERROR(err.Error())
	}
	fmt.Println("")

	/////////NEXT STAGE TEST
	if len(takeFile) == 0 {
		fmt.Println("\rNothing new")
		logger.INFO("No new files found")
		return
	}

	logger.INFO(strconv.Itoa(len(takeFile)) + " new files found")

	runInchecker(takeFile)
	for _, val := range takeFile {
		fmt.Println("Can take", val)
	}
	//os.Exit(2)
	for _, val := range takeFile {
		if strings.Contains(val, ".srt") {
			if err := grabber.CopyFile(val, "d:\\SENDER\\"); err != nil {
				logger.ERROR(err.Error())
			} else {
				logger.TRACE(val + " copied to d:\\SENDER\\")
			}
			continue
		}
		if err := grabber.CopyFile(val, fldr.InPath()); err != nil {
			logger.ERROR(err.Error())
		} else {
			logger.TRACE(val + " copied to " + fldr.InPath())
		}

	}

}

func runInchecker(takeFile []string) []string {
	validFiles := []string{}
	// logger.INFO("Run: " + "inchecker " + strings.Join(takeFile, " "))
	// for _, file := range takeFile {
	// 	_, _, err := cli.RunConsole("inchecker", file)
	// 	if err != nil {
	// 		logger.ERROR(err.Error())
	// 		continue
	// 	}
	// 	logger.TRACE("valid: " + file)
	// 	validFiles = append(validFiles, file)
	// }
	// return validFiles
	_, _, err := cli.RunConsole("inchecker", takeFile...)
	if err != nil {
		logger.ERROR(err.Error())
	}
	return validFiles
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

func argsReceived() {
	for _, val := range os.Args {
		val = strings.ToLower(val)
		switch val {
		case "--incheck", "-c":
			afterCheck = true
		case "--help", "-h":
			printHelp()
		}
	}

}

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

func visit(path string, f os.FileInfo, err error) error {
	if f.IsDir() {
		clearLine()
		fmt.Print("\rSearch: ", path)
	}
	if !strings.Contains(f.Name(), marker) {
		return nil
	}
	dir, base := filepath.Split(path)
	base = strings.TrimSuffix(base, marker)
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, fl := range files {
		if strings.Contains(fl.Name(), base) && !strings.Contains(fl.Name(), marker) {
			takeFile = append(takeFile, dir+fl.Name())
		}
	}
	return nil
}

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
