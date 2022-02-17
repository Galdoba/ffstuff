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
	// rep := reportOnScanningSoundscan(export)
	// if err = addSummary(export, rep); err != nil {
	// 	return err
	// }
	return err
}

func reportOnScanningSoundscan(path string) []string {
	return []string{}
}

func MakeLoudnormReport(path, export string) error {
	listenReport, err := command.New(
		command.CommandLineArguments(fmt.Sprintf("loudnorm %v -scan", path)),
		command.WriteToFile(export),
	)
	if err != nil {
		return err
	}
	if err = listenReport.Run(); err != nil {
		return err
	}
	//rep := reportOnScanningLoudnorm(export)
	// if err = addSummary(export, rep); err != nil {
	// 	return err
	// }
	return err
}

type loudnormResult struct {
	filePath   string
	ra         float64
	channels   []int
	timeLevels []int
}

func LoudnormReportToString(path string) string {
	_, err := os.Stat(path)
	if err != nil {
		return `%data corrupted%`
	}
	//lines := readReport(path)
	rep := reportOnScanningLoudnorm(path)
	repStr := strings.Join(rep, " ")
	repStr = strings.TrimSpace(repStr)
	if repStr == "" {
		return "%%data corrupted%%"
	}
	return repStr
}

func addSummary(path string, rep []string) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err = f.WriteString("\nSUMMARY\n"); err != nil {
		return err
	}
	for _, line := range rep {
		if _, err = f.WriteString(line + "\n"); err != nil {
			return err
		}
	}
	return nil
}

func loudnormReportPath(path string) string {
	return strings.TrimSuffix(path, ".m4a") + "_loudnorm_report.txt"
}

func findRA(input string) string {
	re := regexp.MustCompile(`RA: ([0-9]*\.[0-9]*)`)
	return re.FindString(input)
}

func findChannels(input string) string {
	//channels: -3 -4 0 -10 -10 -10
	//re := regexp.MustCompile(`((channels: (-\d{1,}|0) (-\d{1,}|0) $)|(channels: (-\d{1,}|0) (-\d{1,}|0) (-\d{1,}|0) (-\d{1,}|0) (-\d{1,}|0) (-\d{1,}|0) $))`)
	re := regexp.MustCompile(`channels: (-\d{1,}|0).*$`)
	return re.FindString(input)
}

func findStats(input string) string {
	//channels: -3 -4 0 -10 -10 -10
	re := regexp.MustCompile(`(\d*<\d*<\d*)`)
	return re.FindString(input)
}

func reportOnScanningLoudnorm(path string) []string {
	warns := []string{}
	lines := readReport(path)
	for _, l := range lines {
		if strings.Contains(l, "RA: ") {
			ra := findRA(l)
			warns = append(warns, ra)
		}
		if strings.Contains(l, "channels:") {
			chn := findChannels(l)
			warns = append(warns, chn)
		}
		if strings.Contains(l, "ST stats clean") {
			chn := findStats(l)
			warns = append(warns, "stats: "+chn)
		}

	}
	return warns
}

func prependStr(sl []string, elem string) []string {
	slNew := []string{elem}
	return append(slNew, sl...)
}

func channelWarning(channels []int) string {
	warn := "ok"
	switch len(channels) {
	default:
		return fmt.Sprintf("unexpected number of channels (%v channels)", len(channels))
	case 2:
		left := channels[0]
		right := channels[1]
		if left-right > 2 {
			return fmt.Sprintf("right channel loudness anomaly (%v Db)", left-right)
		}
		if right-left > 2 {
			return fmt.Sprintf("left channel loudness anomaly (%v Db)", right-left)
		}
	case 6:
		l := channels[0]
		r := channels[1]
		c := channels[2]
		//lfe := channels[3]
		ls := channels[4]
		rs := channels[5]
		for _, val := range channels {
			if c-val < 0 {
				return fmt.Sprintf("Center is not the loudest channel %v", channels)
			}
		}
		if l-r > 2 {
			return fmt.Sprintf("right channel loudness anomaly (%v Db)", l-r)
		}
		if r-l > 2 {
			return fmt.Sprintf("right channel loudness anomaly (%v Db)", r-l)
		}
		if ls-rs > 2 {
			return fmt.Sprintf("right channel loudness anomaly (%v Db)", ls-rs)
		}
		if rs-ls > 2 {
			return fmt.Sprintf("right channel loudness anomaly (%v Db)", rs-ls)
		}
	}
	return warn
}

func loudnesRangeWarning(lr float64) string {
	if lr > 20 {
		return fmt.Sprintf("Loudness Range invalid (%v)", lr)
	}
	return "ok"
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
