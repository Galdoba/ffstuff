package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/Galdoba/devtools/cli/command"
	"github.com/Galdoba/devtools/keyval"
)

const (
	IN_PATH  = `\\192.168.31.4\buffer\IN\`
	TRL_TAG  = `--TRL--`
	FILM_TAG = `--FILM--`
	SER_TAG  = `--SER--`
)

func main() {
	entries, err := os.ReadDir(IN_PATH)
	if err != nil {
		log.Fatal(err)
	}

	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		for _, tag := range []string{TRL_TAG, FILM_TAG, SER_TAG} {
			if !strings.Contains(e.Name(), tag) {
				//				fmt.Println("no tag", tag, e.Name())
				continue
			}
			switch tag {
			default:
				//fmt.Println(e.Name(), tag)
			case SER_TAG:
				//se, ep := findSERdata(e.Name())
				//fmt.Println(e.Name(), se, ep, tag)
			}
			key := entryKey(e.Name(), tag)
			fmt.Println(key)
			if !entryExists(key) {
				createEntry(key)
			}
			addTaskInputFile(key, e.Name())
		}
	}

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

func entryExists(name string) bool {
	kv, err := keyval.Load("fftasks_status")
	if err != nil {
		return false
	}
	single, err := kv.GetSingle(name)
	if err != nil {
		return false
	}
	switch single {
	case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F":
		return true
	default:
		return false
	}
}

func createEntry(name string) {
	kv, _ := keyval.NewKVlist(name)
	kv.Save()

}

func readEntry(name string) string {
	comm, err := command.New(command.CommandLineArguments(
		"kval", fmt.Sprintf("read -from fftasks_status -k %v", name),
	),
		command.Set(command.BUFFER_ON),
		command.Set(command.TERMINAL_ON),
	)
	err = comm.Run()
	if err != nil {
		fmt.Println(err.Error())
	}
	return comm.StdOut()
}

func addTaskInputFile(task string, file string) error {
	kv, _ := keyval.NewKVlist(fmt.Sprintf("%v", task))
	kv.Save()
	fmt.Println("==", keyval.MakePathJS(task))

	taskKVL, err := keyval.Load(fmt.Sprintf("%v", task))
	if err != nil {
		panic("+++" + err.Error())
	}
	taskKVL.Add("inputfiles", file)
	return nil
}
