package muxer

/*
-функции мукса много раз повторяются из расчета "пусть будет наглядно"

*/

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Galdoba/ffstuff/fldr"
	"github.com/Galdoba/ffstuff/pkg/cli"
	"github.com/Galdoba/ffstuff/pkg/namedata"
	"github.com/Galdoba/utils"
)

const (
	MuxerAR2    = "ar2"
	MuxerA6     = "ar6"
	MuxerAE2    = "ae2"
	MuxerAE6    = "ae6"
	MuxerAR2E2  = "ar2e2"
	MuxerAR2E6  = "ar2e6"
	MuxerA6E2   = "ar6e2"
	MuxerA6E6   = "ar6e6"
	MuxerAR2s   = "ar2_sr"
	MuxerA6s    = "ar6_sr"
	MuxerAR2E2s = "ar2e2_sr"
	MuxerAR2E6s = "ar2e6_sr"
	MuxerA6E2s  = "ar6e2_sr"
	MuxerA6E6s  = "ar6e6_sr"
	MuxerNA     = ""
	MuxerSKIP   = "ST"
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

func MuxList() ([]string, error) {
	file, err := os.Open(fldr.MuxPath() + "muxlist.txt")
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, err
		}
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
	return list, scanner.Err()
}

func ChooseMuxer(task string) ([]string, string, error) {
	data := strings.Split(task, " ")
	if data[0] == "" {
		return []string{}, MuxerSKIP, nil
	}
	if len(data) < 2 {
		return []string{}, MuxerNA, errors.New("muxer not assigned")
	}
	switch data[1] {
	default:
		return []string{}, MuxerSKIP, errors.New("muxer not recognised")
	case MuxerAR2, MuxerA6, MuxerAE2, MuxerAE6, MuxerAR2E2, MuxerAR2E6, MuxerA6E2, MuxerA6E6, MuxerAR2s, MuxerA6s, MuxerAR2E2s, MuxerAR2E6s, MuxerA6E2s, MuxerA6E6s:
		paths := defineFiles(task)
		return paths, data[1], nil
	}
}

func Run(muxerTask string, files []string) error {
	if err := assertInputFiles(files...); err != nil {
		return err
	}
	switch muxerTask {
	default:
		return errors.New("undefined muxer task")
	case MuxerAR2:
		return MuxA2(files[0], files[1])
	case MuxerA6:
		return MuxA6(files[0], files[1])
	case MuxerAE2:
		return MuxAE2(files[0], files[1])
	case MuxerAE6:
		return MuxAE6(files[0], files[1])
	case MuxerAR2E2:
		return MuxA2E2(files[0], files[1], files[2])
	case MuxerAR2E6:
		return MuxA2E6(files[0], files[1], files[2])
	case MuxerA6E2:
		return MuxA6E2(files[0], files[1], files[2])
	case MuxerA6E6:
		return MuxA6E6(files[0], files[1], files[2])
	case MuxerAR2s:
		return MuxA2s(files[0], files[1], files[2])
	case MuxerA6s:
		return MuxA6s(files[0], files[1], files[2])
	case MuxerAR2E2s:
		return MuxA2E2s(files[0], files[1], files[2], files[3])
	case MuxerAR2E6s:
		return MuxA2E6s(files[0], files[1], files[2], files[3])
	case MuxerA6E2s:
		return MuxA6E2s(files[0], files[1], files[2], files[3])
	case MuxerA6E6s:
		return MuxA6E6s(files[0], files[1], files[2], files[3])
	}

}

func defineFiles(task string) []string {
	data := strings.Split(task, " ")
	paths := []string{fldr.MuxPath() + data[0]}
	base := strings.TrimSuffix(data[0], ".mp4")
	base = fldr.MuxPath() + base
	switch data[1] {
	case MuxerAR2:
		paths = append(paths, base+"_rus20.ac3")
	case MuxerA6:
		paths = append(paths, base+"_rus51.ac3")
	case MuxerAE2:
		paths = append(paths, base+"_eng20.ac3")
	case MuxerAE6:
		paths = append(paths, base+"_eng51.ac3")
	case MuxerAR2E2:
		paths = append(paths, base+"_rus20.ac3")
		paths = append(paths, base+"_eng20.ac3")
	case MuxerAR2E6:
		paths = append(paths, base+"_rus20.ac3")
		paths = append(paths, base+"_eng51.ac3")
	case MuxerA6E2:
		paths = append(paths, base+"_rus51.ac3")
		paths = append(paths, base+"_eng20.ac3")
	case MuxerA6E6:
		paths = append(paths, base+"_rus51.ac3")
		paths = append(paths, base+"_eng51.ac3")
	case MuxerAR2s:
		paths = append(paths, base+"_rus20.ac3")
		paths = append(paths, base+".srt")
	case MuxerA6s:
		paths = append(paths, base+"_rus51.ac3")
		paths = append(paths, base+".srt")
	case MuxerAR2E2s:
		paths = append(paths, base+"_rus20.ac3")
		paths = append(paths, base+"_eng20.ac3")
		paths = append(paths, base+".srt")
	case MuxerAR2E6s:
		paths = append(paths, base+"_rus20.ac3")
		paths = append(paths, base+"_eng51.ac3")
		paths = append(paths, base+".srt")
	case MuxerA6E2s:
		paths = append(paths, base+"_rus51.ac3")
		paths = append(paths, base+"_eng20.ac3")
		paths = append(paths, base+".srt")
	case MuxerA6E6s:
		paths = append(paths, base+"_rus51.ac3")
		paths = append(paths, base+"_eng51.ac3")
		paths = append(paths, base+".srt")
	}
	return paths
}

func MuxA2(video, audio1 string) error {
	base := utils.CommonPrefix(video, audio1)
	base = namedata.RetrieveShortName(base)
	prog := "ffmpeg"
	args := []string{
		"-i", video,
		"-i", audio1,
		"-codec", "copy",
		"-map", "0:v",
		"-map", "1:a", "-metadata:s:a:0", "language=rus",
		fldr.OutPath() + base + "_ar2.mp4",
	}
	_, _, err := cli.RunConsole(prog, args...)
	return err
}

func MuxAE2(video, audio1 string) error {
	base := utils.CommonPrefix(video, audio1)
	base = namedata.RetrieveShortName(base)
	prog := "ffmpeg"
	args := []string{
		"-i", video,
		"-i", audio1,
		"-codec", "copy",
		"-map", "0:v",
		"-map", "1:a", "-metadata:s:a:0", "language=eng",
		fldr.OutPath() + base + "_ae2.mp4",
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
		"-codec", "copy",
		"-map", "0:v",
		"-map", "1:a", "-metadata:s:a:0", "language=rus",
		fldr.OutPath() + base + "_ar6.mp4",
	}
	_, _, err := cli.RunConsole(prog, args...)
	return err
}

func MuxAE6(video, audio1 string) error {
	base := utils.CommonPrefix(video, audio1)
	base = namedata.RetrieveShortName(base)
	prog := "ffmpeg"
	args := []string{
		"-i", video,
		"-i", audio1,
		"-codec", "copy",
		"-map", "0:v",
		"-map", "1:a", "-metadata:s:a:0", "language=eng",
		fldr.OutPath() + base + "_ae6.mp4",
	}
	_, _, err := cli.RunConsole(prog, args...)
	return err
}

func MuxA2E2(video, audio1, audio2 string) error {
	base := utils.CommonPrefix(video, audio1, audio2)
	base = namedata.RetrieveShortName(base)
	prog := "ffmpeg"
	args := []string{
		"-i", video,
		"-i", audio1,
		"-i", audio2,
		"-codec", "copy",
		"-map", "0:v",
		"-map", "1:a", "-metadata:s:a:0", "language=rus",
		"-map", "2:a", "-metadata:s:a:1", "language=eng",
		fldr.OutPath() + base + "_ar2e2.mp4",
	}
	_, _, err := cli.RunConsole(prog, args...)
	return err
}

func MuxA2E6(video, audio1, audio2 string) error {
	base := utils.CommonPrefix(video, audio1, audio2)
	base = namedata.RetrieveShortName(base)
	prog := "ffmpeg"
	args := []string{
		"-i", video,
		"-i", audio1,
		"-i", audio2,
		"-codec", "copy",
		"-map", "0:v",
		"-map", "1:a", "-metadata:s:a:0", "language=rus",
		"-map", "2:a", "-metadata:s:a:1", "language=eng",
		fldr.OutPath() + base + "_ar2e6.mp4",
	}
	_, _, err := cli.RunConsole(prog, args...)
	return err
}

func MuxA6E2(video, audio1, audio2 string) error {
	base := utils.CommonPrefix(video, audio1, audio2)
	base = namedata.RetrieveShortName(base)
	prog := "ffmpeg"
	args := []string{
		"-i", video,
		"-i", audio1,
		"-i", audio2,
		"-codec", "copy",
		"-map", "0:v",
		"-map", "1:a", "-metadata:s:a:0", "language=rus",
		"-map", "2:a", "-metadata:s:a:1", "language=eng",
		fldr.OutPath() + base + "_ar6e2.mp4",
	}
	_, _, err := cli.RunConsole(prog, args...)
	return err
}

func MuxA6E6(video, audio1, audio2 string) error {
	base := utils.CommonPrefix(video, audio1, audio2)
	base = namedata.RetrieveShortName(base)
	prog := "ffmpeg"
	args := []string{
		"-i", video,
		"-i", audio1,
		"-i", audio2,
		"-codec", "copy",
		"-map", "0:v",
		"-map", "1:a", "-metadata:s:a:0", "language=rus",
		"-map", "2:a", "-metadata:s:a:1", "language=eng",
		fldr.OutPath() + base + "_ar6e6.mp4",
	}
	_, _, err := cli.RunConsole(prog, args...)
	return err
}

////////////////////
func MuxA2s(video, audio1, subs string) error {
	base := utils.CommonPrefix(video, audio1)
	base = namedata.RetrieveShortName(base)
	prog := "ffmpeg"
	args := []string{
		"-i", video,
		"-i", audio1,
		"-i", subs,
		"-codec", "copy", "-codec:s", "mov_text",
		"-map", "0:v",
		"-map", "1:a", "-metadata:s:a:0", "language=rus",
		"-map", "2:s", "-metadata:s:s:0", "language=rus",
		fldr.OutPath() + base + "_ar2_sr.mp4",
	}
	_, _, err := cli.RunConsole(prog, args...)
	return err
}

func MuxA6s(video, audio1, subs string) error {
	base := utils.CommonPrefix(video, audio1)
	base = namedata.RetrieveShortName(base)
	prog := "ffmpeg"
	args := []string{
		"-i", video,
		"-i", audio1,
		"-i", subs,
		"-codec", "copy", "-codec:s", "mov_text",
		"-map", "0:v",
		"-map", "1:a", "-metadata:s:a:0", "language=rus",
		"-map", "2:s", "-metadata:s:s:0", "language=rus",
		fldr.OutPath() + base + "_ar6_sr.mp4",
	}
	_, _, err := cli.RunConsole(prog, args...)
	return err
}

func MuxA2E2s(video, audio1, audio2, subs string) error {
	base := utils.CommonPrefix(video, audio1, audio2)
	base = namedata.RetrieveShortName(base)
	prog := "ffmpeg"
	args := []string{
		"-i", video,
		"-i", audio1,
		"-i", audio2,
		"-i", subs,
		"-codec", "copy", "-codec:s", "mov_text",
		"-map", "0:v",
		"-map", "1:a", "-metadata:s:a:0", "language=rus",
		"-map", "2:a", "-metadata:s:a:1", "language=eng",
		"-map", "3:s", "-metadata:s:s:0", "language=rus",
		fldr.OutPath() + base + "_ar2e2_sr.mp4",
	}
	_, _, err := cli.RunConsole(prog, args...)
	return err
}

func MuxA2E6s(video, audio1, audio2, subs string) error {
	base := utils.CommonPrefix(video, audio1, audio2)
	base = namedata.RetrieveShortName(base)
	prog := "ffmpeg"
	args := []string{
		"-i", video,
		"-i", audio1,
		"-i", audio2,
		"-i", subs,
		"-codec", "copy", "-codec:s", "mov_text",
		"-map", "0:v",
		"-map", "1:a", "-metadata:s:a:0", "language=rus",
		"-map", "2:a", "-metadata:s:a:1", "language=eng",
		"-map", "3:s", "-metadata:s:s:0", "language=rus",
		fldr.OutPath() + base + "_ar2e6_sr.mp4",
	}
	_, _, err := cli.RunConsole(prog, args...)
	return err
}

func MuxA6E2s(video, audio1, audio2, subs string) error {
	base := utils.CommonPrefix(video, audio1, audio2)
	base = namedata.RetrieveShortName(base)
	prog := "ffmpeg"
	args := []string{
		"-i", video,
		"-i", audio1,
		"-i", audio2,
		"-i", subs,
		"-codec", "copy", "-codec:s", "mov_text",
		"-map", "0:v",
		"-map", "1:a", "-metadata:s:a:0", "language=rus",
		"-map", "2:a", "-metadata:s:a:1", "language=eng",
		"-map", "3:s", "-metadata:s:s:0", "language=rus",
		fldr.OutPath() + base + "_ar6e2_sr.mp4",
	}
	_, _, err := cli.RunConsole(prog, args...)
	return err
}

func MuxA6E6s(video, audio1, audio2, subs string) error {
	base := utils.CommonPrefix(video, audio1, audio2, subs)
	base = namedata.RetrieveShortName(base)
	prog := "ffmpeg"
	args := []string{
		"-i", video,
		"-i", audio1,
		"-i", audio2,
		"-i", subs,
		"-codec", "copy", "-codec:s", "mov_text",
		"-map", "0:v",
		"-map", "1:a", "-metadata:s:a:0", "language=rus",
		"-map", "2:a", "-metadata:s:a:1", "language=eng",
		"-map", "3:s", "-metadata:s:s:0", "language=rus",
		fldr.OutPath() + base + "_ar6e6_sr.mp4",
	}
	_, _, err := cli.RunConsole(prog, args...)
	return err
}

func assertInputFiles(filePath ...string) error {
	for _, path := range filePath {
		_, err := os.Stat(path)
		if err != nil {
			if os.IsNotExist(err) {
				return errors.New("File not found: " + path)
			}
			return err
		}
	}
	return nil
}

type Task struct {
	instruction string
	video       string
	audio1      string
	audio2      string
	subtitles   string
	err         error
}

func (t *Task) Validate() {
	base := baseOf(t)

	switch {
	default:
		t.err = nil
	case !strings.Contains(t.instruction, "2") && !strings.Contains(t.instruction, "6"):
		t.err = fmt.Errorf("instruction has no valid channel layout (expecting '2' or '6')")
	case !strings.Contains(t.instruction, "r") && !strings.Contains(t.instruction, "e") && !strings.Contains(t.instruction, "qqq"):
		t.err = fmt.Errorf("instruction has no valid language marker (expecting 'r', 'e' or 'qqq')")
	case strings.Contains(t.instruction, "qqq2qqq") || strings.Contains(t.instruction, "qqq6qqq"):
		t.err = fmt.Errorf("instruction can not have same language")
	case strings.Contains(t.instruction, "rus2rus") || strings.Contains(t.instruction, "rus6rus"):
		t.err = fmt.Errorf("instruction can not have same language")
	case strings.Contains(t.instruction, "eng2eng") || strings.Contains(t.instruction, "eng6eng"):
		t.err = fmt.Errorf("instruction can not have same language")
	case strings.Contains(t.instruction, "aqqq") && !strings.Contains(t.audio1, "_qqq"):
		t.err = fmt.Errorf("audio1 does not match with instruction")
	case strings.Contains(t.instruction, "ar") && !strings.Contains(t.audio1, "_rus"):
		t.err = fmt.Errorf("audio1 does not match with instruction")
	case strings.Contains(t.instruction, "ae") && !strings.Contains(t.audio1, "_eng"):
		t.err = fmt.Errorf("audio1 does not match with instruction")
	case strings.Contains(t.instruction, "2qqq") && !strings.Contains(t.audio2, "_qqq"):
		t.err = fmt.Errorf("audio1 does not match with instruction")
	case strings.Contains(t.instruction, "2r") && !strings.Contains(t.audio2, "_rus"):
		t.err = fmt.Errorf("audio1 does not match with instruction")
	case strings.Contains(t.instruction, "2e") && !strings.Contains(t.audio2, "_eng"):
		t.err = fmt.Errorf("audio1 does not match with instruction")
	case strings.Contains(t.instruction, "6qqq") && !strings.Contains(t.audio2, "_qqq"):
		t.err = fmt.Errorf("audio1 does not match with instruction")
	case strings.Contains(t.instruction, "6r") && !strings.Contains(t.audio2, "_rus"):
		t.err = fmt.Errorf("audio1 does not match with instruction")
	case strings.Contains(t.instruction, "6e") && !strings.Contains(t.audio2, "_eng"):
		t.err = fmt.Errorf("audio1 does not match with instruction")
	case strings.Contains(t.instruction, "_sr") && t.subtitles != baseOf(t)+".srt":
		t.err = fmt.Errorf("subtitles does not match with instruction (have '%v', expect '%v'", t.subtitles, base+".srt")
	case !strings.Contains(t.instruction, "_sr") && t.subtitles != "":
		t.err = fmt.Errorf("subtitles defined but not instructed")

	}
}

func baseOf(t *Task) string {
	base := strings.TrimSuffix(t.video, ".mp4")
	base = strings.TrimSuffix(t.video, ".mpeg")
	return base
}

func decodeInstruction(t *Task) (string, string, string) {
	a1, a2, s := "", "", ""
	left := t.instruction
	if strings.TrimSuffix(t.instruction, "_sr") != left {
		left = strings.TrimSuffix(t.instruction, "_sr")
		s = baseOf(t) + ".srt"
	}
	instrMap := make(map[string]string)
	instrMap["ar2"] = "_rus20.ac3"
	instrMap["ae2"] = "_eng20.ac3"
	instrMap["aqqq2"] = "_qqq20.ac3"
	instrMap["ar6"] = "_rus51.ac3"
	instrMap["ae6"] = "_eng51.ac3"
	instrMap["aqqq6"] = "_qqq51.ac3"
	instrMap["r2"] = "_rus20.ac3"
	instrMap["e2"] = "_eng20.ac3"
	instrMap["qqq2"] = "_qqq20.ac3"
	instrMap["r6"] = "_rus51.ac3"
	instrMap["e6"] = "_eng51.ac3"
	instrMap["qqq6"] = "_qqq51.ac3"
	first := []string{"ar2", "ar6", "ae2", "ae6", "aqqq2", "aqqq6"}
	for _, p1 := range first {
		if strings.TrimPrefix(t.instruction, p1) != left {
			left = strings.TrimSuffix(t.instruction, p1)
			a1 = baseOf(t) + instrMap[p1]
		}
	}
	if left != "" {
		a2 = baseOf(t) + instrMap[left]
	}
	return a1, a2, s
}
