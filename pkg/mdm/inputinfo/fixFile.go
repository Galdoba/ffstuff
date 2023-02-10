package inputinfo

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Galdoba/utils"
)

func CleanScanData() {
	fmt.Print("Scan data searching... ")
	list := findDuplicateless(gatherInfo())
	l := len(list)
	fmt.Println("")
	for i, inp := range list {
		fmt.Printf("Adding input %v/%v\r", i+1, l)
		//fmt.Println(inp.String())
		AddNewDataToFile(inp)
	}
	fmt.Println("")
	fmt.Println("Creating backup...")
	backName := strings.TrimSuffix(originalFile, ".txt") + timestamp()
	os.Rename(originalFile, backName)
	fmt.Printf("%v . . . ok\n", backName)
	os.Rename(newFile, originalFile)
	fmt.Println("Scan data is ready to use")
}

func timestamp() string {
	tm := time.Now()
	return fmt.Sprintf("%v.txt", tm.Format("20060102_150405"))
}

func findDuplicateless(inputs []inputdata) []inputdata {
	newList := []inputdata{}
	os.Truncate(newFile, 0)
	unique := 0
	duplicates := 0
	total := len(inputs)
	fmt.Printf("Searching duplicates: \n")
	for i1, data1 := range inputs {
		duplicated := false
		for i2, data2 := range inputs {
			if i2 <= i1 {
				continue
			}

			if duplicated {
				continue
			}
			fmt.Printf("file %v of %v (u=%v/d=%v)                   \r", i1, total, unique, duplicates)
			if dataIsSame(data1, data2) {
				duplicated = true

			}
		}
		if !duplicated {
			newList = append(newList, data1)
			unique++
		} else {
			duplicates++
		}
	}
	fmt.Println("")
	return newList
}

func AddNewDataToFile(input inputdata) {
	appendTextToFile(newFile, "-START--------------------------------------------------------------------------\n")
	appendTextToFile(newFile, input.String())
	// for _, line := range input.data {

	// 	appendTextToFile(newFile, line)
	// }
	appendTextToFile(newFile, "-END----------------------------------------------------------------------------\n")
	appendTextToFile(newFile, " \n")
}

func appendTextToFile(filename, text string) {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	if _, err = f.WriteString(text); err != nil {
		panic(err)
	}
}

func dataIsSame(input1, input2 inputdata) bool {
	if len(input1.data) != len(input2.data) {
		return false
	}
	for i, line1 := range input1.data {

		switch {
		case strings.Contains(line1, "Started: "):
		case strings.Contains(line1, "Duration: N/A"):
		default:
			if line1 != input2.data[i] {
				return false
			}
		}

	}
	return true
}

var originalFile = `\\nas\buffer\IN\ScanData\input\ffmpeg\data.txt`
var newFile = `\\nas\buffer\IN\ScanData\input\ffmpeg\data_new.txt`

func originalFileLines() int {
	return len(utils.LinesFromTXT(originalFile))
}

func gatherInfo() []inputdata {
	allData := []inputdata{}
	data := inputdata{}
	// for _, line := range utils.LinesFromTXT(`C:\Users\a.akkulov\Desktop\cmdLine.go`) {
	// 	if line == `/*` {
	// 		data = inputdata{}
	// 		continue
	// 	}
	// 	if line == `*/` {
	// 		allData = append(allData, data)

	// 		continue
	// 	}
	// 	data.data = append(data.data, line)
	// }
	for _, line := range utils.LinesFromTXT(`\\nas\buffer\IN\ScanData\input\ffmpeg\data.txt`) {
		if line == `-START--------------------------------------------------------------------------` {
			data = inputdata{}
			continue
		}
		if line == `-END----------------------------------------------------------------------------` {
			data = skipTrashData(data)
			allData = append(allData, data)

			continue
		}
		data.data = append(data.data, line+"\n")
	}

	fmt.Println(len(allData), "data poits found")
	return allData
}

//, from 'Chetyre_sezona_v_gavane_s01_05_PRT230109004417_SER_00079_18.mp4':
//Input #0, mov,mp4,m4a,3gp,3g2,mj2, from '\\nas\ROOT\IN\@TRAILERS\_DONE\5_neizvestnyh_TRL\5_neizvestnyh_Trailer_20_rus_eng_HD.mov':
func skipTrashData(original inputdata) inputdata {
	fixed := inputdata{}
	for _, line := range original.data {
		switch {
		case strings.Contains(line, "At least one output file must be specified"):
			continue
		}
		if strings.Contains(line, "from '") && strings.Contains(line, "':") {
			startLine := strings.Split(line, "from '")[0]
			oriName := strings.Split(line, "from '")[1]
			oriName = strings.Split(oriName, "':")[0]
			shortName := filepath.Base(oriName)
			line = startLine + "from '" + shortName + "':\n"
		}
		fixed.data = append(fixed.data, line)
	}
	return fixed
}
