package inputinfo

import (
	"fmt"
	"os"
	"strings"

	"github.com/Galdoba/devtools/cli/command"
)

func ParseFile(path string) (*parseInfo, error) {
	_, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("os.Stat(%v): %v", path, err)
	}
	com, err := command.New(
		command.CommandLineArguments("ffmpeg", "-hide_banner -i "+path),
		command.Set(command.BUFFER_ON),
		//command.Set(command.TERMINAL_ON),
	)
	if err != nil {
		return nil, fmt.Errorf("command.New(%v): %v", path, err)
	}
	com.Run()
	buf := com.StdErr()
	fmt.Println("buffer:", buf)
	input := inputdata{strings.Split(buf, "\n")}
	pi, err := parse(input)
	if err != nil {
		return nil, fmt.Errorf("parse(%v): %v", path, err)
	}
	return pi, err
}
