package actioncombine

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/Galdoba/ffstuff/app/demuxer/handle"
	"github.com/Galdoba/ffstuff/pkg/namedata"
	"github.com/Galdoba/ffstuff/pkg/spreadsheet"
	"github.com/Galdoba/ffstuff/pkg/spreadsheet/tablemanager"
	"github.com/urfave/cli"
)

func Run(c *cli.Context) (string, error) {
	fmt.Println("--------------------------------------------------------------------------------")
	args := c.Args()
	autoMap, err := autoMap(args)
	handle.Error(err)
	sheet, err := spreadsheet.New()
	handle.Error(err)
	taskList := tablemanager.TaskListFrom(sheet)
	readyNames := []string{}
	for _, task := range taskList.ReadyForDemux() {
		readyNames = append(readyNames, task.Name())
	}
	name := handle.SelectionSingle("Введите название из таблицы: ", readyNames...)
	editName := namedata.TransliterateForEdit(name)
	fmt.Println("Базовое имя:", editName)
	task := taskList.ByName(name)
	propose := tablemanager.ProposeTargetDirectory(taskList, task)
	destination := handle.SelectionSingle("Где должен быть файл конечный файл?: ", []string{`\\nas\ROOT\EDIT\` + propose, `\\nas\root\EDIT\@trailers_temp\`, "[LOCAL]"}...)
	if destination == "[LOCAL]" {
		destination = ""
	}
	destination, _ = filepath.Abs(destination)
	fmt.Println("Путь просчета:", destination)
	langTag := handle.SelectionSingle("Какой должен быть звук? :", []string{"AUDIORUS51", "AUDIOENG51"}...)
	atempo := handle.SelectionSingle("C каким atempo считать файл? :", []string{"(25/1)", "(24/1)", "(24000/1001)"}...)
	line := fmt.Sprintf("-i %v -i %v -i %v -i %v -i %v -i %v -filter_complex amerge=inputs=6,channelmap=channel_layout=5.1,aresample=48000,atempo=25/%v[audio] -map [audio] -c:a alac -compression_level 0 -map_metadata -1 -map_chapters -1 %v\\%v", autoMap["L"], autoMap["R"], autoMap["C"], autoMap["LFE"], autoMap["Ls"], autoMap["Rs"], atempo, destination, editName+"_"+langTag+".m4a")
	return line, nil
}

func autoMap(inputs []string) (map[string]string, error) {
	tagMap := make(map[string]string)
	if len(inputs) != 6 {
		return tagMap, fmt.Errorf("input invalid (expect 6 audio files)")
	}
	wavArgsFound := 0
	for _, arg := range inputs {
		trimmed := strings.TrimSuffix(arg, ".wav")
		if trimmed != arg {
			wavArgsFound++
		}
	}
	if wavArgsFound != 6 {
		return tagMap, fmt.Errorf("input invalid (not .wav files received)")
	}
	tags := []string{".L.", ".R.", ".C.", ".LFE.", ".Ls.", ".Rs."}

	for _, path := range inputs {
		switch {
		case strings.Contains(path, tags[0]):
			tagMap["L"] = path
		case strings.Contains(path, tags[1]):
			tagMap["R"] = path
		case strings.Contains(path, tags[2]):
			tagMap["C"] = path
		case strings.Contains(path, tags[3]):
			tagMap["LFE"] = path
		case strings.Contains(path, tags[4]):
			tagMap["Ls"] = path
		case strings.Contains(path, tags[5]):
			tagMap["Rs"] = path
		}
	}
	if len(tagMap) != 6 {
		return tagMap, fmt.Errorf("channels were not assingned properly\n%v:%v\n%v:%v\n%v:%v\n%v:%v\n%v:%v\n%v:%v", tags[0], tagMap["L"], tags[1], tagMap["R"], tags[2], tagMap["C"], tags[3], tagMap["LFE"], tags[4], tagMap["Ls"], tags[5], tagMap["Rs"])
	}
	return tagMap, nil
}
