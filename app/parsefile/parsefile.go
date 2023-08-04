package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/Galdoba/ffstuff/pkg/mdm/inputinfo"
	"github.com/Galdoba/ffstuff/pkg/namedata"
	"github.com/Galdoba/ffstuff/pkg/spreadsheet"
	"github.com/Galdoba/ffstuff/pkg/spreadsheet/tablemanager"
	"github.com/Galdoba/ffstuff/pkg/translit"
)

// func main() {
// 	inputinfo.CleanScanData()
// }

func main() {

	for _, arg := range os.Args {
		pi, err := inputinfo.ParseFile(arg)
		if err != nil {
			//fmt.Println(err.Error())
			continue
		}

		bud := ""
		bud += pi.Buffer()
		bud = strings.ReplaceAll(bud, "At least one output file must be specified\n", "")
		bLines := strings.Split(bud, "\n")
		shortBuf := ""
		for _, line := range bLines {
			switch {
			case strings.Contains(line, "Input #"):
				shortBuf += line + "\n"
			case strings.Contains(line, "Duration:"):
				shortBuf += line + "\n"
			case strings.Contains(line, "Stream #"):
				shortBuf += line + "\n"
			}
		}
		fmt.Println(shortBuf)

		if strings.Contains(pi.FileName(), "--") {
			editname := namedata.EditForm(pi.FileName()).EditName()
			sheet, err := spreadsheet.New()
			if err == nil {
				taskList := tablemanager.TaskListFrom(sheet)
				readyList := taskList.ReadyForDemux()
				for _, check := range readyList {
					if strings.Contains(translit.Transliterate(check.Name()), editname) {
						fmt.Println(fmt.Sprintf(`FILE="%v"`, pi.FileName()))
						fmt.Println(fmt.Sprintf(`OUTBASE="%v"`, editname))
						fmt.Println(fmt.Sprintf(`EDIT_PATH="/mnt/pemaltynov/ROOT/EDIT/%v"`, tablemanager.ProposeTargetDirectory(taskList, check)))
						fmt.Println(fmt.Sprintf(`ARCHIVE_PATH="/mnt/pemaltynov/ROOT/IN/%v"`, tablemanager.ProposeArchiveDirectory(check)))
					}
				}
			}
		}

		for _, w := range pi.Warnings() {
			fmt.Println("warning:", w)
		}
	}
}
