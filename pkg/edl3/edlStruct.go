package edl3

import "github.com/macroblock/imed/pkg/types"

type edlData struct {
	output []media
}

type media struct {
	video  []clip
	audio1 []clip
	audio2 []clip
	audio3 []clip
	audio4 []clip
}

type clip struct {
	id        string
	channel   string
	mixType   string
	sourceA   string
	inPointA  types.Timecode
	durationA types.Timecode
	sourceB   string
	inPointB  types.Timecode
	durationB types.Timecode
}

type sample struct {
	source   string
	inPoint  types.Timecode
	duration types.Timecode
}

func newSample(source string, inPoint, outPoint types.Timecode) (sample, error) {
	s := sample{}
	return s, nil
}

func newClip(channel, mixType string, sources ...sample) (clip, error) {
	c := clip{}
	return c, nil
}
