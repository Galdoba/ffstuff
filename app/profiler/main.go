package main

import (
	"encoding/json"
	"fmt"
	"os"
)

/*
TZ:
Утилита должна оценить файл и составить его рабочий Профиль

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
