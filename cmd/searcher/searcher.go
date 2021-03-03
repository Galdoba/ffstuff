package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/Galdoba/ffstuff/fldr"
	"github.com/Galdoba/ffstuff/pkg/cli"

	"github.com/Galdoba/devtools/cli/user"
	"github.com/Galdoba/ffstuff/pkg/config"
	"github.com/Galdoba/ffstuff/pkg/logfile"
)

var configMap map[string]string
var takeFile []string
var marker string

func init() {
	configMapTemp := make(map[string]string)
	configMapTemp, err := config.Read()
	if err != nil {
		fmt.Println(err)
		if err.Error() == "Config file not found" {
			fmt.Print("Creating config...")
			_, err := config.Construct()
			if err != nil {
				fmt.Println(err)
				os.Exit(3)
			}
			config.SetField("marker", ".ready")
			config.SetField("root", "UNDEFINED")
			fmt.Print("	done\n")
		}
	}
	configMap = configMapTemp
}

func main() {
	fldr.Init()
	logger := logfile.New(fldr.MuxPath()+"logfile.txt", logfile.LogLevelWARN)

	// configMap, configErr := config.Read()
	// if configErr != nil {
	// 	fmt.Println(configErr)
	// 	os.Exit(4)
	// }
	takeFile = []string{}
	// flag.Parse()
	// fmt.Println(configMap)
	// root := flag.Arg(0)
	root := configMap["ROOT"]
	if root == "UNDEFINED" {
		fmt.Println("Enter path to root folder:")
		fmt.Print("Root=")
		str, err := user.InputStr()
		if err != nil {
			logger.WARN(err.Error())
		}
		config.SetField("ROOT", str)
		root = str
	}
	//fmt.Println("root =", root)
	marker = configMap["MARKER"]

	err := filepath.Walk(root, visit)
	//fmt.Printf("filepath.Walk() returned %v\n", err)
	if err != nil {
		logger.ERROR(err.Error())
	}
	clearLine()

	/////////NEXT STAGE TEST
	if len(takeFile) == 0 {
		fmt.Println("\rNothing new")
		logger.INFO("No new files found")
		return
	}

	logger.INFO(strconv.Itoa(len(takeFile)) + " new files found")
	// fmt.Println("\rNew Files Found:")

	// fmt.Println("Checking via 'inchecker':")

	logger.INFO("Run: " + "inchecker " + strings.Join(takeFile, " "))
	_, _, err = cli.RunConsole("inchecker", takeFile...)
	if err != nil {
		logger.ERROR(err.Error())
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
