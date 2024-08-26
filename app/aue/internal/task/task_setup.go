package task

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/Galdoba/devtools/cli/command"
	. "github.com/Galdoba/ffstuff/app/aue/internal/define"
	"github.com/Galdoba/ffstuff/pkg/ump"
)

const (
	unknownParamValue = "[UNKNOWN]"
)

func NewTask(name string) *cliTask {
	ct := cliTask{}

	ct.name = name
	ct.parameter = make(map[string]string)
	for _, key := range expectedParameters(ct.name) {
		ct.parameter[key] = unknownParamValue
	}
	switch name {
	default:
		panic(fmt.Sprintf("task %v not implemented", name))
	case TASK_MoveFile:
		ct.parameter["format"] = fmt.Sprintf("mv {%v} {%v}", TASK_PARAM_OldPath, TASK_PARAM_NewPath)
		ct.execFunc = moveFileFunc
	case TASK_CopyFile:
		ct.parameter["format"] = fmt.Sprintf("cp {%v} {%v}", TASK_PARAM_OldPath, TASK_PARAM_NewPath)
		ct.execFunc = copyFileFunc
	case TASK_Make_Dir:
		ct.parameter["format"] = fmt.Sprintf("mkdir {%v}", TASK_PARAM_NewPath)
		ct.execFunc = makeDirFunc
	case TASK_Encode_v1a1:
		ct.parameter["format"] = fmt.Sprintf("ffmpeg -n -r 25 -i {%v} "+
			"-filter_complex [0:v:0]setsar=(1/1)[vidHD];[0:a:0]aresample=48000,atempo=25/(25/1)[aud1] "+
			"-map [vidHD] -c:v libx264 -preset medium -crf 21 -pix_fmt yuv420p -profile high -g 0 -map_metadata -1 -map_chapters -1 {%v} "+
			"-map [aud1] -c:a alac -compression_level 0 -map_metadata -1 -map_chapters -1 {%v}",
			TASK_PARAM_Encode_input, PURPOSE_Output_Video, PURPOSE_Output_Audio1)
		ct.execFunc = encode_v1a1_Func
	case TASK_Encode_v1a2:
		ct.parameter["format"] = fmt.Sprintf("ffmpeg -n -r 25 -i {%v} "+
			"-filter_complex [0:v:0]setsar=(1/1)[vidHD];[0:a:0]aresample=48000,atempo=25/(25/1)[aud1];[0:a:1]aresample=48000,atempo=25/(25/1)[aud2] "+
			"-map [vidHD] -c:v libx264 -preset medium -crf 21 -pix_fmt yuv420p -profile high -g 0 -map_metadata -1 -map_chapters -1 {%v} "+
			"-map [aud1] -c:a alac -compression_level 0 -map_metadata -1 -map_chapters -1 {%v} "+
			"-map [aud2] -c:a alac -compression_level 0 -map_metadata -1 -map_chapters -1 {%v}",
			TASK_PARAM_Encode_input, PURPOSE_Output_Video, PURPOSE_Output_Audio1, PURPOSE_Output_Audio2)
		ct.execFunc = encode_v1a2_Func
	}

	return &ct
}

// fflite -n -r 25 -i /home/pemaltynov/IN/_IN_PROGRESS/Agenty_vo_vremeni_s01e09--SER--Agenty_vo_vremeni_s01e09_PRT240807071844_SER_04901_18.mp4 \
//   -filter_complex "[0:v:0]split=2[vidHD][inProxy]; [inProxy]scale=iw/2:ih, setsar=(1/1)*2[vidHD_pr]; [0:a:0]aresample=48000,atempo=25/(25/1),asplit=2[aud1][aud1_pr]; [0:a:1]aresample=48000,atempo=25/(25/1),asplit=2[aud2][aud2_pr]"\
//      -map "[vidHD]" -c:v libx264 -preset medium -crf 16 -pix_fmt yuv420p -g 0 -map_metadata -1 -map_chapters -1 /mnt/pemaltynov/ROOT/EDIT/_amedia/Agenty_vo_vremeni_s01/Agenty_vo_vremeni_s01_09_PRT240807071844_HD.mp4\
//      -map "[vidHD_pr]" -c:v libx264 -x264opts interlaced=1 -preset superfast -pix_fmt yuv420p  -b:v 2000k -maxrate 2000k /mnt/pemaltynov/ROOT/EDIT/_amedia/Agenty_vo_vremeni_s01/Agenty_vo_vremeni_s01_09_PRT240807071844_HD_proxy.mp4\
//      -map "[aud1]" -c:a alac -compression_level 0 -map_metadata -1 -map_chapters -1 /mnt/pemaltynov/ROOT/EDIT/_amedia/Agenty_vo_vremeni_s01/Agenty_vo_vremeni_s01_09_PRT240807071844_AUDIORUS20.m4a\
//      -map "[aud1_pr]" -c:a ac3 -b:a 128k /mnt/pemaltynov/ROOT/EDIT/_amedia/Agenty_vo_vremeni_s01/Agenty_vo_vremeni_s01_09_PRT240807071844_AUDIORUS20_proxy.ac3\
//      -map "[aud2]" -c:a alac -compression_level 0 -map_metadata -1 -map_chapters -1 /mnt/pemaltynov/ROOT/EDIT/_amedia/Agenty_vo_vremeni_s01/Agenty_vo_vremeni_s01_09_PRT240807071844_AUDIOCHI20.m4a\
//      -map "[aud2_pr]" -c:a ac3 -b:a 128k /mnt/pemaltynov/ROOT/EDIT/_amedia/Agenty_vo_vremeni_s01/Agenty_vo_vremeni_s01_09_PRT240807071844_AUDIOCHI20_proxy.ac3\

// ffmpeg -n -r 25 -i %v
//		-filter_complex "[0:v:0]setsar=(1/1)[vidHD]; [0:a:0]aresample=48000,atempo=25/(25/1)[aud1]; [0:a:1]aresample=48000,atempo=25/(25/1)[aud2]"\
//      -map "[vidHD]" -c:v libx264 -preset medium -crf 21 -pix_fmt yuv420p -profile high -g 0 -map_metadata -1 -map_chapters -1 %v\
//      -map "[aud1]" -c:a alac -compression_level 0 -map_metadata -1 -map_chapters -1 %v\
//      -map "[aud2]" -c:a alac -compression_level 0 -map_metadata -1 -map_chapters -1 %v\

// ffmpeg -n -r 25 -i %v -filter_complex \"[0:v:0]setsar=(1/1)[vidHD]; [0:a:0]aresample=48000,atempo=25/(25/1)[aud1]\" -map \"[vidHD]\" -c:v libx264 -preset medium -crf 21 -pix_fmt yuv420p -profile high -g 0 -map_metadata -1 -map_chapters -1 %v -map \"[aud1]\" -c:a alac -compression_level 0 -map_metadata -1 -map_chapters -1 %v

func mediaInfo(path string) (*ump.MediaProfile, error) {
	prf := ump.NewProfile()
	err := prf.ConsumeFile(path)
	return prf, err
}

func moveFileFunc(params map[string]string) error {
	for _, key := range expectedParameters(TASK_MoveFile) {
		switch key {
		case TASK_PARAM_OldPath, TASK_PARAM_NewPath:
			if _, ok := params[key]; !ok {
				return fmt.Errorf("parametr '%v' is absent", key)
			}
		default:
			continue
		}
	}
	oldPath := params[TASK_PARAM_OldPath]
	newPath := params[TASK_PARAM_NewPath]
	err := os.Rename(oldPath, newPath)
	if err != nil {
		fmt.Println(err.Error())
		panic(1)
	}
	return nil
}

func makeDirFunc(params map[string]string) error {
	for _, key := range expectedParameters(TASK_Make_Dir) {
		switch key {
		case TASK_PARAM_NewPath:
			if _, ok := params[key]; !ok {
				return fmt.Errorf("parametr '%v' is absent", key)
			}
		default:
			continue
		}
	}
	newPath := params[TASK_PARAM_NewPath]
	return os.MkdirAll(newPath, 0666)
}

func copyFileFunc(params map[string]string) error {
	for _, key := range expectedParameters(TASK_CopyFile) {
		switch key {
		case TASK_PARAM_OldPath, TASK_PARAM_NewPath:
			if _, ok := params[key]; !ok {
				return fmt.Errorf("parametr '%v' is absent", key)
			}
		default:
			continue
		}
	}
	defer os.Rename(params[TASK_PARAM_NewPath]+".tmp", params[TASK_PARAM_NewPath])
	oldPath, err := os.OpenFile(params[TASK_PARAM_OldPath], os.O_RDWR, 0666)
	if err != nil {
		return fmt.Errorf("oldpath creation failed: %v", err)
	}
	defer oldPath.Close()
	newPath, err := os.Create(params[TASK_PARAM_NewPath] + ".tmp")
	if err != nil {
		return fmt.Errorf("newpath creation failed: %v", err)
	}
	_, err = io.Copy(newPath, oldPath)
	if err != nil {
		return fmt.Errorf("copying failed: %v", err)
	}
	defer newPath.Close()
	return nil
}

func encode_v1a1_Func(params map[string]string) error {
	for _, key := range expectedParameters(TASK_Encode_v1a1) {
		switch key {
		case TASK_PARAM_Encode_input, PURPOSE_Output_Video, PURPOSE_Output_Audio1, "format":
			if _, ok := params[key]; !ok {
				return fmt.Errorf("parametr '%v' is absent", key)
			}
		default:
			continue
		}
	}
	fmt.Println(params["format"])
	text := params["format"]
	cmd := exec.Command(text)
	return cmd.Run()
}

func encode_v1a2_Func(params map[string]string) error {
	for _, key := range expectedParameters(TASK_Encode_v1a2) {
		switch key {
		case TASK_PARAM_Encode_input, PURPOSE_Output_Video, PURPOSE_Output_Audio1, PURPOSE_Output_Audio2, "format":
			if _, ok := params[key]; !ok {
				return fmt.Errorf("parametr '%v' is absent", key)
			}
		default:
			continue
		}
	}

	text := params["format"]
	for _, value := range []string{TASK_PARAM_Encode_input, PURPOSE_Output_Video, PURPOSE_Output_Audio1, PURPOSE_Output_Audio2} {

		text = strings.ReplaceAll(text, fmt.Sprintf("{%v}", value), params[value])
	}

	cmnd, err := command.New(command.CommandLineArguments(fmt.Sprintf("%v", text)),
		command.Set(command.TERMINAL_ON),
	)
	if err != nil {
		fmt.Println("create command err:", err.Error())
	}
	err = cmnd.Run()

	// cmd := exec.Command("ffmpeg", argsCreated...)
	// err := cmd.Run()
	if err != nil {
		fmt.Println(err.Error())
		panic("v1a2")
	}
	return nil
}
