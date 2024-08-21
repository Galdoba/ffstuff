package job

import (
	"fmt"

	source "github.com/Galdoba/ffstuff/app/aue/internal/files/sourcefile"
	target "github.com/Galdoba/ffstuff/app/aue/internal/files/targetfile"
	"github.com/Galdoba/ffstuff/app/aue/internal/task"
)

type jobAdmin struct {
	source  []source.SourceFile
	target  []target.TargetFile
	tasks   []task.Task
	options *jobOptions
}

type inputFile struct {
	name string
}

func New(source []source.SourceFile, jobOptions ...JobOptsFunc) (*jobAdmin, error) {
	ja := jobAdmin{}
	ja.source = source
	options := defaultJobOptions()
	for _, enrichWith := range jobOptions {
		enrichWith(&options)
	}
	ja.options = &options
	return &ja, nil
}

func (ja *jobAdmin) DecideType() error {
	return fmt.Errorf("TODO: chosse what job based on input")
}

func (ja *jobAdmin) CompileTasks() error {
	return fmt.Errorf("TODO: generate tasks based on job name")
}

func (ja *jobAdmin) Execute() error {
	return fmt.Errorf("TODO: execute tasks based on job processingMode")
}

/*scenario:
job := job.New(inputPaths, options...)
job.DecideType()
job.CompileTasks()
job.Execute()

*/
