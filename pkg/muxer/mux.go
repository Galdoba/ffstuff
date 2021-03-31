package muxer

import (
	"bufio"
	"log"
	"os"

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
	//Read mux list
	//err: file not found = skip
	//err: end file = end program
	//Choose Muxer
	//err: input files are not in the same folder = skip
	//err: input files not same duration = skip (need inchecker)
	//Run Muxer
	//err: output file not same duration as input = skip (need inchecker)
	//Rename input files
}

func MuxList() []string {
	file, err := os.Open(fldr.MuxPath() + "muxlist.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	list := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		list = append(list, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return list
}

func MuxA2(video, audio1 string) error {
	base := utils.CommonPrefix(video, audio1)
	base = namedata.RetrieveShortName(base)
	//dir := namedata.RetrieveDirectory() -
	prog := "ffmpeg"
	args := []string{
		"-i", video,
		"-i", audio1,
		"-codec", "copy", "-codec:s", "mov_text",
		"-map", "0:v",
		"-map", "1:a", "-metadata:s:a:0", "language=rus",
		fldr.OutPath() + base + "_ar2.mp4",
	}
	_, _, err := cli.RunConsole(prog, args...)
	return err
}

func MuxA6(video, audio1 string) error {
	base := utils.CommonPrefix(video, audio1)
	base = namedata.RetrieveShortName(base)
	prog := "ffmpeg"
	args := []string{
		"-i", video,
		"-i", audio1,
		"-codec", "copy", "-codec:s", "mov_text",
		"-map", "0:v",
		"-map", "1:a", "-metadata:s:a:0", "language=rus",
		fldr.OutPath() + base + "_ar6.mp4",
	}
	_, _, err := cli.RunConsole(prog, args...)
	return err
}
