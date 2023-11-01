package main

import (
	"fmt"

	"github.com/Galdoba/devtools/cli/command"
	"github.com/Galdoba/devtools/keyval"
)

const (
	key_inputfiles = "inputfiles"
	kval_workdirs  = "work_dirs"
)

type managerTasks struct {
	data   keyval.Kval
	status string
}

func main() {
	//проходим по задачам
	//	исходя из текущего статуса предлогаем команду

	status, err := keyval.Load("fftasks_status")
	if err != nil {
		fmt.Println(err.Error())
	}
	projectData := make(map[string]keyval.Kval)
	for _, projectKey := range status.Keys() {
		project, err := keyval.Load("ffprojects/" + projectKey)
		if err != nil {
			fmt.Println("err:", err.Error())
		}
		projectData[projectKey] = project

	}
	for _, currentStatus := range []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F"} {
		switch currentStatus {
		default:
			panic(currentStatus + " not inplemented")
		case "0":
			for projectKey, project := range projectData {
				fmt.Println("")
				currentS, err := status.GetSingle(projectKey)
				if err != nil {
					println("manager fail:", err.Error())
				}
				if currentS != currentStatus {
					continue
				}
				inputFiles, err := project.GetAll(key_inputfiles)
				if err != nil {
					panic("checking input file: " + err.Error())
				}
				if len(inputFiles) > 0 {
					fmt.Println("suggest: set status to '1'", projectKey)
					out, err := command.RunSilent("changestatus", fmt.Sprintf("%v 1", projectKey))
					fmt.Println(out)
					if err != nil {
						panic("change status to 1: " + err.Error())
					}
				}
				fmt.Println(project, "rejected step 1")

			}
		}
	}

}
