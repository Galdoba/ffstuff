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
	fmt.Println("creating targets for", ja.options.jobType)

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

func (ja *jobAdmin) setupTaskList() error {
	input := make(map[string]*source.SourceFile)
	for _, src := range ja.source {
		input[src.Purpose()] = src
	}
	output := make(map[string]*target.TargetFile)
	for _, tgt := range ja.target {
		output[tgt.ClaimedGoal] = tgt
		fmt.Println("TARGET ADDED:")
		fmt.Println(tgt.Details())
	}
	BUFFER_IN := ja.options.inputDir
	IN_PROGRESS := ja.options.processingDir
	EDIT := ja.options.outDir
	DONE := ja.options.doneDir

	switch ja.options.jobType {
	default:
		panic(ja.options.jobType + " is unexpected")
	case JOB_V1A1S0, JOB_V1A2S0, JOB_V1A2S1:
		//move all inputs to progress
		for _, src := range input {
			oldPath := BUFFER_IN + src.Name()
			newPath := IN_PROGRESS + src.Name()
			ja.tasks = append(ja.tasks, taskMove(oldPath, newPath))
		}
		//encode
		inputPaths := setupPaths(IN_PROGRESS, input[PURPOSE_Input_Media].Name())
		outputsPaths := setupOutputPaths(EDIT, output)
		encodeTask := taskEncode(ja.options.jobType, inputPaths[0], outputsPaths...)
		ja.tasks = append(ja.tasks, encodeTask)
		//copy srt
		for _, target := range ja.target {
			if target.ClaimedGoal != PURPOSE_Output_Subs {
				continue
			}
			oldPath := IN_PROGRESS + input[PURPOSE_Input_Subs].Name()
			newPath := EDIT + target.ExpectedName
			ja.tasks = append(ja.tasks, taskCopy(oldPath, newPath))
		}

		//move input to done
		for _, src := range input {
			oldPath := IN_PROGRESS + src.Name()
			newPath := DONE + src.Name()
			ja.tasks = append(ja.tasks, taskMove(oldPath, newPath))
		}
	}
	return nil
}

func taskMove(old, new string) task.Task {
	tskMove := task.NewTask(TASK_MoveFile)
	tskMove.SetParameters(
		task.NewParameterData(TASK_PARAM_OldPath, old),
		task.NewParameterData(TASK_PARAM_NewPath, new),
	)
	return tskMove
}

func taskEncode(job, inputPath string, outputPaths ...string) task.Task {
	encodeType := jobToTaskEncodeType(job)
	encodeTask := task.NewTask(encodeType)
	encodeTask.SetParameters(task.NewParameterData(TASK_PARAM_Encode_input, inputPath))
	params := []string{PURPOSE_Output_Video, PURPOSE_Output_Audio1, PURPOSE_Output_Audio2}
	for i, param := range params {
		encodeTask.SetParameters(task.NewParameterData(param, outputPaths[i]))
	}
	return encodeTask
}

func taskCopy(old, new string) task.Task {
	tskMove := task.NewTask(TASK_CopyFile)
	tskMove.SetParameters(
		task.NewParameterData(TASK_PARAM_OldPath, old),
		task.NewParameterData(TASK_PARAM_NewPath, new),
	)
	return tskMove
}

func jobToTaskEncodeType(job string) string {
	switch job {
	default:
		panic(fmt.Sprintf("job_types undefined job %v", job))
	case JOB_V1A1S0:
		return TASK_Encode_v1a1
	case JOB_V1A2S0, JOB_V1A2S1:
		return TASK_Encode_v1a2
	}
}

func setupPaths(dir string, files ...string) []string {
	paths := []string{}
	for _, file := range files {
		paths = append(paths, dir+file)
	}
	return paths
}

func setupOutputPaths(dir string, output map[string]*target.TargetFile) []string {
	paths := []string{}
	purposes := []string{PURPOSE_Output_Video, PURPOSE_Output_Audio1, PURPOSE_Output_Audio2, PURPOSE_Output_Subs}
	for _, purpose := range purposes {
		file := output[purpose]
		if file != nil {
			paths = append(paths, dir+file.ExpectedName)
		}
	}
	return paths
}
