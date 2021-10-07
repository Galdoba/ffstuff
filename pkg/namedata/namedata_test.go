package namedata

import (
	"fmt"
	"strings"
	"testing"
)

type testData struct {
	input          string
	expectedOutput string
}

func emulatedNames0() []testData {
	return []testData{
		{"some_name_xx21__hd_eng20-ebur128.ac3", "some_name_xx21__hd_eng20.ac3"},
		{"some_name_xx21__hd_eng51-ebur128.ac3", "some_name_xx21__hd_eng51.ac3"},
		{"some_name_xx21__hd_eng51-ebur128-stereo.ac3", "some_name_xx21__sd_eng20.ac3"},
		{"some_name_xx21__hd_eng20-ebur128-stereo.ac3", "some_name_xx21__sd_eng20.ac3"},
		{"some_name_xx21__hd_eng20_ebur128_stereo.ac3", "some_name_xx21__hd_eng20.ac3"},
		{"some_name_xx21_hd_eng20-ebur128.ac3", "some_name_xx21__hd_eng20.ac3"},
	}
}

func emulatedNames() []string {
	var emulatedNames []string
	base := []string{"film_name_0000", "serial_name_s00_00_0000"}
	vidTag := []string{"sd", "hd", "4k"}
	audTag := []string{"rus20", "rus51", "eng20", "eng51", "qqq20", "qqq51"}
	eburTag := []string{"-ebur128", "-ebur128-stereo", ""}
	resolutions := []string{".ac3", ".aac", ".m4a", ".txt", ".mp4"}
	for _, name := range base {
		for _, vt := range vidTag {
			for _, at := range audTag {
				for _, et := range eburTag {
					for _, rt := range resolutions {
						emulatedNames = append(emulatedNames, name+"__"+vt+"_"+at+et+rt)
					}
				}
			}
		}
	}
	emulatedNames = append(emulatedNames, "film_name_0000__sd.mp4")
	return emulatedNames
}

func TestEburTrimmer(t *testing.T) {
	undecided := 0
	for i, oldName := range emulatedNames() {
		newName, err := TrimLoudnormPrefix(oldName)

		if err != nil {
			if strings.Contains(err.Error(), "invalid name [") {
				continue
			}
			t.Errorf("Test %v:	Error: %v", i, err.Error())
			undecided++
		}
		if newName != oldName {
			t.Errorf("oldname > newname (%v > %v)", oldName, newName)
		}

	}
	fmt.Println("Undecided:", undecided)
}
