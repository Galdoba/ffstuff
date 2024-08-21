package action

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Galdoba/ffstuff/pkg/ump"
)

// ParseLayout - парсим код плашки на раскладку каналов и кол-во субтитров
// a62_s1 - распасится на 5.1, stereo и 1 subtitle
func ParseLayout(layout string) ([]string, int, error) {
	audio := []string{}
	parts := strings.Split(layout, "_")
	if len(parts) > 2 {
		return nil, -1, fmt.Errorf("bad layout '%v': to many delimeters", layout)
	}
	parts[0] = strings.TrimPrefix(parts[0], "a")
	if _, err := strconv.Atoi(parts[0]); err != nil {
		return nil, -1, fmt.Errorf("bad layout '%v': can't parse audio channels", layout)
	}
	chans := strings.Split(parts[0], "")
	for _, aud := range chans {
		switch aud {
		default:
			return nil, -1, fmt.Errorf("bad layout '%v': can't parse audio channels: undefined value '%v'", layout, aud)
		case "2", "6":
			audio = append(audio, aud)
		}
	}
	if len(parts) == 1 {
		return audio, 0, nil
	}
	if !strings.HasPrefix(parts[1], "s") {
		return nil, -1, fmt.Errorf("bad layout '%v': can't parse subtitle number", layout)
	}
	srtVal := strings.TrimPrefix(parts[1], "s")
	switch srtVal {
	case "x":
		return audio, 1, nil
	}
	n, err := strconv.Atoi(srtVal)
	if err != nil {
		return nil, -1, fmt.Errorf("bad layout '%v': can't parse subtitle number: undefined value '%v'", layout, srtVal)
	}
	return audio, n, nil
}

func StdDestinationDir(file string) string {
	base := filepath.Base(file)
	sep := string(filepath.Separator)
	root := os.Getenv("AGELOGOPATH")
	if strings.Contains(base, "__") {
		return root + strings.Split(base, "__")[0] + sep
	} else {
		return root + "undefined" + sep
	}
}

type GenerationTask struct {
	videoSource string
	key         string
	destination string
	audioData   []string
	frmt        []string
	srtNum      int
}

func NewTask(videoSource string, destination string, audioData []string, srtNum int) *GenerationTask {
	return &GenerationTask{videoSource: videoSource, destination: destination, audioData: audioData, srtNum: srtNum}
}

type Product struct {
	base   string
	frmt   string
	audio  []string
	srtNum int
}

func (p *Product) OutName() string {
	root := os.Getenv("AGELOGOPATH")
	sep := string(filepath.Separator)
	root += p.base + sep + p.base + "_" + p.frmt
	if len(p.audio) == 0 {
		root += ".mp4"
		return root
	}
	root += "_a"
	for _, a := range p.audio {
		switch a {
		default:
			root += "x"
		case "stereo":
			root += "2"
		case "5.1":
			root += "6"
		}
	}
	if p.srtNum == 0 {
		root += ".mp4"
		return root
	}
	return fmt.Sprintf("%v_s%v.mp4", root, p.srtNum)
}

func (p *Product) VSource() string {
	return p.base
}

func DetectSources() []string {
	root := os.Getenv("AGELOGOPATH")
	sep := string(filepath.Separator)
	sourceDir := root + "originals" + sep
	files := []string{}
	fi, err := ioutil.ReadDir(sourceDir)
	if err != nil {
		fmt.Println(err.Error())
	}
	for _, f := range fi {
		files = append(files, sourceDir+f.Name())
	}
	return files
}

type Sourcevalidation struct {
	Valid bool
	MSG   string
}

func SourceValid(source string, width, height, seconds int) Sourcevalidation {
	vld := Sourcevalidation{}
	mp := ump.NewProfile()
	if err := mp.ConsumeFile(source); err != nil {
		vld.MSG = err.Error()
		return vld
	}
	dur, err := strconv.ParseFloat(mp.Format.Duration, 64)
	if err != nil {
		vld.MSG = err.Error()
		return vld
	}
	if dur < float64(seconds) {
		vld.MSG = fmt.Sprintf("not enough duration [%v] seconds (expect at least %v)", dur, seconds)
		return vld
	}
	vidNum := 0
	for _, vid := range mp.Streams {
		if vid.Codec_type != "video" {
			continue
		}
		if vid.Width < width || vid.Height < height {
			vld.MSG = fmt.Sprintf("expected size at least %vx%v (have %vx%v)", width, height, vid.Width, vid.Height)
			return vld
		}
		vidNum++
	}
	if vidNum != 1 {
		vld.MSG = fmt.Sprintf("1 video stream expected (have %v)", vidNum)
		return vld
	}
	vld.Valid = true
	return vld
}
