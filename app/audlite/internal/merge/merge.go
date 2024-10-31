package merge

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/Galdoba/ffstuff/pkg/ump"
)

const (
	merge2 = "2 mono ==> streo"
	merge6 = "6 mono ==> 5.1"
)

/*
FC_AUD6="amerge=inputs=6,aresample=48000,atempo=25/(__)"
fflite -i [file1] -i [file2] -i [file3] -i [file4] -i [file5] -i [file6] \
 -filter_complex "[0:a:0][1:a:0][2:a:0][3:a:0][4:a:0][5:a:0],${FC_AUD6}[audio_out]" \
 -map "[audio_out]" -c:a alac -compression_level 0 -map_metadata -1 -map_chapters -1 outfile.m4a \

merge		: "amerge=inputs=6"
aresample	: "aresample=48000"
timeFactor  : "%v/%v" [curentDuration, tgDuration]
atempo		: "atempo=25/(%v)" [timeFactror]

fc: `"[0:a:0][1:a:0][2:a:0][3:a:0][4:a:0][5:a:0]%v[audio_out]"` [merge,aresample,timeFactor]

command:
`ffmpeg -i %v -i %v -i %v -i %v -i %v -i %v -filter_complex "%v" -map "[audio_out]" -c:a alac -compression_level 0 -map_metadata -1 -map_chapters -1 %v`




Input #0, wav, from 'Dom_C.wav':
  Duration: 01:31:23.14, bitrate: 1152 kb/s
    Stream #0:0: Audio: pcm_s24le ([1][0][0][0] / 0x0001), 48000 Hz, mono, s32 (24 bit), 1152 kb/s
	 01:31:23.14   ==> 01:27:38.58
	 5,483.14      ==> 5,258.58
	 25/23.976     ==> 1.0427
	 548314/525858 ==> 1.0427

ffmpeg -i Dom_L.wav -i Dom_R.wav -i Dom_C.wav -i Dom_Lfe.wav -i Dom_Ls.wav -i Dom_Rs.wav -filter_complex "[0:a:0][1:a:0][2:a:0][3:a:0][4:a:0][5:a:0]amerge=inputs=6,aresample=48000,atempo=25/(24000/1001)[audio]" -map "[audio]" -c:a alac -compression_level 0 -map_metadata -1 -map_chapters -1 outFile_2397.m4a

ffmpeg -i Dom_L.wav -i Dom_R.wav -i Dom_C.wav -i Dom_Lfe.wav -i Dom_Ls.wav -i Dom_Rs.wav -filter_complex "[0:a:0][1:a:0][2:a:0][3:a:0][4:a:0][5:a:0]amerge=inputs=6,aresample=48000,atempo=(548314/525858)[audio]" -map "[audio]" -c:a alac -compression_level 0 -map_metadata -1 -map_chapters -1 outFile_2397_dur.m4a


*/

type mergeProc struct {
	descr          string
	source         map[string]*ump.MediaProfile
	sourcesOrdered []string
	srcFps         float64
	srcDuration    float64
	tgDuration     float64
	destination    string
	target         string
}

func NewMerge(sources ...string) (*mergeProc, error) {
	mp := mergeProc{}
	mp.source = make(map[string]*ump.MediaProfile)
	durAverage := 0.0
	for _, path := range sources {
		mpR := ump.NewProfile()
		if err := mpR.ConsumeFile(path); err != nil {
			return nil, fmt.Errorf("source injection failed: %v", err)
		}
		mp.source[path] = mpR
		dur, err := strconv.ParseFloat(mpR.Format.Duration, 64)
		if err != nil {
			return nil, fmt.Errorf("duration parsing failed (%v): %v", mpR.Format.Duration, err)
		}
		durAverage += dur
	}
	if err := mp.defineOrder(); err != nil {
		return nil, err
	}
	durAverage = durAverage / float64(len(mp.source))
	mp.srcDuration = durAverage
	return &mp, nil
}

func (mp *mergeProc) SetOriginalFPS(fps float64) {
	mp.srcFps = fps
}

func (mp *mergeProc) SetOriginalDuration(duration string) error {
	dur, err := durationToFl64(duration)
	if err != nil {
		return err
	}
	mp.srcDuration = dur
	return nil
}

func (mp *mergeProc) SetTargetName(name string) error {
	mp.target = name
	return nil
}

func (mp *mergeProc) SetDestination(name string) error {
	fi, err := os.Stat(name)
	if err != nil {
		return err
	}
	if !fi.IsDir() {
		return fmt.Errorf("destination is not directory")
	}
	mp.destination = name
	return nil
}

func durationToFl64(duration string) (float64, error) {
	reg := regexp.MustCompile(`(\d\d\:)?(\d\d\:)?(\d*(\.*\d*)?)`)
	fl := reg.FindString(duration)
	parts := strings.Split(fl, ":")
	if len(parts) == 1 {
		return strconv.ParseFloat(fl, 64)
	}
	f := 0.0
	for i, seg := range parts {
		switch i {
		case len(parts) - 1:
			num, err := strconv.ParseFloat(seg, 64)
			if err != nil {
				return -1, fmt.Errorf("failed to parse '%v' from %v: %v", seg, fl, err)
			}
			f += num
		default:
			adder := 1
			switch i {
			case len(parts) - 2:
				adder = 60
			case len(parts) - 3:
				adder = 60 * 60
			case len(parts) - 4:
				adder = 60 * 60 * 24
			}
			numI, err := strconv.Atoi(seg)
			if err != nil {
				return -1, fmt.Errorf("failed to parse '%v' from %v: %v", seg, fl, err)
			}
			f += float64(numI * adder)
		}
	}
	return f, nil
}

func (mp *mergeProc) defineOrder() error {
	num := len(mp.source)
	switch num {
	default:
		return fmt.Errorf("can't set order of %v sources", num)
	case 6:
		mp.descr = merge6
		l, r, c, lfe, ls, rs := "", "", "", "", "", ""
		for path := range mp.source {
			path_low := strings.ToLower(path)
			switch {
			case strings.HasSuffix(path_low, "01.wav") || strings.HasSuffix(path_low, "l.wav"):
				l = path
			case strings.HasSuffix(path_low, "02.wav") || strings.HasSuffix(path_low, "r.wav"):
				r = path
			case strings.HasSuffix(path_low, "03.wav") || strings.HasSuffix(path_low, "c.wav"):
				c = path
			case strings.HasSuffix(path_low, "04.wav") || strings.HasSuffix(path_low, "lfe.wav"):
				lfe = path
			case strings.HasSuffix(path_low, "05.wav") || strings.HasSuffix(path_low, "ls.wav"):
				ls = path
			case strings.HasSuffix(path_low, "06.wav") || strings.HasSuffix(path_low, "rs.wav"):
				rs = path
			default:
			}
		}
		mp.sourcesOrdered = []string{l, r, c, lfe, ls, rs}
	case 2:
		mp.descr = merge2
		l, r := "", ""
		for path := range mp.source {
			path_low := strings.ToLower(path)
			switch {
			case strings.HasSuffix(path_low, "01.wav") || strings.HasSuffix(path_low, "l.wav"):
				l = path
			case strings.HasSuffix(path_low, "02.wav") || strings.HasSuffix(path_low, "r.wav"):
				r = path
			default:
			}
		}
		mp.sourcesOrdered = []string{l, r}
	}
	for i, channel := range mp.sourcesOrdered {
		if channel == "" {
			return fmt.Errorf("channel %v undefined", i)
		}
	}
	return nil
}

func (mp *mergeProc) Prompt() (string, error) {
	s := "ffmpeg -hide_banner"
	list := []string{}
	for k := range mp.source {
		list = append(list, k)
	}
	for _, fl := range mp.sourcesOrdered {
		s += fmt.Sprintf(" -i %v", fl)
	}
	if mp.tgDuration == 0.0 {
		mp.tgDuration = (mp.srcDuration / mp.srcFps) * 25
	}

	if mp.srcFps == 0 {
		return "", fmt.Errorf("source fps not set")
	}
	if mp.srcDuration == 0 {
		return "", fmt.Errorf("source duration not set")
	}
	if mp.target == "" {
		return "", fmt.Errorf("output name not set")
	}
	fctext := ""
	switch mp.descr {
	case merge2:
		fctext += "amerge=inputs=2,channelmap=channel_layout=stereo"
	case merge6:
		fctext += "amerge=inputs=6"
	default:
		return "", fmt.Errorf("merge scenario unknown '%v'", mp.descr)
	}
	fctext = fmt.Sprintf("%v,aresample=48000,atempo=(%v/%v)", fctext, mp.srcDuration, mp.tgDuration)
	openFC := ""
	for i := range mp.sourcesOrdered {
		openFC += fmt.Sprintf("[%v:a:0]", i)
	}
	fc := fmt.Sprintf("%v%v[out]", openFC, fctext)
	s += fmt.Sprintf(` -filter_complex "%v"`, fc)
	s += fmt.Sprintf(" -map \"[out]\" -c:a alac -compression_level 0 -map_metadata -1 -map_chapters -1 %v.m4a", mp.target)
	return s, nil
}

func fpsToFloat(fps string) float64 {
	data := strings.Split(fps, "/")
	i1, _ := strconv.Atoi(data[0])
	i2, _ := strconv.Atoi(data[1])
	fl := float64(i1) / float64(i2)
	fli := float64(int(fl*1000)) / 1000
	return fli
}
