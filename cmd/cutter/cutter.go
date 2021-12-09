package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/Galdoba/ffstuff/clipmaker"
	"github.com/Galdoba/ffstuff/ediread"
	"github.com/Galdoba/ffstuff/fldr"
	"github.com/Galdoba/ffstuff/pkg/cli"
	"github.com/Galdoba/ffstuff/pkg/glog"
)

func main() {
	fldr.Init()
	//logger := glog.New(fldr.MuxPath()+"logfile.txt", glog.LogLevelWARN)
	logger := glog.New(glog.LogPathDEFAULT, glog.LogLevelINFO)
	logger.INFO("Start cutting")
	edlFile := fldr.SelectEDL()
	edi, err := ediread.NewEdlData(edlFile)
	if err != nil {
		//fmt.Println(err)
		logger.ERROR(err.Error())
	}

	cliTasks := []cli.Task{}
	clipMap := clipmaker.NewClipMap()
	for _, clipData := range edi.Entry() {
		//fmt.Println(clipData)
		//////////////////////////////////
		f, err := os.OpenFile(fldr.MuxPath()+"cutterOutputNames.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		// if _, err = f.WriteString(clipData + "\n"); err != nil {
		// 	fmt.Println(err)
		// }
		//////////////////////////////////
		cl, err := clipmaker.NewClip(clipData)
		if err != nil {
			fmt.Println(err)
		}
		clipMap[cl.Index()] = cl
		newTask := cli.NewTask(clipmaker.CutClip(cl))
		cliTasks = append(cliTasks, newTask)
		if _, err = f.WriteString(newTask.String() + "\n"); err != nil {
			fmt.Println(err)
		}
	}
	cliTasks = sortTasks(cliTasks)

	for _, task := range cliTasks {
		fmt.Print("RUN:", task, "\n")
		logger.INFO("Run: " + task.String())
		taskErr := task.Run()
		if taskErr != nil {
			logger.ERROR(taskErr.Error())
		}
	}
	logger.INFO("Cutting Complete")
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
