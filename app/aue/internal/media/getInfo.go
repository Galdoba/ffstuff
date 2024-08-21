package media

import "github.com/Galdoba/ffstuff/pkg/ump"

const (
	typeVideo    = "video"
	typeAudio    = "audio"
	typeSubtitle = "subtitle"
)

func VideoStreams(prf *ump.MediaProfile) []*ump.Stream {
	video := []*ump.Stream{}
	for _, stream := range prf.Streams {
		if stream.Codec_type == typeVideo {
			video = append(video, stream)
		}
	}
	return video
}

func AudioStreams(prf *ump.MediaProfile) []*ump.Stream {
	audio := []*ump.Stream{}
	for _, stream := range prf.Streams {
		if stream.Codec_type == typeAudio {
			audio = append(audio, stream)
		}
	}
	return audio
}

func SubtitleStreams(prf *ump.MediaProfile) []*ump.Stream {
	audio := []*ump.Stream{}
	for _, stream := range prf.Streams {
		if stream.Codec_type == typeSubtitle {
			audio = append(audio, stream)
		}
	}
	return audio
}

func SteamNumbers(prf *ump.MediaProfile) (int, int, int) {
	numsByType := make(map[string]int)
	for _, stream := range prf.Streams {
		numsByType[stream.Codec_type]++
	}
	return numsByType[typeVideo], numsByType[typeAudio], numsByType[typeSubtitle]
}
