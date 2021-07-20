package edl3

import (
	"testing"

	"github.com/macroblock/imed/pkg/types"
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
	for _, data := range sampleDataList() {
		inPoint, _ := types.ParseTimecode(data[1])
		outPoint, _ := types.ParseTimecode(data[2])

		s, err := newSample(data[0], inPoint, outPoint)
		if !listContains(validWipeCodes(), s.source) {
			t.Errorf("test content error %v %v %v", s.source, s.inPoint, s.duration)
		}
		if err != nil {
			t.Errorf("clip err = '%v', expect nil from %v", err.Error(), s)
		}
	}
}

func TestClip(t *testing.T) {

	cl, err := newClip("channel", "mix", sample{})

	if !listContains(validChannels(), cl.channel) {
		t.Errorf("clip channel = '%v', expect V, A, A2, A3 or A4", cl.channel)
	}
	if !listContains(validWipeCodes(), cl.channel) {
		t.Errorf("clip mixType = '%v', expect C, D or Wxxx", cl.channel)
	}
	if err != nil {
		t.Errorf("clip err = '%v', expect nil", err.Error())
	}

}

func validChannels() []string {
	return []string{"V", "A", "A2", "A3", "A4", "NONE"}
}
