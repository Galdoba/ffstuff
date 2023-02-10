package inputinfo

import (
	"fmt"
	"testing"
)

func TestInputReading(t *testing.T) {

	remember := []videostream{}
	examples := gatherInfo()
	darMap := make(map[string]int)
	//out := []string{}
	metadataKeys := make(map[string]int)
	total := 0
	for i, input := range examples {
		total = i
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
				remember = append(remember, stream)
			}
			// if i == 1400 {
			// 	fmt.Println(pi)
			// 	fmt.Println(pi.video)
			// 	fmt.Println(pi.audio)
			// }
			if len(pi.warnings) > 0 {
				//fmt.Println(pi.filename)
				//for _, wrn := range pi.warnings {
				//	fmt.Print("" + wrn + "\n")
				//}
				//fmt.Println("")
				darMap[fmt.Sprintf("have %v warnings", len(pi.warnings))]++
			} else {
				darMap["possible for automation"]++
				//fmt.Println(pi.String())
			}
		}
		//fmt.Println("----------")
		//fmt.Println(metadataKeys)
		//fmt.Println("----------")

	}
	fmt.Println("cycles", total)
	fmt.Println("===============")
	// for _, data := range remember {
	// 	switch {
	// 	default:
	// 		//fmt.Println(data.sardar)
	// 		// case strings.HasPrefix(data.codecinfo, "prores"):
	// 		darMap[strings.TrimSpace(data.sardar)]++
	// 		darMap["Total_Files"]++
	// 		// case strings.HasPrefix(data.codecinfo, "mpeg2video"):
	// 		// 	darMap["mpeg2video"]++
	// 	}

	// 	// }

	// }
	for i := 0; i < 5000; i++ {
		for k, v := range darMap {
			if v == i {
				fmt.Println(k, v, "--")
			}
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

func TestParseFile(t *testing.T) {
	return
	file := `\\nas\buffer\IN\_DONE\Gorodskie_legendy_s04e02_PRT230208141911_SER_00385_18.mp4`
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
	//fmt.Println(pi.video)
	fmt.Println("data		:", pi.video[0].data)
	//parseVideoLine(pi.video[0].data)
	//fmt.Println("codecinfo	:", pi.video[0].codecinfo)
	//fmt.Println("pix_fmt	:", pi.video[0].pix_fmt)
	//fmt.Println("width		:", pi.video[0].width)
	//fmt.Println("height		:", pi.video[0].height)
	//fmt.Println("bitrate	:", pi.video[0].bitrate)
	//fmt.Println("sardar		:", pi.video[0].sardar)
	//fmt.Println("fps		:", pi.video[0].fps)
	//fmt.Println("tbr		:", pi.video[0].tbr)
	//fmt.Println("tbn		:", pi.video[0].tbn)
	//fmt.Println("tbc		:", pi.video[0].tbc)
	//fmt.Println("lang		:", pi.video[0].lang)
	//fmt.Println("metadata	:", pi.video[0].metadata)
	//fmt.Println("sidedata	:", pi.video[0].sidedata)
	//fmt.Println("warnings	:", pi.video[0].warnings)
	fmt.Println("----------")
	fmt.Println(pi.audio)
	fmt.Println("----------")
	fmt.Println(pi.data)
	fmt.Println("----------")
	fmt.Println(pi.subtitles)
}

func TestDuplicates(t *testing.T) {
	return
	CleanScanData()
	//GenerateNewDataFile(list)
}
