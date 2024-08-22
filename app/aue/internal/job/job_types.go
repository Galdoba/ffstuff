package job

import (
	"fmt"

	"github.com/Galdoba/ffstuff/app/aue/internal/define"
	. "github.com/Galdoba/ffstuff/app/aue/internal/define"
	source "github.com/Galdoba/ffstuff/app/aue/internal/files/sourcefile"
	target "github.com/Galdoba/ffstuff/app/aue/internal/files/targetfile"
	"github.com/Galdoba/ffstuff/app/aue/internal/task"
)

func taskList(jType string) []task.Task {
	tasklist := []task.Task{}
	switch jType {
	default:
		return append(tasklist)
	case define.JOB_V1A2S1:
		tasklist = append(tasklist,
			task.NewTask(TASK_MoveFile),
			task.NewTask(TASK_Make_Dir),
			task.NewTask(TASK_Encode_v1a2),
			task.NewTask(TASK_MoveFile),
		)
	case JOB_V1A2S0:
		tasklist = append(tasklist,
			task.NewTask(TASK_MoveFile),
			task.NewTask(TASK_Make_Dir),
			task.NewTask(TASK_Encode_v1a2),
			task.NewTask(TASK_MoveFile),
		)
	case JOB_V1A1S0:
		tasklist = append(tasklist,
			task.NewTask(TASK_MoveFile),
			task.NewTask(TASK_Make_Dir),
			task.NewTask(TASK_Encode_v1a1),
			task.NewTask(TASK_MoveFile),
		)
	}
	return tasklist
}

func (ja *jobAdmin) setJobCodeAndTargets() error {
	sourceNames := source.Names(ja.source)
	if ja.options.jobType == "" {
		streamsByType := source.MapStreamTypesAll(ja.source)
		code := fmt.Sprintf("V%vA%vS%v", streamsByType[STREAM_VIDEO], streamsByType[STREAM_AUDIO], streamsByType[STREAM_SUBTITLE])
		ja.options.jobType = code

	}
	if err := assertJobCode(ja.options.jobType); err != nil {
		return err
	}

	switch ja.options.jobType {
	case JOB_V1A2S1:
		ja.target = append(ja.target,
			target.New(define.PURPOSE_Output_Video, sourceNames),
			target.New(define.PURPOSE_Output_Audio1, sourceNames),
			target.New(define.PURPOSE_Output_Audio2, sourceNames),
			target.New(define.PURPOSE_Output_Subs, sourceNames),
		)
	case JOB_V1A2S0:
		ja.target = append(ja.target,
			target.New(define.PURPOSE_Output_Video, sourceNames),
			target.New(define.PURPOSE_Output_Audio1, sourceNames),
			target.New(define.PURPOSE_Output_Audio2, sourceNames),
		)
	case JOB_V1A1S0:
		ja.target = append(ja.target,
			target.New(define.PURPOSE_Output_Video, sourceNames),
			target.New(define.PURPOSE_Output_Audio1, sourceNames),
		)
	default:
		return fmt.Errorf("setTargets for job '%v' not implemented", ja.options.jobType)
	}
	return nil
}
