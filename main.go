package main

import (
	"fmt"
	"strings"

	"github.com/Galdoba/ffstuff/clipmaker"
	"github.com/Galdoba/ffstuff/ediread"
	"github.com/Galdoba/ffstuff/fldr"
)

var WorkDate string

func main() {
	fldr.Init()

	edlFile := fldr.SelectEDL()
	edi, err := ediread.NewEdlData(edlFile)
	if err != nil {
		fmt.Println(err)
	}

	cliTasks := []string{}
	clipMap := clipmaker.NewClipMap()
	for _, clipData := range edi.Entry() {
		fmt.Println(clipData)
		cl, err := clipmaker.NewClip(clipData)
		if err != nil {
			fmt.Println(err)
		}
		clipMap[cl.Index()] = cl

		program, arguments := clipmaker.CreateTask(cl)
		cutClip := program + " " + strings.Join(arguments, " ")
		cliTasks = append(cliTasks, cutClip)

	}
	for _, val := range cliTasks {
		fmt.Print("RUN:", val, "\n")

	}
	clipmaker.ConcatClips(clipMap)
}
