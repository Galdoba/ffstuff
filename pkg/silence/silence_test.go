package silence

import (
	"fmt"
	"testing"

	"github.com/Galdoba/ffstuff/pkg/info"
)

func TestDetect(t *testing.T) {
	for _, path := range testPaths() {
		for _, dur := range testDurat() {
			for _, ld := range testLoudness() {
				//fmt.Printf("start test %v\npath: %v\ndur %v\nloud %v\n", n, path, dur, ld)
				durat, err := info.Duration(path)
				fmt.Println(durat, err)
				si, err := Detect(path, ld, dur, true)
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
		"d:\\MUX\\tests\\Ryad_19_TRL_AUDIORUS51.m4a",
		"d:\\MUX\\tests\\Shang-Chi_and_the_Legend_of_the_Ten_Rings_AUDIOENG51.m4a",
		"d:\\MUX\\tests\\Shang-Chi_and_the_Legend_of_the_Ten_Rings_AUDIORUS51.m4a",
	}
}

func testLoudness() []float64 {
	return []float64{
		//LOUDNESS_BELOW_36db,
		//LOUDNESS_BELOW_48db,
		//LOUDNESS_BELOW_60db,
		LOUDNESS_BELOW_72db,
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
		2,
		//1.5,
	}
}
