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
	return ct.execFunc(ct.parameter)
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
	case TASK_Signal_Done:
		parameters = append(parameters, TASK_PARAM_NewPath)
		parameters = append(parameters, TASK_PARAM_Text)
	}
	return parameters
}
