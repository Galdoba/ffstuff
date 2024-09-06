package job

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Galdoba/ffstuff/app/aue/internal/bashgen"
	"github.com/Galdoba/ffstuff/app/aue/internal/define"
	"github.com/Galdoba/ffstuff/app/aue/internal/files/bridge"
	source "github.com/Galdoba/ffstuff/app/aue/internal/files/sourcefile"
	target "github.com/Galdoba/ffstuff/app/aue/internal/files/targetfile"
	"github.com/Galdoba/ffstuff/app/aue/internal/task"
	"github.com/Galdoba/ffstuff/pkg/logman"
)

type jobAdmin struct {
	source  []*source.SourceFile
	target  []*target.TargetFile
	tasks   []task.Task
	options *jobOptions
	name    string
}

type inputFile struct {
	name string
}

func New(sources []*source.SourceFile, targets []*target.TargetFile, jobOptions ...JobOptsFunc) (*jobAdmin, error) {
	ja := jobAdmin{}
	ja.source = sources
	ja.target = targets
	options := defaultJobOptions()
	for _, enrichWith := range jobOptions {
		enrichWith(&options)
	}
	ja.options = &options

	return &ja, nil
}

func (ja *jobAdmin) DecideType() error {

	if ja.options.jobType == "" || len(ja.target) == 0 {
		err := ja.setJobCodeAndTargets()
		if err != nil {
			return logman.Errorf("job targets setup failed: %v", err)
		}
	}

	br := bridge.New()
	if err := br.Connect(ja.source, ja.target); err != nil {
		return err
	}
	ja.source = br.Sources()
	ja.target = br.Targets()
	ja.name = br.ProjectBase()
	if err := approveSources(ja.source, ja.options.jobType); err != nil {
		return fmt.Errorf("job decidion failed: %v", err)
	}
	logman.Ping("end")
	return nil
}

func (ja *jobAdmin) CompileTasks() error {

	return ja.setupTaskList()
}

func (ja *jobAdmin) Execute() error {
	for _, err := range []error{
		generateBashFiles(ja),
		processDirectly(ja),
	} {
		if err != nil {
			return fmt.Errorf("job '%v' execution failed: %v", ja.ProjectName(), err)
		}
	}

	return nil
}

func (ja *jobAdmin) TypeDecided() string {
	return ja.options.jobType
}

func generateBashFiles(ja *jobAdmin) error {
	if ja.options.bashGeneration {
		gen := bashgen.New(ja)
		if err := gen.GenerateBash(ja.tasks); err != nil {
			return logman.Errorf("bash generation failed: %v", err)
		}
		logman.Println("bash generated")
	}
	return nil
}

func processDirectly(ja *jobAdmin) error {
	if ja.options.directProcessing {
		for _, t := range ja.tasks {
			if err := t.Execute(); err != nil {
				return fmt.Errorf("task execution failed: %v", t.String())
			}
		}
	}
	return nil
}

func jobCode(targets []target.TargetFile) string {
	goalPresent := make(map[string]bool)
	for _, trg := range targets {
		goalPresent[trg.ClaimedGoal] = true
	}

	taskCodeNums := []int{0, 0, 0}
	if goalPresent[define.PURPOSE_Output_Video] {
		taskCodeNums[0] = 1
	}
	if goalPresent[define.PURPOSE_Output_Audio1] {
		taskCodeNums[1] = 1
	}
	if goalPresent[define.PURPOSE_Output_Audio2] {
		taskCodeNums[1] = 2
	}
	if goalPresent[define.PURPOSE_Output_Subs] {
		taskCodeNums[2] = 1
	}
	code := fmt.Sprintf("V%vA%vS%v")
	return code
}

func assertJobCode(code string) error {
	switch code {
	case define.JOB_V1A0S0,
		define.JOB_V0A1S0,
		define.JOB_V0A2S0,
		define.JOB_V1A1S0,
		define.JOB_V1A2S0,
		define.JOB_V1A2S1,
		define.JOB_V0A0S1:
		return nil
	default:
		return fmt.Errorf("unknown job code received '%v'", code)
	}
}

func approveSources(sources []*source.SourceFile, job string) error {
	streamTypeMap := jobCodeToStreamTypeMap(job)
	inSourcesTypeMap := make(map[string]int)
	for _, source := range sources {
		profile := source.Profile()
		for _, stream := range profile.Streams {
			inSourcesTypeMap[stream.Codec_type]++
		}
	}
	for k, v := range streamTypeMap {
		if inSourcesTypeMap[k] < v {
			return fmt.Errorf("sources have not enough %v streams", k)
		}
	}
	return nil
}

func jobCodeToStreamTypeMap(job string) map[string]int {
	stMap := make(map[string]int)
	job = strings.ReplaceAll(job, "V", "_")
	job = strings.ReplaceAll(job, "A", "_")
	job = strings.ReplaceAll(job, "S", "_")
	parts := strings.Split(job, "_")
	nums := []int{}
	for _, part := range parts {
		v, err := strconv.Atoi(part)
		if err != nil {
			continue
		}
		nums = append(nums, v)
	}
	for i, n := range nums {
		stType := ""
		switch i {
		case 0:
			stType = define.STREAM_VIDEO
		case 1:
			stType = define.STREAM_AUDIO
		case 2:
			stType = define.STREAM_SUBTITLE
		}
		stMap[stType] = n
	}
	return stMap
}

/*scenario:
job := job.New(inputPaths, options...)
job.DecideType()
job.CompileTasks()
job.Execute()

*/
