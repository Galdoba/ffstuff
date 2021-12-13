package cli

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/Galdoba/utils"
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
	//fmt.Println("Run:", line)

	time.Sleep(time.Millisecond * 2)
	cmd := exec.Command(program, args...)
	cmd.Stdout = os.Stdout //&o
	output, _ := cmd.CombinedOutput()
	sOUT := string(output)
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	return sOUT, cmd.Stderr, err
}

//RunConsole - запускает в дефолтовом терминале cli программу
func RunToFile(file, program string, args ...string) (io.Writer, error) {
	var line []string
	line = append(line, program)
	line = append(line, args...)
	//fmt.Println("Run:", line)
	var o bytes.Buffer
	var e bytes.Buffer
	time.Sleep(time.Millisecond * 2)
	cmd := exec.Command(program, args...)
	cmd.Stdout = &o //os.Stdout
	//output, _ := cmd.CombinedOutput()
	//sOUT := string(output)
	cmd.Stderr = &e //os.Stderr
	err := cmd.Run()
	sOUT := string(o.Bytes()) + "\n" + string(e.Bytes())
	utils.AddLineToFile(file, sOUT)
	return cmd.Stderr, err
}

//RunToBuffer - запускает в дефолтовом терминале cli программу но результат пишет в переменные
func RunToBuffer(program string, args ...string) (string, string, error) {
	var line []string
	line = append(line, program)
	line = append(line, args...)
	var o bytes.Buffer
	var e bytes.Buffer
	time.Sleep(time.Millisecond * 2)
	cmd := exec.Command(program, args...)
	cmd.Stdout = &o //os.Stdout
	cmd.Stderr = &e //os.Stderr

	err := cmd.Run()
	return o.String(), e.String(), err
}

//RunToAll - запускает в дефолтовом терминале cli программу и копирует stdout/stderr в переменные
func RunToAll(program string, args ...string) (string, string, error) {
	var line []string
	line = append(line, program)
	line = append(line, args...)
	var o bytes.Buffer
	var e bytes.Buffer
	time.Sleep(time.Millisecond * 2)
	cmd := exec.Command(program, args...)
	cmd.Stdout = io.MultiWriter(os.Stdout, &o)
	cmd.Stderr = io.MultiWriter(os.Stderr, &e)
	err := cmd.Run()
	return string(o.Bytes()), string(e.Bytes()), err
}

type Task struct {
	program   string
	agruments []string
}

//NewTask - создает объект содержащий инструкции для отправки в командную строку
func NewTask(program string, arguments []string) Task {
	return Task{program, arguments}
}

//String - возвращает строку для консоли
func (t Task) String() string {
	s := t.program
	for _, arg := range t.agruments {
		s = s + " " + arg
	}
	return s
}

//Run - выполняет инструкции
func (t *Task) Run() error {
	_, _, err := RunConsole(t.program, t.agruments...)
	return err
}

//Run - выполняет инструкции
func (t *Task) LastArg() string {
	l := len(t.agruments)
	return t.agruments[l-1]
}
