package task

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/Galdoba/devtools/cli/command"
	. "github.com/Galdoba/ffstuff/app/aue/internal/define"
)

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

func makeFile(params map[string]string) error {
	for _, key := range expectedParameters(TASK_PARAM_NewPath) {
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
	f, err := os.Create(newPath)
	if err != nil {
		return fmt.Errorf("file creation failed: %v")
	}
	defer f.Close()
	return os.MkdirAll(newPath, 0666)
}
