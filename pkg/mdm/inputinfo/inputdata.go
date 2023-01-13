package inputinfo

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

type inputdata struct {
	data []string
}

type parseInfo struct {
	scanTime    string
	filename    string
	metadata    map[string]string
	duration    float64
	start       float64
	globBitrate int
	streams     []stream
	parsedLines int
	parseStage  int
}

type stream struct {
	data     string
	metadata map[string]string
}

const (
	unknown = iota
	stage_ParseStart
	stage_ParseScanTime
	stage_ParseFilename
	stage_ParseGlobalMeta
	stage_ParseDuration
	stage_ParseStreams
	stage_ParseStreamMeta
	parseMethod_FFLITE
	parseMethod_FFMPEG
	prefix_StartTime       = "Started:"
	prefix_FFLITE_input    = "INPUT "
	prefix_FFMPEG_input    = "Input "
	prefix_FFMPEG_Metadata = "Metadata:"
	prefix_FFMPEG_Duration = "Duration:"
)

func hasTrigger(line, trigger string) bool {
	return strings.Contains(line, trigger)
}

func hasStreamInfo(line string) bool {
	//  Stream #0:0(und): Video: h264 (High 4:2:2) (avc1 / 0x31637661), yuv422p, 1920x1080 [SAR 1:1 DAR 16:9], 10687 kb/s, 25 fps, 25 tbr, 12800 tbn, 50 tbc (default)
	//    0:0 Video: prores (HQ) (apch / 0x68637061), yuv422p10le(tv, bt709, progressive), 1920x1080, SAR 1:1 DAR 16:9, 24 tbr, 24 tbn, 24 tbc
	//  Stream #0:0: Audio:
	r := regexp.MustCompile(`Stream\ \#\d+\:`)
	streamTag := r.FindString(line)
	if streamTag != "" {
		return true
	}
	r = regexp.MustCompile(`\ \ \ \ \d+\:\d+`)
	streamTag = r.FindString(line)
	if streamTag != "" {
		return true
	}
	//       \ (\#)?\d+\:\d+
	return false
}

func (inp *inputdata) String() string {
	return strings.Join(inp.data, "")
}

func newStream() stream {
	return stream{"", make(map[string]string)}
}

func parse(input inputdata) (*parseInfo, error) {
	pi := parseInfo{}
	pi.metadata = make(map[string]string)
	pi.start = math.NaN()
	pMethod := unknown
	pStage := stage_ParseScanTime
	metadata := make(map[string]string)
	stream := stream{"NO DATA", nil}
	for _, line := range input.data {
		switch line {
		case "At least one output file must be specified":
			continue
		}
		//fmt.Println("START: ", line)
		switch {
		case hasTrigger(line, prefix_StartTime) && (pStage == stage_ParseScanTime):
			pi.scanTime = grepScantime(line)
			pStage = stage_ParseFilename
		case hasTrigger(line, prefix_FFMPEG_input):
			pMethod = parseMethod_FFMPEG
			pStage = stage_ParseFilename
		case hasTrigger(line, prefix_FFLITE_input):
			pMethod = parseMethod_FFLITE
			pStage = stage_ParseFilename
		case hasTrigger(line, prefix_FFMPEG_Metadata):
			pStage = stage_ParseGlobalMeta
			//continue
		case hasTrigger(line, prefix_FFMPEG_Duration):
			pStage = stage_ParseDuration
		case hasStreamInfo(line):
			pStage = stage_ParseStreams
		}
		switch pStage {
		case stage_ParseStreams:
			if stream.data != "NO DATA" {

				////
				pi.metadata = metadata
				delete(pi.metadata, "Metadata")
				//pi.AddGlobalMetaData(metadata)
			} else {
				pi.AddStreamData(stream, metadata)
			}
			metadata = make(map[string]string)
		}
		//fmt.Println(pStage, line)
		switch pStage {
		case stage_ParseFilename:
			switch pMethod {
			case parseMethod_FFMPEG:
				pi.filename = grepFilenameFFMPEG(line)
			case parseMethod_FFLITE:
				pi.filename = grepFilenameFFLITE(line)
			}
		case stage_ParseGlobalMeta:
			switch pMethod {
			default:
				panic(0)
			case parseMethod_FFMPEG:
				key, val := grepGlobalMetadataFFMPEG(line)
				metadata[key] = val
			}
		case stage_ParseDuration:

			switch pMethod {
			default:
				panic("parse method unknown")
			case parseMethod_FFMPEG:
				//fmt.Println("|||", line)
				pi.duration, pi.start, pi.globBitrate = grepDurationDataFFMPEG(line)
			case parseMethod_FFLITE:
				pi.duration, pi.globBitrate = grepDurationDataFFLITE(line)
			}
			if pi.globBitrate == -1 {
				fmt.Println(pMethod)
				panic(pi.filename)
			}
		case stage_ParseStreams:
			stream = newStream()
			stream.data = line
			pStage = stage_ParseStreamMeta
		case stage_ParseStreamMeta:
			key, val := grepGlobalMetadataFFMPEG(line)
			metadata[key] = val
		}
	}
	pi.AddStreamData(stream, metadata)
	return &pi, nil
}

func (pi *parseInfo) AddStreamData(stream stream, metadata map[string]string) {
	stream.metadata = metadata
	delete(stream.metadata, "Metadata")
	pi.streams = append(pi.streams, stream)
}

/*
^     # start of string
\s*   # optional whitespace
(\w+) # one or more alphanumeric characters, capture the match
\s*   # optional whitespace
\(    # a (
\s*   # optional whitespace
(\d+) # a number, capture the match
\D+   # one or more non-digits
(\d+) # a number, capture the match
\D+   # one or more non-digits
\)    # a )
\s*   # optional whitespace
$     # end of string
`^((?:\d{1,3}\.){3}\d{1,3}) ([a-zA-Z]{3} \d{1,2} \d{4} \d{1,2}:\d{2}:\d{2}) (.*)`
`^(\d{2}\.\d{2}\.\d{2} \d{2}:\d{2}:\d{2},\d{2})(.*)`
"Jan 02 2006 15:04:05"
Started: 09.01.2023 17:27:19,42
*/

//parseScantime - ищет время сканирования. Опционально.
//#запускать если данные уже получены
func grepScantime(line string) string {
	r := regexp.MustCompile(`\d{2}\.\d{2}\.\d{4}\ \d{2}\:\d{2}\:\d{2}\,\d{2}`)
	//scanTime := strings.TrimPrefix(line, prefix_StartTime)
	return r.FindString(line)
}

func grepFilenameFFMPEG(line string) string {
	r := regexp.MustCompile(`\ \'(.*)\'`)
	name := r.FindString(line)
	name = strings.TrimPrefix(name, " '")
	name = strings.TrimSuffix(name, "'")
	return name
}

func grepFilenameFFLITE(line string) string {
	//INPUT 0: TRL_NP_2_2min_6+_vkino.mxf
	r := regexp.MustCompile(`INPUT\ \d*\:\ .*`)
	name := r.FindString(line)
	parts := strings.Split(name, ": ")
	if len(parts) <= 1 {
		return ""
	}
	name = parts[1]
	return name
}

func grepGlobalMetadataFFMPEG(line string) (string, string) {
	key, val := "", ""
	data := strings.Split(line, ":")
	key = data[0]
	if len(data) == 1 {
		val = ""
	} else {
		val = strings.Join(data[1:], ":")
	}
	key = strings.TrimSpace(key)
	val = strings.TrimSpace(val)
	return key, val
}

func parseMetaGlobal(line string) (string, string) {
	if !strings.HasPrefix(line, "    ") {
		return "", ""
	}
	line = strings.TrimSpace(line)
	fields := strings.Split(line, ": ")
	for i, field := range fields {
		fields[i] = strings.TrimSpace(field)
	}
	if len(fields) < 2 {
		return fields[0], ""
	}
	return fields[0], strings.Join(fields[1:], ": ")
}

func grepDurationDataFFMPEG(line string) (float64, float64, int) {
	//  Duration: 00:46:12.02, start: 0.000000, bitrate: 10952 kb/s
	duration := 0.0
	start := 0.0
	bitrate := 0
	data := strings.Split(line, ",")
	for i, ln := range data {
		switch i {
		case 0:
			duration = grepDuration(ln)
		case 1:
			start = grepStart(ln)
		case 2:
			bitrate = grepBitrate(ln)
		}
	}
	return duration, start, bitrate
}

func grepDurationDataFFLITE(line string) (float64, int) {
	duration := 0.0
	bitrate := 0
	data := strings.Split(line, ",")
	for i, ln := range data {
		switch i {
		case 0:
			duration = grepDuration(ln)
		case 1:
			bitrate = grepBitrate(ln)
		}
	}
	return duration, bitrate
}

func durationStrToFl64(dur string) float64 {
	dur = strings.ReplaceAll(dur, ".", ":")
	dt := strings.Split(dur, ":")
	durationFL := 0.0
	for i, d := range dt {
		dInt, err := strconv.Atoi(d)
		if err != nil {
			return -1
		}
		switch i {
		case 0:
			durationFL += float64(dInt) * 3600.0
		case 1:
			durationFL += float64(dInt) * 60.0
		case 2:
			durationFL += float64(dInt)
		case 3:
			durationFL += float64(dInt) * 0.01
		}
	}
	return durationFL
}

func grepDuration(line string) float64 {
	r := regexp.MustCompile(`Duration\:\ \d{2}\:\d{2}\:\d{2}\.\d{2}`)
	dur := r.FindString(line)
	dur = strings.TrimPrefix(dur, "Duration: ")
	return durationStrToFl64(dur)
}

func grepStart(line string) float64 {
	r := regexp.MustCompile(`start\:\ \d*\.\d*`)
	stSTR := r.FindString(line)
	stSTR = strings.TrimPrefix(stSTR, "start: ")
	return startStrToFl64(stSTR)
}

func startStrToFl64(start string) float64 {
	startFL, err := strconv.ParseFloat(start, 64)
	if err != nil {
		return math.NaN()
	}

	return startFL
}

func grepBitrate(line string) int {
	//Duration: 00:40:24.64, start: 0.000000, bitrate: 55084 kb/s
	//
	fmt.Println("grepBitrate", line)
	r := regexp.MustCompile(`\d*\ kb\/s`)
	btr := r.FindString(line)
	btr = strings.TrimSuffix(btr, " kb/s")
	br, err := strconv.Atoi(btr)
	if err != nil {
		return -1
	}
	return br
}

// func ripBetween(line, start, close string) string {
// 	parts := strings.Split(line, start)
// 	if len(parts) < 2 {
// 		return ""
// 	}
// 	startless := strings.Join(parts[1:], "")
// 	parts2 := strings.Split(startless, close)
// 	result := strings.Join(parts2[:len(parts2)-1], "")
// 	return result
// }

// func containsAll(fullString string, sub ...string) bool {
// 	for _, s := range sub {
// 		if !strings.Contains(fullString, s) {
// 			return false
// 		}
// 	}
// 	return true
// }

type videostream struct {
	pix_fmt       string
	width, height int
	bitrate       int
	sar           string
	dar           string
	fps           string
	tbr           string
	tbn           string
	tbc           string
	cc            string
	warnings      []string
}

type audiostream struct {
	hertz          int
	bitrate        string
	channels       int
	channel_layout string
}
