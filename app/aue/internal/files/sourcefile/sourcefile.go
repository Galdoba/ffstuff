package sourcefile

import (
	"fmt"
	"path/filepath"

	"github.com/Galdoba/ffstuff/app/aue/internal/define"
	"github.com/Galdoba/ffstuff/app/aue/internal/media"
	"github.com/Galdoba/ffstuff/pkg/ump"
)

type SourceFile struct {
	path    string
	name    string
	purpose string
	profile *ump.MediaProfile
}

func New(path string, purpose string) *SourceFile {
	sf := SourceFile{}
	sf.path = path
	sf.name = filepath.Base(path)
	sf.purpose = purpose
	return &sf
}

func (sf *SourceFile) FillProfile(prf *ump.MediaProfile) {
	sf.profile = prf
}

func (sf *SourceFile) Validate() error {
	if sf.profile == nil {
		return fmt.Errorf("no media profile")
	}
	vNum, aNum, sNum := media.SteamNumbers(sf.profile)
	switch sf.purpose {
	default:
		return fmt.Errorf("purpose '%v' out of scope of validation func or not implemented", sf.purpose)
	case define.PURPOSE_Input_Media:
		if vNum+aNum == 0 {
			return fmt.Errorf("can't fill purpose '%v': no media streams", sf.purpose)
		}
	case define.PURPOSE_Input_Subs, define.PURPOSE_Input_Hardsub:
		if sNum < 1 {
			return fmt.Errorf("can't fill purpose '%v': no subtitle streams", sf.purpose)
		}
	}
	return nil
}

func (sf *SourceFile) Name() string {
	return sf.name
}

func (sf *SourceFile) Path() string {
	return sf.path
}

func (sf *SourceFile) Purpose() string {
	return sf.purpose
}

func (sf *SourceFile) Profile() *ump.MediaProfile {
	return sf.profile
}

func (sf *SourceFile) FPS() string {
	videoStreams := media.VideoStreams(sf.profile)
	if len(videoStreams) != 1 {
		return ""
	}
	return videoStreams[0].R_frame_rate
}

func MapStreamTypesAll(sources []*SourceFile) map[string]int {
	stMap := make(map[string]int)
	for _, source := range sources {
		profile := source.profile
		for _, stream := range profile.Streams {
			stMap[stream.Codec_type]++
		}
	}
	return stMap
}

func MapByStreamTypes(sources []*SourceFile) map[string][]int {
	stMap := make(map[string][]int)
	for _, source := range sources {
		stMap[source.name] = []int{0, 0, 0}
		profile := source.profile
		for _, stream := range profile.Streams {
			switch stream.Codec_type {
			case define.STREAM_VIDEO:
				stMap[source.name][0]++
			case define.STREAM_AUDIO:
				stMap[source.name][1]++
			case define.STREAM_SUBTITLE:
				stMap[source.name][2]++
			}
		}
	}
	return stMap
}

func Names(sources []*SourceFile) []string {
	names := []string{}
	for _, source := range sources {
		names = append(names, source.name)
	}
	return names
}

func (src *SourceFile) Details() string {
	str := ""
	str += fmt.Sprintf("  path    : %v\n", src.path)
	str += fmt.Sprintf("  name    : %v\n", src.name)
	str += fmt.Sprintf("  purpose : %v\n", src.purpose)
	str += fmt.Sprintf("  profile : %v\n", src.profile)
	return str
}
