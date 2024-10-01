package process

import (
	"fmt"

	"github.com/Galdoba/ffstuff/app/grabber/commands/grabberflag"
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
	copyDecidion          string
	DeleteDecidion        string
	sortDecidion          string
}

func New(opts ...ProcessOption) (*Process, error) {
	logman.Debug(logman.NewMessage("start new process"))
	pr := Process{}
	settings := defaultOptions()
	for _, enrich := range opts {
		enrich(&settings)
	}
	pr.mode = settings.mode
	pr.copyDecidion = settings.copy_decidion
	pr.DeleteDecidion = settings.delete_decidion
	pr.sortDecidion = settings.sort_decidion
	logman.Debug(logman.NewMessage("validate process configuration"))
	err := pr.validate()
	if err == nil {
		logman.Debug(logman.NewMessage("new process started"))
	}
	return &pr, err
}

func (pr *Process) validate() error {
	for _, assert := range []func(*Process) (string, error){
		assertMode,
		assertCopyDecidion,
		assertDeleteDecidion,
		assertSortDecidion,
	} {
		msg, err := assert(pr)
		if err != nil {
			return logman.Errorf("validation failed: %v", err)
		}
		logman.Debug(logman.NewMessage(msg))
	}
	return nil
}

func assertMode(pr *Process) (string, error) {
	switch pr.mode {
	case MODE_GRAB, MODE_TRACK:
	default:
		return "", fmt.Errorf("process mode: %v", pr.mode)
	}
	return "process.mode: " + pr.mode, nil
}

func assertCopyDecidion(pr *Process) (string, error) {
	switch pr.copyDecidion {
	case grabberflag.VALUE_COPY_SKIP, grabberflag.VALUE_COPY_RENAME, grabberflag.VALUE_COPY_OVERWRITE:
		return "process.copyDecidion: " + pr.copyDecidion, nil
	default:
		return "", fmt.Errorf("process copyDecidion invalid: %v", pr.copyDecidion)
	}

}

func assertDeleteDecidion(pr *Process) (string, error) {
	switch pr.DeleteDecidion {
	case grabberflag.VALUE_DELETE_NONE, grabberflag.VALUE_DELETE_MARKER, grabberflag.VALUE_DELETE_ALL:
		return "process.DeleteDecidion: " + pr.DeleteDecidion, nil
	default:
		return "", fmt.Errorf("process DeleteDecidion invalid: %v", pr.DeleteDecidion)
	}

}

func assertSortDecidion(pr *Process) (string, error) {
	switch pr.sortDecidion {
	case grabberflag.VALUE_SORT_PRIORITY, grabberflag.VALUE_SORT_SIZE, grabberflag.VALUE_SORT_NONE:
		return "process.sortDecidion: " + pr.sortDecidion, nil
	default:
		return "", fmt.Errorf("process sortDecidion invalid: %v", pr.sortDecidion)
	}

}

func (pr *Process) ShowOrder() {
	fmt.Println(pr.sortDecidion)
}
