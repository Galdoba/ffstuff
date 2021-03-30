package muxer

import (
	"github.com/Galdoba/ffstuff/fldr"
	"github.com/Galdoba/ffstuff/pkg/cli"
	"github.com/Galdoba/ffstuff/pkg/namedata"
	"github.com/Galdoba/utils"
)

/*
rem ffmpeg ^
rem -i "%f_name%.mp4" ^
rem -i "%f_name%_rus20.ac3" ^
rem -i "%f_name%_eng51.ac3" ^
rem -i "%f_name%.srt" ^
rem -codec copy -codec:s mov_text ^
rem     -map 0:v ^
rem     -map 1:a -metadata:s:a:0 language=rus ^
rem     -map 2:a -metadata:s:a:1 language=eng ^
rem     -map 3:s -metadata:s:s:0 language=rus ^
rem "\\192.168.32.3\ROOT\#PETR\toCheck\%f_name%_ar2e6.mp4"


*/

func Test() {

}

func MuxA2(video, audio1 string) {
	base := utils.CommonPrefix(video, audio1)
	base = namedata.RetrieveShortName(base)
	prog := "ffmpeg"
	args := []string{
		"-i", video,
		"-i", audio1,
		"-codec", "copy", "-codec:s", "mov_text",
		"-map", "0:v",
		"-map", "1:a", "-metadata:s:a:0", "language=rus",
		fldr.OutPath() + base + "_ar2.mp4",
	}
	cli.RunConsole(prog, args...)
}
