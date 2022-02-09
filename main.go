package main

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Galdoba/devtools/cli/command"
	"github.com/Galdoba/ffstuff/pkg/info"
	"github.com/Galdoba/utils"
)

func main() {
	path := "d:\\IN\\IN_2022-02-09\\Eiffel_AUDIORUS51.m4a"
	reportFile := strings.TrimSuffix(path, ".m4a") + "_loudnorm_report.txt"
	fmt.Printf("%v\nChecking...\n", path)
	fmt.Println(info.Duration(path))
	if _, err := os.Stat(reportFile); err == nil {
		fmt.Println("Report file already exist")
		rep := reportOnScanning(reportFile)
		for _, s := range rep {
			fmt.Println(s)
		}
		return
	}

	listenReport, err := command.New(
		command.CommandLineArguments(fmt.Sprintf("loudnorm %v -scan", path)),
		//command.Set(command.BUFFER_ON),
		//command.Set(command.TERMINAL_ON),
		command.WriteToFile(reportFile),
	)
	if err != nil {
		fmt.Println("error:", err.Error())
	}

	timeTaken := 0
	reportComplete := false
	go func() {
		listenReport.Run()
		reportComplete = true
	}() // чтобы трэчить текущий прогресс нужно ставить флаг читать файл? НЕТ. смотри Run_untested()
	for !reportComplete {
		time.Sleep(time.Second)
		timeTaken++
		fmt.Printf("loudnorm scanning... %v\r", timeTaken)
	}
	fmt.Printf("loudnorm scanning completed after %v seconds\n", timeTaken)
	rep := reportOnScanning(reportFile)
	if len(rep) > 0 {
		rep = prependStr(rep, path)
		rep = prependStr(rep, listenReport.StdOut())
	}
	for _, s := range rep {
		fmt.Println(s)
	}
}

func reportOnScanning(path string) []string {
	warns := []string{}
	lines := utils.LinesFromTXT(path)
	channels := []int{}
	for _, l := range lines {
		if strings.Contains(l, "RA: ") {
			data := strings.Split(l, ",")
			for _, dExact := range data {
				if !strings.Contains(dExact, "RA:") {
					continue
				}
				raStr := strings.Fields(dExact)
				lRangeL, err := strconv.ParseFloat(raStr[1], 64)
				if err != nil {
					warns = append(warns, err.Error())
				}
				if lrWarn := loudnesRangeWarning(lRangeL); lrWarn != "ok" {
					warns = append(warns, lrWarn)
				}
			}
		}
		if strings.Contains(l, "channels:") {
			data := strings.Split(l, ",")
			for _, dExact := range data {
				if !strings.Contains(dExact, "channels:") {
					continue
				}
				raStr := strings.Fields(dExact)
				for _, val := range raStr {
					if chanVol, err := strconv.Atoi(val); err == nil {
						channels = append(channels, chanVol)
					}
				}
				if chanWarn := channelWarning(channels); chanWarn != "ok" {
					warns = append(warns, chanWarn)
				}
			}
		}
	}
	warnsCleaned := []string{}
	for _, w := range warns {
		if w != "ok" {
			warnsCleaned = append(warnsCleaned, w)
		}
	}

	return warnsCleaned
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

func PixelColorTest() {
	imageFile, err := os.Open("c:\\Users\\pemaltynov\\go\\src\\github.com\\Galdoba\\ffstuff\\assets\\waveform_test_sqrt_t.png")
	if err != nil {
		panic(1)
	}
	imData, imType, err := image.Decode(imageFile)
	fmt.Println(imData.At(2000, 100))
	fmt.Println(imType)
	if err != nil {
		panic(2)
	}
	imageFile.Seek(0, 0)
	loadedIm, err := png.Decode(imageFile)
	if err != nil {
		panic(3)
	}
	rec := loadedIm.Bounds()
	i := 1
	empty := 0
	filled := 0
	for y := 0; y < rec.Dy(); y++ {
		for x := 0; x < rec.Dx(); x++ {
			if y == rec.Dy()/2 {
				fmt.Printf("Pixel %v	 (%v, %v) is [%v]\n", i, x, y, loadedIm.At(x, y))
			}
			r, g, b, a := loadedIm.At(x, y).RGBA()
			if r+g+b+a == 0 {
				empty++
			} else {

				//fmt.Printf("Pixel %v	 (%v, %v) is [%v]\n", i, x, y, loadedIm.At(x, y))
				filled++
			}

			i++
		}
	}
	fmt.Println("done", empty, filled)
}
