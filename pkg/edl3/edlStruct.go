package edl3

import "github.com/Galdoba/ffstuff/pkg/types"

type edlData struct {
	data []statementData
}

type clip struct {
	id      string
	channel string
	mixType string
	sourceA sample
	sourceB sample
}

type sample struct {
	source   string
	inPoint  types.Timecode
	duration types.Timecode
}

func newSample(source string, inPoint, outPoint types.Timecode) (sample, error) {
	s := sample{}
	s.source = source
	s.inPoint = inPoint
	s.duration = outPoint - inPoint
	return s, nil
}

func newClip(channel, mixType string, sources ...sample) (clip, error) {
	c := clip{}
	c.channel = channel
	c.mixType = mixType
	for i, src := range sources {
		switch i {
		default:
			break
		case 0:
			c.sourceA = src
		case 1:
			c.sourceB = src
		}
	}
	return c, nil
}
