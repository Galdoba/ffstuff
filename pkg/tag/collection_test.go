package tag

import (
	"fmt"
	"testing"
)

var inputsIN = []string{
	//`Zlo--s04e13--SER--HSUB--Zlo_s04e13_PRT240816080000_SER_05005_18.mp4`,
	`Zlo--SER--s04e13--PRT240816080000--HD--Zlo_s04e13_PRT240816080000_SER_05005_18.mp4`,
	`Semion--FILM--4K--durackoe_imya.mp4`,
}

var inputsUnknown = []string{
	`Sem_smertnyh_grehov_s02e02_aad asdfSD.mp4`,
	`Sem_smertnyh_grehov_s02_02_SER515.mp4`,
	`Vosem_smertnyh_grehov_s02_SER515_trl.mp4`,
}

func TestParsing(t *testing.T) {
	for _, input := range inputsIN {
		col, err := NewCollection(input)
		if err != nil {
			fmt.Println("ERROR:   ", err)
		}
		fmt.Println(col)
		fmt.Println(col.FileName)
		fmt.Println(ProjectOutfileName(col, New(EXT, "mp4")))
		fmt.Println(ProjectOutfileName(col, New(AUDIO, "AUDIORUS51"), New(EXT, "m4a")))
		fmt.Println(ProjectOutfileName(col, New(AUDIO, "AUDIOENG20"), New(EXT, "m4a")))
		fmt.Println(ProjectOutfileName(col, New(EXT, "srt")))

	}
	fmt.Println("/////////////////")
	for _, input := range inputsUnknown {
		fmt.Println("----")
		col, err := NewCollection(input)
		if err != nil {
			fmt.Println("ERROR:   ", err)
			continue
		}
		fmt.Println(col)
		fmt.Println(col.FileName)
		fmt.Println(ProjectInfileName(col))

	}
}

// type collectionResult struct {
// 	assighnment map[string][]int
// }

// func (cr *collectionResult) Analize(input string, c *collection) string {
// 	last := 0
// 	detections := make(map[int]int)
// 	for _, slice := range cr.assighnment {
// 		for _, i := range slice {
// 			if last < i {
// 				last = i
// 			}
// 			detections[i]++
// 		}
// 	}
// 	s := ""
// 	for i := 0; i < last; i++ {
// 		s += fmt.Sprintf("%v", detections[i])
// 	}
// 	return s
// }

// func findPositions(str, substr string) []int {
// 	parts := strings.Split(str, substr)
// 	switch len(parts) {
// 	default:
// 		return nil
// 	case 2:
// 	}
// 	offset := len(parts[0])
// 	//leng := len(substr)
// 	res := []int{}
// 	for i, v := range strings.Split(substr, "") {
// 		if v == strings.Split(str, "")[offset+i] {
// 			fmt.Println(v, offset+i)
// 			res = append(res, offset+i)
// 		}
// 	}
// 	return res
// }
