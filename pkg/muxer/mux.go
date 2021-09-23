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
	"github.com/Galdoba/utils"
)

const (
	MuxerAR2 = "ar2"
	MuxerAR6 = "ar6"
	MuxerAE2 = "ae2"
	MuxerAE6 = "ae6"
	MuxerAQ2 = "ae2"
	MuxerAQ6 = "ae6"

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

func MuxListV2(path string) ([]*Task, error) {
	var tl []*Task
	file, err := os.Open(path + "muxlist.txt")
	if err != nil {
		if os.IsNotExist(err) {
			return tl, err
		}
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		tl = append(tl, NewTask(path, scanner.Text()))
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return tl, scanner.Err()
}

func AssertTasks(tl []*Task) []error {
	var errList []error
	for i, tsk := range tl {
		if tsk.err != nil {
			errList = append(errList, fmt.Errorf("Task %v error: %v", i, tsk.err.Error()))
		}
		errList = append(errList, checkFileExistiense(tsk)...) //проверка наличия файлов
		//errList = append(errList, checkOutputFileExistience(tsk))
	}
	return errList
}

func checkFileExistiense(t *Task) []error {
	var errList []error
	checkFiles := []string{t.video, t.audio1, t.audio2, t.subtitles}
	for _, f := range checkFiles {
		if f != "" {
			if _, err := os.Stat(f); os.IsNotExist(err) {
				errList = append(errList, fmt.Errorf("%v: not found", f))
			}
		}
	}
	return errList
}

func checkOutputFileExistience(t *Task) error {
	body := strings.TrimPrefix(baseOf(t), t.path)
	f := fldr.OutPath() + body + "_" + t.instruction + ".mp4"
	fmt.Println(f)
	_, err := os.Stat(f)
	if err == nil {
		return fmt.Errorf("%v: detected", f)
	}
	return nil
}

func MuxV2(t *Task) error {
	prog := "ffmpeg"
	args2 := []string{
		"-i", t.video,
		"-i", t.audio1,
	}
	if t.audio2 != "" {
		args2 = append(args2, "-i")
		args2 = append(args2, t.audio2)
	}
	if t.subtitles != "" {
		args2 = append(args2, "-i")
		args2 = append(args2, t.subtitles)
	}
	args2 = append(args2, []string{
		"-codec", "copy", "-codec:s", "mov_text",
		"-map", "0:v",
		"-map", "1:a", "-metadata:s:a:0", "language=" + t.l1,
	}...)
	if t.audio2 != "" {
		args2 = append(args2, "-map", "2:a", "-metadata:s:a:1", "language="+t.l2)
	}
	if t.subtitles != "" {
		args2 = append(args2, "-map", "3:s", "-metadata:s:s:0", "language=rus")
	}
	//fldr.OutPath() + body + "_" + t.instruction + ".mp4"
	args2 = append(args2, fldr.OutPath()+OutputFile(t))
	_, _, err := cli.RunConsole(prog, args2...)
	return err
}

func OutputFile(t *Task) string {
	body := strings.TrimPrefix(baseOf(t), t.path)
	return body + "_" + t.instruction + ".mp4"
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

////////

type Task struct {
	input       string
	instruction string
	video       string
	audio1      string
	audio2      string
	subtitles   string
	err         error
	l1          string
	l2          string
	path        string
}

func NewTask(path, instructionData string) *Task {
	t := Task{}
	t.path = path
	t.input = instructionData
	data := strings.Split(instructionData, " ")
	if len(data) < 2 {
		t.err = fmt.Errorf("Can not create instruction with line '%v'\n", instructionData)
		return &t
	}
	if !utils.ListContains(validInstructions(), data[1]) {
		t.err = fmt.Errorf("instruction '%v' is invalid\n", data)
		return &t
	}
	t.instruction = data[1]
	t.video = t.path + data[0]
	a1, a2, s := decodeInstruction(&t)
	base := baseOf(&t)
	t.audio1 = base + a1
	if a2 != "" {
		t.audio2 = base + a2
	}
	if s != "" {
		t.subtitles = base + s
	}
	t.l1 = a1
	t.l1 = strings.TrimPrefix(t.l1, "_")
	t.l1 = strings.TrimSuffix(t.l1, "51.ac3")
	t.l1 = strings.TrimSuffix(t.l1, "20.ac3")
	t.l2 = a2
	t.l2 = strings.TrimPrefix(t.l2, "_")
	t.l2 = strings.TrimSuffix(t.l2, "51.ac3")
	t.l2 = strings.TrimSuffix(t.l2, "20.ac3")
	return &t
}

func (t *Task) Line() string {
	return t.input
}

func (t *Task) Validate() {
	//base := baseOf(t)
	a1, a2, s := decodeInstruction(t)
	switch {
	default:
		t.err = fmt.Errorf("Undecided")
	case t.instruction == "":
		t.err = fmt.Errorf("instruction not set")
	case t.video == "":
		t.err = fmt.Errorf("video not set")
	case strings.Contains(t.video, ".mpeg"):
		t.err = fmt.Errorf("[WARNING] '.mpeg' is not used anymore")
	case !strings.Contains(t.audio1, a1):
		t.err = fmt.Errorf("audio1 is not correct: have '%v', expect '%v'", t.audio1, a1)
	case !strings.Contains(t.audio2, a2):
		t.err = fmt.Errorf("audio1 is not correct: have '%v', expect '%v'", t.audio2, a2)
	case !strings.Contains(t.subtitles, s):
		t.err = fmt.Errorf("subtitles is not correct: have '%v', expect '%v'", t.subtitles, s)
	case a2 == "" && t.audio2 != a2:
		t.err = fmt.Errorf("audio2 is not correct: have '%v', expect '%v'", t.audio2, a2)
	case s == "" && t.subtitles != s:
		t.err = fmt.Errorf("subtitles is not correct: have '%v', expect '%v'", t.subtitles, s)
	}
	if t.err.Error() == "Undecided" && utils.ListContains(validInstructions(), t.instruction) {
		t.err = nil
	}
}

func baseOf(t *Task) string {
	base := strings.TrimSuffix(t.video, ".mp4")
	return base
}

func decodeInstruction(t *Task) (string, string, string) {
	a1, a2, s := "", "", ""
	left := t.instruction
	prefixes := []string{"ar2", "ar6", "ae2", "ae6", "aqqq2", "aqqq6"}
	middlefixes := []string{"", "r2", "r6", "e2", "e6", "qqq2", "qqq6"}
	if strings.TrimSuffix(t.instruction, "_sr") != left {
		left = strings.TrimSuffix(t.instruction, "_sr")
		s = ".srt"
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
	mid := "???"
	for _, p1 := range prefixes {
		mid = strings.TrimPrefix(t.instruction, p1)
		if mid != t.instruction {
			a1 = instrMap[p1]
			break
		}
	}
	mid = strings.TrimSuffix(mid, "_sr")
	for _, p2 := range middlefixes {
		if p2 == mid {
			a2 = instrMap[mid]
			break
		}
	}
	//a2 = instrMap[left]
	return a1, a2, s
}

func validInstructions() []string {
	prefix := "a"
	l1 := []string{"r", "e", "qqq"}
	c1 := []string{"2", "6"}
	l2 := []string{"", "r", "e", "qqq"}
	c2 := []string{"", "2", "6"}
	postfix := []string{"", "_sr", " "}
	var instructions []string
	for _, a := range l1 {
		for _, b := range c1 {
			for _, c := range l2 {
				for _, d := range c2 {
					for _, e := range postfix {
						//fmt.Println(prefix + a + b + c + d + e)
						instructions = append(instructions, prefix+a+b+c+d+e)
					}
				}
			}
		}
	}
	return instructions
}

func ShowTaskList(tl []*Task) {
	fmt.Println("================================================================================")
	for _, t := range tl {
		fmt.Println(t.video, t.instruction)
	}
	fmt.Println("================================================================================")
}

///////////////////////////////////////////
///////////////////////////////////////////
///////////////LEGACY//////////////////////
///////////////////////////////////////////
///////////////////////////////////////////

// func Run(muxerTask string, files []string) error {
// 	if err := assertInputFiles(files...); err != nil {
// 		return err
// 	}
// 	switch muxerTask {
// 	default:
// 		return errors.New("undefined muxer task")
// 	case MuxerAR2:
// 		return MuxA2(files[0], files[1])
// 	case MuxerAR6:
// 		return MuxA6(files[0], files[1])
// 	case MuxerAE2:
// 		return MuxAE2(files[0], files[1])
// 	case MuxerAE6:
// 		return MuxAE6(files[0], files[1])
// 	case MuxerAR2E2:
// 		return MuxA2E2(files[0], files[1], files[2])
// 	case MuxerAR2E6:
// 		return MuxA2E6(files[0], files[1], files[2])
// 	case MuxerA6E2:
// 		return MuxA6E2(files[0], files[1], files[2])
// 	case MuxerA6E6:
// 		return MuxA6E6(files[0], files[1], files[2])
// 	case MuxerAR2s:
// 		return MuxA2s(files[0], files[1], files[2])
// 	case MuxerA6s:
// 		return MuxA6s(files[0], files[1], files[2])
// 	case MuxerAR2E2s:
// 		return MuxA2E2s(files[0], files[1], files[2], files[3])
// 	case MuxerAR2E6s:
// 		return MuxA2E6s(files[0], files[1], files[2], files[3])
// 	case MuxerA6E2s:
// 		return MuxA6E2s(files[0], files[1], files[2], files[3])
// 	case MuxerA6E6s:
// 		return MuxA6E6s(files[0], files[1], files[2], files[3])
// 	}

// }

// func defineFiles(task string) []string {
// 	data := strings.Split(task, " ")
// 	paths := []string{fldr.MuxPath() + data[0]}
// 	base := strings.TrimSuffix(data[0], ".mp4")
// 	base = fldr.MuxPath() + base
// 	switch data[1] {
// 	case MuxerAR2:
// 		paths = append(paths, base+"_rus20.ac3")
// 	case MuxerAR6:
// 		paths = append(paths, base+"_rus51.ac3")
// 	case MuxerAE2:
// 		paths = append(paths, base+"_eng20.ac3")
// 	case MuxerAE6:
// 		paths = append(paths, base+"_eng51.ac3")
// 	case MuxerAR2E2:
// 		paths = append(paths, base+"_rus20.ac3")
// 		paths = append(paths, base+"_eng20.ac3")
// 	case MuxerAR2E6:
// 		paths = append(paths, base+"_rus20.ac3")
// 		paths = append(paths, base+"_eng51.ac3")
// 	case MuxerA6E2:
// 		paths = append(paths, base+"_rus51.ac3")
// 		paths = append(paths, base+"_eng20.ac3")
// 	case MuxerA6E6:
// 		paths = append(paths, base+"_rus51.ac3")
// 		paths = append(paths, base+"_eng51.ac3")
// 	case MuxerAR2s:
// 		paths = append(paths, base+"_rus20.ac3")
// 		paths = append(paths, base+".srt")
// 	case MuxerA6s:
// 		paths = append(paths, base+"_rus51.ac3")
// 		paths = append(paths, base+".srt")
// 	case MuxerAR2E2s:
// 		paths = append(paths, base+"_rus20.ac3")
// 		paths = append(paths, base+"_eng20.ac3")
// 		paths = append(paths, base+".srt")
// 	case MuxerAR2E6s:
// 		paths = append(paths, base+"_rus20.ac3")
// 		paths = append(paths, base+"_eng51.ac3")
// 		paths = append(paths, base+".srt")
// 	case MuxerA6E2s:
// 		paths = append(paths, base+"_rus51.ac3")
// 		paths = append(paths, base+"_eng20.ac3")
// 		paths = append(paths, base+".srt")
// 	case MuxerA6E6s:
// 		paths = append(paths, base+"_rus51.ac3")
// 		paths = append(paths, base+"_eng51.ac3")
// 		paths = append(paths, base+".srt")
// 	}
// 	return paths
// }

// func MuxA2(video, audio1 string) error {
// 	base := utils.CommonPrefix(video, audio1)
// 	base = namedata.RetrieveShortName(base)
// 	prog := "ffmpeg"
// 	args := []string{
// 		"-i", video,
// 		"-i", audio1,
// 		"-codec", "copy",
// 		"-map", "0:v",
// 		"-map", "1:a", "-metadata:s:a:0", "language=rus",
// 		fldr.OutPath() + base + "_ar2.mp4",
// 	}
// 	_, _, err := cli.RunConsole(prog, args...)
// 	return err
// }

// func MuxAE2(video, audio1 string) error {
// 	base := utils.CommonPrefix(video, audio1)
// 	base = namedata.RetrieveShortName(base)
// 	prog := "ffmpeg"
// 	args := []string{
// 		"-i", video,
// 		"-i", audio1,
// 		"-codec", "copy",
// 		"-map", "0:v",
// 		"-map", "1:a", "-metadata:s:a:0", "language=eng",
// 		fldr.OutPath() + base + "_ae2.mp4",
// 	}
// 	_, _, err := cli.RunConsole(prog, args...)
// 	return err
// }

// func MuxA6(video, audio1 string) error {
// 	base := utils.CommonPrefix(video, audio1)
// 	base = namedata.RetrieveShortName(base)
// 	prog := "ffmpeg"
// 	args := []string{
// 		"-i", video,
// 		"-i", audio1,
// 		"-codec", "copy",
// 		"-map", "0:v",
// 		"-map", "1:a", "-metadata:s:a:0", "language=rus",
// 		fldr.OutPath() + base + "_ar6.mp4",
// 	}
// 	_, _, err := cli.RunConsole(prog, args...)
// 	return err
// }

// func MuxAE6(video, audio1 string) error {
// 	base := utils.CommonPrefix(video, audio1)
// 	base = namedata.RetrieveShortName(base)
// 	prog := "ffmpeg"
// 	args := []string{
// 		"-i", video,
// 		"-i", audio1,
// 		"-codec", "copy",
// 		"-map", "0:v",
// 		"-map", "1:a", "-metadata:s:a:0", "language=eng",
// 		fldr.OutPath() + base + "_ae6.mp4",
// 	}
// 	_, _, err := cli.RunConsole(prog, args...)
// 	return err
// }

// func MuxA2E2(video, audio1, audio2 string) error {
// 	base := utils.CommonPrefix(video, audio1, audio2)
// 	base = namedata.RetrieveShortName(base)
// 	prog := "ffmpeg"
// 	args := []string{
// 		"-i", video,
// 		"-i", audio1,
// 		"-i", audio2,
// 		"-codec", "copy",
// 		"-map", "0:v",
// 		"-map", "1:a", "-metadata:s:a:0", "language=rus",
// 		"-map", "2:a", "-metadata:s:a:1", "language=eng",
// 		fldr.OutPath() + base + "_ar2e2.mp4",
// 	}
// 	_, _, err := cli.RunConsole(prog, args...)
// 	return err
// }

// func MuxA2E6(video, audio1, audio2 string) error {
// 	base := utils.CommonPrefix(video, audio1, audio2)
// 	base = namedata.RetrieveShortName(base)
// 	prog := "ffmpeg"
// 	args := []string{
// 		"-i", video,
// 		"-i", audio1,
// 		"-i", audio2,
// 		"-codec", "copy",
// 		"-map", "0:v",
// 		"-map", "1:a", "-metadata:s:a:0", "language=rus",
// 		"-map", "2:a", "-metadata:s:a:1", "language=eng",
// 		fldr.OutPath() + base + "_ar2e6.mp4",
// 	}
// 	_, _, err := cli.RunConsole(prog, args...)
// 	return err
// }

// func MuxA6E2(video, audio1, audio2 string) error {
// 	base := utils.CommonPrefix(video, audio1, audio2)
// 	base = namedata.RetrieveShortName(base)
// 	prog := "ffmpeg"
// 	args := []string{
// 		"-i", video,
// 		"-i", audio1,
// 		"-i", audio2,
// 		"-codec", "copy",
// 		"-map", "0:v",
// 		"-map", "1:a", "-metadata:s:a:0", "language=rus",
// 		"-map", "2:a", "-metadata:s:a:1", "language=eng",
// 		fldr.OutPath() + base + "_ar6e2.mp4",
// 	}
// 	_, _, err := cli.RunConsole(prog, args...)
// 	return err
// }

// func MuxA6E6(video, audio1, audio2 string) error {
// 	base := utils.CommonPrefix(video, audio1, audio2)
// 	base = namedata.RetrieveShortName(base)
// 	prog := "ffmpeg"
// 	args := []string{
// 		"-i", video,
// 		"-i", audio1,
// 		"-i", audio2,
// 		"-codec", "copy",
// 		"-map", "0:v",
// 		"-map", "1:a", "-metadata:s:a:0", "language=rus",
// 		"-map", "2:a", "-metadata:s:a:1", "language=eng",
// 		fldr.OutPath() + base + "_ar6e6.mp4",
// 	}
// 	_, _, err := cli.RunConsole(prog, args...)
// 	return err
// }

// ////////////////////
// func MuxA2s(video, audio1, subs string) error {
// 	base := utils.CommonPrefix(video, audio1)
// 	base = namedata.RetrieveShortName(base)
// 	prog := "ffmpeg"
// 	args := []string{
// 		"-i", video,
// 		"-i", audio1,
// 		"-i", subs,
// 		"-codec", "copy", "-codec:s", "mov_text",
// 		"-map", "0:v",
// 		"-map", "1:a", "-metadata:s:a:0", "language=rus",
// 		"-map", "2:s", "-metadata:s:s:0", "language=rus",
// 		fldr.OutPath() + base + "_ar2_sr.mp4",
// 	}
// 	_, _, err := cli.RunConsole(prog, args...)
// 	return err
// }

// func MuxA6s(video, audio1, subs string) error {
// 	base := utils.CommonPrefix(video, audio1)
// 	base = namedata.RetrieveShortName(base)
// 	prog := "ffmpeg"
// 	args := []string{
// 		"-i", video,
// 		"-i", audio1,
// 		"-i", subs,
// 		"-codec", "copy", "-codec:s", "mov_text",
// 		"-map", "0:v",
// 		"-map", "1:a", "-metadata:s:a:0", "language=rus",
// 		"-map", "2:s", "-metadata:s:s:0", "language=rus",
// 		fldr.OutPath() + base + "_ar6_sr.mp4",
// 	}
// 	_, _, err := cli.RunConsole(prog, args...)
// 	return err
// }

// func MuxA2E2s(video, audio1, audio2, subs string) error {
// 	base := utils.CommonPrefix(video, audio1, audio2)
// 	base = namedata.RetrieveShortName(base)
// 	prog := "ffmpeg"
// 	args := []string{
// 		"-i", video,
// 		"-i", audio1,
// 		"-i", audio2,
// 		"-i", subs,
// 		"-codec", "copy", "-codec:s", "mov_text",
// 		"-map", "0:v",
// 		"-map", "1:a", "-metadata:s:a:0", "language=rus",
// 		"-map", "2:a", "-metadata:s:a:1", "language=eng",
// 		"-map", "3:s", "-metadata:s:s:0", "language=rus",
// 		fldr.OutPath() + base + "_ar2e2_sr.mp4",
// 	}
// 	_, _, err := cli.RunConsole(prog, args...)
// 	return err
// }

// func MuxA2E6s(video, audio1, audio2, subs string) error {
// 	base := utils.CommonPrefix(video, audio1, audio2)
// 	base = namedata.RetrieveShortName(base)
// 	prog := "ffmpeg"
// 	args := []string{
// 		"-i", video,
// 		"-i", audio1,
// 		"-i", audio2,
// 		"-i", subs,
// 		"-codec", "copy", "-codec:s", "mov_text",
// 		"-map", "0:v",
// 		"-map", "1:a", "-metadata:s:a:0", "language=rus",
// 		"-map", "2:a", "-metadata:s:a:1", "language=eng",
// 		"-map", "3:s", "-metadata:s:s:0", "language=rus",
// 		fldr.OutPath() + base + "_ar2e6_sr.mp4",
// 	}
// 	_, _, err := cli.RunConsole(prog, args...)
// 	return err
// }

// func MuxA6E2s(video, audio1, audio2, subs string) error {
// 	base := utils.CommonPrefix(video, audio1, audio2)
// 	base = namedata.RetrieveShortName(base)
// 	prog := "ffmpeg"
// 	args := []string{
// 		"-i", video,
// 		"-i", audio1,
// 		"-i", audio2,
// 		"-i", subs,
// 		"-codec", "copy", "-codec:s", "mov_text",
// 		"-map", "0:v",
// 		"-map", "1:a", "-metadata:s:a:0", "language=rus",
// 		"-map", "2:a", "-metadata:s:a:1", "language=eng",
// 		"-map", "3:s", "-metadata:s:s:0", "language=rus",
// 		fldr.OutPath() + base + "_ar6e2_sr.mp4",
// 	}
// 	_, _, err := cli.RunConsole(prog, args...)
// 	return err
// }

// func MuxA6E6s(video, audio1, audio2, subs string) error {
// 	base := utils.CommonPrefix(video, audio1, audio2, subs)
// 	base = namedata.RetrieveShortName(base)
// 	prog := "ffmpeg"
// 	args := []string{
// 		"-i", video,
// 		"-i", audio1,
// 		"-i", audio2,
// 		"-i", subs,
// 		"-codec", "copy", "-codec:s", "mov_text",
// 		"-map", "0:v",
// 		"-map", "1:a", "-metadata:s:a:0", "language=rus",
// 		"-map", "2:a", "-metadata:s:a:1", "language=eng",
// 		"-map", "3:s", "-metadata:s:s:0", "language=rus",
// 		fldr.OutPath() + base + "_ar6e6_sr.mp4",
// 	}
// 	_, _, err := cli.RunConsole(prog, args...)
// 	return err
// }

// func ChooseMuxer(task string) ([]string, string, error) {
// 	data := strings.Split(task, " ")
// 	if data[0] == "" {
// 		return []string{}, MuxerSKIP, nil
// 	}
// 	if len(data) < 2 {
// 		return []string{}, MuxerNA, errors.New("muxer not assigned")
// 	}
// 	switch data[1] {
// 	default:
// 		return []string{}, MuxerSKIP, errors.New("muxer not recognised")
// 	case MuxerAR2, MuxerAR6, MuxerAE2, MuxerAE6, MuxerAR2E2, MuxerAR2E6, MuxerA6E2, MuxerA6E6, MuxerAR2s, MuxerA6s, MuxerAR2E2s, MuxerAR2E6s, MuxerA6E2s, MuxerA6E6s:
// 		paths := defineFiles(task)
// 		return paths, data[1], nil
// 	}
// }
