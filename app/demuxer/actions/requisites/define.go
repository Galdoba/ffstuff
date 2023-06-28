package main

import (
	"fmt"

	"github.com/Galdoba/ffstuff/app/demuxer/handle"
	"github.com/Galdoba/ffstuff/pkg/mdm/inputinfo"
	"github.com/Galdoba/ffstuff/pkg/spreadsheet"
	"github.com/Galdoba/ffstuff/pkg/spreadsheet/tablemanager"
)

const (
	TRAILER = "TRAILER"
	FILM    = "FILM"
	SEASON  = "SEASON"
)

type demuxTask struct {
	sources    []string
	taskType   string
	agent      string
	targetOS   string
	editPath   string
	archvePath string
	outputs    []string
	sourceData map[string]inputinfo.ParseInfo
}

func inputList() []string {
	return []string{
		`\\192.168.31.4\buffer\IN\_DONE\Razboi_prores.mov`,
	}
}

func main() {
	dt, err := ConstructTask(inputList())
	if dt.agent == "" {
		fmt.Println("task.agent unset")
	}
	if dt.archvePath == "" {
		fmt.Println("task.archvePath unset")
	}
	if dt.editPath == "" {
		fmt.Println("task.editPath unset")
	}
	if dt.targetOS == "" {
		fmt.Println("task.targetOS unset")
	}
	if dt.taskType == "" {
		fmt.Println("task.taskType unset")
	}
	if len(dt.outputs) == 0 {
		fmt.Println("task.outputs unset")
	}
	if len(dt.sources) == 0 {
		fmt.Println("task.sources unset")
	}
	if len(dt.sourceData) == 0 {
		fmt.Println("task.sourceData unset")
	}
	fmt.Println(dt)
	fmt.Println(err)
}

func ConstructTask(list []string) (demuxTask, error) {
	dt := demuxTask{}
	dt.WithTaskType()
	taskList := filterTasks(dt.taskType)
	if len(taskList) == 0 {
		return dt, fmt.Errorf("Подходящих задач не найдено")
	}
	taskData := SelectTaskData(dt.taskType, taskList)
	fmt.Println("==========")
	fmt.Println(taskData)
	fmt.Println("||||||||||")
	return dt, nil
}

func (dt *demuxTask) WithTaskType() *demuxTask {
	taskType := handle.SelectionSingle("Что хотим получить?", []string{TRAILER, FILM, SEASON}...)
	dt.taskType = taskType
	return dt
}

func SelectTaskData(taskType string, taskList []tablemanager.TaskData) tablemanager.TaskData {
	options := []string{}
	answerMap := make(map[string]tablemanager.TaskData)
	for i, task := range taskList {
		switch taskType {
		case TRAILER:
			options = append(options, task.StringAsTrailer())
		case FILM:
			options = append(options, task.StringAsFilm())
		case SEASON:
			options = append(options, task.StringAsSeason())
		}
		answerMap[options[i]] = task
	}
	fmt.Println("")
	selected := handle.SelectionSingle("Что готовим?", options...)
	return answerMap[selected]
}

func filterTasks(taskType string) []tablemanager.TaskData {
	sp, _ := spreadsheet.New()
	if err := sp.Update(); err != nil {
		return nil
	}
	tlist := tablemanager.TaskListFrom(sp)
	taskList := []tablemanager.TaskData{}
	switch taskType {
	case TRAILER:
		taskList = tlist.ChooseTrailer()
	case FILM:
		taskList = tlist.ChooseFilm()
	case SEASON:
		taskList = tlist.ChooseSeason()
	}
	return taskList
}
