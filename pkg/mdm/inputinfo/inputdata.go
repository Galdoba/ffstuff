package inputinfo

import (
	"fmt"
	"math"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/Galdoba/devtools"
	"github.com/Galdoba/ffstuff/constant"
)

type inputdata struct {
	data []string
}

type ParseInfo struct {
	scanTime    string
	filename    string
	metadata    map[string]string
	duration    float64
	start       float64
	globBitrate int
	streams     []stream
	parsedLines int
	parseStage  int
	comment     string
	Video       []Videostream
	Audio       []Audiostream
	Data        []datastream
	Subtitles   []subtitlestream
	warnings    []string
	buffer      string
}

func (pi *ParseInfo) String() string {

	str := ""
	str += fmt.Sprintf("name: %v\n", pi.filename)
	str += fmt.Sprintf("GlobMeta: %v\n", len(pi.metadata))
	for k, v := range pi.metadata {
		str += fmt.Sprintf("  %v: %v\n", k, v)
	}
	str += fmt.Sprintf("durdata: %v - %v - %v\n", pi.duration, pi.start, pi.globBitrate)
	str += fmt.Sprintf("Streams: %v\n", len(pi.streams))
	str += "Video:\n"
	for _, v := range pi.Video {
		str += fmt.Sprintf("[vid]  %v\n", v.String())
	}

	str += "------------\n"
	return str
}

func (pi *ParseInfo) FileName() string {
	return pi.filename
}

func (pi *ParseInfo) Warnings() []string {
	return pi.warnings
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
	stage_ParseSideData
	stage_ParseDuration
	stage_ParseStreams
	stage_ParseStreamMeta
	parseMethod_FFLITE
	parseMethod_FFMPEG
	prefix_StartTime             = "Started:"
	prefix_FFLITE_input          = "INPUT "
	prefix_FFMPEG_input          = "Input "
	prefix_FFMPEG_Metadata       = "Metadata:"
	prefix_FFMPEG_Duration       = "Duration:"
	prefix_FFMPEG_Side_data      = "Side data"
	prefix_FFMPEG_Video_data     = ": Video:"
	prefix_FFMPEG_Audio_data     = ": Audio:"
	prefix_FFMPEG_Data_data      = ": Data:"
	prefix_FFMPEG_Subtitles_data = ": Subtitles:"
	warning_FFMPEG_Invalid_data  = "Invalid data found when processing input"
)

func hasTrigger(line, trigger string) bool {
	return strings.Contains(line, trigger)
}

func hasStreamInfo(line string) bool {
	if strings.Contains(line, "Video:") {
		return true
	}
	if strings.Contains(line, "Audio:") {
		return true
	}
	if strings.Contains(line, "Data:") {
		return true
	}
	if strings.Contains(line, "Subtitle:") {
		return true
	}
	return false
}

func (inp *inputdata) String() string {
	return strings.Join(inp.data, "")
}

func newStream() stream {
	return stream{"", make(map[string]string)}
}

func parse(input inputdata) (*ParseInfo, error) {
	pi := ParseInfo{}
	pi.metadata = make(map[string]string)
	pi.start = math.NaN()
	pMethod := unknown
	pStage := stage_ParseScanTime
	dsbMap := make(map[string]string)
	for _, line := range input.data {

		switch {
		default:
			switch len(pi.streams) {
			case 0:
				pStage = stage_ParseGlobalMeta
			default:
				pStage = stage_ParseStreamMeta
			}
		case pStage == stage_ParseSideData:
			if hasStreamInfo(line) {
				pStage = stage_ParseStreams
			}
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
			switch len(pi.streams) {
			case 0:
				pStage = stage_ParseGlobalMeta
			default:
				pStage = stage_ParseStreamMeta
			}
		case hasTrigger(line, prefix_FFMPEG_Side_data):
			pStage = stage_ParseSideData
			continue
		case hasTrigger(line, prefix_FFMPEG_Duration):
			pStage = stage_ParseDuration
		case hasStreamInfo(line):
			pStage = stage_ParseStreams
		case hasTrigger(line, warning_FFMPEG_Invalid_data):
			pi.warnings = append(pi.warnings, warning_FFMPEG_Invalid_data)
		}
		//fmt.Println("Stage", pStage, line)
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
				//panic("Unknown panic")
			case parseMethod_FFMPEG:
				key, val := grepGlobalMetadataFFMPEG(line)
				pi.injectMetadata(key, val)
			}
		case stage_ParseDuration:
			dsbMap = parseDSB(line)
			pi.injectDurationInfo(dsbMap)
			if pi.globBitrate == -1 {
				fmt.Println(pMethod)
				panic(pi.filename)
			}
		case stage_ParseStreams:
			stream := newStream()
			stream.data = line
			stream.metadata = make(map[string]string)
			pi.streams = append(pi.streams, stream)
			pStage = stage_ParseStreamMeta
		case stage_ParseStreamMeta:
			key, val := grepGlobalMetadataFFMPEG(line)
			pi.injectMetadata(key, val)
		case stage_ParseSideData:
			last := len(pi.streams) - 1
			pi.streams[last].metadata["Side data"] = line
		}
	}
	pi.parseStreams()
	pi.mergeWarnings()
	return &pi, nil
}

func (pi *ParseInfo) mergeWarnings() {
	for i, v := range pi.Video {
		if len(v.warnings) == 0 {
			continue
		}
		pi.warnings = append(pi.warnings, fmt.Sprintf("video %v:", i))
		for _, w := range v.warnings {
			pi.warnings = append(pi.warnings, "  "+w)
		}
	}
	for i, a := range pi.Audio {
		if len(a.warnings) == 0 {
			continue
		}
		pi.warnings = append(pi.warnings, fmt.Sprintf("audio %v:", i))
		for _, w := range a.warnings {
			pi.warnings = append(pi.warnings, "  "+w)
		}
	}
	if len(pi.Video) > 1 {
		pi.warnings = append(pi.warnings, fmt.Sprintf("file have %v video streams", len(pi.Video)))
	}
	if len(pi.Audio) > 2 {
		pi.warnings = append(pi.warnings, fmt.Sprintf("file have %v audio streams", len(pi.Audio)))
	}
	if len(pi.Subtitles) > 0 && (len(pi.Audio)+len(pi.Video)) > 0 {
		pi.warnings = append(pi.warnings, fmt.Sprintf("file have %v subtitle streams", len(pi.Subtitles)))
	}
}

func (pi *ParseInfo) injectDurationInfo(dsb map[string]string) {
	for k, v := range dsb {
		switch k {
		case "Duration":
			pi.duration = durationStrToFl64(v)
		case "start":
			pi.start = startStrToFl64(v)
		case "bitrate":
			x, _ := strconv.Atoi(strings.TrimSuffix(v, " kb/s"))
			pi.globBitrate = x
		}
	}
}

func (pi *ParseInfo) injectMetadata(key, val string) {
	switch key {
	case "At least one output file must be specified", "Metadata", "":
		return
	default:
		strNum := len(pi.streams)
		switch strNum {
		case 0:
			pi.metadata[key] = val
		default:
			pi.streams[strNum-1].metadata[key] = val
		}
	}
}

func parseDSB(line string) map[string]string {
	parsed := make(map[string]string)
	segments := strings.Split(line, ",")
	for _, seg := range segments {
		dtpts := strings.Split(seg, ": ")
		key := strings.TrimSpace(dtpts[0])
		val := strings.TrimSpace(dtpts[1])
		parsed[key] = val
	}
	return parsed
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
	name = filepath.Base(name)
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
	name = filepath.Base(name)
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
	if key == "Metadata" && val == "" {
		return "", ""
	}
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

func startStrToFl64(start string) float64 {
	startFL, err := strconv.ParseFloat(start, 64)
	if err != nil {
		return math.NaN()
	}

	return startFL
}

func (pi *ParseInfo) parseStreams() {
	//fmt.Println("Parse", len(pi.streams), "streams")
	for i, stream := range pi.streams {
		data := stream.data
		switch {
		default:
			panic("unknown stream type: " + data)
		case strings.Contains(data, "Video:"):
			vs := parseVideoData(pi.streams[i].data)
			vs.metadata = stream.metadata
			pi.Video = append(pi.Video, vs)
		case strings.Contains(data, "Audio:"):
			as := parseAudioData(pi.streams[i].data)
			as.metadata = stream.metadata
			pi.Audio = append(pi.Audio, as)
		case strings.Contains(data, "Data:"):
			dt := parseDataData(pi.streams[i].data)
			dt.metadata = stream.metadata
			pi.Data = append(pi.Data, dt)
		case strings.Contains(data, "Subtitle:"):
			st := parseSubtitleData(pi.streams[i].data)
			st.metadata = stream.metadata
			pi.Subtitles = append(pi.Subtitles, st)
		}
	}

}

func parseDataData(data string) datastream {
	ds := datastream{}
	bra := deBracketSplit(data)
	for i, aud := range bra {
		switch i {
		default:
			ds.coments = append(ds.coments, aud)
		}
	}
	ds.lang = grepLang(data)
	return ds
}

func parseSubtitleData(data string) subtitlestream {
	ss := subtitlestream{}
	bra := deBracketSplit(data)
	for i, aud := range bra {
		switch i {
		default:
			ss.coments = append(ss.coments, aud)
		}
	}
	ss.lang = grepLang(data)
	return ss
}

func parseAudioData(data string) Audiostream {
	as := Audiostream{}
	//0:1 (eng) Audio: pcm_s24le (lpcm / 0x6D63706C), 48000 Hz, 5.1, s32 (24 bit), 6912 kb/s (default)
	//0:1 (rus) Audio: aac (LC), 48000 Hz, stereo, fltp (default)
	//Stream #0:1(rus): Audio: aac (LC) (mp4a / 0x6134706D), 48000 Hz, stereo, fltp, 127 kb/s (default)
	bra := deBracketSplit(data)
	for i, aud := range bra {
		switch i {
		case 0:
			data := strings.Split(aud, "Audio: ")
			as.codec = devtools.AliasByPrefixFromFile(constant.AudioCodecAliasFile(), data[1])

		case 1:
			as.hertz = grepFreq(aud)
		case 2:
			as.channel_layout = aud
		case 3:
			as.sample_format = aud
		case 4:
			dt := strings.Split(aud, " kb/s")
			btrt, err := strconv.Atoi(strings.TrimSpace(dt[0]))
			if err != nil {
				as.warnings = append(as.warnings, "can't parse bitrate: "+aud)
				continue
			}
			as.bitrate = btrt
		default:
			as.warnings = append(as.warnings, "unknown_data: "+aud)
		}
	}
	as.lang = grepLang(data)
	as.assessAudiostream()
	return as
}

func (as *Audiostream) assessAudiostream() {
	switch {
	default:
		as.warnings = append(as.warnings, fmt.Sprintf("channel layout:%v", as.channel_layout))
	case strings.HasPrefix(as.channel_layout, " 5.1"):
		as.channel_layout = "5.1"
	case strings.HasPrefix(as.channel_layout, " mono"):
		as.channel_layout = "mono"
	case strings.HasPrefix(as.channel_layout, " stereo"):
		as.channel_layout = "stereo"
	case strings.HasPrefix(as.channel_layout, " 1 channels"):
		as.channel_layout = "mono"
	}
	if as.bitrate < 80 && as.bitrate != 0 {
		as.warnings = append(as.warnings, fmt.Sprintf("low bitrate: %v kb/s", as.bitrate))
	}
}

func grepWH(data string) (int, int) {
	for _, fld := range strings.Fields(data) {
		r := regexp.MustCompile(`\d\dx\d\d`)
		if r.FindString(fld) == "" {
			continue
		}
		whStr := strings.Split(fld, "x")
		w, _ := strconv.Atoi(whStr[0])
		h, _ := strconv.Atoi(whStr[1])
		return w, h
	}
	return -1, -1
}

func grepBitrate(data string) int {
	r := regexp.MustCompile(`\d* kb`)
	found := r.FindString(data)
	if found == "" {
		return -1
	}
	bt, _ := strconv.Atoi(strings.TrimSuffix(found, " kb"))
	return bt
}

func grepFreq(data string) int {
	r := regexp.MustCompile(`\d* Hz`)
	found := r.FindString(data)
	if found == "" {
		return -1
	}
	bt, _ := strconv.Atoi(strings.TrimSuffix(found, " Hz"))
	return bt
}

func parseVideoData(line string) Videostream {
	vs := Videostream{}
	vs.data = line
	bra := deBracketSplit(line)
	for i, data := range bra {
		switch i {
		case 0:
			dt := strings.Split(data, "Video: ")
			codec := dt[len(dt)-1]
			codec = strings.TrimSpace(codec)
			vs.codecinfo = codec
		case 1:
			vs.pix_fmt = data
		case 2:
			vs.width, vs.height = grepWH(data)
			if strings.Contains(data, "[SAR") {
				vs.sardar = ripBetween(data, "[", "]")
			}
		default:
			if strings.Contains(data, "kb/s") {
				vs.bitrate = data
			}
			if strings.Contains(data, "SAR") {
				if vs.sardar != "" {
					vs.warnings = append(vs.warnings, vs.sardar)
					vs.warnings = append(vs.warnings, data)
				}
				vs.sardar = data
			}
			if strings.Contains(data, " fps") {
				vs.fps = data
			}
			if strings.Contains(data, " tbr") {
				vs.tbr = data
			}
			if strings.Contains(data, " tbn") {
				vs.tbn = data
			}
			if strings.Contains(data, " tbc") {
				vs.tbc = strings.TrimSuffix(data, "\n")
			}
		}
	}
	vs.lang = grepLang(line)
	vs.assessVideostream()
	return vs
}

func (vid *Videostream) assessVideostream() {
	switch vid.fps {
	default:
		vid.warnings = append(vid.warnings, "bad fps:"+strings.TrimSuffix(vid.fps, " fps"))
	case " 25 fps", " 24 fps", " 23.98 fps":
	//case " 24.01 fps", " 24.02 fps", " 29.97 fps", " 24.97 fps", " 30 fps", " 50 fps", " 24.96 fps", " 23.99 fps", " 25.01 fps", " 25.02 fps": //bad
	case "":
		vid.warnings = append(vid.warnings, "no fps detected")
	}
	if strings.TrimSpace(vid.sardar) != "SAR 1:1 DAR 16:9" && strings.TrimSpace(vid.sardar) != "" {
		vid.warnings = append(vid.warnings, "atypical SAR/DAR: '"+vid.sardar+"'")
	}
	switch {
	default:
		vid.warnings = append(vid.warnings, "atypical width/height: "+fmt.Sprintf("%v:%v", vid.width, vid.height))
	case vid.width == 1920 && vid.height == 1080:
	}
}

func deBracketSplit(str string) []string {
	sl := strings.Split(str, "")
	buf := ""
	bracketed := []string{}
	closed := true
	for _, s := range sl {
		switch s {
		case "(", "[":
			closed = false
		case ")", "]":
			closed = true
		case ",":
			if closed {
				bracketed = append(bracketed, buf)
			}
			buf = ""
			continue
		}
		buf += s
	}
	bracketed = append(bracketed, buf)
	return bracketed
}

func grepLang(streamdata string) string {
	pts := strings.Split(streamdata, ":")
	return ripBetween(pts[1], "(", ")")
}

type Videostream struct {
	//  Stream #0:0(und): Video: h264 (High 4:2:2) (avc1 / 0x31637661), yuv422p, 1920x1080 [SAR 1:1 DAR 16:9], 38375 kb/s, 25 fps, 25 tbr, 12800 tbn, 50 tbc (default)
	data      string
	codecinfo string
	pix_fmt   string
	width     int
	height    int
	bitrate   string
	sardar    string
	fps       string
	tbr       string
	tbn       string
	tbc       string
	lang      string
	metadata  map[string]string
	sidedata  string
	warnings  []string
}

func (vs *Videostream) String() string {
	sw, sh := vs.SAR()

	s := fmt.Sprintf("%vx%v [%v:%v]  %v ", vs.width, vs.height, sw, sh, vs.fps)
	return s
}

func (vs *Videostream) Dimentions() (int, int) {
	return vs.width, vs.height
}

func (vs *Videostream) SAR() (int, int) {
	w, h := vs.Dimentions()
	primes := []int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59, 61, 67, 71, 73, 79, 83, 89, 97}
	sort.Reverse(sort.IntSlice(primes))
	for _, prime := range primes {
		check := true
		for check {
			if w%prime == 0 && h%prime == 0 {
				w, h = w/prime, h/prime
			} else {
				check = false
			}
		}
	}
	return w, h
}

func (vs *Videostream) DAR() float64 {
	return float64(vs.width) / float64(vs.height)
}

/*
Video: 1920x1080; 25 fps
*/

type Audiostream struct {
	codec          string
	hertz          int
	channel_layout string
	sample_format  string
	bitrate        int
	lang           string
	metadata       map[string]string
	warnings       []string
}

func (a *Audiostream) String() string {
	s := fmt.Sprintf("%v", a.codec)
	if a.channel_layout != "" {
		s += "  " + a.channel_layout
	}
	if a.bitrate != 0 {
		s += "  " + fmt.Sprintf("bitrate: %v", a.bitrate)
	}

	return s
}

type datastream struct {
	coments  []string
	lang     string
	metadata map[string]string
}

type subtitlestream struct {
	coments  []string
	lang     string
	metadata map[string]string
}

func ripBetween(data, open, close string) string {
	d := strings.Split(data, "")
	opened := false
	buf := ""
	for _, s := range d {
		switch s {
		case open:
			opened = true
		case close:
			opened = false
		default:
			if opened {
				buf += s
			}
		}
	}
	return buf
}

/*
Stream #0:0(und): Video: h264 (High 4:2:2) (avc1 / 0x31637661), yuv422p, 1920x1080 [SAR 1:1 DAR 16:9], 25333 kb/s, 25 fps, 25 tbr, 12800 tbn, 50 tbc (default)

000
111
101

10N        3

101  100
 01   00   5
      00   7
	  00
	  00
	  00
2

2
1
3


*/
