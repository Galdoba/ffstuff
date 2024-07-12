package scan

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/Galdoba/ffstuff/pkg/command"
)

var ErrRWCheck = errors.New("read-write error detected")

func ReadWrite(path string) error {
	process, err := command.New(
		command.CommandLineArguments("ffmpeg", fmt.Sprintf("-i %v -map 0 -scodec srt -dcodec copy -f null NUL", path)),
		command.Set(command.TERMINAL_OFF),
		command.AddBuffer("buf"),
	)
	if err != nil {
		fmt.Fprintf(os.Stdout, "process creation error: %v", err.Error())
		return err
	}

	done := false
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		err = process.Run()
		if err != nil {
			fmt.Fprintf(os.Stderr, "process run error: %v", err.Error())
		}
		done = true
		wg.Done()
	}()
	interupted := false
	fmt.Printf("Read-Write Check: %v\n", path)
	bufText := process.Buffer("buf").String()
	for !done {
		bufText = process.Buffer("buf").String()
		out := progress(bufText)
		fmt.Printf("%v   \r", out)
		if containError(bufText) {
			process.Interrupt()
			interupted = true
		}
	}
	wg.Wait()
	switch interupted {
	case false:
		fmt.Printf("progress: done   \r")
	case true:
		fmt.Printf("                                                                                     \r")
		fmt.Println("read-write error detected: process terminated")
		return ErrRWCheck
	}

	return nil
}

func progress(report string) string {
	ln := strings.Split(report, "\n")
	last := len(ln) - 1
	if last < 0 {
		last = 0
	}
	lines := strings.Split(report, "\n")
	durationTotal := -1.0
	durationScanned := -1.0

	for _, line := range lines {
		if strings.Contains(line, "Duration:") {
			durationTotal = extractDuration(line)
			continue
		}
		if strings.Contains(line, "time=") {
			for _, segment := range strings.Split(line, "\r") {
				if segment != "" && strings.Contains(segment, "time=") {
					durationScanned = extractDuration(segment)
				}
			}
		}

	}
	if durationScanned == -1 {
		return "progress: 0%   "
	}
	prctInt := int(durationScanned / durationTotal * 10000)
	prctFl := float64(prctInt) / 100

	return fmt.Sprintf("progress: %v%v     ", prctFl, "%")
}

func extractDuration(line string) float64 {
	re := regexp.MustCompile(`\d+\:\d+\:\d+\.\d+`)
	str := re.FindString(line)
	part := strings.Split(str, ":")
	durFl := 0.0
	for i, p := range part {
		fl, err := strconv.ParseFloat(p, 64)
		if err != nil {
			return -1
		}
		switch i {
		case 0:
			durFl = fl * 3600
		case 1:
			durFl += fl * 60
		case 2:
			durFl += fl
		}
	}
	return durFl
}

func containError(buf string) bool {
	ln := strings.Split(buf, "\n")
	for _, l := range ln {
		if strings.Contains(l, strings.ToLower("Error")) {
			return true
		}
	}
	return false
}
