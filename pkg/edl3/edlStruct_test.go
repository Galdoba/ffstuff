package edl3

import (
	"testing"

	"github.com/Galdoba/ffstuff/pkg/types"
)

type sampleData struct {
	s   string
	in  types.Timecode
	dur types.Timecode
}

func sampleDataList() [][]string {
	return [][]string{
		{"name.mp4", "00:00:00:00", "00:00:00:00"},
	}
}

func TestSample(t *testing.T) {
	///CREATION
	expectedSource := "clip.mp4"
	expectedIN := 123
	expectedDUR := 456
	smpl, err := newSample(expectedSource, types.Timecode(expectedIN), types.Timecode(expectedDUR))
	if err != nil {
		t.Errorf("newSample(string, time, time) func error: %v", err.Error())
	}
	if smpl.source != expectedSource {
		t.Errorf("Cannot set sample name (have '%v'; expect '%v'", smpl.source, expectedSource)
	}
	if smpl.source != expectedSource {
		t.Errorf("Cannot set sample InPoint (have '%v'; expect '%v'", smpl.inPoint, expectedIN)
	}
	if smpl.source != expectedSource {
		t.Errorf("Cannot set sample Duration (have '%v'; expect '%v'", smpl.duration, expectedDUR)
	}
	//EXAMPLES BY DATA
	for _, data := range sampleDataList() {
		inPoint, errIN := types.ParseTimecode(data[1])
		if errIN != nil {
			t.Errorf("errIN = '%v', expect nil from %v", errIN.Error(), data[1])
		}
		outPoint, errOUT := types.ParseTimecode(data[2])
		if errOUT != nil {
			t.Errorf("errOUT = '%v', expect nil from %v", errOUT.Error(), data[2])
		}
		s, err := newSample(data[0], inPoint, outPoint)
		// if !listContains(validWipeCodes(), s.source) {
		// 	t.Errorf("test content error '%v' '%v' '%v', expect WipeCode and two Timecodes", s.source, s.inPoint, s.duration)
		// }
		if err != nil {
			t.Errorf("clip err = '%v', expect nil from %v", err.Error(), s)
		}
	}
}

func TestClip(t *testing.T) {

}

// func validChannels() []string {
// 	return []string{"V", "A", "A2", "A3", "A4", "NONE"}
// }
