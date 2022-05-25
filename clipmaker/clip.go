package clipmaker

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/Galdoba/ffstuff/fldr"
	"github.com/Galdoba/utils"
)

func NewClipMap() map[int]clip {
	cm := make(map[int]clip)
	return cm
}

type clip struct {
	index               int
	sourceFileName      string
	sourceFileType      string
	targetFileName      string
	clipStart           float64
	clipDuration        float64
	seqPosStartTimeCode string
	seqPosEndTimeCode   string
	nextClipIndex       int
	ffmpegArgs          []string
}

func NewClip(clipData string) (clip, error) {
	c := clip{}
	timestamps, err := parseTimeCodes(clipData)
	if err != nil {
		return c, err
	}
	c.index = parseFileIndex(clipData)
	c.sourceFileName = parseFileName(clipData)
	c.sourceFileType = extention(c.sourceFileName)
	//c.sourceFileFolder = fldr.InPath()
	c.clipStart = prmTime2Seconds(timestamps[0])
	c.clipDuration = utils.RoundFloat64(prmTime2Seconds(timestamps[1])-prmTime2Seconds(timestamps[0]), 3)
	c.seqPosStartTimeCode = timestamps[2]
	c.seqPosEndTimeCode = timestamps[3]
	c.formArgs()
	return c, nil
}

func (cl *clip) Index() int {
	return cl.index
}

func parseTimeCodes(clipData string) ([]string, error) {
	rawData := strings.Split(clipData, " ")
	timestamps := []string{}
	for _, val := range rawData {
		if len(val) != 11 {
			continue
		}
		valBt := []byte(val)
		if string(valBt[2]) != ":" || string(valBt[5]) != ":" || string(valBt[8]) != ":" {
			continue
		}
		timestamps = append(timestamps, val)
	}
	if len(timestamps) != 4 {
		err := errors.New("Timestamps parse Error: [" + strings.Join(timestamps, " ") + "]")
		return []string{}, err
	}
	return timestamps, nil
}

func prmTime2Seconds(s string) float64 {
	data := strings.Split(s, ":")
	hh, mm, ss, fr := 0, 0, 0, 0
	err := errors.New("No data")
	for i, val := range data {
		switch i {
		case 0:
			hh, err = strconv.Atoi(val)
		case 1:
			mm, err = strconv.Atoi(val)
		case 2:
			ss, err = strconv.Atoi(val)
		case 3:
			fr, err = strconv.Atoi(val)
		}
		if err != nil {
			errOut := errors.New(err.Error() + "\n Error: prmTime2Seconds(" + s + ")")
			fmt.Println(errOut.Error())
			return 0.0
		}
	}
	ms := utils.RoundFloat64(float64((fr*40)+(ss*1000)+(mm*60*1000)+(hh*3600*1000))*0.001, 3)
	return ms
}

func parseFileName(clipData string) string {
	rawData := strings.Split(clipData, " ")
	for _, val := range rawData {
		if containsAny(val, validFileExtetions()...) {
			return val
		}
	}
	return ""
}

func parseFileIndex(clipData string) int {
	rawData := strings.Split(clipData, " ")
	index, err := strconv.Atoi(rawData[0])
	if err != nil {
		fmt.Println(err)
	}
	return index
}

func containsAny(s string, subs ...string) bool {
	for _, sub := range subs {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}

func shortName(fileName string) string {
	p := strings.Split(fileName, ".")
	return strings.Join(p[0:len(p)-1], ".")
}

func extention(fileName string) string {
	p := strings.Split(fileName, ".")
	return "." + p[len(p)-1]
}

func validFileExtetions() []string {
	return []string{
		".m4a",
		".mp4",
	}
}

func indexStr(i int) string {
	s := strconv.Itoa(i)
	if i < 100 {
		s = "0" + s
	}
	if i < 10 {
		s = "0" + s
	}
	return s
}

func CutClip(cl clip) (string, []string) {
	program := "ffmpeg"
	argums := cl.formArgs()
	fmt.Println(argums)
	//cli.RunConsole(program, argums...)
	//"ffmpeg", "-i", file, "-map", "0:0", "-vcodec", "copy", "-an", "-t", premToFF(timeLen), "-ss", premToFF(timeStart), outputFile
	return program, argums
}

func (cl *clip) formArgs() []string {
	var argums []string
	//sdfsdf
	ssStamp := strconv.FormatFloat(cl.clipStart, 'f', 3, 64)
	tStamp := strconv.FormatFloat(cl.clipDuration, 'f', 3, 64)
	switch cl.sourceFileType {
	case ".mp4":
		cl.targetFileName = fldr.MuxPath() + shortName(cl.sourceFileName) + "_VCLIP_" + indexStr(cl.index) + extention(cl.sourceFileName)
		argums = []string{"-i", fldr.InPath() + cl.sourceFileName, "-an", "-map", "0:0", "-vcodec", "copy", "-ss", ssStamp, "-t", tStamp, cl.targetFileName}

	case ".m4a":
		cl.targetFileName = fldr.MuxPath() + shortName(cl.sourceFileName) + "_ACLIP_" + indexStr(cl.index) + extention(cl.sourceFileName)
		argums = []string{"-i", fldr.InPath() + cl.sourceFileName, "-vn", "-acodec", "copy", "-ss", ssStamp, "-t", tStamp, cl.targetFileName}
	default:
		fmt.Print("----------" + cl.sourceFileType + "------\n")
	}
	//fmt.Print(ssStamp, tStamp)
	return argums
}

func CutClipD(cl clip, sourceDir, targetDir string) (string, []string) {
	program := "ffmpeg"
	argums := cl.formArgsD(sourceDir, targetDir)
	fmt.Println(argums)
	//cli.RunConsole(program, argums...)
	//"ffmpeg", "-i", file, "-map", "0:0", "-vcodec", "copy", "-an", "-t", premToFF(timeLen), "-ss", premToFF(timeStart), outputFile
	return program, argums
}

func (cl *clip) formArgsD(sourceDir, targetDir string) []string {
	var argums []string

	//sdfsdf
	ssStamp := strconv.FormatFloat(cl.clipStart, 'f', 3, 64)
	tStamp := strconv.FormatFloat(cl.clipDuration, 'f', 3, 64)
	switch cl.sourceFileType {
	case ".mp4":
		cl.targetFileName = targetDir + shortName(cl.sourceFileName) + "_VCLIP_" + indexStr(cl.index) + extention(cl.sourceFileName)
		argums = []string{"-i", sourceDir + cl.sourceFileName, "-an", "-map", "0:0", "-vcodec", "copy", "-ss", ssStamp, "-t", tStamp, cl.targetFileName}

	case ".m4a":
		cl.targetFileName = targetDir + shortName(cl.sourceFileName) + "_ACLIP_" + indexStr(cl.index) + extention(cl.sourceFileName)
		argums = []string{"-i", sourceDir + cl.sourceFileName, "-vn", "-acodec", "copy", "-ss", ssStamp, "-t", tStamp, cl.targetFileName}
	default:
		fmt.Print("----------" + cl.sourceFileType + "------\n")
	}
	//fmt.Print(ssStamp, tStamp)
	return argums
}
