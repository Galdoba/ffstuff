package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/Galdoba/devtools/cli/command"
)

func main() {
	fmt.Println("START locator")
	//
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
	list := strings.Fields(buf)
	pth, _ := command.RunSilent("kval", fmt.Sprintf("read -from work_dirs -k IN"))
	pth = strings.TrimSuffix(pth, "\n")
	fmt.Printf("'%v'\n", pth)
	for _, file := range list {
		fmt.Println(file)
	}

}
