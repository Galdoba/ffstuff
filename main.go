package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/Galdoba/ffstuff/pkg/namedata"
	"github.com/Galdoba/ffstuff/pkg/scanner"
)

func main() {

	// fmt.Println("")
	// fmt.Print(string(rune(9617))) // ░
	// fmt.Print(string(rune(9618))) //   ▒
	// fmt.Print(string(rune(9608))) //█
	// fmt.Print(string(rune(9612))) //▌
	// fmt.Print(string(rune(9619))) // ▓
	// os.Exit(6)
	fmt.Print(string(rune(9617)))
	fmt.Println("Start")
	//list, err := scanner.ListContent("\\\\nas\\ROOT\\EDIT\\_sony\\Breaking_Bad_s01\\")
	list, err := scanner.Scan("\\\\192.168.31.4\\root\\EDIT\\", ".ready")
	fmt.Println("List:")
	for i := range list {
		fmt.Println(i, list[i])
		name := namedata.RetrieveShortName(list[i])
		name = strings.TrimSuffix(name, ".ready")
		dir := namedata.RetrieveDirectory(list[i])
		sList, err2 := scanner.ListContent(dir)
		if err2 != nil {
			fmt.Println(err.Error())
		}
		for f := range sList {
			if strings.Contains(sList[f], name) {
				fmt.Print(sList[f])
				fmt.Print("\n")
			}
		}
	}
	fmt.Println("Err:", err)
	for now := 0; now < 100; now++ {
		downloadbar(now)
		time.Sleep(time.Millisecond * 50)
	}
	fmt.Println("")
	fmt.Println("end")
}

func downloadbar(now int) {
	if now > 100 || now < 0 {
		return
	}
	s := ""
	for i := 0; i < 100; i++ {
		if i <= now {
			s += string(rune(9608))
			continue
		}
		s += string(rune(9617))
	}
	s += "\r"
	fmt.Print(s)
}
