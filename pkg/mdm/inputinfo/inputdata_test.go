package inputinfo

import (
	"fmt"
	"sort"
	"testing"

	"github.com/Galdoba/utils"
)

func TestInputReading(t *testing.T) {
	examples := gatherInfo()
	darMap := make(map[string]int)
	out := []string{}
	for i, input := range examples {
		fmt.Printf("example %v: ", i)
		pi, err := parse(input)
		if err != nil {
			t.Errorf(err.Error())
		}
		fmt.Printf("%v                  \n", pi)
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
	fmt.Println("//////////////")
	//for i := 0; i < 1500; i++ {
	for k, v := range darMap {
		//		if v == i {
		//			fmt.Println(k, v)
		//		}
		out = append(out, fmt.Sprintf("%v = %v", k, v))
	}
	sort.Strings(out)
	//for _, s := range out {
	//fmt.Println(s)
	//}
	//}

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
