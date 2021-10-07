package inchecker

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"strconv"
	"strings"

	"github.com/Galdoba/ffstuff/pkg/namedata"
	"github.com/Galdoba/ffstuff/pkg/scanner"
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
	ffinfoHasBFrames    = "has_b_frames"
)

//Checker - mounts InChecker interface
type Checker struct {
	flagVocal bool
	pathList  []string
	data      map[string]ffinfo.File
	groups    map[string][]string
	errorLog  map[string][]error
}

//NewChecker -
func NewChecker() Checker {
	ch := Checker{}
	ch.data = make(map[string]ffinfo.File)
	ch.errorLog = make(map[string][]error)
	return ch
}

func (ch *Checker) AddTask(path string) {
	//исключаем файлы которые не проверяем точно
	usr, _ := user.Current()

	if stringsContainsAnyOf(path, ".srt", ".ready", "."+usr.Name) {
		return
	}
	ch.pathList = append(ch.pathList, path)
	_, err := fileExists(path)
	if err != nil {
		ch.errorLog[path] = append(ch.errorLog[path], err)
		return
	}
	f, err := ffinfo.Probe(path)
	if err != nil {
		ch.errorLog[path] = append(ch.errorLog[path], err)
		return
	}
	ch.data[path] = *f
}

//Check - проверяет файлы на тему всех косяков о которых я додумался
func (ch *Checker) Check() []error {
	var allErrors []error
	for _, path := range ch.pathList {
		//f := ch.data[path]				DEBUG: принтует все сожержимое файла
		//fmt.Println(f.String())
		if len(ch.errorLog[path]) != 0 {
			for _, err := range ch.errorLog[path] {
				allErrors = append(allErrors, errors.New(path+" - "+err.Error()))
			}
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
			ch.checkBFrames(path),
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

func (ch *Checker) checkDuration(path string) error {
	// если наш файл видео - проверка не делается
	if collectInfo(ch.data[path], 0, ffinfoCodecType) == codecTypeVideo {
		return nil
	}
	//Ищем файлы в текущей папке с базой нашего файла
	base := namedata.RetrieveBase(path)
	folder := namedata.RetrieveDirectory(path)
	relatedFiles, scanErr := scanner.Scan(folder, base)
	if scanErr != nil {
		return scanErr
	}
	baseDuration := "0.0"
	//прозваниваем все найденные
	for _, p := range relatedFiles {
		if stringsContainsAnyOf(p, ".ready", ".srt") {
			continue
		}
		f, err := ffinfo.Probe(p)
		if err != nil {
			ch.errorLog[path] = append(ch.errorLog[path], err)
			return err
		}
		//если не видео - пропускаем
		if collectInfo(*f, 0, ffinfoCodecType) != codecTypeVideo {
			continue
		}
		data := *f
		//узнаем базовую длинну
		baseDuration = collectInfo(data, 0, ffinfoDuration)
	}
	fileDuration := collectInfo(ch.data[path], 0, ffinfoDuration)
	//сравниваем длинну нашего файла и базового
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
			// //Проверка ошибки отсуствия видео тэга. TODO: перенсти отдельно в проверку имен
			// if stringsContainsAnyOf(path, ".mp4") {
			// 	err := errors.New("video tag absent or unknown")
			// 	if stringsContainsAnyOf(path, "_sd") || stringsContainsAnyOf(path, "_hd") {
			// 		err = nil
			// 	}
			// 	if err != nil {
			// 		return err
			// 	}
			// }
			// ////////////////////////////////////////
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
	data := expectedFromVideo(path)
	expWH := data[ffinfoWidth] + "/" + data[ffinfoHeight]
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
	data := expectedFromVideo(path)
	expPixFmt := data[ffinfoPixFmt]
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
	data := expectedFromVideo(path)
	expFPS := data[ffinfoFPS]
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
	data := expectedFromVideo(path)
	expSar := data[ffinfoSAR]
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

func (ch *Checker) checkBFrames(path string) error {
	if collectInfo(ch.data[path], 0, ffinfoCodecType) != codecTypeVideo {
		return nil
	}
	data := expectedFromVideo(path)
	expBF := data[ffinfoHasBFrames]
	for stream := 0; stream < len(ch.data[path].Streams); stream++ {
		bf := collectInfo(ch.data[path], stream, ffinfoHasBFrames)
		if bf != expBF {
			if bf == "" {
				return errors.New("Bframes: no data")
			}
			return errors.New("Bframes: " + bf + " (expect " + expBF + ")")
		}
	}
	return nil
}

func compareDuration(baseDuration, fileDuration string) error {
	if baseDuration == "0.0" {
		return errors.New("no video provided to compare duration with")
	}
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

// func knownTags() []string {
// 	return []string{
// 		"HD",
// 		"SD",
// 		"43",
// 		"AUDIOENG20",
// 		"AUDIORUS20",
// 		"AUDIOENG51",
// 		"AUDIORUS51",
// 		"TRL",
// 		"Proxy",
// 	}
// }

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
	case ffinfoHasBFrames:
		return strconv.Itoa(f.Streams[stream].HasBFrames)

	}
	return "collect info key not added!!!"
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
	if strings.Contains(fileName, ".ac3") {
		data[ffinfoCodecName] = "ac3"
		//	data[ffinfoChannelLayout] = "5.1(side)" TODO: решить где делать эту проверку
	}
	if stringsContainsAnyOf(fileName, "_AUDIORUS51", "_AUDIOENG51", "_AUDIO51", "rus51", "eng51") {
		data[ffinfoChannels] = "6"
		data[ffinfoChannelLayout] = "5.1"
	}
	if stringsContainsAnyOf(fileName, "_AUDIORUS20", "_AUDIOENG20", "_AUDIO20", "rus20", "eng20") {
		data[ffinfoChannels] = "2"
		data[ffinfoChannelLayout] = "stereo"
	}
	return data
} //TODO: переписать иак чтобы оно собирало тэги из имени файла

func stringsContainsAnyOf(s string, substr ...string) bool {
	s = strings.ToLower(s)
	for _, val := range substr {
		val = strings.ToLower(val)
		if strings.Contains(s, val) {
			return true
		}
	}
	return false
}

func expectedFromVideoOLD(fileName string) (wh string, pixFmt string, fps string, sar string) { //TODO: переписать иак чтобы оно собирало тэги из имени файла
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

func expectedFromVideo(fileName string) map[string]string { //TODO: переписать иак чтобы оно собирало тэги из имени файла
	data := make(map[string]string)
	data[ffinfoPixFmt] = "unknown (video tag not found)"
	data[ffinfoFPS] = "unknown (video tag not found)"
	data[ffinfoWidth] = "unknown (video tag not found)"
	data[ffinfoHeight] = "unknown (video tag not found)"
	data[ffinfoSAR] = "unknown (video tag not found)"
	data[ffinfoHasBFrames] = "unknown (video tag not found)"
	if stringsContainsAnyOf(fileName, ".mp4") { // для всего видео
		data[ffinfoPixFmt] = "yuv420p"
		data[ffinfoFPS] = "25/1"
		data[ffinfoHasBFrames] = "0"
	}
	if stringsContainsAnyOf(fileName, "_4k") {
		data[ffinfoWidth] = "3840"
		data[ffinfoHeight] = "2160"
		data[ffinfoSAR] = "1:1"
	}
	if stringsContainsAnyOf(fileName, "_hd") {
		data[ffinfoWidth] = "1920"
		data[ffinfoHeight] = "1080"
		data[ffinfoSAR] = "1:1"
	}
	if stringsContainsAnyOf(fileName, "_hd__proxy") {
		data[ffinfoWidth] = "480"
		data[ffinfoHeight] = "270"
		data[ffinfoSAR] = "1:1"
	}
	if stringsContainsAnyOf(fileName, "_sd") {
		data[ffinfoWidth] = "720"
		data[ffinfoHeight] = "576"
		data[ffinfoSAR] = "64:45"
	}
	if stringsContainsAnyOf(fileName, "_sd_43") {
		data[ffinfoSAR] = "16:15"
	}
	return data
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

//Report - выводит результат проверки
func (ch *Checker) Report(errs []error) {
	//color.Cyan("TEXT")
	nameLen := 0
	errorsFound := 0
	for _, val := range ch.pathList {
		nameLen = maxFrom(nameLen, len(val))
	}
	for i, val := range ch.pathList {
		originalName := val
		if i == 0 {
			head := "===INCHECKER REPORT"
			for len(head) < nameLen {
				head += "="
			}
			head += "======"
			fmt.Println(head)

		}
		for len(val) < nameLen {
			val += "."
		}
		fmt.Print(val, "..")
		if len(ch.errorLog[originalName]) == 0 {
			color.Green("ok")
			continue
		}
		color.Yellow("warning!")
		for _, err := range ch.errorLog[originalName] {
			//fmt.Print("\n	")
			err = errors.New(originalName + " - " + err.Error())
			errorsFound++
		}
	}
	tail := "======"
	for len(tail) < nameLen {
		tail += "="
	}
	tail += "======"
	fmt.Println(tail)
	fmt.Println(strconv.Itoa(len(errs)) + " errors found")
}

func maxFrom(a, b int) int {
	if a > b {
		return a
	}
	return b
}
