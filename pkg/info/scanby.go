package info

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/Galdoba/devtools/cli/command"
)

const (
	SoundscanFAST = iota
	SoundscanFULL
)

func MakeSoundscanReport(path, export string, checkDepth int) error {
	lb := ""
	dur := ""
	switch checkDepth {
	default:
		return fmt.Errorf("checkDepth option unrecognized")
	case SoundscanFAST:
		lb = "-lb 76"
		dur = "-d 2"
	case SoundscanFULL:
		lb = "-lb 84"
		dur = "-d 0.04"
	}
	listenReport, err := command.New(
		command.CommandLineArguments(fmt.Sprintf("soundscan %v %v %v -scan", path, lb, dur)),
		command.WriteToFile(export),
	)
	if err != nil {
		return err
	}
	if err = listenReport.Run(); err != nil {
		return err
	}
	return err
}

func reportOnScanningSoundscan(path string) []string {
	return []string{}
}

func MakeLoudnormReport(path, export string) error {
	listenReport, err := command.New(
		command.Set(command.BUFFER_ON),
		command.CommandLineArguments(fmt.Sprintf("loudnorm %v -scan", path)),
		command.WriteToFile(export),
	)
	if err != nil {
		return err
	}
	if err = listenReport.Run(); err != nil {
		return err
	}
	if err := exportTest("---OUT:\n"+listenReport.StdOut()+"---ERROR:\n"+listenReport.StdErr(), export+"2"); err != nil {
		return err
	}
	return nil
}

func exportTest(buffer, path string) error {
	file, fileErr := os.Create(path)
	if fileErr != nil {
		fmt.Println(fileErr)
		return fmt.Errorf("exportTest: %v", fileErr.Error())
	}
	fmt.Fprintf(file, "%s\n", buffer)
	return nil
}

func LoudnormData(path string) []string {
	_, err := os.Stat(path)
	if err != nil {
		return []string{"data corrupted", "data corrupted", "data corrupted", "data corrupted"}
	}
	rep := reportOnScanningLoudnorm(path)
	return rep
}

func loudnormReportPath(path string) string {
	return strings.TrimSuffix(path, ".m4a") + "_loudnorm_report.txt"
}

func findIL(input string) string {
	re := regexp.MustCompile(`I: -([0-9]*\.[0-9]*)`)
	str := re.FindString(input)
	str = strings.TrimPrefix(str, "I:")
	str = strings.TrimSpace(str)
	return str
}

func findRA(input string) string {
	re := regexp.MustCompile(`RA: ([0-9]*\.[0-9]*)`)
	str := re.FindString(input)
	str = strings.TrimPrefix(str, "RA:")
	str = strings.TrimSpace(str)
	return str
}

func findChannels(input string) string {
	re := regexp.MustCompile(`channels: (-\d{1,}|0).*$`)
	str := re.FindString(input)
	str = strings.TrimPrefix(str, "channels: ")
	str = strings.TrimSpace(str)
	return str
}

func findStats(input string) string {
	re := regexp.MustCompile(`(\d*<\d*<\d*)`)
	str := re.FindString(input)
	str = strings.TrimSpace(str)
	return re.FindString(input)
}

func reportOnScanningLoudnorm(path string) []string {
	warns := []string{}
	searchData := []string{"IL", "RA", "CHN", "STATS"}
	foundData := make(map[string]string)
	lines := readReport(path)
	for _, d := range searchData {
		foundData[d] = "data corrupted"
	}
	for _, l := range lines {
		switch {
		case strings.Contains(l, "RA: "):
			il := findIL(l)
			foundData["IL"] = il
			ra := findRA(l)
			foundData["RA"] = ra
		case strings.Contains(l, "channels:"):
			chn := findChannels(l)
			foundData["CHN"] = chn
		case strings.Contains(l, "ST stats clean"):
			chn := findStats(l)
			foundData["STATS"] = chn
		}

	}
	for _, d := range searchData {
		warns = append(warns, foundData[d])
	}
	return warns
}

func prependStr(sl []string, elem string) []string {
	slNew := []string{elem}
	return append(slNew, sl...)
}

func readReport(path string) []string {
	var lines []string
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return lines
}
