package copyprocess

import (
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/Galdoba/ffstuff/app/grabber/internal/origin"
	"github.com/Galdoba/ffstuff/pkg/logman"
	"github.com/fatih/color"
)

const (
	Status_Pending = 0
	sc             = "\u001B7"
	rc             = "\u001B8"
)

type copyActionState struct {
	sourcePaths          []origin.Origin
	destination          string
	markerExt            string
	creationErrors       []error
	deleteAll            bool
	deleteMarker         bool
	sourceTargetMap      map[origin.Origin]string
	namesLen             int
	sourcesVolume        int64
	downloadedVolume     int64
	downloadedVolumeLast int64
	completed            []origin.Origin
	errors               []error
}

func NewCopyAction(sourceTargetMap map[origin.Origin]string, opts ...Option) *copyActionState {
	cas := copyActionState{}
	// cas.source = src

	settings := defaultCopyOptions()
	for _, modify := range opts {
		modify(&settings)
	}
	cas.destination = settings.destination
	cas.markerExt = settings.markerExt
	cas.sourceTargetMap = sourceTargetMap
	cas.downloadedVolumeLast = -1
	for _, src := range settings.sourcePaths {
		// fmt.Println("---")
		// originFile := origin.New(src)

		cas.sourcePaths = append(cas.sourcePaths, src)
	}
	return &cas
}

func (cas *copyActionState) Start() error {
	startTime := time.Now()
	logman.Info("start grabbing")
	copyErrors := newErrorCollector()
	for _, src := range cas.sourcePaths {
		if cas.namesLen < len(src.Name()) {
			tgt := cas.sourceTargetMap[src]
			tgt = filepath.Base(tgt)
			cas.namesLen = len(tgt)
		}
		si, err := os.Stat(src.Path())
		if err != nil {
			return fmt.Errorf("failed to get %v size", src.Name())
		}
		cas.sourcesVolume += si.Size()
	}
	fmt.Println("target directory:", cas.destination)
	fmt.Print(sc)
	fmt.Println(cas.copyProcessData())

	for _, src := range cas.sourcePaths {
		source := src.Path()
		srcInfo, errS := os.Stat(source)
		if errS != nil {
			cas.errors = append(cas.errors, logman.Errorf("failed to get info on %v: %v", source, errS))
			continue
		}
		if !srcInfo.Mode().IsRegular() { // cannot copy non-regular files (e.g., directories, symlinks, devices, etc.)
			cas.errors = append(cas.errors, fmt.Errorf("cannot copy non-regular files (directories, symlinks, devices, etc.): "+source+" ("+srcInfo.Mode().String()+")"))
			continue
		}
		sourceSize := srcInfo.Size()
		target := cas.sourceTargetMap[src]

		go copyContent(src.Path(), target, copyErrors)
		doneCopying := false
		goodDone := false
		time.Sleep(time.Second)
		for !doneCopying {
			copyFile, err := os.Stat(target)
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
			cas.completed = append(cas.completed, src)
		} else {
			cas.errors = append(cas.errors, logman.Errorf("failed to grab %v", source))
		}

	}
	fmt.Print(rc + sc)
	fmt.Println(cas.copyProcessData() + "\n")
	dur := time.Since(startTime)
	tookTime := dur.String()

	allErrors := append(copyErrors.collected, cas.errors...)
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
	logman.Info("done grabbing (process took %v)", tookTime)
	if len(cas.completed) < len(cas.sourcePaths) {
		return logman.Warn("grabbed %v of %v files", len(cas.completed), len(cas.sourcePaths))
	}
	if len(cas.errors) > 0 {
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
	cas.downloadedVolumeLast = cas.downloadedVolume
	cas.downloadedVolume = 0
	out := "process report:\n"
	for _, src := range cas.sourcePaths {
		tgt := cas.sourceTargetMap[src]
		progress, trueSize := currentProgress(src.Path(), tgt)
		progressString := formatProgress(progress)
		tgtName := filepath.Base(tgt)
		out += fmt.Sprintf("  %v    %v           \n", wideName(tgtName, cas.namesLen), progressString)
		cas.downloadedVolume += trueSize
	}
	if len(cas.sourcePaths) == 0 {
		out += "  nothing to grab"
		return out
	}

	out += "summary:\n"
	out += fmt.Sprintf("  %v        ", summaryString(cas.sourcesVolume, cas.downloadedVolume, cas.downloadedVolumeLast))
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
		progressString = color.GreenString("ok")
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
	if len(cas.errors) == 0 {
		return nil
	}
	text := "processing error(s) detected:\n"
	for i, err := range cas.errors {
		text += fmt.Sprintf("%v: %v\n", i+1, err)
	}
	return errors.New(text)
}

/*
create
*/
