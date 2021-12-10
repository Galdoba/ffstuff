package silence

import (
	"fmt"
	"math"
	"os"

	"github.com/Galdoba/ffstuff/pkg/cli"
)

type silence struct {
	coords []timeCoord
}

//Detect - слушает файл и возвращает координаты тишины
func Detect(path string) (*silence, error) {
	//debugMsg("START: Detect(f *os.File) (*silence, error)")
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	sOUT, sERR, cERR := cli.RunConsole("ffmpeg", "-i", f.Name(), "-loglevel", "verbose") //, "-loglevel", "quiet")

	if sOUT != "" {
		return nil, fmt.Errorf("sOUT != ''")
	}
	if sERR != nil {

	}
	if cERR != nil {
		//return nil, cERR
	}

	//fmt.Println(f.Name(), "correct\n ")
	//debugMsg("END: Detect(f *os.File) (*silence, error)")
	return nil, nil
}

//timeCoord - описывает кусок пустоты на треке. Пустотой считаем громкость ниже -104.5 Db
type timeCoord struct {
	start    float64
	duration float64
}

//end - предпологаймая точка конца пустоты
func (tc *timeCoord) end() float64 {
	return round(tc.start + tc.duration)
}

///////////Helpers
//round - округляет до 6 разрядов для единообразия с ffmpeg
func round(f float64) float64 {
	return math.Trunc(f/0.000001) * 0.000001
}

func debugMsg(s string) {
	fmt.Printf("debug Message: " + s + "\n")
}
