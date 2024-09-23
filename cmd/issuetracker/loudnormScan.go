package main

// import (
// 	"fmt"
// 	"strconv"
// 	"strings"

// 	"github.com/Galdoba/utils"
// )

// func loudnormReportPath(path string) string {
// 	return strings.TrimSuffix(path, ".m4a") + "_loudnorm_report.txt"
// }

// func reportOnScanning(path string) []string {
// 	warns := []string{}
// 	lines := utils.LinesFromTXT(path)
// 	channels := []int{}
// 	for _, l := range lines {
// 		if strings.Contains(l, "RA: ") {
// 			data := strings.Split(l, ",")
// 			for _, dExact := range data {
// 				if !strings.Contains(dExact, "RA:") {
// 					continue
// 				}
// 				raStr := strings.Fields(dExact)
// 				lRangeL, err := strconv.ParseFloat(raStr[1], 64)
// 				if err != nil {
// 					warns = append(warns, err.Error())
// 				}
// 				if lrWarn := loudnesRangeWarning(lRangeL); lrWarn != "ok" {
// 					warns = append(warns, lrWarn)
// 				}
// 			}
// 		}
// 		if strings.Contains(l, "channels:") {
// 			data := strings.Split(l, ",")
// 			for _, dExact := range data {
// 				if !strings.Contains(dExact, "channels:") {
// 					continue
// 				}
// 				raStr := strings.Fields(dExact)
// 				for _, val := range raStr {
// 					if chanVol, err := strconv.Atoi(val); err == nil {
// 						channels = append(channels, chanVol)
// 					}
// 				}
// 				if chanWarn := channelWarning(channels); chanWarn != "ok" {
// 					warns = append(warns, chanWarn)
// 				}
// 			}
// 		}
// 	}
// 	warnsCleaned := []string{}
// 	for _, w := range warns {
// 		if w != "ok" {
// 			warnsCleaned = append(warnsCleaned, w)
// 		}
// 	}
// 	if len(warnsCleaned) == 0 {
// 		warnsCleaned = append(warnsCleaned, "ok")
// 	}

// 	return warnsCleaned
// }

// func prependStr(sl []string, elem string) []string {
// 	slNew := []string{elem}
// 	return append(slNew, sl...)
// }

// func channelWarning(channels []int) string {
// 	warn := "ok"
// 	switch len(channels) {
// 	default:
// 		return fmt.Sprintf("unexpected number of channels (%v channels)", len(channels))
// 	case 2:
// 		left := channels[0]
// 		right := channels[1]
// 		if left-right > 2 {
// 			return fmt.Sprintf("right channel loudness anomaly (%v Db)", left-right)
// 		}
// 		if right-left > 2 {
// 			return fmt.Sprintf("left channel loudness anomaly (%v Db)", right-left)
// 		}
// 	case 6:
// 		l := channels[0]
// 		r := channels[1]
// 		c := channels[2]
// 		//lfe := channels[3]
// 		ls := channels[4]
// 		rs := channels[5]
// 		for _, val := range channels {
// 			if c-val < 0 {
// 				return fmt.Sprintf("Center is not the loudest channel %v", channels)
// 			}
// 		}
// 		if l-r > 2 {
// 			return fmt.Sprintf("right channel loudness anomaly (%v Db)", l-r)
// 		}
// 		if r-l > 2 {
// 			return fmt.Sprintf("right channel loudness anomaly (%v Db)", r-l)
// 		}
// 		if ls-rs > 2 {
// 			return fmt.Sprintf("right channel loudness anomaly (%v Db)", ls-rs)
// 		}
// 		if rs-ls > 2 {
// 			return fmt.Sprintf("right channel loudness anomaly (%v Db)", rs-ls)
// 		}
// 	}
// 	return warn
// }

// func loudnesRangeWarning(lr float64) string {
// 	if lr > 20 {
// 		return fmt.Sprintf("Loudness Range invalid (%v)", lr)
// 	}
// 	return "ok"
// }
