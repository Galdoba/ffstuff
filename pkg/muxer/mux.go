package muxer

import (
	"github.com/Galdoba/ffstuff/fldr"
	"github.com/Galdoba/ffstuff/pkg/cli"
	"github.com/Galdoba/ffstuff/pkg/namedata"
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

func commonPrefix(sl ...string) string {
	leastLen := 10000000
	for _, v := range sl {
		if len(v) < leastLen {
			leastLen = len(v)
		}
	}
	common := ""
	dif := false
	for i := 0; i < leastLen; i++ {
		if dif {
			break
		}
		examine := ""
		for in, v := range sl {
			if string(v[in]) != string(sl[0][in]) {
				dif = true
				break
			}
			examine = string(sl[0][i])
		}
		if !dif {
			common += examine
		}
	}
	return common
}

func MuxA2(video, audio1 string) {
	base := commonPrefix(video, audio1)
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
