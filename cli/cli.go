package cli

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

//RunConsole - запускает в дефолтовом терминале cli программу
func RunConsole(program string, args ...string) {
	var line []string
	line = append(line, program)
	line = append(line, args...)
	fmt.Println("Run:", line)
	time.Sleep(time.Millisecond * 150)
	cmd := exec.Command(program, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Run()
}
