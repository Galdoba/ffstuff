package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Galdoba/devtools/cli/command"
)

func main() {
	//ffmpeg -i c:\Users\Public\IMG_1988.jpg -vf "drawtext=text='AAA\:AAA/AAA AAA.AAAAAA':fontcolor=red:fontsize=75:x=200:y=100, drawtext=text='BBBBBBBBB':fontcolor=red:fontsize=75:x=400:y=200" -y c:\Users\Public\IMG_1988_2.jpg
	args := []string{}
	folder := `c:\Users\Public\`
	file, err := os.Open(folder + `Table.csv`)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	line := -1
	for scanner.Scan() {
		line++
		if line == 0 {
			continue
		}
		fmt.Println(scanner.Text())
		scanned := scanner.Text()
		//		scanned = strings.ReplaceAll(scanned, ":", " ")
		data := strings.Split(scanned, ",")
		input := folder + data[0]
		fmt.Println(data)
		args = append(args, fmt.Sprintf(`-i %v -vf drawtext=textfile=text1.txt:fontcolor=red:fontsize=75:x=1002:y=100:,drawtext=textfile=text2.txt:fontcolor=red:fontsize=75:x=1002:y=200: -y aaa.jpg`, input))
		//"drawtext=textfile=text.txt:x=640:y=360:fontsize=24:fontcolor=white"
		fmt.Println(args)
		command, err := command.New(
			command.CommandLineArguments("ffmpeg", args...),
			command.Set(command.TERMINAL_ON),
		)
		if err != nil {
			panic(err.Error())
		}
		command.Run()
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	panic("stop")
	command, err := command.New(
		command.CommandLineArguments("ffmpeg"),
		command.Set(command.TERMINAL_ON),
	)
	if err != nil {
		panic(err.Error())
	}
	command.Run()
}
