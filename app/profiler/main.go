package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

/*
TZ:
Утилита должна оценить файл и составить его рабочий Профиль

//
Input #0, mov,mp4,m4a,3gp,3g2,mj2, from 'Barri_4s_treyler_a_teka.mp4':
  Duration: 00:01:26.08, start: 0.000000, bitrate: 15338 kb/s
    Stream #0:0(eng): Video: h264 (Main) (avc1 / 0x31637661), yuv420p, 1920x1080 [SAR 1:1 DAR 16:9], 14831 kb/s, 25 fps, 25 tbr, 25k tbn, 50 tbc (default)
    Stream #0:1(eng): Audio: aac (LC) (mp4a / 0x6134706D), 48000 Hz, stereo, fltp, 317 kb/s (default)


//
[1100]
[1{#1920x1080#25#5822#[]#[[SAR 1:1 DAR 16:9]]#5.3%};2{#5.1#48.0#341}{#stereo#48.0#129};0;1]

ffprobe -v quiet -of json -show_format -show_streams Barri_4s_treyler_a_teka.mp4

Структура Профиля:
краткий:
[ABCDE-F], где
a = количество видеопотоков
b = количество аудиопотоков
c = количество потоков данных
d = количество srt потоков
e = количество хз чего
f = количество замечаний от библиотеки парса

Сепараторы:
- = логический раздел (Абзац)
; = потоки (Строка)
# = данные внутри потоков (Слово)

развернутый:
[A#a#b#c#d#e#f;B#g#h#i;C          ], где
A = количество видеопотоков     формат: eHex                (eHex = целое число с базой 32)          int
a = размер видео                формат: IxI                 (I = целое число)                        int
b = fps видео                   формат: F                   (F = десятичное число)                   float64
c = битрейт                     формат: I                   (I = целое число)                        int
d = SARDAR внешний              формат: [SAR I:I DAR I:I]   (I = целое число)                        []int
e = SARDAR внутренний           формат: [d]                 (d = формат внутреннего SARDAR)          []int
f = интерлейс                   формат: F                   (F = средний % не прогрессивных кадров)  float64
B = количество аудиопотоков     формат: eHex                (eHex = целое число с базой 32)          int
g = раскладка каналов           формат: S                   (S = текст ключ-значение по таблице)     string
h = герцовка                    формат: F                   (F = десятичное число kHz)               float64
i = битрейт                     формат: I                   (I = целое число значение kbit/s)        int
C = количество потоков данных   формат: eHex                (eHex = целое число с базой 32)          int??????????






e = битрейт      формат: [I]   (I = целое число)
d = количество srt потоков
e = количество хз чего
f = количество замечаний от библиотеки парса

Определения:
1.Профиль - форматированная информация о внутренней медиа структуре файла.
2.Аргументы = файлы к которым нужно составить профайл


*/

func main() {
	fmt.Println("Here be profiler")
}

func newProfile() *Profile {
	sr := &Profile{}
	//	sr.Format = &Format{}
	return sr
}

func ConsumeJSON(path string) (*Profile, error) {
	if !strings.HasSuffix(path, ".json") {
		return nil, fmt.Errorf("file is nit json")
	}
	sr := &Profile{}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("can't read json: %v", err.Error())
	}
	err = json.Unmarshal(data, &sr)
	if err != nil {
		return nil, fmt.Errorf("can't unmarshal json: %v (%v)", err.Error(), path)
	}
	err = sr.validate()
	if err != nil {
		return sr, fmt.Errorf("validation error: %v", err.Error())
	}
	return sr, nil
}

func assertNoError(err error) {
	if err != nil {
		panic(err.Error())
	}
}

type Profile struct {
	Format   *Format   `json:"format"`
	Streams  []*Stream `json:"streams,omitempty"`
	warnings []string
	short    string
	long     string
}

func (p *Profile) validate() error {
	vSNum := 0
	aSNum := 0
	dSNum := 0
	sSNum := 0
	warnings := 0
	for _, stream := range p.Streams {
		switch stream.Codec_type {
		default:
			return fmt.Errorf("unknown codec type: %v", stream.Codec_type)
		case "video":

			vSize := fmt.Sprintf("%vx%v", stream.Width, stream.Height)
			switch vSize {
			case "720x576", "1920x1080", "3840x2160":
			default:
				p.warnings = append(p.warnings, fmt.Sprintf("stream [0:v:%v] ==> bad video size: [%v]", vSNum, vSize))
			}
			fps := stream.R_frame_rate
			switch fps {
			case "24/1", "25/1", "24000/1001":
			case "2997/125": //about valid
				fps = "24000/1001"
			default:
				p.warnings = append(p.warnings, fmt.Sprintf("stream [0:v:%v] ==> bad fps: [%v]", vSNum, fps))
			}
			sardar := "SAR=" + stream.Sample_aspect_ratio + " DAR=" + stream.Display_aspect_ratio
			switch sardar {
			case "SAR=1:1 DAR=16:9", "":
			case "SAR= DAR=":
				p.warnings = append(p.warnings, fmt.Sprintf("stream [0:v:%v] ==> blank SAR DAR data present", vSNum))
				warnings++
			case "SAR=1:1 DAR=1024:429", "SAR=1:1 DAR=37:20", "SAR=1:1 DAR=160:67":
				//p.warnings = append(p.warnings, fmt.Sprintf("stream [0:v:%v] ==> bad (but possible) SAR DAR: [%v]", vSNum, sardar))
			default:
				p.warnings = append(p.warnings, fmt.Sprintf("stream [0:v:%v] ==> bad SAR DAR: [%v]", vSNum, sardar))
			}
			btrate := stream.Bit_rate
			switch btrate {
			case "":
				p.warnings = append(p.warnings, fmt.Sprintf("stream [0:v:%v] ==> blank bitRate data present", vSNum))
			default:
				btr, err := strconv.Atoi(btrate)
				if err != nil {
					panic(err.Error())
				}
				if btr == 0 {
					panic(fmt.Sprintf("stream [0:v:%v] ==> bitRate: [%v]", vSNum, btr))
				}
				//p.warnings = append(p.warnings, fmt.Sprintf("stream [0:v:%v] ==> bitRate: [%v]", vSNum, btr))
			}
			vSNum++
		case "audio":
			aSNum++
		case "data":
			dSNum++
		case "subtitle":
			sSNum++
		}
	}
	if vSNum > 1 {
		p.warnings = append(p.warnings, fmt.Sprintf("file ==> %v video streams present", vSNum))
	}
	p.short = fmt.Sprintf("%v%v%v%v-%v", ehex(vSNum), ehex(aSNum), ehex(dSNum), ehex(sSNum), ehex(len(p.warnings)))
	return nil
}

func ehex(i int) string {
	switch i {
	default:
		return "?"
	case 0, 1, 2, 3, 4, 5, 6, 7, 8, 9:
		return fmt.Sprintf("%v", i)
	case 10:
		return "A"
	case 11:
		return "B"
	case 12:
		return "C"
	case 13:
		return "D"
	case 14:
		return "E"
	case 15:
		return "F"
	case 16:
		return "G"
	case 17:
		return "H"
	case 18:
		return "J"
	case 19:
		return "K"
	case 20:
		return "L"
	case 21:
		return "M"
	case 22:
		return "N"
	case 23:
		return "P"
	case 24:
		return "Q"
	case 25:
		return "R"
	case 26:
		return "S"
	case 27:
		return "T"
	case 28:
		return "U"
	case 29:
		return "V"
	case 30:
		return "W"
	case 31:
		return "X"
	case 32:
		return "Y"
	case 33:
		return "Z"
	}
}

func (p *Profile) Warnings() []string {
	return p.warnings
}

func (p *Profile) Short() string {
	return p.short
}

// func (sr *Profile) String() string {
// 	s := "file: " + sr.Format.Filename + "\n"
// 	s += fmt.Sprintf("streams total: %v", len(sr.Streams))
// 	return s
// }

type Format struct {
	Bit_rate         string            `json:"bit_rate,omitempty"`
	Duration         string            `json:"duration,omitempty"`
	Filename         string            `json:"filename,omitempty"`
	Format_long_name string            `json:"format_long_name,omitempty"`
	Format_name      string            `json:"format_name,omitempty"`
	Nb_programs      int               `json:"nb_programs,omitempty"`
	Nb_streams       int               `json:"nb_streams,omitempty"`
	Probe_score      int               `json:"probe_score,omitempty"`
	Size             string            `json:"size,omitempty"`
	Start_time       string            `json:"start_time,omitempty"`
	Tags             map[string]string `json:"tags,omitempty"`
}

type Stream struct {
	Avg_frame_rate       string                  `json:"avg_frame_rate,omitempty"`
	Bit_rate             string                  `json:"bit_rate,omitempty"`
	Bits_per_raw_sample  string                  `json:"bits_per_raw_sample,omitempty"`
	Bits_per_sample      int                     `json:"bits_per_sample,omitempty"`
	Channel_layout       string                  `json:"channel_layout,omitempty"`
	Channels             int                     `json:"channels,omitempty"`
	Chroma_location      string                  `json:"chroma_location,omitempty"`
	Closed_captions      int                     `json:"closed_captions,omitempty"`
	Codec_long_name      string                  `json:"codec_long_name,omitempty"`
	Codec_name           string                  `json:"codec_name,omitempty"`
	Codec_tag            string                  `json:"codec_tag,omitempty"`
	Codec_tag_string     string                  `json:"codec_tag_string,omitempty"`
	Codec_time_base      string                  `json:"codec_time_base,omitempty"`
	Codec_type           string                  `json:"codec_type,omitempty"`
	Coded_height         int                     `json:"coded_height,omitempty"`
	Coded_width          int                     `json:"coded_width,omitempty"`
	Color_primaries      string                  `json:"color_primaries,omitempty"`
	Color_range          string                  `json:"color_range,omitempty"`
	Color_space          string                  `json:"color_space,omitempty"`
	Color_transfer       string                  `json:"color_transfer,omitempty"`
	Display_aspect_ratio string                  `json:"display_aspect_ratio,omitempty"`
	Divx_packed          string                  `json:"divx_packed,omitempty"`
	Dmix_mode            string                  `json:"dmix_mode,omitempty"`
	Duration             string                  `json:"duration,omitempty"`
	Duration_ts          int                     `json:"duration_ts,omitempty"`
	Field_order          string                  `json:"field_order,omitempty"`
	Has_b_frames         int                     `json:"has_b_frames,omitempty"`
	Height               int                     `json:"height,omitempty"`
	Id                   string                  `json:"id,omitempty"`
	Index                int                     `json:"index,omitempty"`
	Is_avc               string                  `json:"is_avc,omitempty"`
	Level                int                     `json:"level,omitempty"`
	Loro_cmixlev         string                  `json:"loro_cmixlev,omitempty"`
	Loro_surmixlev       string                  `json:"loro_surmixlev,omitempty"`
	Ltrt_cmixlev         string                  `json:"ltrt_cmixlev,omitempty"`
	Ltrt_surmixlev       string                  `json:"ltrt_surmixlev,omitempty"`
	Max_bit_rate         string                  `json:"max_bit_rate,omitempty"`
	Nal_length_size      string                  `json:"nal_length_size,omitempty"`
	Nb_frames            string                  `json:"nb_frames,omitempty"`
	Pix_fmt              string                  `json:"pix_fmt,omitempty"`
	Profile              string                  `json:"profile,omitempty"`
	Quarter_sample       string                  `json:"quarter_sample,omitempty"`
	R_frame_rate         string                  `json:"r_frame_rate,omitempty"`
	Refs                 int                     `json:"refs,omitempty"`
	Sample_aspect_ratio  string                  `json:"sample_aspect_ratio,omitempty"`
	Sample_fmt           string                  `json:"sample_fmt,omitempty"`
	Sample_rate          string                  `json:"sample_rate,omitempty"`
	Start_pts            int                     `json:"start_pts,omitempty"`
	Start_time           string                  `json:"start_time,omitempty"`
	Time_base            string                  `json:"time_base,omitempty"`
	Width                int                     `json:"width,omitempty"`
	Side_data_list       []Side_data_list_struct `json:"side_data_list,omitempty"`
	Tags                 map[string]string       `json:"tags,omitempty"`
	Disposition          map[string]int          `json:"disposition,omitempty"`
}

type Side_data_list_struct struct {
	Side_data map[string]string
}

/*



_streams___disposition 20117
_streams___disposition____attached_pic 20117
_streams___disposition____clean_effects 20117
_streams___disposition____comment 20117
_streams___disposition____default 20117
_streams___disposition____dub 20117
_streams___disposition____forced 20117
_streams___disposition____hearing_impaired 20117
_streams___disposition____karaoke 20117
_streams___disposition____lyrics 20117
_streams___disposition____original 20117
_streams___disposition____timed_thumbnails 23343
_streams___disposition____visual_impaired 20117
*/
