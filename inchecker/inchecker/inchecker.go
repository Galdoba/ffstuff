package inchecker

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/malashin/ffinfo"
)

const (
	ffinfoCodecType     = "codec_type"
	ffinfoChannels      = "channels"
	ffinfoChannelLayout = "channel_layout"
)

type InChecker interface {
	CheckValidity(string) error
}

func Test() {
	for i, val := range os.Args {
		if i == 0 {
			continue
		}
		//		val = strings.ReplaceAll(val, "\\\\", "\\")
		repFile, err := ffinfo.Probe(val)
		if err != nil {
			fmt.Println(err)
		}
		if err := CheckInput(repFile); err != nil {
			fmt.Print("\n")
			fmt.Println(err.Error())
			continue
		}
		fmt.Println("...ok")
		// cli.RunConsole("ffprobe", "-i", val)
		// fmt.Println("----------------------------------------")
		// fmt.Println(repFile.String())
		// fmt.Println("----------------------------------------")

	}
}

func CheckInput(f *ffinfo.File) error {
	fmt.Print("check file: ", f.Format.Filename)
	err := errors.New("Unchecked")
	switch CollectInfo(f, ffinfoCodecType) {
	default:
		err = errors.New("WARNING: Codec Type '" + CollectInfo(f, ffinfoCodecType) + "' unimplemented")
	case "audio":
		err = checkAudio(f)
	case "video":
		fmt.Println(f.String())
		err = checkVideo(f)
		fmt.Println("========================================")
		fmt.Println(f.Streams[0].SampleAspectRatio)
		fmt.Println(f.Streams[0].DisplayAspectRatio)
	}
	return err
}

func checkAudio(repFile *ffinfo.File) error {
	fileName := repFile.Format.Filename
	report := ""
	expChan, expLayout := expectedFromAudio(fileName)
	if expChan != CollectInfo(repFile, ffinfoChannels) {
		report += "Warning! Channels:\n"
		report += "Expect: '" + expChan + "', have '" + CollectInfo(repFile, ffinfoChannels) + "'\n"
	}
	if expLayout != CollectInfo(repFile, ffinfoChannelLayout) {
		report += "Warning! Channel Layout:\n"
		report += "Expect: '" + expLayout + "', have '" + CollectInfo(repFile, ffinfoChannelLayout) + "'\n"
	}
	if report != "" {
		return errors.New(report)
	}
	return nil
}

func checkVideo(f *ffinfo.File) error {
	fileName := f.Format.Filename
	report := ""
	fmt.Println(fileName, report)
	return nil
}

func CollectInfo(f *ffinfo.File, key string) string {
	key = strings.ToLower(key)
	key = strings.ReplaceAll(key, " ", "_")
	switch key {
	case ffinfoCodecType:
		return f.Streams[0].CodecType
	case ffinfoChannels:
		return strconv.Itoa(f.Streams[0].Channels)
	case ffinfoChannelLayout:
		return f.Streams[0].ChannelLayout

	}
	return "UNKNOWN KEY"
}

func expectedFromAudio(fileName string) (string, string) { //TODO: переписать иак чтобы оно собирало тэги из имени файла
	if strings.Contains(fileName, "_AUDIORUS20") {
		return "2", "stereo"
	}
	if strings.Contains(fileName, "_AUDIORUS51") {
		return "6", "5.1"
	}
	if strings.Contains(fileName, "_AUDIOENG20") {
		return "2", "stereo"
	}
	if strings.Contains(fileName, "_AUDIOENG51") {
		return "6", "5.1"
	}
	return "unknown audio tags", "unknown audio tags"
}

func expectedFromVideo(fileName string) (string, string) { //TODO: переписать иак чтобы оно собирало тэги из имени файла
	if strings.Contains(fileName, "_HD") {
		wh := "1920/1080"
		pixFmt := "yuv420"
		fps := "25/1"
		fmt.Println(wh, pixFmt, fps)

		return "2", "stereo"
	}

	return "unknown audio tags", "unknown audio tags"
}

//go run inchecker.go f:\Work\petr_proj\___IN\IN_2021-02-05\Zloy_duh_AUDIORUS51.m4a f:\Work\petr_proj\___IN\IN_2021-02-05\Zvonok_poslednyaya_glava_AUDIORUS51.m4a f:\Work\petr_proj\___IN\IN_2021-02-05\Daniel_Isnt_Real_HD.mp4 f:\Work\petr_proj\___IN\IN_2021-02-05\Daniel_Isnt_Real_SD.mp4
