package inchecker

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/malashin/ffinfo"
)

const (
	ffinfoCodecType     = "codec_type"
	ffinfoChannels      = "channels"
	ffinfoChannelLayout = "channel_layout"
	ffinfoWidth         = "width"
	ffinfoHeight        = "height"
	ffinfoPixFmt        = "pix_fmt"
	ffinfoFPS           = "r_frame_rate"
	ffinfoSAR           = "sample_aspect_ratio"
)

//InChecker - checks video and audio files to match input base format
type InChecker interface {
	CheckValidity(string) error
}

//Checker - mounts InChecker interface
type Checker struct{}

//NewChecker -
func NewChecker() Checker {
	return Checker{}
}

//CheckValidity - Checks File for valid format
func (ch *Checker) CheckValidity(path string) error {
	repFile, err := ffinfo.Probe(path)
	if err != nil {
		return errors.New("\n" + "ffinfo.Probe(string): " + err.Error())
	}
	if err := checkInput(repFile); err != nil {
		return err
	}
	return nil
}

func knownTags() []string {
	return []string{
		"HD",
		"SD",
		"43",
		"AUDIOENG20",
		"AUDIORUS20",
		"AUDIOENG51",
		"AUDIORUS51",
		"TRL",
		"Proxy",
	}
}

func checkInput(f *ffinfo.File) error {
	err := errors.New("Initial Error (MUST NOT HAPPEN)")
	base, ext, tags := decodeName(f.Format.Filename)
	if ext == "srt" {
		return errors.New("\nFile is subtitles")
	}
	for stream := 0; stream < len(f.Streams); stream++ {
		switch collectInfo(f, 0, ffinfoCodecType) {
		default:
			err = errors.New("WARNING: Codec Type '" + collectInfo(f, 0, ffinfoCodecType) + "' unknown")
			fmt.Println(base, ext)
		case "audio":
			err = checkAudio(f, stream, tags)
		case "video":
			err = checkVideo(f, stream)
		}
	}
	return err
}

func checkAudio(repFile *ffinfo.File, stream int, tags []string) error {
	fileName := repFile.Format.Filename
	report := "\n"
	expChan, expLayout := expectedFromAudio(fileName)
	if expChan != collectInfo(repFile, stream, ffinfoChannels) {
		report += "Channels:"
		report += " expect '" + expChan + "', have '" + collectInfo(repFile, stream, ffinfoChannels) + "'\n"
	}
	if expLayout != collectInfo(repFile, stream, ffinfoChannelLayout) {
		report += "Channel Layout:"
		report += " expect '" + expLayout + "', have '" + collectInfo(repFile, stream, ffinfoChannelLayout) + "'\n"
	}
	if report != "\n" {
		return errors.New(report)
	}
	return nil
}

func checkVideo(f *ffinfo.File, stream int) error {
	fileName := f.Format.Filename
	report := "\n"
	expWH, expPixFmt, expFPS, expSAR := expectedFromVideo(fileName)
	trueWH := collectInfo(f, stream, ffinfoWidth) + "/" + collectInfo(f, stream, ffinfoHeight)
	truePixFmt := collectInfo(f, stream, ffinfoPixFmt)
	trueFPS := collectInfo(f, stream, ffinfoFPS)
	trueSAR := collectInfo(f, stream, ffinfoSAR)
	if expWH != trueWH {
		report += "Width/Height:"
		report += " expect '" + expWH + "', have '" + trueWH + "'\n"
	}
	if expPixFmt != truePixFmt {
		report += "PixFmt:"
		report += " expect '" + expPixFmt + "', have '" + truePixFmt + "'\n"
	}
	if trueFPS != expFPS {
		report += "FPS:"
		report += " expect '" + expFPS + "', have '" + trueFPS + "'\n"
	}
	if trueSAR != expSAR && trueSAR != "" {
		report += "SAR:"
		report += " expect '" + expSAR + "', have '" + trueSAR + "'\n"
	}
	if report != "\n" {
		return errors.New(report)
	}
	return nil
}

func compareStrData(expected, real, dataType string) error {
	if real == "" {
		fmt.Println("!have no data on SAR!")
		return nil
	}
	if expected != real {
		report := "File data do not match expected:"
		report += "Have '" + dataType + "' : " + real + "\n"
		report += "Expect '" + dataType + "' : " + expected + "\n"
		return errors.New(report)
	}
	return nil
}

func collectInfo(f *ffinfo.File, stream int, key string) string {
	key = strings.ToLower(key)
	key = strings.ReplaceAll(key, " ", "_")
	switch key {
	case ffinfoCodecType:
		return f.Streams[stream].CodecType
	case ffinfoChannels:
		return strconv.Itoa(f.Streams[stream].Channels)
	case ffinfoChannelLayout:
		return f.Streams[stream].ChannelLayout
	case ffinfoWidth:
		return strconv.Itoa(f.Streams[stream].Width)
	case ffinfoHeight:
		return strconv.Itoa(f.Streams[stream].Height)
	case ffinfoPixFmt:
		return f.Streams[stream].PixFmt
	case ffinfoFPS:
		return f.Streams[stream].RFrameRate
	case ffinfoSAR:
		return f.Streams[stream].SampleAspectRatio

	}
	return "UNKNOWN KEY"
}

func expectedFromAudio(fileName string) (string, string) { //TODO: переписать иак чтобы оно собирало тэги из имени файла

	if strings.Contains(fileName, "_AUDIORUS51") {
		return "6", "5.1"
	}
	if strings.Contains(fileName, "_AUDIOENG51") {
		return "6", "5.1"
	}
	if strings.Contains(fileName, "_AUDIO51") {
		return "6", "5.1"
	}
	if strings.Contains(fileName, "_AUDIORUS20") {
		return "2", "stereo"
	}
	if strings.Contains(fileName, "_AUDIOENG20") {
		return "2", "stereo"
	}
	if strings.Contains(fileName, "_AUDIO20") {
		return "2", "stereo"
	}

	return "unknown audio tags", "unknown audio tags"
}

func expectedFromVideo(fileName string) (string, string, string, string) { //TODO: переписать иак чтобы оно собирало тэги из имени файла
	if strings.Contains(fileName, "_HD") {
		wh := "1920/1080"
		pixFmt := "yuv420p"
		fps := "25/1"
		sar := "1:1"
		return wh, pixFmt, fps, sar
	}
	if strings.Contains(fileName, "_SD_43") {
		wh := "720/576"
		pixFmt := "yuv420p"
		fps := "25/1"
		sar := "16:15"
		return wh, pixFmt, fps, sar
	}
	if strings.Contains(fileName, "_SD") {
		wh := "720/576"
		pixFmt := "yuv420p"
		fps := "25/1"
		sar := "64:45"
		return wh, pixFmt, fps, sar
	}
	return "unknown video tags", "unknown video tags", "unknown video tags", "unknown video tags"
}

func decodeName(path string) (string, string, []string) {
	fileName := shortFileName(path)
	tags, ext := splitName(fileName)
	base, tags2 := nameBase(tags)
	return base, ext, tags2
}

func shortFileName(path string) string {
	data := strings.Split(path, "\\")
	fileName := path
	if len(data) > 1 {
		fileName = data[len(data)-1]
	}
	return fileName
}

func splitName(fileName string) ([]string, string) {
	data := strings.Split(fileName, "_")
	tags := []string{}
	ext := ""
	for index, val := range data {
		if index == len(data)-1 {
			p := strings.Split(val, ".")
			ext = p[len(p)-1]
			tags = append(tags, p[0])
			continue
		}
		tags = append(tags, val)
	}
	return tags, ext
}

func nameBase(tags []string) (string, []string) {
	base := ""
	tags2 := []string{}
	for i, val := range tags {
		for _, tg := range knownTags() {
			if tg == val {
				tags2 = append(tags2, tg)
			}
			if tg != val {
				continue
			}
		}
		if i != len(tags)-1 {
			base += val + "_"
		}
	}
	base = strings.TrimSuffix(base, "_")
	return base, tags2
}
