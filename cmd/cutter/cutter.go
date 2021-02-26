package main

import (
	"fmt"
	"strings"

	"github.com/Galdoba/ffstuff/clipmaker"
	"github.com/Galdoba/ffstuff/ediread"
	"github.com/Galdoba/ffstuff/fldr"
	"github.com/Galdoba/ffstuff/pkg/cli"
)

func main() {
	fldr.Init()

	edlFile := fldr.SelectEDL()
	edi, err := ediread.NewEdlData(edlFile)
	if err != nil {
		fmt.Println(err)
	}

	cliTasks := []cli.Task{}
	clipMap := clipmaker.NewClipMap()
	for _, clipData := range edi.Entry() {
		fmt.Println(clipData)
		cl, err := clipmaker.NewClip(clipData)
		if err != nil {
			fmt.Println(err)
		}
		clipMap[cl.Index()] = cl
		cliTasks = append(cliTasks, cli.NewTask(clipmaker.CutClip(cl)))
	}
	cliTasks = sortTasks(cliTasks)
	for _, task := range cliTasks {
		fmt.Print("RUN:", task, "\n")
		task.Run()

	}
	clipmaker.ConcatClips(clipMap)
}

//ставит резку аудио перед резкой видео.
func sortTasks(unsorted []cli.Task) []cli.Task {
	sorted := []cli.Task{}
	for _, task := range unsorted {
		if strings.Contains(task.String(), "_ACLIP_") {
			sorted = append(sorted, task)
		}
	}
	for _, task := range unsorted {
		if strings.Contains(task.String(), "_VCLIP_") {
			sorted = append(sorted, task)
		}
	}
	if len(sorted) != len(unsorted) {
		fmt.Println("Can not sort")
		return unsorted
	}
	return sorted
}
