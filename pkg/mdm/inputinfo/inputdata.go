package inputinfo

import (
	"fmt"
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

// inputinfo.Parse(filepath string) (*Info, error)

const (
	stage_ParseScanTime = iota
	stage_ParseFilename
	stage_ParseGlobalMeta
	stage_ParseDuration
	stage_ParseStreams
	scanTimePrefix = "Started: "
)

func parse(input inputdata) (*parseInfo, error) {
	pi := parseInfo{}
	pi.metadata = make(map[string]string)
	for _, line := range input.data {
		switch pi.parseStage {
		case stage_ParseScanTime:
			//fmt.Printf("Stage %v: Line %v| \n%v\n", pi.parseStage, ln, pi)
			pi.parseScanTime(line)
		case stage_ParseFilename:
			//fmt.Printf("Stage %v: Line %v| \n%v\n", pi.parseStage, ln, pi)
			pi.parseName(line)
		case stage_ParseGlobalMeta:
			//fmt.Printf("Stage %v: Line %v| \n%v\n", pi.parseStage, ln, pi)
			if err := pi.parseMetaGlobal(line); err != nil {
				return &pi, err
			}
		case stage_ParseDuration:
			//fmt.Printf("Stage %v: Line %v| \n%v\n", pi.parseStage, ln, pi)
			//panic("must not happen")
			if err := pi.parseDuration(line); err != nil {
				return &pi, err
			}
		case stage_ParseStreams:
			//fmt.Printf("Stage %v: Line %v| \n", pi.parseStage, ln)
		}
		pi.parsedLines++
	}
	return &pi, nil
}

//parseScantime - ищет время сканирования. Опционально.
//#запускать если данные уже получены
func parseScantime(line string) string {
	scanTime := strings.TrimPrefix(line, scanTimePrefix)
	return scanTime
}

func (pi *parseInfo) parseScanTime(line string) {
	scanTime := parseScantime(line)
	switch scanTime {
	default:
		pi.scanTime = scanTime
		pi.parseStage = stage_ParseFilename
	case "":
		pi.parseName(line)
	}
}

func parseName(line string) string {
	switch {
	case strings.Contains(line, ffmpegInputPrefix): //ffmpeg
		name := ripBetween(line, " '", "':")
		return name
	case strings.Contains(line, ffliteInputPrefix): //fflite
		line = strings.TrimSpace(line)
		name := strings.TrimPrefix(line, ffliteInputPrefix)
		return name
	default: //No input Prefix
		return ""
	}
}

func (pi *parseInfo) parseName(line string) {
	name := parseName(line)
	if name != "" {
		pi.filename = name
		pi.parseStage = stage_ParseGlobalMeta
	}
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

func (pi *parseInfo) parseMetaGlobal(line string) error {
	if strings.HasPrefix(line, "  Duration:") {
		pi.parseStage = stage_ParseDuration
		return pi.parseDuration(line)
	}
	key, val := parseMetaGlobal(line)
	switch {
	case key+val == "":
	case key == "" && val != "":
		return fmt.Errorf("no key but have value")
	default:
		pi.metadata[key] = val
	}
	return nil
}

func parseDuration(line string) (duration float64, start float64, bitrate int, err error) {
	data := strings.Split(line, ", ")
	hms := 0
	for i, val := range data {
		val = strings.TrimSpace(val)
		switch i {
		case 0:
			d := strings.Split(val, ": ")
			timeStr := strings.ReplaceAll(d[1], ":", ".")
			timeSegments := strings.Split(timeStr, ".")
			for j, ts := range timeSegments {
				switch j {
				case 0:
					hms, err = strconv.Atoi(ts)
					duration += 3600.0 * float64(hms)
				case 1:
					hms, err = strconv.Atoi(ts)
					duration += 60.0 * float64(hms)
				case 2:
					hms, err = strconv.Atoi(ts)
					duration += 1.0 * float64(hms)
				case 3:
					hms, err = strconv.Atoi(ts)
					duration += 0.001 * float64(hms)
				}
				if err != nil {
					return
				}
			}
		case 1:
			d := strings.Split(val, ": ")
			start, err = strconv.ParseFloat(d[1], 64)
		case 2:
			d := strings.Split(val, ": ")
			bitrt := strings.Fields(d[1])
			bitrate, err = strconv.Atoi(bitrt[0])
		}
		if err != nil {
			return
		}
	}
	return
}

func (pi *parseInfo) parseDuration(line string) error {
	duration, start, bitrate, err := parseDuration(line)
	pi.duration = duration
	pi.start = start
	pi.globBitrate = bitrate
	pi.parseStage = stage_ParseStreams
	return err
}

func parseStreamData(line string) []string {
	stdt := []string{}
	if !strings.Contains(line, " Video: ") {
		return nil
	}
	//  Stream #0:0: Video: mpeg2video (4:2:2), yuv422p(tv, unknown/bt709/bt709, progressive), 1920x1080 [SAR 1:1 DAR 16:9], 50000 kb/s, 25 fps, 25 tbr, 25 tbn, 50 tbc
	//    0:0 (und) Video: prores (HQ) (apch / 0x68637061), yuv422p10le(tv, bt709, progressive), 1920x1080, 165748 kb/s, SAR 1:1 DAR 16:9, 23.98 fps, 23.98 tbr, 24k tbn, 24k tbc (default)
	//отделяем анкер
	//mpeg2video (4:2:2), yuv422p(tv, unknown/bt709/bt709, progressive), 1920x1080 [SAR 1:1 DAR 16:9], 50000 kb/s, 25 fps, 25 tbr, 25 tbn, 50 tbc
	//prores (HQ) (apch / 0x68637061), yuv422p10le(tv, bt709, progressive), 1920x1080, 165748 kb/s, SAR 1:1 DAR 16:9, 23.98 fps, 23.98 tbr, 24k tbn, 24k tbc (default)
	//исключаем скобки
	//mpeg2video , yuv422p, 1920x1080 , 50000 kb/s, 25 fps, 25 tbr, 25 tbn, 50 tbc
	//prores  , yuv422p10le, 1920x1080, 165748 kb/s, SAR 1:1 DAR 16:9, 23.98 fps, 23.98 tbr, 24k tbn, 24k tbc

	return stdt
}

func splitVideoStreamGroups(line string) []string {
	stdt := []string{}
	if !strings.Contains(line, " Video: ") {
		continue
	}
	return nil
}

/*
parserSearch:
error		- opt //sdfsdf.asda: Invalid data found when processing input
started		- opt
name		- mand
metadata	- opt
streams		- mand
*/

const (
	ffliteInputPrefix = "INPUT 0: "
	ffmpegInputPrefix = "Input #0, "
)

func ripBetween(line, start, close string) string {
	parts := strings.Split(line, start)
	if len(parts) < 2 {
		return ""
	}
	startless := strings.Join(parts[1:], "")
	parts2 := strings.Split(startless, close)
	result := strings.Join(parts2[:len(parts2)-1], "")
	return result
}

/*
func collectVideoData(feed string) videostream {
	vs := videostream{}
	data := strings.Split(feed, ",")
	for _, dataPart := range data {
		if w, h := extractDimention(dataPart); w != 0 && h != 0 {
			vs.width, vs.height = w, h
		}
		for _, ex := range []extracted{
			extractSar(dataPart),
			extractDAR(dataPart),
			extractFPS(dataPart),
			extractTBR(dataPart),
			extractTBN(dataPart),
			extractTBC(dataPart),
			extractBITRATE(dataPart),
			extractCaptionsInfo(dataPart),
		} {
			if ex.val == "" {
				continue
			}
			if !ex.norm {
				vs.warnings = append(vs.warnings, "abnormal_"+ex.key+": "+strings.TrimSpace(ex.val))
			}
			switch ex.key {
			case "sar":
				vs.sar = ex.val
			case "dar":
				vs.dar = ex.val
			case "fps":
				vs.fps = ex.val
			case "tbr":
				vs.tbr = ex.val
			case "tbn":
				vs.tbn = ex.val
			case "tbc":
				vs.tbc = ex.val
			case "cc":
				vs.cc = ex.val
			case "bitrate":
				vs.bitrate, _ = strconv.Atoi(ex.val)
			}
		}

	}
	return vs
}

func extractDimention(str string) (int, int) {
	str = strings.TrimSpace(str)
	strArr := strings.Split(str, " ")
	width := 0
	height := 0
	for _, s := range strArr {
		dim := strings.Split(s, "x")
		if len(dim) != 2 {
			continue
		}
		w, errW := strconv.Atoi(dim[0])
		if errW != nil {
			continue
		}
		h, errH := strconv.Atoi(dim[1])
		if errH != nil {
			continue
		}
		width = w
		height = h
	}
	return width, height
}

func extractSar(str string) extracted {
	key := "sar"
	if strings.Contains(str, "SAR") {
		validSARs := []string{"1:1", "4:3", "64:45"}
		for _, sar := range validSARs {
			if strings.Contains(str, "SAR "+sar) {
				return extracted{key, sar, true}
			}
		}
		return extracted{key, str, false}
	}
	return extracted{key, "", true}
}

func extractDAR(str string) extracted {
	key := "dar"
	if strings.Contains(str, "DAR ") {
		pt := strings.Split(str, " ")
		for i, p := range pt {
			if p == "DAR" {
				foundDar := strings.TrimSuffix(pt[i+1], "]")
				if foundDar == "16:9" {
					return extracted{key, foundDar, true}
				} else {
					return extracted{key, foundDar, false}
				}
			}
		}
	}
	return extracted{key, "", true}
}

func extractFPS(str string) extracted {
	key := "fps"
	if strings.Contains(str, " fps") {
		switch strings.TrimSpace(str) {
		case "25 fps", "24 fps", "23.98 fps":
			return extracted{key, strings.TrimSuffix(str, " fps"), true}
		default:
			return extracted{key, strings.TrimSuffix(str, " fps"), false}
		}
	}
	return extracted{key, "", true}
}

func extractTBR(str string) extracted {
	key := "tbr"
	if strings.Contains(str, " tbr") {
		switch strings.TrimSpace(str) {
		case "25 tbr", "24 tbr", "23.98 tbr":
			return extracted{key, strings.TrimSuffix(str, " tbr"), true}
		default:
			return extracted{key, strings.TrimSuffix(str, " tbr"), false}
		}

	}
	return extracted{key, "", true}
}

func extractTBN(str string) extracted {
	key := "tbn"
	if strings.Contains(str, " tbn") {
		switch strings.TrimSpace(str) {
		//case "12800 tbn":
		//	return extracted{key, strings.TrimSuffix(str, " tbn"), true}
		default:
			return extracted{key, strings.TrimSuffix(str, " tbn"), true}
		}

	}
	return extracted{key, "", true}
}

func extractBITRATE(str string) extracted {
	key := "bitrate"
	if strings.Contains(str, "kb/s") {
		switch strings.TrimSpace(str) {
		//case "12800 tbn":
		//	return extracted{key, strings.TrimSuffix(str, " tbn"), true}
		default:
			return extracted{key, strings.TrimSpace(strings.TrimSuffix(str, "kb/s")), true}
		}

	}
	return extracted{key, "", true}
}

func extractTBC(str string) extracted {
	key := "tbc"
	if strings.Contains(str, " tbc") {
		switch strings.TrimSpace(str) {
		//case "12800 tbc":
		//	return extracted{key, strings.TrimSuffix(str, " tbc"), true}
		default:
			return extracted{key, strings.TrimSuffix(str, " tbc"), true}
		}

	}
	return extracted{key, "", true}
}
*/

func containsAll(fullString string, sub ...string) bool {
	for _, s := range sub {
		if !strings.Contains(fullString, s) {
			return false
		}
	}
	return true
}

/*
Input #0, mov,mp4,m4a,3gp,3g2,mj2, from 'SeaHorses_FTR_1080p_RU-XX_20_24fps_PR422HQ.mov'

> ffmpeg -i S_LUBOVYU_I_YAROSTYU_235_24FPS_RUS51.mov
  INPUT 0: S_LUBOVYU_I_YAROSTYU_235_24FPS_RUS51.mov
  Duration: 01:56:46.00, start: 0.000000, bitrate: 173412 kb/s
    0:0 (eng) Video: prores (HQ) (apch / 0x68637061), yuv422p10le(tv, bt709/unknown/unknown, progressive), 1920x1080, 168410 kb/s, SAR 1:1 DAR 16:9, 24 fps, 24 tbr, 24 tbn, 24 tbc (default)
    0:1 (eng) Data: none (tmcd / 0x64636D74), 0 kb/s (default)
      handler_name    : Apple Time Code Media Handler
    0:2 (eng) Audio: pcm_s16le (sowt / 0x74776F73), 48000 Hz, 5.1(side), s16, 4608 kb/s (default)
    S_lyubovyu_i_yarostyu
$ ffmpeg -i Spletnica_s02e07__PRT230106115407_SER_00087_18.mp4
ffmpeg version 4.4.1-full_build-www.gyan.dev Copyright (c) 2000-2021 the FFmpeg developers
  built with gcc 11.2.0 (Rev1, Built by MSYS2 project)
  configuration: --enable-gpl --enable-version3 --enable-static --disable-w32threads --disable-autodetect --enable-fontconfig --enable-iconv --enable-gnutls --enable-libxml2 --enable-gmp --enable-lzma --enable-libsnappy --enable-zlib --enable-librist --enable-libsrt --enable-libssh --enable-libzmq --enable-avisynth --enable-libbluray --enable-libcaca --enable-sdl2 --enable-libdav1d --enable-libzvbi --enable-librav1e --enable-libsvtav1 --enable-libwebp --enable-libx264 --enable-libx265 --enable-libxvid --enable-libaom --enable-libopenjpeg --enable-libvpx --enable-libass --enable-frei0r --enable-libfreetype --enable-libfribidi --enable-libvidstab --enable-libvmaf --enable-libzimg --enable-amf --enable-cuda-llvm --enable-cuvid --enable-ffnvcodec --enable-nvdec --enable-nvenc --enable-d3d11va --enable-dxva2 --enable-libmfx --enable-libglslang --enable-vulkan --enable-opencl --enable-libcdio --enable-libgme --enable-libmodplug --enable-libopenmpt --enable-libopencore-amrwb --enable-libmp3lame --enable-libshine --enable-libtheora --enable-libtwolame --enable-libvo-amrwbenc --enable-libilbc --enable-libgsm --enable-libopencore-amrnb --enable-libopus --enable-libspeex --enable-libvorbis --enable-ladspa --enable-libbs2b --enable-libflite --enable-libmysofa --enable-librubberband --enable-libsoxr --enable-chromaprint
  libavutil      56. 70.100 / 56. 70.100
  libavcodec     58.134.100 / 58.134.100
  libavformat    58. 76.100 / 58. 76.100
  libavdevice    58. 13.100 / 58. 13.100
  libavfilter     7.110.100 /  7.110.100
  libswscale      5.  9.100 /  5.  9.100
  libswresample   3.  9.100 /  3.  9.100
  libpostproc    55.  9.100 / 55.  9.100
Input #0, mov,mp4,m4a,3gp,3g2,mj2, from 'Spletnica_s02e07__PRT230106115407_SER_00087_18.mp4':
  Metadata:
    major_brand     : isom
    minor_version   : 512
    compatible_brands: isomiso2avc1mp41
    encoder         : Lavf58.45.100
  Duration: 00:54:07.02, start: 0.000000, bitrate: 8376 kb/s
  Stream #0:0(und): Video: h264 (High 4:2:2) (avc1 / 0x31637661), yuv422p, 1920x1080 [SAR 1:1 DAR 16:9], 7897 kb/s, 25 fps, 25 tbr, 12800 tbn, 50 tbc (default)
    Metadata:
      handler_name    : VideoHandler
      vendor_id       : [0][0][0][0]
      timecode        : 01:00:00:00
  Stream #0:1(rus): Audio: aac (LC) (mp4a / 0x6134706D), 48000 Hz, 5.1, fltp, 341 kb/s (default)
    Metadata:
      handler_name    : SoundHandler
      vendor_id       : [0][0][0][0]
  Stream #0:2(eng): Audio: aac (LC) (mp4a / 0x6134706D), 48000 Hz, stereo, fltp, 128 kb/s
    Metadata:
      handler_name    : SoundHandler
      vendor_id       : [0][0][0][0]
  Stream #0:3(eng): Data: none (tmcd / 0x64636D74)
    Metadata:
      handler_name    : TimeCodeHandler
      timecode        : 01:00:00:00


*/

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
