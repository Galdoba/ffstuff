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
	os.Remove(clipData)
	f, err := os.OpenFile(clipData, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		println(err.Error())
		return
	}
	defer f.Close()
	data := clipboard.Read(clipboard.FmtText)
	_, err = f.WriteString(string(data))
	if err != nil {
		println(err.Error())
	}
	//fmt.Println("bytes written to file", wr)
	for i, arg := range os.Args {
		if i == 0 {
			continue
		}
		_, err := f.WriteString("\n" + string(arg))
		if err != nil {
			println(err.Error())
		}
		//fmt.Println("bytes written to file", wr)
	}
	clipBytes, err := os.ReadFile(clipData)
	if err != nil {
		println(err.Error())
		return
	}
	clipboard.Write(clipboard.FmtText, clipBytes)

}
