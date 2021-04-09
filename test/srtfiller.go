package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/Galdoba/utils"
)

func main() {
	lines := utils.LinesFromTXT("c:\\Users\\pemaltynov\\go\\src\\github.com\\Galdoba\\ffstuff\\test\\Uoker_s01e08_SER_10517.srt")
	file := "newFile.srt"
	newFile, err := os.Create(file)
	if err != nil {
		panic(err)
	}
	defer newFile.Close()

	utils.AddLineToFile(file, "1")
	utils.AddLineToFile(file, "00:00:00,000 --> 00:00:00,040")
	utils.AddLineToFile(file, " ")
	utils.AddLineToFile(file, "")

	for i, line := range lines {
		fmt.Print(i, " of ", len(lines), "\r")
		if line == "" {
			utils.AddLineToFile(file, line+"\n")
			continue
		}
		num, err := strconv.Atoi(line)
		if err != nil {
			utils.AddLineToFile(file, line)
			continue
		}
		num++
		utils.AddLineToFile(file, strconv.Itoa(num))
	}

}
