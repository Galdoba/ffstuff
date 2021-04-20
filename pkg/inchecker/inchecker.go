package inchecker

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Galdoba/ffstuff/pkg/namedata"
	"github.com/fatih/color"

	"github.com/malashin/ffinfo"
)

const (
	codecTypeVideo      = "video"
	codecTypeAudio      = "audio"
	ffinfoCodecType     = "codec_type"
	ffinfoCodecName     = "codec_name"
	ffinfoChannels      = "channels"
	ffinfoChannelLayout = "channel_layout"
	ffinfoWidth         = "width"
	ffinfoHeight        = "height"
	ffinfoPixFmt        = "pix_fmt"
	ffinfoFPS           = "r_frame_rate"
	ffinfoSAR           = "sample_aspect_ratio"
	ffinfoDuration      = "duration"
)

//Checker - mounts InChecker interface
type Checker struct {
	flagVocal bool
	pathList  []string
	//logger   logfile.Logger
	data     map[string]ffinfo.File
	groups   map[string][]string
	errorLog map[string][]error
}

//NewChecker -
func NewChecker() Checker {
	ch := Checker{}
	ch.groups = make(map[string][]string)
	ch.data = make(map[string]ffinfo.File)
	ch.errorLog = make(map[string][]error)
	//ch.logger = logfile.New(fldr.MuxPath()+"logfile.txt", logfile.LogLevelINFO)
	return ch
}

func (ch *Checker) AddTask(path string) {
	//fmt.Println("Adding", path)
	ch.pathList = append(ch.pathList, path)
	_, err := fileExists(path)
	if err != nil {
		ch.errorLog[path] = append(ch.errorLog[path], err)
		return
	}
	f, err := ffinfo.Probe(path)
	if err != nil {
		ch.errorLog[path] = append(ch.errorLog[path], err)
		//fmt.Println(f)
		return
	}
	ch.data[path] = *f
	//ch.pathList = append(ch.pathList, path)
	base := namedata.RetrieveBase(path)
	ch.groups[base] = append(ch.groups[base], path)

}

//Check - проверяет файлы на тему всех косяков о которых я додумался
func (ch *Checker) Check() []error {
	var allErrors []error
	for _, path := range ch.pathList {
		//f := ch.data[path]				DEBUG: принтует все сожержимое файла
		//fmt.Println(f.String())
		if len(ch.errorLog[path]) != 0 {
			fmt.Println(ch.errorLog[path])
			fmt.Println(ch.errorLog)
			continue
		}
		ch.errorLog[path] = addError(
			ch.checkCodecName(path),
			ch.checkLayout(path),
			ch.checkChannels(path),
			ch.checkDuration(path),
			ch.checkWidthHeight(path),
			ch.checkPixFmt(path),
			ch.checkFPS(path),
			ch.checkSAR(path),
		)
		for _, err := range ch.errorLog[path] {
			allErrors = append(allErrors, errors.New(path+" - "+err.Error()))
			// 	fmt.Println(ch.errorLog[path])
		}
	}
	return allErrors
}

func addError(allErrors ...error) []error {
	errLog := []error{}
	for _, err := range allErrors {
		if err != nil {
			errLog = append(errLog, err)
		}
	}
	return errLog
}

//Report - выводит результат проверки
func (ch *Checker) Report() {
	//color.Cyan("TEXT")
	for _, val := range ch.pathList {
		fmt.Print(val, ": ")
		if len(ch.errorLog[val]) == 0 {

			color.Green("		ok")
			continue
		}
		color.Yellow("		Warning!")
		//fmt.Print(val, ": ")
		for _, err := range ch.errorLog[val] {
			//fmt.Print("\n	")
			err = errors.New(val + " - " + err.Error())
		}
	}
}

func (ch *Checker) checkDuration(path string) error {
	if collectInfo(ch.data[path], 0, ffinfoCodecType) == codecTypeVideo {
		return nil
	}
	base := namedata.RetrieveBase(path)
	baseDuration := "0.0"
	if len(ch.groups[base]) < 2 {
		return nil
	}
	for _, p := range ch.groups[base] {
		data := ch.data[p]
		if collectInfo(data, 0, ffinfoCodecType) != codecTypeVideo {
			continue
		}
		baseDuration = collectInfo(data, 0, ffinfoDuration)
	}
	fileDuration := collectInfo(ch.data[path], 0, ffinfoDuration)
	if err := compareDuration(baseDuration, fileDuration); err != nil {
		return err
	}
	return nil
}

func (ch *Checker) checkCodecName(path string) error {
	switch collectInfo(ch.data[path], 0, ffinfoCodecType) {
	case codecTypeAudio:
		for stream := 0; stream < len(ch.data[path].Streams); stream++ {
			data := expectedFromAudio(path)
			codecName := collectInfo(ch.data[path], stream, ffinfoCodecName)
			if codecName != data[ffinfoCodecName] {
				return errors.New("Codec_Name: " + codecName + " (expect " + data[ffinfoCodecName] + ")")
			}
		}
	case codecTypeVideo:
		for stream := 0; stream < len(ch.data[path].Streams); stream++ {
			codecName := collectInfo(ch.data[path], stream, ffinfoCodecName)
			if codecName != "h264" {
				return errors.New("Codec_Name: " + codecName + " (expect " + "h264" + ")")
			}
		}
	}
	return nil
}

func (ch *Checker) checkChannels(path string) error {
	if collectInfo(ch.data[path], 0, ffinfoCodecType) != codecTypeAudio {
		return nil
	}
	for stream := 0; stream < len(ch.data[path].Streams); stream++ {
		data := expectedFromAudio(path)
		channnels := collectInfo(ch.data[path], stream, ffinfoChannels)
		if channnels != data[ffinfoChannels] {
			return errors.New("Channels: " + channnels + " (expect " + data[ffinfoChannels] + ")")
		}
	}
	return nil
}

func (ch *Checker) checkLayout(path string) error {
	if collectInfo(ch.data[path], 0, ffinfoCodecType) != codecTypeAudio {
		return nil
	}
	for stream := 0; stream < len(ch.data[path].Streams); stream++ {
		data := expectedFromAudio(path)
		layout := collectInfo(ch.data[path], stream, ffinfoChannelLayout)
		if layout != data[ffinfoChannelLayout] {
			return errors.New("Channel layout: " + layout + " (expect " + data[ffinfoChannelLayout] + ")")
		}
	}
	return nil
}

func (ch *Checker) checkWidthHeight(path string) error {
	if collectInfo(ch.data[path], 0, ffinfoCodecType) != codecTypeVideo {
		return nil
	}
	expWH, _, _, _ := expectedFromVideo(path)
	for stream := 0; stream < len(ch.data[path].Streams); stream++ {
		whData := collectInfo(ch.data[path], stream, ffinfoWidth) + "/" + collectInfo(ch.data[path], stream, ffinfoHeight)
		if whData != expWH {
			return errors.New("Width/Height: " + whData + " (expect " + expWH + ")")
		}
	}
	return nil
}

func (ch *Checker) checkPixFmt(path string) error {
	if collectInfo(ch.data[path], 0, ffinfoCodecType) != codecTypeVideo {
		return nil
	}
	_, expPixFmt, _, _ := expectedFromVideo(path)
	for stream := 0; stream < len(ch.data[path].Streams); stream++ {
		pixFmt := collectInfo(ch.data[path], stream, ffinfoPixFmt)
		if pixFmt != expPixFmt {
			return errors.New("PixFmt: " + pixFmt + " (expect " + expPixFmt + ")")
		}
	}
	return nil
}

func (ch *Checker) checkFPS(path string) error {
	if collectInfo(ch.data[path], 0, ffinfoCodecType) != codecTypeVideo {
		return nil
	}
	_, _, expFPS, _ := expectedFromVideo(path)
	for stream := 0; stream < len(ch.data[path].Streams); stream++ {
		fps := collectInfo(ch.data[path], stream, ffinfoFPS)
		if fps != expFPS {
			return errors.New("FPS: " + fps + " (expect " + expFPS + ")")
		}
	}
	return nil
}

func (ch *Checker) checkSAR(path string) error {
	if collectInfo(ch.data[path], 0, ffinfoCodecType) != codecTypeVideo {
		return nil
	}
	_, _, _, expSar := expectedFromVideo(path)
	for stream := 0; stream < len(ch.data[path].Streams); stream++ {
		sar := collectInfo(ch.data[path], stream, ffinfoSAR)
		if sar != expSar {
			if sar == "" {
				return errors.New("SAR: no data")
			}
			return errors.New("SAR: " + sar + " (expect " + expSar + ")")
		}
	}
	return nil
}

func compareDuration(baseDuration, fileDuration string) error {
	fDur, err := strconv.ParseFloat(fileDuration, 'f')
	if err != nil {
		return err
	}
	bDur, err := strconv.ParseFloat(baseDuration, 'f')
	if err != nil {
		return err
	}
	if fDur-bDur > 0.385 || fDur-bDur < -0.385 {
		fl := fDur - bDur
		flStr := strconv.FormatFloat(fl, 'f', 3, 64)
		return errors.New("Duration mismatch: " + flStr + " seconds")
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

func collectInfo(f ffinfo.File, stream int, key string) string {
	key = strings.ToLower(key)
	key = strings.ReplaceAll(key, " ", "_")
	switch key {
	case ffinfoCodecName:
		return f.Streams[stream].CodecName
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
	case ffinfoDuration:
		return f.Format.Duration

	}
	return "UNKNOWN KEY"
}

// func expectedFromAudio(fileName string) (string, string) { //TODO: переписать иак чтобы оно собирало тэги из имени файла

// 	if strings.Contains(fileName, "_AUDIORUS51") {
// 		return "6", "5.1"
// 	}
// 	if strings.Contains(fileName, "_AUDIOENG51") {
// 		return "6", "5.1"
// 	}
// 	if strings.Contains(fileName, "_AUDIO51") {
// 		return "6", "5.1"
// 	}
// 	if strings.Contains(fileName, "_AUDIORUS20") {
// 		return "2", "stereo"
// 	}
// 	if strings.Contains(fileName, "_AUDIOENG20") {
// 		return "2", "stereo"
// 	}
// 	if strings.Contains(fileName, "_AUDIO20") {
// 		return "2", "stereo"
// 	}

// 	return "unknown audio tags", "unknown audio tags"
// }

func expectedFromAudio(fileName string) map[string]string {
	data := make(map[string]string)
	data[ffinfoCodecName] = "alac"
	if stringsContainsAnyOf(fileName, "_AUDIORUS51", "_AUDIOENG51", "_AUDIO51") {
		data[ffinfoChannels] = "6"
		data[ffinfoChannelLayout] = "5.1"
	}
	if stringsContainsAnyOf(fileName, "_AUDIORUS20", "_AUDIOENG20", "_AUDIO20") {
		data[ffinfoChannels] = "2"
		data[ffinfoChannelLayout] = "stereo"
	}
	return data
} //TODO: переписать иак чтобы оно собирало тэги из имени файла

func stringsContainsAnyOf(s string, substr ...string) bool {
	for _, val := range substr {
		if strings.Contains(s, val) {
			return true
		}
	}
	return false
}

func expectedFromVideo(fileName string) (wh string, pixFmt string, fps string, sar string) { //TODO: переписать иак чтобы оно собирало тэги из имени файла
	if strings.Contains(fileName, "_HD__Proxy__") {
		wh := "480/270"
		pixFmt := "yuv420p"
		fps := "25/1"
		sar := "1:1"
		return wh, pixFmt, fps, sar
	}
	if strings.Contains(fileName, "_4K") {
		wh := "3840/2160"
		pixFmt := "yuv420p"
		fps := "25/1"
		sar := "1:1"
		return wh, pixFmt, fps, sar
	}
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
		for _, tg := range namedata.KnownTags() {
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

func fileExists(path string) (bool, error) {
	if _, err := os.Stat(path); err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, errors.New("File not exists ")

	} else {
		// Schrodinger: file may or may not exist. See err for details.
		return false, errors.New("Schrodinger: file may or may not exist. See err for details")
	}
}

func MediaDuration(path string) (float64, error) {
	ch := NewChecker()
	ch.AddTask(path)
	duration := ch.data[path].Format.Duration
	fDur, err := strconv.ParseFloat(duration, 'f')
	if err != nil {
		return 0.0, err
	}
	return fDur, nil
}
