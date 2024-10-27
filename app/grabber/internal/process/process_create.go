package process

import (
	"fmt"

	"github.com/Galdoba/ffstuff/app/grabber/commands/grabberflag"
	"github.com/Galdoba/ffstuff/app/grabber/internal/origin"
	"github.com/Galdoba/ffstuff/app/grabber/internal/validation"
	"github.com/Galdoba/ffstuff/pkg/logman"
)

/*
Stages:
0	Setup
1	CheckTrigger
2	WaitTrigger
3	ScanRoot
4	CollectPotential
5	CreateList
6	SortList
7	FilterList
8	Grab
9	Report

*/

type Process struct {
	activeStage           int
	nextStage             int
	timeoutTriggerSeconds int
	cronShedule           string
	nextSheduleTrigger    int
	mode                  string
	CopyDecidion          string
	DeleteDecidion        string
	SortDecidion          string
	KeepMarkerGroups      bool
	DestinationDir        string
	SourceTargetMap       map[origin.Origin]string
}

func New(opts ...ProcessOption) (*Process, error) {
	logman.Debug(logman.NewMessage("start new process"))
	pr := Process{}
	settings := defaultOptions()
	for _, enrich := range opts {
		enrich(&settings)
	}
	pr.SourceTargetMap = make(map[origin.Origin]string)
	pr.mode = settings.mode
	pr.CopyDecidion = settings.copy_decidion
	pr.DeleteDecidion = settings.delete_decidion
	pr.SortDecidion = settings.sort_decidion
	pr.KeepMarkerGroups = settings.keepmarkerGroups
	pr.DestinationDir = settings.destination
	logman.Debug(logman.NewMessage("validate process configuration"))
	err := pr.validate()
	if err == nil {
		logman.Debug(logman.NewMessage("new process started"))
	}
	return &pr, err
}

func (pr *Process) validate() error {
	for _, assert := range []func(*Process) (logman.Message, error){
		assertMode,
		assertCopyDecidion,
		assertDeleteDecidion,
		assertSortDecidion,
		assertDestination,
	} {
		msg, err := assert(pr)
		if err != nil {
			logman.ProcessMessage(msg, logman.DEBUG)
			return fmt.Errorf("validation failed: %v", err)
		}
		logman.Debug(msg)
	}
	return nil
}

func assertMode(pr *Process) (logman.Message, error) {
	switch pr.mode {
	case MODE_GRAB, MODE_TRACK:
	default:
		return nil, fmt.Errorf("process mode: %v", pr.mode)
	}
	return logman.NewMessage("process.mode: %v", pr.mode), nil
}

func assertCopyDecidion(pr *Process) (logman.Message, error) {
	switch pr.CopyDecidion {
	case grabberflag.VALUE_COPY_SKIP, grabberflag.VALUE_COPY_RENAME, grabberflag.VALUE_COPY_OVERWRITE:
		return logman.NewMessage("process.CopyDecidion: %v", pr.CopyDecidion), nil
	default:
		return nil, fmt.Errorf("process CopyDecidion invalid: %v", pr.CopyDecidion)
	}

}

func assertDeleteDecidion(pr *Process) (logman.Message, error) {
	switch pr.DeleteDecidion {
	case grabberflag.VALUE_DELETE_NONE, grabberflag.VALUE_DELETE_MARKER, grabberflag.VALUE_DELETE_ALL:
		return logman.NewMessage("process.DeleteDecidion: %v", pr.DeleteDecidion), nil
	default:
		return nil, fmt.Errorf("process DeleteDecidion invalid: %v", pr.DeleteDecidion)
	}

}

func assertSortDecidion(pr *Process) (logman.Message, error) {
	switch pr.SortDecidion {
	case grabberflag.VALUE_SORT_PRIORITY, grabberflag.VALUE_SORT_SIZE, grabberflag.VALUE_SORT_NONE:
		return logman.NewMessage("process.SortDecidion: %v", pr.SortDecidion), nil
	default:
		return nil, fmt.Errorf("process SortDecidion invalid: %v", pr.SortDecidion)
	}

}

func assertDestination(pr *Process) (logman.Message, error) {
	if err := validation.DirectoryValidation(pr.DestinationDir); err != nil {
		return nil, err
	}
	return logman.NewMessage("process.DestinationDir: %v", pr.DestinationDir), nil

}

func (pr *Process) ShowOrder() {
	fmt.Println(pr.SortDecidion)
}
