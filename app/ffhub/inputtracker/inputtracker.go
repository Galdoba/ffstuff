package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/Galdoba/devtools/cli/command"
)

/*
ПРАВИЛА:
1. Всегда печатай что происходит
2. Используем только утилиту kval
*/

const (
	//IN_PATH     = `\\192.168.31.4\buffer\IN\`
	TRL_TAG     = `--TRL--`
	FILM_TAG    = `--FILM--`
	SER_TAG     = `--SER--`
	status_file = "fftasks_status"
	input_key   = "inputfiles"
)

/*
checklist:
подтверди статусный файл
	если нет - создай
собери список файлов в папке Входящее
	если ошибка - заверши работу
	ПО СПИСКУ ФАЙЛОВ
	если не линкованый - игнорируем
	подтверди отсуствие статус проекта
		если нет - создать файл проекта
	подтверди статус проекта равный 0
		если нет - игнорируем
	добавить имя файла в исходные файлы проекта
заверши работу
*/

func main() {
	// подтверди статусный файл
	// 	если нет - создай
	if !confirm(status_file, "") {
		fmt.Println("creating status file")
		command.RunSilent("kval", fmt.Sprintf("new %v", status_file))
	}
	// собери список файлов в папке Входящее
	// 	если ошибка - заверши работу
	IN_PATH, err := command.RunSilent("kval", fmt.Sprintf("read -from work_dirs -k IN"))
	IN_PATH = strings.TrimSuffix(IN_PATH, "\n")
	entries, err := os.ReadDir(IN_PATH)
	if err != nil {
		fmt.Println("не могу создать список файлов в директории %v:\n%v", IN_PATH, err.Error())
		os.Exit(1)
	}
	// 	ПО СПИСКУ ФАЙЛОВ
	for _, e := range entries {
		//fmt.Print(e.Name())
		if e.IsDir() {
			//fmt.Print(" игнорируем (директория)\n")
			continue
		}
		tag := ""
		fileName := e.Name()
		switch {
		case strings.Contains(fileName, TRL_TAG):
			tag = TRL_TAG
		case strings.Contains(fileName, FILM_TAG):
			tag = FILM_TAG
		case strings.Contains(fileName, SER_TAG):
			tag = SER_TAG
		}
		if tag == "" {
			//fmt.Print(" игнорируем (не привязан)\n")
			continue
		}
		key := entryKey(fileName, tag)
		if !confirm("fftasks_status", key) {
			command.RunSilent("kval", fmt.Sprintf("write -to %v -k %v 0", status_file, key))

		}
		if !confirm(key, "") {
			command.RunSilent("kval", fmt.Sprintf("new ffprojects/%v", key))
		}
		addTaskInputFile(key, e.Name())

	}
	// 	если не линкованый - игнорируем
	// 	подтверди отсуствие статус проекта
	// 		если нет - создать файл проекта
	// 	подтверди статус проекта равный 0
	// 		если нет - игнорируем
	// 	добавить имя файла в исходные файлы проекта
	// заверши работу

}

func entryKey(name, tag string) string {
	data := strings.Split(name, tag)
	return data[0] + tag
}

func findSERdata(s string) (int, int) {
	re := regexp.MustCompile(`s[0-9]+e[0-9]+`)
	match := re.FindString(s)
	if match != "" {
		data := strings.TrimPrefix(match, "s")
		parts := strings.Split(data, "e")
		se, _ := strconv.Atoi(parts[0])
		ep, _ := strconv.Atoi(parts[1])
		return se, ep
	}

	return -1, -1
}

func confirm(list, key string) bool {
	keyStr := ""
	if key != "" {
		keyStr = fmt.Sprintf("-k %v", key)
	}
	comm, err := command.New(
		command.CommandLineArguments("kval", fmt.Sprintf("confirm -page %v %v", list, keyStr)),
		command.Set(command.BUFFER_ON),
		command.Set(command.TERMINAL_OFF),
	)
	err = comm.Run()
	if err != nil {
		fmt.Println(err.Error())
		panic(2)
		return false
	}
	out := comm.StdOut()

	switch out {
	default:
	case "1\n":
		return true
	case "0\n":
		return false
	}
	return false
}

func addTaskInputFile(task string, file string) string {
	_, err := command.RunSilent("kval", fmt.Sprintf("append -to ffprojects/%v -k inputfiles -u %v", task, file))

	if err != nil {
		return err.Error()
	}
	return fmt.Sprintf("input file '%v' added to %v", file, task)
}
