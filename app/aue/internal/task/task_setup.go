package task

import (
	"fmt"

	. "github.com/Galdoba/ffstuff/app/aue/internal/define"
)

const (
	unknownParamValue = "[UNKNOWN]"
	bashFormat        = "format"
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
		ct.parameter[bashFormat] = fmt.Sprintf("mv {%v} {%v}", TASK_PARAM_OldPath, TASK_PARAM_NewPath)
		ct.execFunc = moveFileFunc
	case TASK_CopyFile:
		ct.parameter[bashFormat] = fmt.Sprintf("cp {%v} {%v}", TASK_PARAM_OldPath, TASK_PARAM_NewPath)
		ct.execFunc = copyFileFunc
	case TASK_Make_Dir:
		ct.parameter[bashFormat] = fmt.Sprintf("mkdir {%v}", TASK_PARAM_NewPath)
		ct.execFunc = makeDirFunc
	case TASK_Encode_v1a1:
		ct.parameter[bashFormat] = fmt.Sprintf("ffmpeg -n -r 25 -i {%v} "+
			"-filter_complex [0:v:0]setsar=(1/1)[vidHD];[0:a:0]aresample=48000,atempo=25/(25/1)[aud1] "+
			"-map [vidHD] -c:v libx264 -preset medium -crf 21 -pix_fmt yuv420p -profile high -g 0 -map_metadata -1 -map_chapters -1 {%v} "+
			"-map [aud1] -c:a alac -compression_level 0 -map_metadata -1 -map_chapters -1 {%v}",
			TASK_PARAM_Encode_input, PURPOSE_Output_Video, PURPOSE_Output_Audio1)
		ct.execFunc = encode_v1a1_Func
	case TASK_Encode_v1a2:
		ct.parameter[bashFormat] = fmt.Sprintf("ffmpeg -n -r 25 -i {%v} "+
			"-filter_complex [0:v:0]setsar=(1/1)[vidHD];[0:a:0]aresample=48000,atempo=25/(25/1)[aud1];[0:a:1]aresample=48000,atempo=25/(25/1)[aud2] "+
			"-map [vidHD] -c:v libx264 -preset medium -crf 21 -pix_fmt yuv420p -profile high -g 0 -map_metadata -1 -map_chapters -1 {%v} "+
			"-map [aud1] -c:a alac -compression_level 0 -map_metadata -1 -map_chapters -1 {%v} "+
			"-map [aud2] -c:a alac -compression_level 0 -map_metadata -1 -map_chapters -1 {%v}",
			TASK_PARAM_Encode_input, PURPOSE_Output_Video, PURPOSE_Output_Audio1, PURPOSE_Output_Audio2)
		ct.execFunc = encode_v1a2_Func
	case TASK_Signal_Done:
		//создать ready file
		ct.parameter[bashFormat] = fmt.Sprintf("touch %v",
			TASK_PARAM_NewPath)
		ct.execFunc = makeFile

	}

	return &ct
}

// func mediaInfo(path string) (*ump.MediaProfile, error) {
// 	prf := ump.NewProfile()
// 	err := prf.ConsumeFile(path)
// 	return prf, err
// }
