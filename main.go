package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Galdoba/ffstuff/pkg/namedata"
	"github.com/Galdoba/ffstuff/pkg/scanner"
)

func main() {
	for i := 0; i < 30000; i++ {
		dn := int64(i)
		var size int64
		size = 25000
		fmt.Print(downloadbar(dn, size, true))
		time.Sleep(time.Millisecond * 500)
	}

	// logger := glog.New(fldr.MuxPath()+"testLog.txt", 2)
	// tasks, err := muxer.MuxList()
	// if err != nil {
	// 	logger.ERROR(err.Error())
	// 	fmt.Println(err)
	// }
	// for i, task := range tasks {
	// 	fmt.Print("Task ", i, "/", len(tasks), ":\n")
	// 	files, muxTask, err := muxer.ChooseMuxer(task)
	// 	if err != nil {
	// 		logger.ERROR(err.Error())
	// 		fmt.Println(err)
	// 		continue
	// 	}
	// 	err = muxer.Run(muxTask, files)
	// 	if err != nil {
	// 		logger.ERROR(err.Error())
	// 		fmt.Println(err)
	// 	}
	// 	logger.TRACE("Task Complete: " + task)
	// }
	// logger.INFO("Muxig Complete")

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

		time.Sleep(time.Millisecond * 50)
	}
	fmt.Println("")
	fmt.Println("end")
}

func downloadbar(bts, size int64, percentage bool) string {
	str := ""
	if size == 0 {
		size = 1
	}
	if percentage {
		prc := float64(bts) / float64(size/100)
		prcStr := strconv.FormatFloat(prc, 'f', 3, 64)
		str = "[ " + prcStr + "% ]\r"
		return str
	}
	str = "Downloaded: " + size2GbString(bts) + "/" + size2GbString(size) + " Gb\r"
	return str

}

func size2GbString(bts int64) string {
	gbt := float64(bts) / 1073741824.0
	gbtStr := strconv.FormatFloat(gbt, 'f', 2, 64)
	return gbtStr
}

// func progressBar(now, all int64) string {

// }

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
[ progress: 100.000% ]
[>>>>>>              ]
[++++++--------------]

if now%4 >= 2 {
				s += string(rune(9612))
				continue
*/

/*
[]

*/

//DurationStamp - return duration (float64 - seconds) as a string in format: [HH:MM:SS.ms]
func DurationStamp(dur float64) string {
	if dur < 0 {
		return "NEGATIVE"
	}
	stamp := ""
	hh := int(dur) / int(3600)
	hLeft := int(dur) % int(3600)
	mm := hLeft / 60
	ss := hLeft % 60
	sLeft := dur - (float64(hh*3600) + float64(mm*60) + float64(ss))
	ms := int(sLeft * 1000)
	////////
	hhs := strconv.Itoa(int(hh))
	for len(hhs) < 2 {
		hhs = "0" + hhs
	}
	mms := strconv.Itoa(int(mm))
	for len(mms) < 2 {
		mms = "0" + mms
	}
	sss := strconv.Itoa(int(ss))
	for len(sss) < 2 {
		sss = "0" + sss
	}
	mss := strconv.Itoa(int(ms))
	for len(mss) < 3 {
		mss = "0" + mss
	}
	stamp = hhs + ":" + mms + ":" + sss + "." + mss
	return stamp

}
