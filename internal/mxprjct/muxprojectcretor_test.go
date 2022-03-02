package mxprjct

import (
	"fmt"
	"testing"
)

func pathInput() [][]string {
	return [][]string{
		{"d:\\IN\\IN_testInput\\Agnes_AUDIORUS51.m4a", "d:\\IN\\IN_testInput\\Agnes_HD.mp4}"},
		{"d:\\IN\\IN_testInput\\agnes_2015__hd.mp4", "d:\\IN\\IN_testInput\\agnes_2015__hd_rus51.m4a"},
		{"d:\\IN\\IN_testInput\\agnes_2015__hd.mp4", "d:\\IN\\IN_testInput\\agnes_2015__hd_eng20.m4a"},
		{"d:\\IN\\IN_testInput\\agnes_2015__hd.mp4", "d:\\IN\\IN_testInput\\agnes_2015__hd_eng51.m4a", "d:\\IN\\IN_testInput\\agnes_2015__hd_rus51.m4a"},
		{"d:\\IN\\IN_testInput\\issues.txt"},
		{"d:\\IN\\IN_testInput\\issues - Copy.srt"},
	}
}

func Test_CreateMuxProjectStruct(t *testing.T) {
	for _, inputSet := range pathInput() {
		fmt.Printf("start input set: %v\n \n", inputSet)
		mp, err := Create(inputSet)
		if err != nil {
			t.Errorf("creation error: %v", err.Error())
			continue
		}
		if mp.inputV == "" {
			t.Errorf("can not have project with absent video")
			continue
		}
		if len(mp.inputA) == 0 {
			t.Errorf("can not have project with absent audio")
			continue
		}
		if mp.expectedResult == "" {
			t.Errorf("expected result not calculated")
			continue
		}
		fmt.Printf("\nend input set --------------\n")
	}
}
