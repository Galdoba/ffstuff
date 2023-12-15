package main

import (
	"encoding/json"
	"fmt"
	"os"
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
	concept()
}

func concept() {
	sr := &ScanResult{}
	sr.Frmt = &Format{Filename: "testName", Nb_streams: 5}
	data, err := os.ReadFile("/home/galdoba/workbench/flscnr/ffprf/scan.json")
	//data0, err2 := json.MarshalIndent(&sr, "", "  ")

	assertNoError(err)
	if len(data) == 0 {
		data, err = json.MarshalIndent(&sr, "", "  ")
		if err != nil {
			println(err.Error())
			//os.Exit(1)
		}
	}
	err = json.Unmarshal(data, &sr)
	if err != nil {
		errText := fmt.Sprintf("can't unmarshal config data: %v", err.Error())
		println(errText)
		os.Exit(1)
	}
	//fmt.Println(string(data))
	//fmt.Println(string(data0))
	// if err2 != nil {
	// 	fmt.Println(err.Error())
	// }
	fmt.Println(sr.Frmt.Filename)
	fmt.Println(sr.Streams[2].Codec_name)
	fmt.Println(sr.String())
}

func assertNoError(err error) {
	if err != nil {
		panic(err.Error())
	}
}

type ScanResult struct {
	//TestStr `json:"Test"`
	Frmt    *Format   `json:"format"`
	Streams []*Stream `json:"streams,omitempty"`
}

func (sr *ScanResult) String() string {
	s := "file: " + sr.Frmt.Filename + "\n"
	s += fmt.Sprintf("streams total: %v", len(sr.Streams))
	return s
}

type Format struct {
	Filename         string            `json:"filename,omitempty"`
	Nb_streams       int               `json:"nb_streams,omitempty"`
	Nb_programs      int               `json:"nb_programs,omitempty"`
	Format_name      string            `json:"format_name,omitempty"`
	Format_long_name string            `json:"format_long_name,omitempty"`
	Start_time       string            `json:"start_time,omitempty"`
	Duration         string            `json:"duration,omitempty"`
	Size             string            `json:"size,omitempty"`
	Bit_rate         string            `json:"bit_rate,omitempty"`
	Probe_score      int               `json:"probe_score,omitempty"`
	Tags             map[string]string `json:"tags,omitempty"`
}

type Stream struct {
	Index               int    `json:"index"`
	Codec_name          string `json:"codec_name"`
	Codec_long_name     string `json:"codec_long_name"`
	Profile             string `json:"profile"`
	Codec_type          string `json:"codec_type"`
	Codec_time_base     string `json:"codec_time_base"`
	Codec_tag_string    string `json:"codec_tag_string"`
	Codec_tag           string `json:"codec_tag"`
	Width               int    `json:"width"`
	Height              int    `json:"height"`
	Codec_width         int    `json:"coded_width"`
	Codec_height        int    `json:"coded_height"`
	Closed_captions     int    `json:"closed_captions"`
	Has_b_frames        int    `json:"has_b_frames"`
	SAR                 string `json:"sample_aspect_ratio"`
	DAR                 string `json:"display_aspect_ratio"`
	Pix_fmt             string `json:"pix_fmt"`
	Level               int    `json:"level"`
	Color_Range         string `json:"color_range"`
	Color_space         string `json:"color_space"`
	Field_order         string `json:"field_order"`
	Refs                int    `json:"refs"`
	R_frame_rate        string `json:"r_frame_rate"`
	Avg_frame_rate      string `json:"avg_frame_rate"`
	Time_base           string `json:"time_base"`
	Start_pts           int    `json:"start_pts"`
	Start_time          string `json:"start_time"`
	Duration_ts         int    `json:"duration_ts"`
	Duration            string `json:"duration"`
	Bit_rate            string `json:"bit_rate"`
	Bits_per_raw_sample string `json:"bits_per_raw_sample"`
	Nb_frames           string `json:"nb_frames"`
}
