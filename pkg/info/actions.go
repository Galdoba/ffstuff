package info

import (
	"fmt"
	"strconv"

	"github.com/malashin/ffinfo"
)

/*
info.Duration(filepath string) float64
fileDuration := collectInfo(ch.data[path], 0, ffinfoDuration)
*/

//Duration - TODO: вернуть значение длинны видео в секундах
func Duration(filepath string) (float64, error) {
	err := fmt.Errorf("Err value not touched")
	f, err := ffinfo.Probe(filepath)
	durations := []float64{}
	for i := 0; i < len(f.Streams); i++ {
		dur, err := strconv.ParseFloat(f.Streams[i].Duration, 'f')
		if err != nil {
			return 0, fmt.Errorf("duration parse failed: %v", err.Error())
		}
		//d := strconv.FormatFloat(dur, 'f', 3, 64)
		durations = append(durations, dur)
	}
	dSum := 0.0
	for _, d := range durations {
		dSum += d
	}
	dAverage := dSum / float64(len(durations))

	return dAverage, err
}
