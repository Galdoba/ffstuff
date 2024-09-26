package origin

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type origin struct {
	path       string
	message    string
	killOnDone bool
	score      int
	err        error
}

var originConstructor *constructor

type ConstructorOption func(*constructorOptions)

type constructorOptions struct {
	filePriorityMap map[string]int
	dirPriorityMap  map[string]int
	killMarker      bool
	killAll         bool
	markerExt       string
}

func defaultConstructor() constructorOptions {
	return constructorOptions{
		filePriorityMap: make(map[string]int),
		dirPriorityMap:  make(map[string]int),
	}
}

func ConstructorSetup(options ...ConstructorOption) error {
	originConstructor = &constructor{}
	settngs := defaultConstructor()
	for _, enrich := range options {
		enrich(&settngs)
	}
	originConstructor.filePriorityMap = settngs.filePriorityMap
	originConstructor.dirPriorityMap = settngs.dirPriorityMap
	originConstructor.killMarker = settngs.killMarker
	originConstructor.killAll = settngs.killAll
	originConstructor.markerExt = settngs.markerExt
	if originConstructor.markerExt == "" {
		return fmt.Errorf("marker extention was not provided")
	}
	return nil
}

type constructor struct {
	filePriorityMap map[string]int
	dirPriorityMap  map[string]int
	killMarker      bool
	killAll         bool
	markerExt       string
}

type Origin interface {
	Path() string
	Message() string
	MustDie() bool
}

func New(path string) *origin {
	o := origin{}
	o.path = path
	if strings.HasSuffix(path, originConstructor.markerExt) {
		if originConstructor.killMarker {
			o.killOnDone = true
		}
		bt, err := os.ReadFile(o.path)
		if err != nil {
			o.err = fmt.Errorf("failed to read message")
		}
		switch len(bt) > 10000 {
		case true:
			o.err = fmt.Errorf("dev: message is to long (10k+ bytes)")
		case false:
			o.message = string(bt)
		}
	}
	if originConstructor.killAll {
		o.killOnDone = true
	}
	for key, score := range originConstructor.filePriorityMap {
		if strings.Contains(filepath.Base(o.path), key) {
			o.score += score
		}
	}
	for key, score := range originConstructor.dirPriorityMap {
		if strings.Contains(filepath.Dir(o.path), key) {
			o.score += score
		}
	}
	return &o
}

func (or *origin) Path() string {
	return or.path
}
func (or *origin) Message() string {
	return or.message
}
func (or *origin) MustDie() bool {
	return or.killOnDone
}

func WithFilePriority(priorityMap map[string]int) ConstructorOption {
	return func(co *constructorOptions) {
		for k, v := range priorityMap {
			co.filePriorityMap[k] = v
		}
	}
}

func WithDirectoryPriority(priorityMap map[string]int) ConstructorOption {
	return func(co *constructorOptions) {
		for k, v := range priorityMap {
			co.dirPriorityMap[k] = v
		}
	}
}

func KillMarkers(killSignal bool) ConstructorOption {
	return func(co *constructorOptions) {
		co.killMarker = killSignal
	}
}

func KillAll(killSignal bool) ConstructorOption {
	return func(co *constructorOptions) {
		co.killMarker = killSignal
	}
}

func WithMarkerExt(ext string) ConstructorOption {
	return func(co *constructorOptions) {
		co.markerExt = ext
	}
}
