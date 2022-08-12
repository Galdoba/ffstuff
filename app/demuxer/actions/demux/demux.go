package actiondemux

import (
	"fmt"
	"strings"

	"github.com/Galdoba/ffstuff/pkg/mdm/format"

	"github.com/Galdoba/ffstuff/pkg/spreadsheet/tablemanager"

	"github.com/Galdoba/devtools/cli/command"

	"github.com/Galdoba/ffstuff/app/demuxer/handle"

	"github.com/urfave/cli"
)

/*
ПРИМЕРЫ ПРИМЕНЕНИЯ
demuxer -tofile file.txt -update demux -i film.mp4
	-tofile file.txt - терминал будет писаться в указанный файл

Для демукса требуются:
1. Исходник(и)
2. Информация по заданию (данные из таблицы)
3. Задание (ввод в ручную что это фильм/трейлер/сериал)

ПЛАН:
1. Собираем данные:
	1.1 Подтверждаем исходник(и)
	1.2 Запрашиваем задание
	1.3 Читаем таблицу.
	1.4 ДЕБАГ: Выводим имена и пути предпологаемых результатов.


*/

var inputBuffer []string
var inputPaths []string

func Run(c *cli.Context) error {
	fmt.Println("RUN Precheck")
	if err := Precheck(c); err != nil {
		return err
	}
	fmt.Println("Precheck complete")
	taskType := handle.SelectionSingle("Что в исходнике?", []string{"Фильм", "Трейлер", "Сериал"}...)
	task := tablemanager.TaskData{}
	fmt.Println("в исходнике: ", taskType)
	switch taskType {
	case "Фильм":
		filmtask, err := DefineFilmTask()
		if err != nil {
			return err
		}
		task = filmtask
	}
	////////
	//fmt.Println(task)
	nameBase := task.OutputBaseName()
	//setpts=N/(25*TB)
	fmt.Printf(`DEMUX TO: \\nas\ROOT\EDIT\%v%v`+"\n", tablemanager.ProposeTargetDirectory(handle.TaskListFull(), task), nameBase)
	archive := tablemanager.ProposeArchiveDirectory(task)
	fmt.Printf("ARCHIVE TO: %v\n", archive)
	tFormat := &format.TargetFormat{}
	switch {
	default:
		tFormat, _ = format.SetAs(format.FilmHD)
	case strings.Contains(task.Name(), " SD"):
		tFormat, _ = format.SetAs(format.FilmSD)
	case strings.Contains(task.Name(), " 4K"):
		tFormat, _ = format.SetAs(format.Film4K)
	}
	fmt.Println(" ")
	fmt.Println("Check Input:")
	for _, issue := range Issues(tFormat, task) {
		fmt.Println("WARNING: " + issue)
	}
	return nil
}

func Precheck(c *cli.Context) error {
	args := c.Args()
	if len(args) == 0 {
		return fmt.Errorf("no arguments provided")
	}
	for _, arg := range args {
		fmt.Println(arg)
	}
	selected := handle.SelectionSingle("Перечень исходных файлов корректен?", []string{"ДА", "НЕТ"}...)
	if selected != "ДА" {
		return fmt.Errorf("User abort")
	}
	for _, arg := range args {
		com, err := command.New(
			command.CommandLineArguments("fflite", "-i "+arg),
			command.Set(command.TERMINAL_ON),
			command.Set(command.BUFFER_ON),
		)
		if err != nil {
			return err
		}
		fmt.Println(" ")
		com.Run()
		inputBuffer = strings.Split(com.StdOut()+"\n"+com.StdErr(), "\n")
	}
	fmt.Println(" ")
	return nil
}

func DefineFilmTask() (tablemanager.TaskData, error) {
	task := tablemanager.TaskData{}
	taskList := handle.SelectFromTable("Фильм")
	s := []string{}
	for _, t := range taskList {
		s = append(s, t.String())
	}
	taskStr := handle.SelectionSingle("Данные из таблицы", s...)
	for _, t := range taskList {
		if t.Match(taskStr) {
			task = t
		}
	}
	return task, nil
}

func Issues(tFormat *format.TargetFormat, task tablemanager.TaskData) []string {
	fmt.Println("DEBUG: Checking Issues")
	videoSizeValid := false
	issues := []string{}
	videoFound := 0
	audioFound := 0
	soundMap := make(map[string]int)
	for _, data := range inputBuffer {
		switch {
		case strings.Contains(data, `Video: `):
			videoFound++
			if strings.Contains(data, ` 1920x1080`) && (!strings.Contains(task.Name(), " SD") && !strings.Contains(task.Name(), " 4K")) {
				videoSizeValid = true
			}
			if strings.Contains(data, ` 720x576`) && strings.Contains(task.Name(), " SD") {
				videoSizeValid = true
			}
			if strings.Contains(data, ` 3840x2160`) && strings.Contains(task.Name(), " 4K") {
				videoSizeValid = true
			}
			if !videoSizeValid {
				issues = append(issues, "scaling needed")
			}
		case strings.Contains(data, `Audio: `):
			audioFound++
			if strings.Contains(data, ` stereo`) {
				soundMap["stereo"]++
			}
			if strings.Contains(data, ` 5.1`) {
				soundMap["5.1"]++
			}
			if strings.Contains(data, ` 5.1(side)`) {
				soundMap["5.1"]++
			}
			if strings.Contains(data, ` mono`) {
				soundMap["mono"]++
			}
			if audioFound != mapSum(soundMap) {
				soundMap["warning"]++
			}
		}
	}
	if soundMap["warning"] > 0 {
		issues = append(issues, fmt.Sprintf("%v audio streams require attention", soundMap["warning"]))
	}
	return issues
}

func mapSum(sm map[string]int) int {
	s := 0
	for _, v := range sm {
		s += v
	}
	return s
}
