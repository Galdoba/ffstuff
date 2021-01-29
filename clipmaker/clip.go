package clipmaker

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/Galdoba/ffstuff/cli"
)

type ClipMap map[int]clip

type clip struct {
	index               int
	sourceFileName      string
	sourceFileFolder    string
	sourceFileType      string
	clipStart           float64
	clipDuration        float64
	seqPosStartTimeCode string
	seqPosEndTimeCode   string
	nextClipIndex       int
}

func NewClip(clipData, folder string) (clip, error) {
	c := clip{}
	timestamps, err := parseTimeCodes(clipData)
	if err != nil {
		return c, err
	}
	c.clipStart = prmTime2Seconds(timestamps[2])
	c.clipDuration = prmTime2Seconds(timestamps[3])
	c.seqPosStartTimeCode = timestamps[0]
	c.seqPosEndTimeCode = timestamps[1]
	c.sourceFileName = parseFileName(clipData)
	c.sourceFileFolder = folder
	return c, nil
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
	ms := (fr * 40) + (ss * 1000) + (mm * 60 * 1000) + (hh * 3600 * 1000)
	sec := float64(ms / 1000)
	return sec
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

func Create(cl clip) {
	program := "ffmpeg"
	ssStamp := strconv.FormatFloat(cl.clipStart, 'f', 3, 64)
	tStamp := strconv.FormatFloat(cl.clipDuration, 'f', 3, 64)
	args := []string{"-i", cl.sourceFileFolder + cl.sourceFileName, "-an", "-map", "0:0", "-vcodec", "copy", "-ss", ssStamp, "-t", tStamp,
		shortName(cl.sourceFileName) + extention(cl.sourceFileName)}
	cli.RunConsole(program, args...)
	//"ffmpeg", "-i", file, "-map", "0:0", "-vcodec", "copy", "-an", "-t", premToFF(timeLen), "-ss", premToFF(timeStart), outputFile
}
