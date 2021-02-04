package cli

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

//RunConsole - запускает в дефолтовом терминале cli программу
func RunConsole(program string, args ...string) error {
	var line []string
	line = append(line, program)
	line = append(line, args...)
	fmt.Println("Run:", line)
	time.Sleep(time.Millisecond * 100)
	cmd := exec.Command(program, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	return err
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
	return RunConsole(t.program, t.agruments...)
}
