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

func New(path string, purpose string) (SourceFile, error) {
	sf := SourceFile{}
	sf.path = path
	sf.name = filepath.Base(path)
	sf.profile = ump.NewProfile()
	if err := sf.profile.ConsumeFile(path); err != nil {
		return sf, fmt.Errorf("profiling failed: %v", err)
	}
	if err := validate(sf); err != nil {
		return sf, fmt.Errorf("source invalid: %v", err)
	}
	return sf, nil
}

func validate(sf SourceFile) error {
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
