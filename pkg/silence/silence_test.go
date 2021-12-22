package silence

import (
	"fmt"
	"testing"
)

func TestDetect(t *testing.T) {
	for _, path := range testPaths() {
		for _, dur := range testDurat() {
			for _, ld := range testLoudness() {
				//fmt.Printf("start test %v\npath: %v\ndur %v\nloud %v\n", n, path, dur, ld)
				si, err := Detect(path, ld, dur)
				if err != nil {
					t.Errorf("Detect returned error:\nfile = '%v'\nerror = '%v'\n", si, err)
				}
				if si != nil {
					fmt.Println(si)
					fmt.Println(si.Timings())
				}

				fmt.Println(" ")
			}
		}
	}
}

func testPaths() []string {
	return []string{
		"d:\\IN\\IN_2021-12-21\\The_Big_Bang_Theory_s01e01_AUDIO20_FOR_SYNC.m4a",
		"d:\\IN\\IN_2021-12-21\\The_Big_Bang_Theory_s01e01_KURAZH20.m4a",
		"d:\\IN\\IN_2021-12-21\\The_Big_Bang_Theory_s01e02_AUDIO20_FOR_SYNC.m4a",
		"d:\\IN\\IN_2021-12-21\\The_Big_Bang_Theory_s01e02_KURAZH20.m4a",
		"d:\\IN\\IN_2021-12-21\\The_Big_Bang_Theory_s01e03_AUDIO20_FOR_SYNC.m4a",
		"d:\\IN\\IN_2021-12-21\\The_Big_Bang_Theory_s01e03_KURAZH20.m4a",
		"d:\\IN\\IN_2021-12-21\\The_Big_Bang_Theory_s01e04_AUDIO20_FOR_SYNC.m4a",
		"d:\\IN\\IN_2021-12-21\\The_Big_Bang_Theory_s01e04_KURAZH20.m4a",
		"d:\\IN\\IN_2021-12-21\\The_Big_Bang_Theory_s01e05_AUDIO20_FOR_SYNC.m4a",
		"d:\\IN\\IN_2021-12-21\\The_Big_Bang_Theory_s01e05_KURAZH20.m4a",
		"d:\\IN\\IN_2021-12-21\\The_Big_Bang_Theory_s01e06_AUDIO20_FOR_SYNC.m4a",
		"d:\\IN\\IN_2021-12-21\\The_Big_Bang_Theory_s01e06_KURAZH20.m4a",
		"d:\\IN\\IN_2021-12-21\\The_Big_Bang_Theory_s01e07_AUDIO20_FOR_SYNC.m4a",
		"d:\\IN\\IN_2021-12-21\\The_Big_Bang_Theory_s01e07_KURAZH20.m4a",
		"d:\\IN\\IN_2021-12-21\\The_Big_Bang_Theory_s01e08_AUDIO20_FOR_SYNC.m4a",
		"d:\\IN\\IN_2021-12-21\\The_Big_Bang_Theory_s01e08_KURAZH20.m4a",
		"d:\\IN\\IN_2021-12-21\\The_Big_Bang_Theory_s01e09_AUDIO20_FOR_SYNC.m4a",
		"d:\\IN\\IN_2021-12-21\\The_Big_Bang_Theory_s01e09_KURAZH20.m4a",
		"d:\\IN\\IN_2021-12-21\\The_Big_Bang_Theory_s01e10_AUDIO20_FOR_SYNC.m4a",
		"d:\\IN\\IN_2021-12-21\\The_Big_Bang_Theory_s01e10_KURAZH20.m4a",
		"d:\\IN\\IN_2021-12-21\\The_Big_Bang_Theory_s01e11_AUDIO20_FOR_SYNC.m4a",
		"d:\\IN\\IN_2021-12-21\\The_Big_Bang_Theory_s01e11_KURAZH20.m4a",
		"d:\\IN\\IN_2021-12-21\\The_Big_Bang_Theory_s01e12_AUDIO20_FOR_SYNC.m4a",
		"d:\\IN\\IN_2021-12-21\\The_Big_Bang_Theory_s01e12_KURAZH20.m4a",
		"d:\\IN\\IN_2021-12-21\\The_Big_Bang_Theory_s01e13_AUDIO20_FOR_SYNC.m4a",
		"d:\\IN\\IN_2021-12-21\\The_Big_Bang_Theory_s01e13_KURAZH20.m4a",
		"d:\\IN\\IN_2021-12-21\\The_Big_Bang_Theory_s01e14_AUDIO20_FOR_SYNC.m4a",
		"d:\\IN\\IN_2021-12-21\\The_Big_Bang_Theory_s01e14_KURAZH20.m4a",
		"d:\\IN\\IN_2021-12-21\\The_Big_Bang_Theory_s01e15_AUDIO20_FOR_SYNC.m4a",
		"d:\\IN\\IN_2021-12-21\\The_Big_Bang_Theory_s01e15_KURAZH20.m4a",
		"d:\\IN\\IN_2021-12-21\\The_Big_Bang_Theory_s01e16_AUDIO20_FOR_SYNC.m4a",
		"d:\\IN\\IN_2021-12-21\\The_Big_Bang_Theory_s01e16_KURAZH20.m4a",
		"d:\\IN\\IN_2021-12-21\\The_Big_Bang_Theory_s01e17_AUDIO20_FOR_SYNC.m4a",
		"d:\\IN\\IN_2021-12-21\\The_Big_Bang_Theory_s01e17_KURAZH20.m4a",
	}
}

func testLoudness() []float64 {
	return []float64{
		//LOUDNESS_BELOW_36db,
		//LOUDNESS_BELOW_48db,
		LOUDNESS_BELOW_60db,
		//LOUDNESS_BELOW_72db,
		// LOUDNESS_BELOW_84db,
		// LOUDNESS_BELOW_96db,
		// LOUDNESS_BELOW_108db,
	}
}

func testDurat() []float64 {
	return []float64{
		// -0.5,
		// 0,
		// 0.5,
		1,
		//1.5,
	}
}
