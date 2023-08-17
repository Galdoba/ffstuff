package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Galdoba/ffstuff/pkg/namedata"
)

func visit(path string, f os.FileInfo, err error) error {
	if strings.Contains(path, "--") {
		ed := namedata.EditForm(path)
		dir := namedata.RetrieveDirectory(path)
		//fmt.Println(dir)
		short := namedata.RetrieveShortName(path)
		edName := ed.EditName()
		cntnt := ed.ContentType()
		newName := dir + strings.TrimPrefix(short, edName+"--"+cntnt+"--")
		nm := strings.Split(newName, "--")
		if len(nm) == 2 {
			newName = dir + nm[1]
		}

		err := os.Rename(path, newName)
		if err == nil {
			fmt.Println("path corrected:", newName)
			visited++
		} else {
			fmt.Printf("%v: %v", path, err.Error())
		}
	}

	fmt.Printf("%v                         \r", visited)
	return nil
}

var visited int

func main() {
	path := `\\192.168.31.4\root\IN\`
	filepath.Walk(path, visit)
}
