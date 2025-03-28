package origin

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Galdoba/ffstuff/app/grabber/commands/grabberflag"
)

type origin struct {
	path       string
	message    string
	group      string
	killOnDone bool
	score      int
	isMarker   bool
	err        error
}

var originConstructor *constructor

type ConstructorOption func(*constructorOptions)

type constructorOptions struct {
	filePriorityMap map[string]int
	dirPriorityMap  map[string]int
	killSignal      string
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
	originConstructor.killSignal = settngs.killSignal
	originConstructor.markerExt = settngs.markerExt
	if originConstructor.markerExt == "" {
		return fmt.Errorf("marker extention was not provided")
	}
	return nil
}

type constructor struct {
	filePriorityMap map[string]int
	dirPriorityMap  map[string]int
	killSignal      string
	markerExt       string
}

type Origin interface {
	Path() string
	Name() string
	Message() string
	Group() string
	MustDie() bool
	Score() int
	IsMarker() bool
	MarshalJSON() ([]byte, error)
	UnmarshalJSON([]byte) error
}

func New(path string, group ...string) *origin {
	o := origin{}
	o.path = path

	if strings.HasSuffix(path, originConstructor.markerExt) {
		o.isMarker = true
		if originConstructor.killSignal == grabberflag.VALUE_DELETE_MARKER {
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
	if originConstructor.killSignal == grabberflag.VALUE_DELETE_ALL {
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
	for _, gr := range group {
		o.group = gr
	}
	return &o
}

func (or *origin) Path() string {
	return or.path
}

func (or *origin) Name() string {
	return filepath.Base(or.path)
}

func (or *origin) Message() string {
	return or.message
}
func (or *origin) Group() string {
	return or.group
}
func (or *origin) MustDie() bool {
	return or.killOnDone
}
func (or *origin) Score() int {
	return or.score
}
func (or *origin) IsMarker() bool {
	return or.isMarker
}

func (or *origin) MarshalJSON() ([]byte, error) {
	e := Export(or)
	return json.MarshalIndent(e, "", "  ")
}

func (or *origin) UnmarshalJSON(data []byte) error {
	e := Export(or)
	if err := json.Unmarshal(data, &e); err != nil {
		return err
	}
	i := Inject(*e)
	or = i
	//*or = origin(s)
	return nil
}

/*
func (receiver *Mode) UnmarshalJSON(data []byte) error {
    *receiver = Mode(data[1 : len(data)-1])
    return nil
}
*/

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

func KillSignal(killSignal string) ConstructorOption {
	return func(co *constructorOptions) {
		co.killSignal = killSignal
	}
}

func WithMarkerExt(ext string) ConstructorOption {
	return func(co *constructorOptions) {
		co.markerExt = ext
	}
}

type ExportOrigin struct {
	Path       string
	Message    string
	Group      string
	KillOnDone bool
	Score      int
	IsMarker   bool
	Err        error
}

func Export(o Origin) *ExportOrigin {
	e := ExportOrigin{}
	e.Path = o.Path()
	e.Message = o.Message()
	e.Group = o.Group()
	e.KillOnDone = o.MustDie()
	e.Score = o.Score()
	e.IsMarker = o.IsMarker()
	return &e
}

func Inject(e ExportOrigin) *origin {
	o := origin{}
	o.path = e.Path
	o.message = e.Message
	o.group = e.Group
	o.killOnDone = e.KillOnDone
	o.score = e.Score
	o.isMarker = e.IsMarker
	return &o
}
