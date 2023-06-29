package inputinfo

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Galdoba/devtools"
)

func TestInputReading(t *testing.T) {
	//return
	remember := []Audiostream{}
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
			remember = append(remember, pi.Audio...)
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
				//				darMap[fmt.Sprintf("have %v warnings", len(pi.warnings))]++
			} else {
				//				darMap["possible for automation"]++
				//fmt.Println(pi.String())
			}
		}
		//fmt.Println("----------")
		//fmt.Println(metadataKeys)
		//fmt.Println("----------")

	}
	fmt.Println("cycles", total)
	fmt.Println("===============")
	for _, data := range remember {
		codec := devtools.AliasByPrefixFromFile(`c:\Users\pemaltynov\.ffstuff\data\alias_AudioCodec`, data.codec)
		// switch {
		// default:
		// 	//fmt.Println(data.sardar)
		// case strings.HasPrefix(data.codec, "aac"):
		// 	codec = "aac"

		// case strings.HasPrefix(data.codec, "ac3"):
		// 	codec = "ac3"
		// case strings.HasPrefix(data.codec, "pcm_s24le"):
		// 	codec = "pcm_s24le"
		// case strings.HasPrefix(data.codec, "pcm_s16le"):
		// 	codec = "pcm_s16le"
		// }
		darMap[codec]++
		darMap["Total_Files"]++
		// }

	}
	for i := 0; i < 15000; i++ {
		for k, v := range darMap {
			if v == i {
				fmt.Println(k, v, "--")
			}
		}
	}
	fmt.Println("Stop")

}

func hasAnyPrefix(str string, preflist []string) string {
	for _, prefix := range preflist {
		if strings.HasPrefix(str, prefix) {
			return prefix
		}
	}
	return str
}

func TestParseFile(t *testing.T) {
	return
	file := `\\192.168.31.4\buffer\IN\_REJECTED\Пустая_Церковь_Ростелеком_R2.mp4`
	pi, err := ParseFile(file)
	fmt.Println(err)
	if err != nil {
		fmt.Println(err.Error())
	}
	for _, w := range pi.warnings {
		fmt.Println(w)
	}

}

func TestDuplicates(t *testing.T) {
	return
	CleanScanData()
	//GenerateNewDataFile(list)
}
