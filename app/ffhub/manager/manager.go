package main

import (
	"fmt"

	"github.com/Galdoba/devtools/cli/command"
	"github.com/Galdoba/devtools/keyval"
)

const (
	key_inputfiles = "inputfiles"
)

func main() {
	//проходим по задачам
	//	исходя из текущего статуса предлогаем команду

	kv, err := keyval.Load("fftasks_status")
	if err != nil {
		fmt.Println(err.Error())
	}

	for i, projectName := range kv.Keys() {
		projectStatus, e := kv.GetSingle(projectName)
		if e != nil {
			panic(e.Error())
		}

		projectData, err := keyval.Load(fmt.Sprintf("ffprojects/%v", projectName))
		if err != nil {
			continue
		}
		fmt.Println(i, projectName, projectStatus, e)
		inputs, err := projectData.GetAll(key_inputfiles)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		switch projectStatus {
		case "0": //собираем в проект список исходных файлов
			if len(inputs) > 0 {
				command.RunSilent("changestatus", fmt.Sprintf("%v 1", projectName))
				//fmt.Println("Suggest: changestatus", fmt.Sprintf("%v 1", projectName))
				return
			}
		case "1": //смотрим чтобы все исходные файлы лежали в IN

		}
		// for _, inp := range inputs {
		// 	//добавляем в проектный файл первичные данные
		// 	//текущая локация "current_dir: "
		// 	//текущая локация "current_dir: "
		// 	fmt.Println(inp)

		// }

	}
}
