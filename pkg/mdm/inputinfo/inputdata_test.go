package inputinfo

import (
	"fmt"
	"testing"

	"github.com/Galdoba/utils"
)

func TestInputReading(t *testing.T) {

	remember := []string{}
	examples := gatherInfo()
	//darMap := make(map[string]int)
	//out := []string{}
	metadataKeys := make(map[string]int)
	for i, input := range examples {

		if i < 1800 {
			continue
		}
		pi, err := parse(input)
		if err != nil {
			t.Errorf(err.Error())
		}
		if pi.filename != "" {
			//fmt.Printf("example %v: \n", i)

			//fmt.Printf("%v", pi.String())

			//for k, _ := range pi.metadata {
			metadataKeys[pi.comment]++
			if pi.comment == "---" {
				panic(pi.filename)
			}
			for _, stream := range pi.video {
				remember = append(remember, stream.data)
			}
			// if i == 1400 {
			// 	fmt.Println(pi)
			// 	fmt.Println(pi.video)
			// 	fmt.Println(pi.audio)
			// }

		}
		//fmt.Println("----------")
		//fmt.Println(metadataKeys)
		//fmt.Println("----------")
		for _, data := range remember {
			fmt.Print(data)

		}
	}
	fmt.Println("Stop")
}

func (pi *parseInfo) String() string {

	str := ""
	str += fmt.Sprintf("name: %v\n", pi.filename)
	str += fmt.Sprintf("scanned: %v\n", pi.scanTime)
	str += fmt.Sprintf("GlobMeta: %v\n", len(pi.metadata))
	for k, v := range pi.metadata {
		str += fmt.Sprintf("  %v: %v\n", k, v)
	}
	str += fmt.Sprintf("durdata: %v - %v - %v\n", pi.duration, pi.start, pi.globBitrate)
	str += fmt.Sprintf("Streams: %v\n", len(pi.streams))
	str += "Video:\n"
	for _, v := range pi.video {
		fmt.Sprintf("%v++\n", v.fps)
	}
	// for i, s := range pi.streams {
	// 	str += s.data + "||"
	// 	if len(s.metadata) > 0 {
	// 		str += fmt.Sprintf(" | stream %v has %v metadata\n", i, len(s.metadata))
	// 		for k, v := range s.metadata {
	// 			str += fmt.Sprintf("    %v: %v\n", k, v)
	// 		}
	// 	}
	// 	//for k, v := range s.metadata {
	// 	//str += fmt.Sprintf("%v|---|%v\n", k, v)
	// 	//}
	// }
	// for _, s := range pi.video {
	// 	str += s.codecinfo + "\n"
	// 	str += s.pix_fmt + "\n"
	// 	str += s.sardar + "\n"
	// 	str += s.fps + "\n"

	// }

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
		data.data = append(data.data, line+"\n")
	}

	fmt.Println(len(allData), "examples found")
	return allData
}

func TestParseFile(t *testing.T) {
	return
	file := `\\nas\buffer\IN\_DONE\Poema_o_shanyan_s01e10_PRT230119004502_SER_12536.mp4`
	pi, err := ParseFile(file)
	fmt.Println(pi)
	fmt.Println("----------")
	//fmt.Println(pi.String())
	fmt.Println("----------")
	fmt.Println(err)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("----------")

}
