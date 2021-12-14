package silence

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/Galdoba/devtools/cli/command"
)

type silence struct {
	coords []timeCoord
}

//Detect - слушает файл и возвращает координаты тишины
func Detect(path string) (*silence, error) {
	fmt.Printf("Scanning file: %v\n", path)
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	if !audioStreamContained(f.Name()) {
		fmt.Printf("Warning: no audio stream detected\n %v ", f.Name())
		return nil, fmt.Errorf("file have no audio stream")
	}
	listenReport, err := command.New(
		command.CommandLineArguments(fmt.Sprintf("ffmpeg -i %v -vn -af silencedetect=n=-60dB:d=0.5 -f null - -loglevel info", f.Name())),
		//command.Set(command.TERMINAL_ON),
		command.Set(command.BUFFER_ON),
	)
	listenReport.Run()
	coords, err := parseSilenceData(listenReport.StdErr())
	if err != nil {
		return nil, err
	}
	fmt.Printf("Completed\n")
	return &silence{coords}, nil
}

func parseSilenceData(data string) ([]timeCoord, error) {
	lines := strings.Split(data, "\n")
	info := []string{}
	for _, l := range lines {
		if strings.Contains(l, "silencedetect @ ") {
			info = append(info, l)
		}
	}
	start := []float64{}
	dur := []float64{}
	tCords := []timeCoord{}
	for _, i := range info {
		switch {
		case strings.Contains(i, "silence_start"):
			str := strings.Split(i, "silence_start: ")
			s := strings.TrimSpace(str[1])
			flt, _ := strconv.ParseFloat(s, 64)
			start = append(start, flt)
		case strings.Contains(i, "silence_duration"):
			str := strings.Split(i, "silence_duration:")
			s := strings.TrimSpace(str[1])
			flt, _ := strconv.ParseFloat(s, 64)
			dur = append(dur, flt)
		}
	}
	if len(start) != len(dur) {
		return tCords, fmt.Errorf("silence coordinated parsing failed")
	}
	for i := 0; i < len(start); i++ {
		tCords = append(tCords, timeCoord{start[i], dur[i]})
	}
	return tCords, nil
}

//timeCoord - описывает кусок пустоты на треке. Пустотой считаем громкость ниже -104.5 Db
type timeCoord struct {
	start    float64
	duration float64
}

//end - предпологаймая точка конца пустоты
func (tc *timeCoord) end() float64 {
	return round(tc.start + tc.duration)
}

///////////Helpers
//round - округляет до 6 разрядов для единообразия с ffmpeg
func round(f float64) float64 {
	return math.Trunc(f/0.000001) * 0.000001
}

func debugMsg(s string) {
	fmt.Printf("debug Message: " + s + "\n")
}

func audioStreamContained(path string) bool {
	//	out, errors, err := cli.RunToAll("ffprobe", "-i", path, "-show_entries", "format : stream=codec_type")
	com, err := command.New(
		command.CommandLineArguments(fmt.Sprintf("ffprobe -i %v", path)),
		command.Set(command.BUFFER_ON),
	)
	if err != nil {
		return false
	}
	err = com.Run()
	out := com.StdErr()
	lines := strings.Split(out, "\n")
	for _, l := range lines {
		if strings.Contains(l, "Stream") && strings.Contains(l, "Audio") {
			return true
		}
	}
	return false
}

/*
Task:
1.оценить содержит ли файл звук
1а.оценить длительность звука
2.общая оценка звука
	если пустот более 3 секунд нет ==> конец
3 создаем warning файл
	если есть пустоты в первых или последних 3 минутах
	3а. отрезаем первые и последние 3 минуты
		слушаем внимательнее
		дополняем warning файл с описанием пустот и предположение откуда отрезать



*/

//consoleFileName := strings.ReplaceAll(f.Name(), "\\", "\\\\")
//out, errors, err := cli.RunToAll("ffprobe", "-i", consoleFileName, "-show_entries", "format : stream=codec_type")
//	out, errors, err := cli.RunToAll("ffmpeg", "-i", consoleFileName, "-af", "silencedetect=n=-75dB:d=2", "-f", "null", "-", "-loglevel", "info")
//sERR, cERR := cli.RunToFile("d:\\MUX\\tests\\log2.txt", "ffmpeg", "-i", consoleFileName, "-af", "silencedetect=n=-90dB:d=2", "-f", "null", "-", "-loglevel", "info")

//fmt.Println(f.Name(), "correct\n ")
//debugMsg("END: Detect(f *os.File) (*silence, error)")
