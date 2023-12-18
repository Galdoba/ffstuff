package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/Galdoba/ffstuff/app/profiler/ump"
)

func Test_ConsumeJSON(t *testing.T) {

	pr, err := ump.New(`\\192.168.31.4\buffer\IN\ScanData\input\Shifter_5.1_RUS.mov`)
	fmt.Println("test struct:", pr)
	if pr != nil {
		fmt.Println(pr.Short())
	}
	if err != nil {
		fmt.Println(err.Error())
	}
	bitMap := make(map[string]int)
	//	return
	dir := `\\192.168.31.4\buffer\IN\ScanData\input\files\`
	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	for testNum, e := range entries {
		if testNum > 0 {
			break
		}
		fmt.Printf("%v\r", testNum)
		path := dir + e.Name()
		sr, err := ump.ConsumeJSON(path)
		if err != nil {
			t.Errorf("%v", err.Error())
		}
		warnMsg := ""
		//fmt.Println(sr.Format.Filename)
		//fmt.Println(sr.Short())
		for _, warn := range sr.Warnings() {
			warnMsg += warn + "\n"
		}

		// 	fmt.Println(sr.Short())
		if warnMsg != "" {
			fmt.Println(testNum, "-------")
			fmt.Println(sr.Format.Filename)
			fmt.Println("SHORT:", sr.Short())
			fmt.Println("LONG :", sr.Long())
			fmt.Println("WARNINGS :")
			fmt.Println(warnMsg)

		}

	}
	for k, v := range bitMap {
		fmt.Println(k, v)
	}
}

// func TestFields(t *testing.T) {
// 	dir := `\\192.168.31.4\buffer\IN\ScanData\input\files\`
// 	entries, err := os.ReadDir(dir)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	//depth1Keys := []string{"_streams", "_format"}
// 	datamap := make(map[string]int)
// 	ln := len(entries)
// 	for eNum, e := range entries {
// 		// if eNum > 10000 {
// 		// 	break
// 		// }
// 		// //fmt.Println(e.Name())
// 		// bt, err := os.ReadFile(dir + e.Name())
// 		// if err != nil {
// 		// 	t.Errorf("OpenFile: %v", err.Error())
// 		// }
// 		// s := string(bt)
// 		// lines := strings.Split(s, "\n")
// 		// lastDepth := 0
// 		// //lastQoutes := ""
// 		// last1 := ""
// 		// last2 := ""
// 		// last3 := ""
// 		// last4 := ""
// 		// last5 := ""
// 		// last6 := ""
// 		// skip := 0
// 		// for l, line := range lines {

// 		// 	qt, dp := parseQuoted(line)

// 		// 	if dp <= lastDepth {
// 		// 		switch dp {
// 		// 		case 1:
// 		// 			last2 = ""
// 		// 			last3 = ""
// 		// 			last4 = ""
// 		// 			last5 = ""
// 		// 			last6 = ""
// 		// 		case 2:
// 		// 			last2 = ""
// 		// 			last3 = ""
// 		// 			last4 = ""
// 		// 			last5 = ""
// 		// 			last6 = ""
// 		// 		case 3:
// 		// 			last3 = ""
// 		// 			last4 = ""
// 		// 			last5 = ""
// 		// 			last6 = ""
// 		// 		case 4:
// 		// 			last4 = ""
// 		// 			last5 = ""
// 		// 			last6 = ""
// 		// 		case 5:
// 		// 			last5 = ""
// 		// 			last6 = ""
// 		// 		case 6:
// 		// 			last6 = ""

// 		// 		}
// 		// 	}
// 		// 	if strings.Contains(line, " }") {
// 		// 		skip++
// 		// 		continue
// 		// 	}
// 		// 	skip = 0
// 		// 	if qt != "" {
// 		// 		switch dp {
// 		// 		case 1:
// 		// 			last1 = qt
// 		// 			last2 = ""
// 		// 			last3 = ""
// 		// 			last4 = ""
// 		// 			last5 = ""
// 		// 			last6 = ""
// 		// 		case 2:
// 		// 			last2 = qt
// 		// 			last3 = ""
// 		// 			last4 = ""
// 		// 			last5 = ""
// 		// 			last6 = ""
// 		// 		case 3:
// 		// 			last3 = qt
// 		// 			last4 = ""
// 		// 			last5 = ""
// 		// 			last6 = ""
// 		// 		case 4:
// 		// 			last4 = qt
// 		// 			last5 = ""
// 		// 			last6 = ""
// 		// 		case 5:
// 		// 			last5 = qt
// 		// 			last6 = ""
// 		// 		case 6:
// 		// 			last6 = qt

// 		// 		}
// 		// 	}

// 		// 	//fmt.Printf("%v||\n", []byte(line))

// 		// 	//fmt.Println(qt, dp)
// 		// 	if eNum%84 == 0 {
// 		// 		fmt.Printf("%v/%v\r", eNum, ln)
// 		// 	}
// 		// 	datamap[last1+last2+last3+last4+last5+last6]++
// 		// 	if last1+last2+last3+last4+last5+last6 == "_format___disposition____timed_thumbnails" {
// 		// 		fmt.Println(line, l)
// 		// 		fmt.Println(last1)
// 		// 		fmt.Println(last2)
// 		// 		fmt.Println(last3)
// 		// 		fmt.Println(last4)
// 		// 		fmt.Println(last5)
// 		// 		fmt.Println(last6)
// 		// 		fmt.Println(e.Name())

// 		// 		panic("===================")
// 		// 	}
// 		// 	lastDepth = dp
// 		//}
// 		sr := &ScanResult{}
// 		sr.Format = &Format{Filename: "testName", Nb_streams: 5}
// 		data, err := os.ReadFile(dir + e.Name())
// 		//data0, err2 := json.MarshalIndent(&sr, "", "  ")

// 		assertNoError(err)
// 		if len(data) == 0 {
// 			data, err = json.MarshalIndent(&sr, "", "  ")
// 			if err != nil {
// 				println(err.Error())
// 				//os.Exit(1)
// 			}
// 		}
// 		err = json.Unmarshal(data, &sr)
// 		if err != nil {
// 			errText := fmt.Sprintf("can't unmarshal config data: %v", err.Error())
// 			println(errText)
// 			os.Exit(1)
// 		}
// 		//fmt.Println(string(data))
// 		//fmt.Println(string(data0))
// 		// if err2 != nil {
// 		// 	fmt.Println(err.Error())
// 		// }
// 		fmt.Println(sr.Format.Filename)
// 		//fmt.Println(sr.Streams[2].Codec_name)
// 		fmt.Println(sr.String(), ln, eNum)
// 		//time.Sleep(time.Second)

// 	}
// 	fmt.Println("")
// 	keys := []string{}
// 	for k := range datamap {
// 		keys = append(keys, k)
// 	}
// 	sort.Strings(keys)
// 	for _, k := range keys {
// 		fmt.Println(k, datamap[k])
// 	}
// 	//fmt.Println(datamap)

// }

func parseQuoted(line string) (string, int) {
	prsed := ""
	depth := 0
	for strings.HasPrefix(line, "    ") {
		prsed += "_"
		depth++
		line = strings.TrimPrefix(line, "    ")
	}
	if strings.HasPrefix(line, `"`) {
		parts := strings.Split(line, `"`)
		prsed += parts[1]
		return prsed, depth
	}
	return "", -1
}
