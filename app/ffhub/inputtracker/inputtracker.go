package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/Galdoba/devtools/cli/command"
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
				//fmt.Println("No ENTRY", key)
				createEntry(key)
			} else {
				fmt.Println("HAVE ENTRY", key)
			}
			//fmt.Println("readEntry", key, readEntry(key))

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
	comm, _ := command.New(command.CommandLineArguments(
		"kval", "read ", "-from ", "fftasks_status ", "-k ", name,
	),
		command.Set(command.BUFFER_ON),
	)
	comm.Run()

	return false
}

func createEntry(name string) {
	comm, err := command.New(command.CommandLineArguments(
		"kval", fmt.Sprintf("write -to fftasks_status -k %v 0", name),
	),
		command.Set(command.BUFFER_ON),
	)
	err = comm.Run()
	if err != nil {
		fmt.Println(err.Error())
	}

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
