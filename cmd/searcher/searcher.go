package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/Galdoba/ffstuff/pkg/cli"
	"github.com/Galdoba/ffstuff/pkg/config"
	"github.com/Galdoba/utils"
)

func visit(path string, f os.FileInfo, err error) error {
	marker := ".ready"
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

func main() {
	config.Load("READYMARKER")

	os.Exit(1)
	for {
		config.Load("READYMARKER")
		utils.ClearScreen()
		takeFile = []string{}
		flag.Parse()
		root := flag.Arg(0)
		err := filepath.Walk(root, visit)
		//fmt.Printf("filepath.Walk() returned %v\n", err)
		if err != nil {
			fmt.Print(err)
		}
		clearLine()
		if len(takeFile) > 0 {
			fmt.Println("\rNew Files Found:")
		} else {
			fmt.Println("\rNothing new")
		}
		for _, val := range takeFile {
			//fmt.Println(val)
			cli.RunConsole("inchecker", val)
		}
		wait := time.Second * 20
		time.Sleep(wait)
	}
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

var takeFile []string
