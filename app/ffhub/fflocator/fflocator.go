package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/Galdoba/devtools/cli/command"
)

func main() {
	args := os.Args
	if len(args) > 2 {
		fmt.Println("fflocator failed: to many arguments")
		return
	}
	if len(args) < 2 {
		fmt.Println("fflocator failed: not enough arguments")
		return
	}

	buf, err := command.RunSilent("kval", fmt.Sprintf("read -from ffprojects/%v -k inputfiles", args[1]))
	if err != nil {
		fmt.Printf("fflocator failed: %v\n", err.Error())
		return
	}
	if strings.Contains(buf, "kval returned error") {
		fmt.Printf("fflocator failed: %v", strings.TrimPrefix(buf, "kval returned error: "))
		return
	}
	list := strings.Fields(buf)

	buf, err = command.RunSilent("kval", fmt.Sprintf("keys work_dirs"))
	if err != nil {
		fmt.Printf("fflocator failed: %v\n", err.Error())
		return
	}
	key_fold := strings.Fields(buf)
	key_fold = key_fold[2:]

	paths := []string{}
	for _, key := range key_fold {
		buf, err = command.RunSilent("kval", fmt.Sprintf("read -from work_dirs -k %v", key))
		if err != nil {
			fmt.Printf("fflocator failed: %v\n", err.Error())
			return
		}
		pth := strings.Fields(buf)
		paths = append(paths, pth...)
	}
	filesLeft := len(list)

	loc := ""
	for _, dir := range paths {
		if filesLeft == 0 {
			break
		}
		for _, file := range list {
			_, err := os.Stat(dir + file)
			if err != nil {
				continue
			}
			if loc == "" {
				loc = dir
			}
			filesLeft--
		}

	}
	if loc == "" {
		fmt.Println("fflocator failed: project files were not located")
		return
	}
	fmt.Println(loc)
}
