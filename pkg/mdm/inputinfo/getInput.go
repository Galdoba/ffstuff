package inputinfo

import (
	"fmt"
	"os"
	"strings"

	"github.com/Galdoba/devtools/cli/command"
)

func ParseFile(path string) (*ParseInfo, error) {
	_, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("os.Stat(%v): %v", path, err)
	}
	//fmt.Println(path)
	com, err := command.New(
		command.CommandLineArguments(`ffmpeg`, "-hide_banner -i "+path),
		command.Set(command.BUFFER_ON),
		command.Set(command.TERMINAL_OFF),
	)
	if err != nil {
		return nil, fmt.Errorf("command.New(%v): %v", path, err)
	}
	err = com.Run()
	parseWarn := ""
	if err != nil {
		if err.Error() != "exit status 1" {
			parseWarn += err.Error()
		} else {

		}

	}
	buf := com.StdErr()
	//fmt.Println("ffmpeg output:\n", buf)
	input := inputdata{strings.Split(buf, "\n")}
	pi, err := parse(input)
	pi.buffer = buf
	if parseWarn != "" {
		pi.buffer += "\nparse: " + parseWarn
	}
	if err != nil {
		return nil, fmt.Errorf("parse(%v): %v", path, err)
	}
	return pi, err
}

func (pi *ParseInfo) NumAudio() int {
	return len(pi.Audio)
}

func (pi *ParseInfo) Buffer() string {
	return pi.buffer
}
