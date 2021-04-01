package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Galdoba/ffstuff/fldr"
	"github.com/Galdoba/ffstuff/pkg/glog"
	"github.com/Galdoba/ffstuff/pkg/muxer"
	"github.com/Galdoba/ffstuff/pkg/namedata"
	"github.com/Galdoba/ffstuff/pkg/scanner"
)

func main() {
	logger := glog.New(fldr.MuxPath()+"testLog.txt", 2)
	tasks, err := muxer.MuxList()
	if err != nil {
		logger.ERROR(err.Error())
		fmt.Println(err)
	}
	for i, task := range tasks {
		fmt.Print("Task ", i, "/", len(tasks), ":\n")
		files, muxTask, err := muxer.ChooseMuxer(task)
		if err != nil {
			logger.ERROR(err.Error())
			fmt.Println(err)
			continue
		}
		err = muxer.Run(muxTask, files)
		if err != nil {
			logger.ERROR(err.Error())
			fmt.Println(err)
		}
		logger.TRACE("Task Complete: " + task)
	}
	logger.INFO("Muxig Complete")

	// str, rtr, err := cli.RunConsole("inchecker", "\\\\nas\\ROOT\\EDIT\\21_02_20\\Ya_podaryu_tebe_pobedu_AUDIORUS51.m4a")
	// fmt.Println(str)
	// fmt.Println(rtr)
	// fmt.Println(err)

	// fmt.Println("")
	// fmt.Print(string(rune(9617))) // ░
	// fmt.Print(string(rune(9618))) //   ▒
	// fmt.Print(string(rune(9608))) //█
	// fmt.Print(string(rune(9612))) //▌
	// fmt.Print(string(rune(9619))) // ▓
	os.Exit(6)
	fmt.Print(string(rune(9617)))
	fmt.Println("Start")
	//list, err := scanner.ListContent("\\\\nas\\ROOT\\EDIT\\_sony\\Breaking_Bad_s01\\")
	list, err := scanner.Scan("\\\\192.168.31.4\\root\\EDIT\\", ".ready")
	fmt.Println("List:")
	for i := range list {
		fmt.Println(i, list[i])
		name := namedata.RetrieveShortName(list[i])
		name = strings.TrimSuffix(name, ".ready")
		dir := namedata.RetrieveDirectory(list[i])
		sList, err2 := scanner.ListContent(dir)
		if err2 != nil {
			fmt.Println(err.Error())
		}
		for f := range sList {
			if strings.Contains(sList[f], name) {
				fmt.Print(sList[f])
				fmt.Print("\n")
			}
		}
	}
	fmt.Println("Err:", err)
	for now := 0; now < 100; now++ {
		downloadbar(now)
		time.Sleep(time.Millisecond * 50)
	}
	fmt.Println("")
	fmt.Println("end")
}

func downloadbar(now int) {
	if now > 100 || now < 0 {
		return
	}
	s := "["
	for i := 0; i < 25; i++ {
		if i < now/4 {
			s += string(rune(9608))
			continue
		}
		if now == 100 {
			s += string(rune(9608))
			continue
		}
		if i == now/4 {
			s += string(rune(9612))
			continue
		}

		s += " "
	}
	s += "]"
	fmt.Print(s)
}

/*
10
██████████ 10%

20
████████████████████ 5%
25
█████████████████████████ 4%
30
██████████████████████████████ 4%
40
[█████████████████████████████████████▌  ] 2.5%
50
██████████████████████████████████████████████████ 2%

[1234567890123456789012345678901234567890]
[progress: 100.0%]

if now%4 >= 2 {
				s += string(rune(9612))
				continue
*/
