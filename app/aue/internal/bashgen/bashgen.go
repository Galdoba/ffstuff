package bashgen

import (
	"fmt"
	"os"
	"strings"

	"github.com/Galdoba/ffstuff/app/aue/internal/task"
)

type generator struct {
	base           string
	destination    string
	translationMap map[string]string
}

type BashGenOptions interface {
	ProjectName() string
	BashDestination() string
	BashTranslationMap() map[string]string
}

func New(bgOpts BashGenOptions) *generator {
	gen := generator{}
	gen.base = bgOpts.ProjectName()
	gen.translationMap = bgOpts.BashTranslationMap()
	gen.destination = bgOpts.BashDestination()
	return &gen
}

type Bashgen interface {
	GenerateBash([]task.Task) error
}

func (gen *generator) GenerateBash(allTasks []task.Task) error {
	bash := bashHeader()
	bash += fmt.Sprintf("PRIORITY=8\n")
	if len(gen.translationMap) == 0 {
		panic("no translation map")
	}
	if gen.destination == "" {
		panic("no dest")
	}
	if gen.base == "" {
		panic("no base")
	}
	for _, tsk := range allTasks {
		bash += fmt.Sprintf("%v\n", lineTranslated(tsk.String(), gen.translationMap))
	}

	f, err := os.Create(gen.destination + gen.base + ".sh")
	if err != nil {
		return fmt.Errorf("can't create bash file")
	}
	if _, err := f.WriteString(bash); err != nil {
		return fmt.Errorf("can't write bash file")
	}
	defer f.Close()
	return nil
}

func bashHeader() string {
	bash := ""
	bash += fmt.Sprintf("#!/bin/bash\n")
	bash += fmt.Sprintf("#\n")
	bash += fmt.Sprintf("set -o nounset    # error when referensing undefined variable\n")
	bash += fmt.Sprintf("set -o errexit    # exit when command fails\n")
	bash += fmt.Sprintf("shopt -s extglob\n")
	bash += fmt.Sprintf("shopt -s nullglob\n")
	bash += fmt.Sprintf("#\n")
	return bash
}

func lineTranslated(origin string, translationMap map[string]string) string {
	score := make(map[string]int)
	maxScore := 0
	for k, _ := range translationMap {
		if strings.Contains(origin, k) {
			score[k] = len(k)
			maxScore = max(len(k), maxScore)
		}
	}
	for currentScore := maxScore; currentScore > 0; currentScore-- {
		for k, s := range score {
			if currentScore != s {
				continue
			}
			if !strings.Contains(origin, k) {
				continue
			}
			origin = strings.ReplaceAll(origin, k, translationMap[k])
		}
	}
	origin = strings.ReplaceAll(origin, " [", ` "[`)
	origin = strings.ReplaceAll(origin, "] ", `]" `)
	origin = strings.ReplaceAll(origin, "ffmpeg", "fflite")
	return origin
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
