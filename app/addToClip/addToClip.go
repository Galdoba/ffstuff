package main

import (
	"os"
	"path/filepath"

	"golang.design/x/clipboard"
)

func main() {
	uhd, err := os.UserHomeDir()
	sep := string(filepath.Separator)
	if err != nil {
		println(err.Error())
		return
	}
	clipData := uhd + sep + "clip.txt"
	err = os.Remove(clipData)
	if err != nil {
		println(err.Error())
		return
	}
	f, err := os.OpenFile(clipData, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		println(err.Error())
		return
	}
	defer f.Close()
	data := clipboard.Read(clipboard.FmtText)
	f.WriteString(string(data))
	for i, arg := range os.Args {
		if i == 0 {
			continue
		}
		f.WriteString("\n" + string(arg))
	}
	clipBytes, err := os.ReadFile(clipData)
	if err != nil {
		println(err.Error())
		return
	}
	clipboard.Write(clipboard.FmtText, clipBytes)
}
