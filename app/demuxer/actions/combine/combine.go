package combine

import (
	"fmt"
	"strings"

	"github.com/Galdoba/ffstuff/pkg/spreadsheet"
	"github.com/Galdoba/ffstuff/pkg/spreadsheet/tablemanager"
	"github.com/macroblock/imed/pkg/translit"
	"github.com/urfave/cli"
	"gopkg.in/AlecAivazis/survey.v1"
)

func Combine(c *cli.Context) (string, error) {
	args := c.Args()
	autoMap, err := autoMap(args)
	if err != nil {
		return "", err
	}
	validator := survey.ComposeValidators()
	sheet, err := spreadsheet.New()
	if err != nil {
		return "", err
	}
	taskList := tablemanager.TaskListFrom(sheet)
	readyNames := []string{}
	for _, task := range taskList.ReadyForDemux() {
		readyNames = append(readyNames, task.Name())
	}
	name, err := askSelection(validator, "Введите имя конечного файла: ", readyNames)
	if err != nil {
		return "", err
	}
	task := taskList.ByName(name)
	propose := tablemanager.ProposeTargetDirectory(taskList, task)
	destination, err := askSelection(validator, "Где должен быть файл конечный файл?: ", []string{`\\nas\ROOT\EDIT\` + propose})
	if err != nil {
		return "", err
	}
	atempo, err := askSelection(validator, "C каким atempo считать файл? :", []string{"(25/1)", "(24/1)", "(24000/1001)"})
	if err != nil {
		return "", err
	}
	line := fmt.Sprintf("-i %v -i %v -i %v -i %v -i %v -i %v -filter_complex amerge=inputs=6,channelmap=channel_layout=5.1,aresample=48000,atempo=25/%v[audio] -map [audio] -c:a alac -compression_level 0 -map_metadata -1 -map_chapters -1 %v%v.m4a", autoMap["L"], autoMap["R"], autoMap["C"], autoMap["LFE"], autoMap["Ls"], autoMap["Rs"], atempo, destination, preapareFileName(name)+"_AUDIORUS51")
	return line, nil
}

func preapareFileName(name string) string {
	nameS := strings.SplitN(name, " (", 1)[0]
	nameT, err := translit.Do(nameS)
	if err != nil {
		return name
	}
	return nameT
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
		return tagMap, fmt.Errorf("channels were not assingned properly")
	}
	return tagMap, nil // SUPERSTAR_TRL_5.1_MIX_86.0dB(Int).C.wav 	SUPERSTAR_TRL_5.1_MIX_86.0dB(Int).L.wav 	SUPERSTAR_TRL_5.1_MIX_86.0dB(Int).LFE.wav 	SUPERSTAR_TRL_5.1_MIX_86.0dB(Int).Ls.wav  	SUPERSTAR_TRL_5.1_MIX_86.0dB(Int).R.wav	SUPERSTAR_TRL_5.1_MIX_86.0dB(Int).Rs.wav
}

func askSelection(validator survey.Validator, message string, options []string) (string, error) {
	atempo := ""
	promptSelect := &survey.Select{
		Message: message,
		Options: append(options, "Иное (Ввести руками)"),
	}
	if err := survey.AskOne(promptSelect, &atempo, validator); err != nil {
		return atempo, err
	}
	if atempo != "Иное (Ввести руками)" {
		return atempo, nil
	}

	return askInput(validator, "Введите atempo: ")
}

func askInput(val survey.Validator, message string) (string, error) {
	result := ""
	promptInput := &survey.Input{
		Message: message,
	}
	if err := survey.AskOne(promptInput, &result, val); err != nil {
		return result, err
	}
	return result, nil
}
