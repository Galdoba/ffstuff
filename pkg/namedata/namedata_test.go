package namedata

import (
	"fmt"
	"testing"
)

func emulatedNames() []string {
	return []string{
		"amerikanskaya_istoriya_prestupleniy_s03_03_2021__hd_eng20-ebur128.ac3",
		"amerikanskaya_istoriya_prestupleniy_s03_03_2021__hd_eng20-ebur128-stereo.ac3",
		"amerikanskaya_istoriya_prestupleniy_s03_03_2021_hd_eng20-ebur128.ac3",
	}
}

func TestEburTrimmer(t *testing.T) {
	for _, oldName := range emulatedNames() {
		newName, err := TrimLoudnormPrefix(oldName)
		if err != nil {
			t.Errorf("Error: %v", err.Error())
		}
		fmt.Println(oldName)
		fmt.Println(newName)
		fmt.Println("")
	}
}
