package task

import (
	"fmt"

	is "github.com/Galdoba/ffstuff/app/aue/internal/define"
	"github.com/Galdoba/ffstuff/pkg/ump"
)

const (
	unknownParamValue = "[UNKNOWN]"
)

func NewTask(name string) (*cliTask, error) {
	ct := cliTask{}
	switch name {
	default:
		return nil, fmt.Errorf("task %v not implemented", name)
	case is.TASK_MoveFile:
	}
	ct.name = name
	ct.parameter = make(map[string]string)
	for _, key := range expectedParameters(ct.name) {
		ct.parameter[key] = unknownParamValue
	}
	return &ct, nil
}

func mediaInfo(path string) (*ump.MediaProfile, error) {
	prf := ump.NewProfile()
	err := prf.ConsumeFile(path)
	return prf, err
}
