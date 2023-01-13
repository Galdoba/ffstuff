package inputinfo

import (
	"fmt"
	"testing"

	"github.com/Galdoba/utils"
)

func TestInputReading(t *testing.T) {
	examples := gatherInfo()
	//darMap := make(map[string]int)
	//out := []string{}
	metadataKeys := make(map[int]int)
	for i, input := range examples {

		if i < 0 || i > 2000 {
			continue
		}
		pi, err := parse(input)
		if err != nil {
			t.Errorf(err.Error())
		}
		//if pi.filename != "" {
		fmt.Printf("example %v: \n", i)

		fmt.Printf("%v", pi)
		//for k, _ := range pi.metadata {
		metadataKeys[pi.globBitrate]++

		if pi.globBitrate == -1 {
			fmt.Println(pi.filename, "------")

		}
		//}

		//}
		// mf := CollectData(input)
		// name, foundName := input.parseName()
		// if foundName == -1 {
		// 	t.Errorf("expecting name to be foudn in: \n%v", input.String())
		// 	panic(1)
		// }
		// mf.filename = name
		// //if strings.Contains(name, " ") {
		// fmt.Println("")
		// fmt.Println(i, mf)
		// strDataAll := input.parseStreamData()
		// for _, v := range strDataAll {
		// 	fmt.Println(v)
		// }
		// if len(strDataAll) == 0 {
		// 	t.Errorf("expecting streamData to be foudn in: \n%v", input.String())
		// 	//panic(2)
		// }
		// video, _, _ := sortStreamData(strDataAll)

		// for _, v := range video {
		// 	vs := collectVideoData(v)
		// 	fmt.Println(vs)
		// 	darMap[vs.tbn+" "+vs.fps]++
		// 	if len(vs.warnings) > 0 {
		// 		fmt.Printf("%v", strings.Join(vs.warnings, "\n"))

		// 		//time.Sleep(time.Second * 2)
		// 	}
		// }

		// // for _, strData := range strDataAll {
		// // 	fmt.Println(strData)
		// // }
		// //}

	}
	//out := []string{}
	// fmt.Println("//////////////")
	for i := 0; i < 1900; i++ {
		for k, v := range metadataKeys {
			if v == i {
				fmt.Println(k, v)
			}
			//out = append(out, fmt.Sprintf("%v = %v", k, v))
		}
	}
	// sort.Strings(out)
	// fmt.Println("----------------")
	// for _, s := range out {
	// 	fmt.Println(s)
	// 	//}
	// 	//}

	// }
}

func (pi *parseInfo) String() string {
	str := ""
	str += fmt.Sprintf("name: %v\n", pi.filename)
	str += fmt.Sprintf("scanned: %v\n", pi.scanTime)
	str += fmt.Sprintf("GlobMeta: %v\n", len(pi.metadata))
	//for k, v := range pi.metadata {
	//str += fmt.Sprintf("  %v: %v\n", k, v)
	//}
	str += fmt.Sprintf("durdata: %v - %v - %v\n", pi.duration, pi.start, pi.globBitrate)
	str += fmt.Sprintf("Streams: %v\n", len(pi.streams))
	for i, s := range pi.streams {
		if len(s.metadata) > 0 {
			str += fmt.Sprintf("  stream %v has %v metadata\n", i, len(s.metadata))
		}
	}

	str += "------------\n"
	return str
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
			allData = append(allData, data)

			continue
		}
		data.data = append(data.data, line)
	}

	fmt.Println(len(allData), "examples found")
	return allData
}
