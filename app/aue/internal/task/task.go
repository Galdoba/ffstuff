package task

import (
	"fmt"
	"strings"

	. "github.com/Galdoba/ffstuff/app/aue/internal/define"
)

// const ()

// type job struct {
// 	name     string
// 	sequanse []Task
// }

type cliTask struct {
	cliLine   string
	name      string
	parameter map[string]string
	execFunc  func(map[string]string) error
}

func (ct *cliTask) SetParameters(params ...parameterData) error {
	for _, param := range params {
		if err := ct.setParameter(param); err != nil {
			return err
		}
	}
	return nil
}

func (ct *cliTask) setParameter(pInput parameterData) error {
	if _, ok := ct.parameter[pInput.key]; !ok {
		return fmt.Errorf("parameter '%v' is not expected in task '%v'", pInput.key, ct.name)
	}
	if ct.parameter[pInput.key] != unknownParamValue {
		return fmt.Errorf("parameter '%v' was set to '%v': modification is forbidden", pInput.key, ct.parameter[pInput.key])
	}
	ct.parameter[pInput.key] = pInput.val
	return nil
}

type parameterData struct {
	key string
	val string
}

func NewParameterData(key string, val string) parameterData {
	return parameterData{key, val}
}

func (ct *cliTask) Execute() error {
	return fmt.Errorf("TODO")
}

func (ct *cliTask) String() string {
	str := ct.parameter["format"]
	for key, value := range ct.parameter {
		str = strings.ReplaceAll(str, fmt.Sprintf("{%v}", key), value)
	}
	return str
}

type Task interface {
	SetParameters(...parameterData) error
	Execute() error
	String() string
}

func expectedParameters(name string) []string {
	parameters := append([]string{}, "format")
	switch name {
	default:
		return nil
	case TASK_MoveFile, TASK_CopyFile:
		parameters = append(parameters, TASK_PARAM_OldPath)
		parameters = append(parameters, TASK_PARAM_NewPath)
	case TASK_Make_Dir:
		parameters = append(parameters, TASK_PARAM_NewPath)
	case TASK_Encode_v1a1:
		parameters = append(parameters, TASK_PARAM_Encode_input)
		parameters = append(parameters, PURPOSE_Output_Video)
		parameters = append(parameters, PURPOSE_Output_Audio1)
	case TASK_Encode_v1a2:
		parameters = append(parameters, TASK_PARAM_Encode_input)
		parameters = append(parameters, PURPOSE_Output_Video)
		parameters = append(parameters, PURPOSE_Output_Audio1)
		parameters = append(parameters, PURPOSE_Output_Audio2)
	}
	return parameters
}

// func command() {
// 	//declaration
// 	inputFile := ""
// 	fps := ""
// 	d := fmt.Sprintf("ffmpeg -n -r 25 -i %v ", inputFile)
// 	//filter_complex 1aud
// 	fc1 := fmt.Sprintf("-filter_complex \"[0:v:0]setsar=1/1[video]; [0:a:0]aresample=48000,atempo=25/(%v)[audio1]\"", fps)
// 	fc2 := fmt.Sprintf("-filter_complex \"[0:v:0]setsar=1/1[video]; [0:a:0]aresample=48000,atempo=25/(%v)[audio1]; [0:a:1]aresample=48000,atempo=25/(%v) [audio2]\"", fps, fps)
// 	//mappings
// 	outVideo := ""
// 	outAudio1 := ""
// 	outAudio2 := ""
// 	mpVid := fmt.Sprintf("-map \"[video]\" -c:v libx264 -preset medium -crf 21 -pix_fmt yuv422p -profile high -g 0 -map_metadata -1 -map_chapters -1 %v", outVideo)
// 	mpAud1 := fmt.Sprintf("-map \"[aud1]\" -c:a alac -compression_level 0 -map_metadata -1 -map_chapters -1 %v", outAudio1)
// 	mpAud2 := fmt.Sprintf("-map \"[aud2]\" -c:a alac -compression_level 0 -map_metadata -1 -map_chapters -1 %v", outAudio2)
// 	return "ffmpeg -n -r 25 -i %v "
// }

/*
mkdir -p /mnt/pemaltynov/ROOT/IN/_AMEDIA/_DONE/Agenty_vo_vremeni_s01/
mkdir -p /mnt/pemaltynov/ROOT/EDIT/_amedia/Agenty_vo_vremeni_s01/
mv /home/pemaltynov/IN/Agenty_vo_vremeni_s01e09--SER--Agenty_vo_vremeni_s01e09_PRT240807071844_SER_04901_18.mp4 /home/pemaltynov/IN/_IN_PROGRESS/ || exit
fflite -n -r 25 -i /home/pemaltynov/IN/_IN_PROGRESS/Agenty_vo_vremeni_s01e09--SER--Agenty_vo_vremeni_s01e09_PRT240807071844_SER_04901_18.mp4 \
  -filter_complex "[0:v:0]split=2[vidHD][inProxy]; [inProxy]scale=iw/2:ih, setsar=(1/1)*2[vidHD_pr]; [0:a:0]aresample=48000,atempo=25/(25/1),asplit=2[aud1][aud1_pr]; [0:a:1]aresample=48000,atempo=25/(25/1),asplit=2[aud2][aud2_pr]"\
     -map "[vidHD]" -c:v libx264 -preset medium -crf 16 -pix_fmt yuv420p -g 0 -map_metadata -1 -map_chapters -1 /mnt/pemaltynov/ROOT/EDIT/_amedia/Agenty_vo_vremeni_s01/Agenty_vo_vremeni_s01_09_PRT240807071844_HD.mp4\
     -map "[vidHD_pr]" -c:v libx264 -x264opts interlaced=1 -preset superfast -pix_fmt yuv420p  -b:v 2000k -maxrate 2000k /mnt/pemaltynov/ROOT/EDIT/_amedia/Agenty_vo_vremeni_s01/Agenty_vo_vremeni_s01_09_PRT240807071844_HD_proxy.mp4\
     -map "[aud1]" -c:a alac -compression_level 0 -map_metadata -1 -map_chapters -1 /mnt/pemaltynov/ROOT/EDIT/_amedia/Agenty_vo_vremeni_s01/Agenty_vo_vremeni_s01_09_PRT240807071844_AUDIORUS20.m4a\
     -map "[aud1_pr]" -c:a ac3 -b:a 128k /mnt/pemaltynov/ROOT/EDIT/_amedia/Agenty_vo_vremeni_s01/Agenty_vo_vremeni_s01_09_PRT240807071844_AUDIORUS20_proxy.ac3\
     -map "[aud2]" -c:a alac -compression_level 0 -map_metadata -1 -map_chapters -1 /mnt/pemaltynov/ROOT/EDIT/_amedia/Agenty_vo_vremeni_s01/Agenty_vo_vremeni_s01_09_PRT240807071844_AUDIOCHI20.m4a\
     -map "[aud2_pr]" -c:a ac3 -b:a 128k /mnt/pemaltynov/ROOT/EDIT/_amedia/Agenty_vo_vremeni_s01/Agenty_vo_vremeni_s01_09_PRT240807071844_AUDIOCHI20_proxy.ac3\
    \
    \
 2>&1 \
*/
