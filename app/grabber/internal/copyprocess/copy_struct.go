package copyprocess

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/Galdoba/ffstuff/app/grabber/internal/origin"
	"github.com/Galdoba/ffstuff/pkg/logman"
	"github.com/Galdoba/ffstuff/pkg/stdpath"
	"github.com/fatih/color"
)

const (
	Status_Pending = 0
	sc             = "\u001B7"
	rc             = "\u001B8"
)

type copyActionState struct {
	SourcePaths          []origin.Origin
	Destination          string
	MarkerExt            string
	CreationErrors       []error
	DeleteAll            bool
	DeleteMarker         bool
	SourceTargetMap      map[origin.Origin]string `json:"-"`
	AdapterSources       map[int]origin.ExportOrigin
	AdapterTargets       map[int]string
	NamesLen             int
	SourcesVolume        int64
	DownloadedVolume     int64
	DownloadedVolumeLast int64
	Completed            []string
	KillList             []string
	Errors               []error
}

type OriginData struct {
	Path string
	Name string
}

type CopyProcess interface {
	Start() error
	ErrorReport() error
	AddToQueue() error
}

func NewCopyAction(sourceTargetMap map[origin.Origin]string, opts ...Option) *copyActionState {
	cas := copyActionState{}
	// cas.source = src

	settings := defaultCopyOptions()
	for _, modify := range opts {
		modify(&settings)
	}
	cas.Destination = settings.destination
	cas.MarkerExt = settings.markerExt
	cas.DeleteAll = settings.deleteAll
	cas.DeleteMarker = settings.deleteMarkers
	cas.SourceTargetMap = sourceTargetMap
	cas.DownloadedVolumeLast = -1
	for _, src := range settings.sourcePaths {
		// fmt.Println("---")
		// originFile := origin.New(src)

		cas.SourcePaths = append(cas.SourcePaths, src)
	}

	return &cas
}

func (cas *copyActionState) Start() error {
	startTime := time.Now()
	logman.Info("begin transfert")
	copyErrors := newErrorCollector()
	for _, src := range cas.SourcePaths {
		if cas.NamesLen < len(src.Name()) {
			tgt := cas.SourceTargetMap[src]
			tgt = filepath.Base(tgt)
			cas.NamesLen = len(tgt)
		}
		si, err := os.Stat(src.Path())
		if err != nil {
			return fmt.Errorf("failed to get %v size", src.Name())
		}
		cas.SourcesVolume += si.Size()
	}
	fmt.Println("target directory:", cas.Destination)
	fmt.Print(sc)
	fmt.Println(cas.copyProcessData())

	for _, src := range cas.SourcePaths {
		source := src.Path()
		srcInfo, errS := os.Stat(source)
		if errS != nil {
			cas.Errors = append(cas.Errors, logman.Errorf("failed to get info on %v: %v", source, errS))
			continue
		}
		if !srcInfo.Mode().IsRegular() { // cannot copy non-regular files (e.g., directories, symlinks, devices, etc.)
			cas.Errors = append(cas.Errors, fmt.Errorf("cannot copy non-regular files (directories, symlinks, devices, etc.): "+source+" ("+srcInfo.Mode().String()+")"))
			continue
		}
		sourceSize := srcInfo.Size()
		target := ""
		for k, v := range cas.SourceTargetMap {
			if k.Path() == src.Path() {
				target = v
				break
			}
		}

		go copyContent(src.Path(), target, copyErrors)
		doneCopying := false
		goodDone := false
		time.Sleep(time.Second)
		for !doneCopying {
			copyFile, err := os.Stat(target)
			if err != nil {
				fmt.Println(err, "----------------")
			}
			copySize := copyFile.Size()
			fmt.Print(rc + sc)
			fmt.Println(cas.copyProcessData())
			if err != nil {
				logman.Error(err)
				doneCopying = true
			}
			time.Sleep(time.Millisecond * 1000)
			if copySize >= sourceSize {
				doneCopying = true
				goodDone = true
			}
		}
		if goodDone {
			cas.Completed = append(cas.Completed, src.Path())
			if src.MustDie() {
				cas.KillList = append(cas.KillList, src.Path())
			}
		} else {
			cas.Errors = append(cas.Errors, logman.Errorf("failed to grab %v", source))
		}

	}
	fmt.Print(rc + sc)
	fmt.Println(cas.copyProcessData() + "\n")
	for _, path := range cas.KillList {
		switch os.Remove(path) {
		case nil:
			logman.Info("deleted: %v", path)
		default:
			cas.Errors = append(cas.Errors, fmt.Errorf("failed to delete %v", path))
		}
	}
	dur := time.Since(startTime)
	tookTime := dur.String()

	allErrors := append(copyErrors.collected, cas.Errors...)
	if len(allErrors) > 0 {
		fmt.Println(color.RedString("grabbing errors:"))
		for _, err := range allErrors {
			fmt.Printf("  %v\n", err)
		}

		for _, err := range allErrors {
			logman.Debug(logman.NewMessage("processing error intercepted: %v", err))
		}
		// return logman.Errorf("%v errors intecepted, while processing", len(copyErrors.collected))
	}
	logman.Info("transfert complete (process took %v)", tookTime)
	if len(cas.Completed) < len(cas.SourcePaths) {
		return logman.Warn("grabbed %v of %v files", len(cas.Completed), len(cas.SourcePaths))
	}
	if len(cas.Errors) > 0 {
		return cas.ErrorReport()
	}

	return nil
}

type errorCollector struct {
	collected []error
}

func newErrorCollector() *errorCollector {
	return &errorCollector{}
}

func (cas *copyActionState) copyProcessData() string {
	cas.DownloadedVolumeLast = cas.DownloadedVolume
	cas.DownloadedVolume = 0
	out := "process report:\n"
	for _, src := range cas.SourcePaths {
		//tgt := cas.SourceTargetMap[src]
		for k, _ := range cas.SourceTargetMap {
			if k.Path() == src.Path() {
				src = k
			}
		}
		tgt := cas.SourceTargetMap[src]
		progress, trueSize := currentProgress(src.Path(), tgt)
		progressString := formatProgress(progress)
		tgtName := filepath.Base(tgt)
		mustDie := " "
		if src.MustDie() {
			mustDie = "*"
			if progressString == "NaN%" || progressString == color.GreenString("ok") {
				progressString = "dead"
			}
		}
		out += fmt.Sprintf("  %v    %v%v           \n", wideName(tgtName, cas.NamesLen), progressString, mustDie)
		cas.DownloadedVolume += trueSize
	}
	if len(cas.SourcePaths) == 0 {
		out += "  nothing to grab"
		return out
	}

	out += "summary:\n"
	out += fmt.Sprintf("  %v        ", summaryString(cas.SourcesVolume, cas.DownloadedVolume, cas.DownloadedVolumeLast))
	return out
}

func wideName(name string, width int) string {
	for len(name) < width {
		name += " "
	}
	return name
}

func currentProgress(src, tgt string) (float64, int64) {
	si, err := os.Stat(src)
	if err != nil {
		return math.NaN(), 0
	}
	sSize := si.Size()
	ti, err := os.Stat(tgt)
	if err != nil {
		return 0, 0
	}
	tSize := ti.Size()
	if sSize == tSize {
		return 100, sSize
	}
	prgr := calculateProgress(sSize, tSize)
	return prgr * 100, tSize
}

func calculateProgress(sSize, tSize int64) float64 {
	prgr := float64(tSize) / float64(sSize)
	return prgr
}

func formatProgress(progress float64) string {
	progressString := strconv.FormatFloat(progress, 'f', 2, 64) + `%`
	switch progress {
	case 0:
		progressString = "wait"
	case 100:
		progressString = color.GreenString("done")
	case math.NaN():
		progressString = color.RedString("N/A")
	}
	return progressString
}

func summaryString(sourceSize, downloadedNow, downloadedLast int64) string {
	prog := calculateProgress(sourceSize, downloadedNow) * 100
	progStr := formatProgress(prog)
	downStr := size2GbString(downloadedNow) + "/" + size2GbString(sourceSize) + " Gb"
	speed := fmt.Sprintf("%v kb/s", (downloadedNow-downloadedLast)/1024)
	return fmt.Sprintf("total progress: %v  %v  speed=%v", progStr, downStr, speed)
}

func size2GbString(bts int64) string {
	gbt := float64(bts) / 1073741824.0
	gbtStr := strconv.FormatFloat(gbt, 'f', 2, 64)
	return gbtStr
}

func (cas *copyActionState) ErrorReport() error {
	if len(cas.SourcePaths) == 0 {
		return fmt.Errorf("process is empty")
	}
	if len(cas.Errors) == 0 {
		return nil
	}
	text := "processing error(s) detected:\n"
	for i, err := range cas.Errors {
		text += fmt.Sprintf("%v: %v\n", i+1, err)
	}
	return errors.New(text)
}

/*
create
*/

func (cas *copyActionState) AddToQueue() error {
	id := time.Now().UnixNano()
	name := fmt.Sprintf("%v.json", id)
	procFile := stdpath.ProgramDir("queue") + name
	f, err := os.Create(procFile)

	cas.AdapterSources = make(map[int]origin.ExportOrigin)
	cas.AdapterTargets = make(map[int]string)
	i := 0
	for k, v := range cas.SourceTargetMap {
		e := origin.Export(k)
		cas.AdapterSources[i] = *e
		cas.AdapterTargets[i] = v
		i++
	}
	defer f.Close()
	if err != nil {
		return logman.Errorf("process file creation failed: %v", err)
	}
	bt, err := json.MarshalIndent(cas, "", "  ")
	if err != nil {
		return logman.Errorf("process marshaling failed: %v", err)
	}
	err = os.WriteFile(procFile, bt, 0660)
	if err != nil {
		return logman.Errorf("process file writing failed: %v", err)
	}
	return nil
}

func Reconstruct(path string) (CopyProcess, error) {
	bt, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read process file: %v", err)
	}
	cp := &copyActionState{}
	cp.AdapterSources = make(map[int]origin.ExportOrigin)
	cp.AdapterTargets = make(map[int]string)
	cp.SourceTargetMap = make(map[origin.Origin]string)

	if err := cp.UnmarshalJSON(bt); err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %v", err)
	}
	cp.AdapterSources = nil
	cp.AdapterTargets = nil

	return cp, nil
}

type adapter struct {
	SourcePaths          []origin.ExportOrigin
	Destination          string
	MarkerExt            string
	CreationErrors       []error
	DeleteAll            bool
	DeleteMarker         bool
	AdapterSources       map[int]origin.ExportOrigin
	AdapterTargets       map[int]string
	NamesLen             int
	SourcesVolume        int64
	DownloadedVolume     int64
	DownloadedVolumeLast int64
	Completed            []string
	KillList             []string
	Errors               []error
}

func (receiver *copyActionState) UnmarshalJSON(data []byte) error {
	//var s string
	//fmt.Println("start")
	//type Alias *copyActionState
	cp := adapter{}
	cp.AdapterSources = make(map[int]origin.ExportOrigin)
	cp.AdapterTargets = make(map[int]string)
	if err := json.Unmarshal(data, &cp); err != nil {
		return err
	}
	for i := 0; i <= len(cp.AdapterSources)-1; i++ {
		e := origin.Inject(cp.AdapterSources[i])

		receiver.SourceTargetMap[e] = cp.AdapterTargets[i]
	}
	for _, e := range cp.SourcePaths {
		receiver.SourcePaths = append(receiver.SourcePaths, origin.Inject(e))
	}
	//receiver.AdapterSources = cp.AdapterSources
	//receiver.AdapterTargets = cp.AdapterTargets
	receiver.Destination = cp.Destination
	receiver.MarkerExt = cp.MarkerExt
	receiver.CreationErrors = cp.CreationErrors
	receiver.DeleteAll = cp.DeleteAll
	receiver.DeleteMarker = cp.DeleteMarker
	receiver.NamesLen = cp.NamesLen
	receiver.SourcesVolume = cp.SourcesVolume
	receiver.DownloadedVolume = cp.DownloadedVolume
	receiver.DownloadedVolumeLast = cp.DownloadedVolumeLast
	receiver.Completed = cp.Completed
	receiver.KillList = cp.KillList
	receiver.Errors = cp.Errors

	return nil
}
