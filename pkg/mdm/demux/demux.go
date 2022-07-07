package demux

import (
	"fmt"
	"strings"

	"github.com/Galdoba/devtools/cli/command"
	"github.com/malashin/ffinfo"
)

type Demuxer interface {
	Demux() error //Демуксит файл исходя из заданных параметров
}

/*
fmpeg -i Сквозь_огонь_Through_the_fire.mkv
  INPUT 0: Сквозь_огонь_Through_the_fire.mkv
  Duration: 01:52:36.12, start: 0.000000, bitrate: 18720 kb/s
    0:0 (eng) Video: h264 (Main), yuv420p(progressive), 1920x1080 [SAR 1:1 DAR 16:9], 25 fps, 25 tbr, 1k tbn, 50 tbc (default)
    0:1 (rus) Audio: mp3, 48000 Hz, stereo, fltp, 320 kb/s (default)
    0:2 (rus) Audio: ac3, 48000 Hz, 5.1(side), fltp, 448 kb/s
    0:3 (fre) Audio: mp3, 48000 Hz, stereo, fltp, 320 kb/s
    0:4 (fre) Audio: ac3, 48000 Hz, 5.1(side), fltp, 448 kb/s
*/
/*
clear
&& mkdir -p /mnt/aakkulov/ROOT/IN/_MEGO_DISTRIBUTION/_DONE/Skvoz_ogon
&& mkdir -p /mnt/aakkulov/ROOT/EDIT/_mego_distribushn/
&& mv /home/aakkulov/IN/Сквозь_огонь_Through_the_fire.mkv /home/aakkulov/IN/_IN_PROGRESS/
&& fflite -r 25 -i /home/aakkulov/IN/_IN_PROGRESS/Сквозь_огонь_Through_the_fire.mkv
	-filter_complex "[0:a:1]aresample=48000,atempo=25/(25)[arus]"
	-map [arus] -c:a alac -compression_level 0 -map_metadata -1 -map_chapters -1 /mnt/aakkulov/ROOT/EDIT/_mego_distribushn/Skvoz_ogon_AUDIORUS51.m4a
	-map 0:v:0 -c:v libx264 -preset medium -crf 16 -pix_fmt yuv420p -g 0 -map_metadata -1 -map_chapters -1 /mnt/aakkulov/ROOT/EDIT/_mego_distribushn/Skvoz_ogon_HD.mp4
&& touch /mnt/aakkulov/ROOT/EDIT/_mego_distribushn/Skvoz_ogon.ready
&& mv /home/aakkulov/IN/_IN_PROGRESS/Сквозь_огонь_Through_the_fire.mkv /home/aakkulov/IN/_DONE/
&& at now + 10 hours <<< "mv /home/aakkulov/IN/_DONE/Сквозь_огонь_Through_the_fire.mkv /mnt/aakkulov/ROOT/IN/_MEGO_DISTRIBUTION/_DONE/Skvoz_ogon"
&& clear
&& touch /home/aakkulov/IN/TASK_COMPLETE_Сквозь_огонь_Through_the_fire.mkv.txt
*/

type demuxValuesPreset struct {
	manual            bool
	yadif             bool
	downscale         bool
	crf               int
	atempo            int
	asetrate          int
	acodec            string
	compression_level string
	map_metadata      string
	map_chapters      string
}

//   decode.Video(path string) (string, error)

/*
АЛГОРИТМ
Узнать базу имени для выходных файлов
Узнать путь куда отправить выходные файлы - ok
Выбрать видеопоток который будет использоваться
    определить Видео Фильтр
        определить нужно ли резать видео
        определить нужно ли скалировать видео
        определить нужно ли менять пиксельность [SAR/DAR]
Выбрать аудиопотоки которые будут использоваться
    уточнить язык
    уточнить количество дорожек
    определить степень сжатия звука
Выбрать поток с Субтитрами который будет использоваться
*/

//AllAsIs - возвращает строку аргументов для того чтобы вытащить все потоки "как есть"
func AllAsIs(path string) (string, error) {
	f, e := ffinfo.Probe(path)
	if e != nil {
		return "", e
	}
	output := "-i " + path
	sl := strings.Split(path, ".")
	baseOutputPath := strings.Join(sl[0:len(sl)-1], "")
	vStr := -1
	aStr := -1
	for _, stream := range f.Streams {
		switch stream.CodecType {
		case "video":
			vStr++
			output += fmt.Sprintf(" -map 0:v:%v -vcodec copy %v_RAW_%v.mp4 ", vStr, baseOutputPath, vStr)
		case "audio":
			aStr++
			output += fmt.Sprintf(" -map 0:a:%v -acodec copy %v_RAW_%v.wav ", aStr, baseOutputPath, aStr)
		}
	}
	return output, nil
}

//VideoFCLine - должен дать параметры которые вставляются в filter_complex для видео
func VideoFCLine(path string) (string, error) {
	f, e := ffinfo.Probe(path)
	if e != nil {
		return "", e
	}
	com, _ := command.New(
		command.CommandLineArguments(fmt.Sprintf("ffmpeg -i %v", path)),
		command.Set(command.TERMINAL_ON),
	)
	com.Run()
	fmt.Println(f)
	return "FCLine", nil
}
