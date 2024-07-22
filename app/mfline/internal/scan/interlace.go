package scan

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/Galdoba/ffstuff/pkg/command"
)

func Interlace(path string) (string, error) {

	devnull := ""
	switch runtime.GOOS {
	case "linux":
		devnull = "/dev/null"
	case "windows":
		devnull = "NUL"
	}
	process, err := command.New(
		command.CommandLineArguments("ffmpeg", fmt.Sprintf("-hide_banner -filter:v idet -frames:v 9999 -an -f rawvideo -y %v -i %v", devnull, path)),
		command.Set(command.TERMINAL_OFF),
		command.AddBuffer("buf"),
	)
	if err != nil {
		fmt.Fprintf(os.Stdout, "process creation error: %v", err.Error())
		return "", err
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

	bufText := process.Buffer("buf").String()
	for !done {
		bufText = process.Buffer("buf").String()
		out := interlaceProgress(bufText)
		fmt.Printf("%v   \r", out)

	}
	wg.Wait()
	reportText := reportInterlace(process.Buffer("buf").String())
	fmt.Printf("%v     \n", reportText)

	return reportText, nil
}

func reportInterlace(bufText string) string {
	lines := strings.Split(bufText, "\n")
	replines := []string{}
	for _, line := range lines {
		if strings.Contains(line, "[Parsed_idet_") {
			replines = append(replines, line)
		}
	}
	data := strings.Join(replines, "|")
	fields := strings.Fields(data)
	numbers := []int{}
	for _, field := range fields {
		if n, err := strconv.Atoi(field); err == nil {
			numbers = append(numbers, n)
		}
	}
	if len(numbers) < 11 {
		return "no video data"
	}
	if len(numbers) > 11 {
		panic("multiple video data")
	}
	progressive := numbers[5] + numbers[9] + numbers[6] + numbers[10]
	total := 0
	for i := 2; i < len(numbers); i++ {
		total += numbers[i]
		switch i {
		case 6, 10:
			// total += numbers[i]
		}
	}
	pct := float64(progressive) / float64(total)
	val := int(pct * 10000)
	pctClean := float64(val) / 100
	text := fmt.Sprintf("progressive fields: %v%v", pctClean, "%")

	return text
	/*
		GARBAGE...
		[Parsed_idet_0 @ 000001eb2b7def00] Repeated Fields: Neither:   659 Top:    46 Bottom:    47
		[Parsed_idet_0 @ 000001eb2b7def00] Single frame detection: TFF:   421 BFF:     0 Progressive:   324 Undetermined:     7
		[Parsed_idet_0 @ 000001eb2b7def00] Multi frame detection: TFF:   425 BFF:     0 Progressive:   326 Undetermined:     1
	*/
}

func interlaceProgress(report string) string {
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
			if durationTotal > 400.0 {
				durationTotal = 400.0
			}
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
