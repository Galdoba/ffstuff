package silence

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/Galdoba/devtools/cli/command"
)

const (
	LOUDNESS_BELOW_36db  = -36.0
	LOUDNESS_BELOW_48db  = -48.0
	LOUDNESS_BELOW_60db  = -60.0
	LOUDNESS_BELOW_72db  = -72.0 // предположителный дефолт
	LOUDNESS_BELOW_84db  = -84.0
	LOUDNESS_BELOW_96db  = -96.0
	LOUDNESS_BELOW_108db = -108.0
)

type silence struct {
	coords        []timeCoord
	totalDuration float64
}

//Detect - слушает файл и возвращает координаты тишины
func Detect(path string, loudnessBorder, duration float64) (*silence, error) {
	//Pre-check
	for _, err := range []error{
		assertPath(path),
		assertLounessBorder(loudnessBorder),
		assertDuration(duration),
	} {
		if err != nil {
			return nil, err
		}
	}
	//Body:
	fmt.Printf("File: %v\nScanning for silence segments below %v Db with duration >= %v+ seconds...\n", path, loudnessBorder, duration)
	listenReport, err := command.New(
		command.CommandLineArguments(fmt.Sprintf("ffmpeg -i %v -vn -af silencedetect=n=%vdB:d=%v -f null - -loglevel info", path, loudnessBorder, duration)),
		//command.Set(command.TERMINAL_ON),
		command.Set(command.BUFFER_ON),
	)
	listenReport.Run() // чтобы трэчить текущий прогресс нужно ставить флаг читать файл?
	coords, dr, err := parseSilenceData(listenReport.StdErr())
	if err != nil {
		return nil, err
	}
	fmt.Printf("Completed...\n")
	return &silence{coords, dr}, nil
}

func assertPath(path string) error {
	p, err := os.Stat(path)
	if err != nil {
		fmt.Printf("Error: %v\n", err.Error())
		return fmt.Errorf("stat error: %v", err.Error())
	}
	mode := p.Mode()
	switch {
	case mode.IsDir():
		fmt.Printf("Error: path is Directory (all directory scanning is not yet implemented)\n")
		return fmt.Errorf("%v is a directory", path)
	}
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	if err := audioStreamContained(f.Name()); err != nil {
		fmt.Printf("Error: %v\n", err.Error())
		return fmt.Errorf(err.Error())
	}
	return nil
}

func assertLounessBorder(loudnessBorder float64) error {
	if loudnessBorder >= 0 {
		fmt.Printf("Warning: Loudness border must be negative (have '%v')\n", loudnessBorder)
		return fmt.Errorf("loudness parameter have incorect value: %v", loudnessBorder)
	}
	if loudnessBorder <= -120 {
		fmt.Printf("Warning: Loudness border is below void\n")
		return fmt.Errorf("loudness parameter have incorect value: %v", loudnessBorder)
	}
	return nil
}

func assertDuration(duration float64) error {
	if duration < 0.04 {
		fmt.Printf("Warning: Duration is to low\n")
		return fmt.Errorf("duration parameter have incorect value: %v", duration)
	}
	if duration < 1.50 {
		fmt.Printf("Warning: Duration is very low - application may take long time to scan the file\n")
	}
	if duration > 8 {
		fmt.Printf("Warning: Duration is very long - application may miss all silence segments\n")
	}
	return nil
}

func parseSilenceData(data string) ([]timeCoord, float64, error) {
	lines := strings.Split(data, "\n")
	info := []string{}
	start := []float64{}
	dur := []float64{}
	var dr float64
	tCords := []timeCoord{}
	for _, l := range lines {
		if strings.Contains(l, "silencedetect @ ") {
			info = append(info, l)
		}
		if strings.Contains(l, "  Duration:") {
			dur := strings.TrimPrefix(l, "  Duration: ")
			dt := strings.Split(dur, ",")
			dr = ffTimeStampToFloat64(dt[0])
		}
	}
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
		return tCords, dr, fmt.Errorf("silence coordinated parsing failed")
	}
	for i := 0; i < len(start); i++ {
		tCords = append(tCords, timeCoord{start[i], dur[i]})
	}
	return tCords, dr, nil
}

func ffTimeStampToFloat64(stamp string) float64 {
	dt := strings.Split(stamp, ":")
	dtt := strings.Split(dt[2], ".")
	dt[0] = strings.TrimPrefix(dt[0], "0")
	h, _ := strconv.Atoi(dt[0])
	m, _ := strconv.Atoi(dt[1])
	s, _ := strconv.Atoi(dtt[0])
	f, _ := strconv.Atoi(dtt[1])
	sec := h*3600 + m*60 + s
	fl := float64(sec) + (float64(f) / 1000.0)
	return fl

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

func audioStreamContained(path string) error {
	//	out, errors, err := cli.RunToAll("ffprobe", "-i", path, "-show_entries", "format : stream=codec_type")
	com, err := command.New(
		command.CommandLineArguments(fmt.Sprintf("ffprobe -i %v", path)),
		command.Set(command.BUFFER_ON),
	)
	if err != nil {
		return err
	}
	err = com.Run()
	out := com.StdErr()
	lines := strings.Split(out, "\n")
	for _, l := range lines {
		if strings.Contains(l, "Stream") && strings.Contains(l, "Audio") {
			return nil
		}
	}
	return fmt.Errorf("file have no audio stream")
}

func (si *silence) String() string {
	duration := si.totalDuration / 100.0
	pts := []float64{}
	for i := 0; i < 100; i++ {
		pts = append(pts, duration*(float64(i)+1))
	}
	glph := []string{}
	silDet := 0
	for _, p := range pts {
		switch isSilent(p, si) {
		case false:
			glph = append(glph, "-")
		case true:
			glph = append(glph, "_")
			silDet++
		}
	}

	rep := ""
	if silDet > 0 {
		rep = "[" + strings.Join(glph, "") + "]"
		rep += fmt.Sprintf("\nSilence detected: %v", silDet) + "%"
	} else {
		rep += "No silence detected"
	}
	return rep
}

func (si *silence) Timings() string {
	rep := ""
	for i, s := range si.coords {
		rep += fmt.Sprintf("Segment %v: from %v to %v (%v)\n", i+1, secondToTimecode(s.start), secondToTimecode(s.end()), secondToTimecode(s.duration))
	}
	return rep
}

func isSilent(p float64, si *silence) bool {
	for _, seg := range si.coords {
		if p >= seg.start && p <= seg.end() {
			return true
		}
	}
	return false
}

func secondToTimecode(fl float64) string {
	h := 0
	m := 0
	s := 0
	fr := 0
	for fl > 3600.0 {
		h++
		fl = fl - 3600
	}
	for fl > 60.0 {
		m++
		fl = fl - 60
	}
	for fl > 1.0 {
		s++
		fl = fl - 1
	}
	fl = fl * 100
	fr = int(fl) / 4
	tmcd := ""
	if h < 10 {
		tmcd += "0"
	}
	tmcd += strconv.Itoa(h) + ":"
	if m < 10 {
		tmcd += "0"
	}
	tmcd += strconv.Itoa(m) + ":"
	if s < 10 {
		tmcd += "0"
	}
	tmcd += strconv.Itoa(s) + ":"
	if fr < 10 {
		tmcd += "0"
	}
	tmcd += strconv.Itoa(fr)
	return tmcd
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
