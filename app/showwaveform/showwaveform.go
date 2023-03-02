package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/Galdoba/devtools/cli/command"
	"github.com/Galdoba/devtools/cli/user"
	"github.com/Galdoba/ffstuff/pkg/mdm/inputinfo"
)

func main() {
	/*
		входные данные - путь
		выходные данные - НЕТ
		алгоритм
		1 - узнать количество аудио
		2 - составить строку для картинок по каждой форме
		3 - составить строку для объединения картинок
	*/
	paths, err := checkInput()
	if err != nil {
		fmt.Println(err.Error())
		user.InputStr()
		os.Exit(1)
	}
	for i, path := range paths {
		fmt.Printf("Start file %v/%v (%v):\n", i+1, len(paths), path)
		// изучаем файл
		pi, err := inputinfo.ParseFile(path)
		if err != nil {
			fmt.Printf("file parsing: %v\n", err.Error())
			continue
		}
		audChans := pi.NumAudio()
		if audChans == 0 {
			fmt.Printf("file has no audio streams\n")
			continue
		}
		// рисуем индивидуальные дорожки
		arguments1 := makeIndividualWFs(path, audChans)
		fmt.Println("Команда для индивидуальных графиков:")
		fmt.Print("ffmpeg ")
		for _, a := range arguments1 {
			fmt.Print(a + " ")
		}
		fmt.Println("")
		if errCommence := commenseCommand(arguments1); errCommence != nil {
			fmt.Println(errCommence.Error())
			continue
		}
		// рисуем индивидуальные дорожки на одном файле
		arguments2 := makeJointWFs(path, audChans)
		fmt.Println("Команда для объединенных графиков:")
		fmt.Print("ffmpeg ")
		for _, a := range arguments2 {
			fmt.Print(a + " ")
		}
		fmt.Println("")
		if errCommence := commenseCommand(arguments2); errCommence != nil {
			fmt.Println(errCommence.Error())
			continue
		}
	}
}

func commenseCommand(argumentLine []string) error {
	//fmt.Println("ffmpeg " + argumentLine)
	for i, _ := range argumentLine {
		argumentLine[i] = argumentLine[i] + " "
	}
	comm1, errCrt := command.New(
		command.CommandLineArguments("ffmpeg", argumentLine...),
		command.Set(command.BUFFER_OFF),
		command.Set(command.TERMINAL_OFF),
	)
	if errCrt != nil {
		return fmt.Errorf("Не могу создать задачу: %v", errCrt.Error())
	}
	errRun := comm1.Run()
	if errRun != nil {
		return fmt.Errorf("Не могу выполнить задачу: %v", errRun.Error())
	}
	return nil
}

func checkInput() ([]string, error) {
	args := os.Args
	if len(args) < 2 {
		return nil, fmt.Errorf("found no input")
	}
	_, err := os.Stat(args[1])
	if err != nil {
		return nil, fmt.Errorf("file argument: (%v)\n%v", args[1], err.Error())
	}
	return args[1:], nil
}

/*
любопытная схема:
ffmpeg -i INPUT.mov -filter_complex "gradients=s=vga:c0=303030:c1=7c84cc:x0=0:x1=0:y0=0:y1=640 [bg];[0:a]showwaves=mode=p2p:s=vga:split_channels=1:scale=lin:draw=scale:colors=Yellow|White[v];[bg][v]overlay=shortest=1:format=rgb:alpha=premultiplied [vid]" -map "[vid]" -map 0:a  -vsync vfr -preset ultrafast -f matroska - | ffplay -autoexit -
*/

func makeIndividualWFs(filename string, audChan int) []string {
	//line := ` -y -i ` + filename + ` -filter_complex `
	fcArg := ""
	for i := 0; i < audChan; i++ {
		fcArg += fmt.Sprintf("[0:a:%v]compand,showwavespic=s=800x600:split_channels=1[a%v];", i, i)
	}
	fcArg = strings.TrimSuffix(fcArg, ";")
	outArg := ""
	for i := 0; i < audChan; i++ {
		outArg += fmt.Sprintf("-map [a%v] -frames:v 1 %v_wf%v.png ", i, filename, i)
	}
	line := fmt.Sprintf("-y -i %v -filter_complex %v %v", filename, fcArg, outArg)
	return strings.Fields(line)
}

func makeJointWFs(filename string, audChan int) []string {
	//line := ` -y `
	inputArgs := ""
	for i := 0; i < audChan; i++ {
		inputArgs += fmt.Sprintf("-i %v_wf%v.png ", filename, i)
	}
	//line += `-filter_complex "`
	fcArgs := ""
	for i := 0; i < audChan; i++ {
		fcArgs += fmt.Sprintf("[%v]", i)
	}
	fcArgs += fmt.Sprintf(`vstack=inputs=%v`, audChan)
	//line += fmt.Sprintf("%v_output.png", filename)
	line := fmt.Sprintf("-y %v -filter_complex %v %v_MERGED.png", inputArgs, fcArgs, filename)
	return strings.Fields(line)
}
