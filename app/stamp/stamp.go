package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/Galdoba/devtools/cli/command"
	"github.com/Galdoba/utils"
)

func main() {
	//ffmpeg -i c:\Users\Public\IMG_1988.jpg -vf "drawtext=text='AAA\:AAA/AAA AAA.AAAAAA':fontcolor=red:fontsize=75:x=200:y=100, drawtext=text='BBBBBBBBB':fontcolor=red:fontsize=75:x=400:y=200" -y c:\Users\Public\IMG_1988_2.jpg
	args := []string{}
	folder := `d:\tests\`
	file, err := os.Open(folder + `Foto.csv`)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// text1, e := os.Create("/home/galdoba/workbench/text1.txt")
	// if e != nil {
	// 	log.Fatal(e)
	// }
	// text2, e := os.Create("/home/galdoba/workbench/text2.txt")
	// if e != nil {
	// 	log.Fatal(e)
	// }

	scanner := bufio.NewScanner(file)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	lines := utils.LinesFromTXT(folder + `photos.csv`)
	line := -1
	for scanner.Scan() {
		clear, err := command.New(
			command.CommandLineArguments("clear"),
			command.Set(command.TERMINAL_ON),
		)
		if err != nil {
			panic(err.Error())
		}

		clear.Run()
		line++
		if line == 0 {
			continue
		}
		text1, e := os.Create(`d:\tests\text1.txt`)
		if e != nil {
			log.Fatal(e)
		}
		text2, e := os.Create(`d:\tests\text2.txt`)
		if e != nil {
			log.Fatal(e)
		}

		//fmt.Println("\n========================================\n")
		found := false
		scanned := scanner.Text()
		//		scanned = strings.ReplaceAll(scanned, ":", " ")
		data := strings.Split(scanned, ",")

		input := folder + data[0]
		//fmt.Println(data)
		files, err := IOReadDir(folder)
		if err != nil {
			panic(err)
		}
		outFolder := ""
		for _, inList := range files {
			//	fmt.Println(inList)
			if strings.HasSuffix(inList, data[0]) {
				input = inList
				found = true
				outFolder = strings.TrimSuffix(input, data[0])
			}
			if found {
				break
			}
		}
		//base := strings.Split(data[0], ".")[0]
		output := outFolder + "Stamped_" + data[0]
		//fmt.Println("Write line 1:", data[1])
		text1.Write([]byte(data[1]))
		//fmt.Println("Write line 2:", data[2])
		text2.Write([]byte(data[2]))
		//text1.WriteString(data[1])
		//text2.WriteString(data[2])
		fmt.Printf("Line %v/%v\n", line, len(lines))
		args = append(args, fmt.Sprintf(`-i %v -vf drawtext=textfile=/home/galdoba/workbench/text1.txt:fontcolor=orange:fontsize=150:x=((w/5)*3):y=((h/100)*85):,drawtext=textfile=/home/galdoba/workbench/text2.txt:fontcolor=orange:fontsize=150:x=((w/5)*3):y=((h/100)*85)+160: -y %v`, input, output))
		//"drawtext=textfile=text.txt:x=640:y=360:fontsize=24:fontcolor=white"     fontsize=150:x=1002:y=100     5152x3864                             fontsize=h/30: x=(w-text_w)/2: y=(h-text_h*2)
		command, err := command.New(
			command.CommandLineArguments("ffmpeg", args...),
			command.Set(command.TERMINAL_ON),
		)
		if err != nil {
			panic(err.Error())
		}

		command.Run()
		fmt.Println(input)
		fmt.Println(output)
		//time.Sleep(time.Microsecond * 1000000)
		args = []string{}
		fmt.Println("Delete line 1:", data[1])
		text1.Truncate(0)
		fmt.Println("Delete line 2:", data[2])
		text2.Truncate(0)

	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	// panic("stop")
	// command, err := command.New(
	// 	command.CommandLineArguments("ffmpeg"),
	// 	command.Set(command.TERMINAL_ON),
	// )
	// if err != nil {
	// 	panic(err.Error())
	// }
	// command.Run()

}

func IOReadDir(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}
