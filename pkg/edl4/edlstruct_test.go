package edl4

// import (
// 	"fmt"
// 	"testing"

// 	"github.com/Galdoba/ffstuff/pkg/types"
// 	"github.com/Galdoba/utils"
// )

// func testInput() []string {
// 	return []string{
// 		// "d:\\IN\\IN_2021-09-14\\Test.edl",
// 		// "d:\\IN\\IN_2021-09-14\\Test2.edl",
// 		// "d:\\IN\\IN_2021-09-14\\test3.edl",
// 		// "d:\\IN\\IN_2021-09-14\\Test4.edl",
// 		// "d:\\IN\\IN_2021-09-15\\Simple.edl",
// 		// "d:\\IN\\IN_2021-09-15\\Simple2.edl",
// 		// "d:\\IN\\IN_2021-09-15\\Simple3.edl",
// 		// "d:\\IN\\IN_2021-09-15\\hard.edl",
// 		"d:\\IN\\IN_2021-09-16\\simple.edl",
// 	}
// }

// func TestEDLstruct(t *testing.T) {
// 	for _, file := range testInput() {
// 		parced, err := ParseFile(file)
// 		if err != nil {
// 			t.Errorf("enexpected error: %v", err.Error())
// 		}
// 		for i, val := range parced.statements {
// 			fmt.Println(i, val)
// 		}
// 		fmt.Println("/////////")
// 		cData, err := BuildClips(parced)
// 		for i, v := range cData {
// 			fmt.Println(v)
// 			if !clipDataComplete(v) {
// 				t.Errorf("File: %v\nClip %v data incomplete: %v", file, i, v)
// 			}
// 		}
// 		fmt.Println("/////////")

// 		// if len(cData) < 1 {
// 		// 	t.Errorf("enexpected len(cData) = %v", len(cData))
// 		// }
// 		//fmt.Println("/////////")
// 		//fmt.Println(parced, err)
// 	}
// }

// func InputClips() []clipdata {
// 	var clData []clipdata
// 	wipeCodes := append(validWipeCodes(), []string{"", " ", "|"}...)
// 	sources := []string{"", "file1", "file2", "BL"}
// 	testDurations := []types.Timecode{-5, 0, 5}
// 	testInpoints := []types.Timecode{-5, 0, 5}
// 	testChennels := []string{"V", "A", "A2", "3", "4", "?", "NONE"}

// 	//////////////////////////////////////////////////////////////
// 	for _, wc := range wipeCodes {
// 		for _, s1 := range sources {
// 			for _, s2 := range sources {
// 				for _, i1 := range testInpoints {
// 					for _, i2 := range testInpoints {
// 						for _, d := range testDurations {
// 							for _, c := range testChennels {
// 								clData = append(clData, clipdata{nil, wc, s1, s2, i1, i2, d, c})
// 							}
// 						}
// 					}
// 				}
// 			}
// 		}
// 	}
// 	return clData
// }

// func TestClipdata(t *testing.T) {
// 	for index, val := range InputClips() {
// 		if !clipDataComplete(val) {
// 			continue
// 		}
// 		switch {
// 		default:
// 			//t.Errorf("testObject %v - undefined, %v", index, val)
// 			//сюда проходят все неопределенные (то есть условно завершенные clipData)
// 		case val.transitionType == "":
// 			t.Errorf("val.transitionType == '', it must not -- testObject %v, %v", index, val)
// 		case val.transitionType == "C" && val.sourceB != "":
// 			t.Errorf("val.transitionType == 'C' and val.sourceB != '', it must not -- testObject %v, %v", index, val)
// 		case val.sourceA == "":
// 			t.Errorf("Cannot have absent sourceA, testObject %v, %v", index, val)
// 		case val.transitionType != "C" && val.sourceB == "":
// 			t.Errorf("val.transitionType != 'C' and val.sourceB == '', it must not -- testObject %v, %v", index, val)
// 		case (val.sourceA == val.sourceB) && (val.inPointA == val.inPointB):
// 			t.Errorf("Sources duplicated, it must not -- testObject %v, %v", index, val)
// 		case val.duration <= 0:
// 			t.Errorf("Duration cannot be <= 0, testObject %v, %v", index, val)
// 		case val.inPointA < 0:
// 			t.Errorf("inPointA cannot be <= 0, testObject %v, %v", index, val)
// 		case val.inPointB < 0:
// 			t.Errorf("inPointB cannot be <= 0, testObject %v, %v", index, val)
// 		case val.sourceB == "" && val.inPointB > 0:
// 			t.Errorf("Cannot have inPoint on absent source, testObject %v, %v", index, val)
// 		case !utils.ListContains(validWipeCodes(), val.transitionType):
// 			t.Errorf("Invalid transition Code '%v', testObject %v, %v", val.transitionType, index, val)
// 		// case val.sourceA == "BL" || val.sourceB == "BL":
// 		// 	t.Errorf("testing 'BL': testObject %v, %v", index, val)
// 		case val.sourceA == "BL" && val.inPointA != 0:
// 			t.Errorf("BL can not have inPoint != 0: testObject %v, %v", index, val)
// 		case !utils.ListContains(validChannels(), val.channel):
// 			t.Errorf("unknown channel: testObject %v, %v", index, val)
// 		}
// 	}
// }
