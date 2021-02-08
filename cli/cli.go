package cli

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"
)

//TODO: перестроить так чтобы запускальщик работал по типу фабрики
/*
Пример:
cli.Run(                                                           //-имя явно поменять что-то типа RunProgram или Setup
	runner.SetProgram("ffmpeg"),
	runner.SetArguments(arguments...),
	runner.SetStdOutOption(int)
	runner.SetStdErrOption(int)
)

*/

//RunConsole - запускает в дефолтовом терминале cli программу
func RunConsole(program string, args ...string) (string, io.Writer, error) {
	var line []string
	line = append(line, program)
	line = append(line, args...)
	fmt.Println("Run:", line)
	time.Sleep(time.Millisecond * 2)
	cmd := exec.Command(program, args...)
	cmd.Stdout = os.Stdout
	output, _ := cmd.CombinedOutput()
	sOUT := string(output)
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return sOUT, cmd.Stderr, err
}

type Task struct {
	program   string
	agruments []string
}

//NewTask - создает объект содержащий инструкции для отправки в командную строку
func NewTask(program string, arguments []string) Task {
	return Task{program, arguments}
}

//Run - выполняет инструкции
func (t *Task) Run() error {
	_, _, err := RunConsole(t.program, t.agruments...)
	return err
}
