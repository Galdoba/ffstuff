package pic

import (
	"fmt"
	"image/png"
	"log"
	"os"
	"time"

	"github.com/Galdoba/utils"
)

func Run() {
	t0 := time.Now().Nanosecond()

	file, _ := os.Open("c:\\Users\\pemaltynov\\go\\src\\github.com\\Galdoba\\ffstuff\\assets\\testscreen15_240.png")
	img, err := png.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	bpix := 0
	totalPix := img.Bounds().Max.Y * img.Bounds().Max.X
	for y := 0; y < img.Bounds().Max.Y; y++ {
		for x := 0; x < img.Bounds().Max.X; x++ {
			//color := img.At(x, y)
			r, g, b, _ := img.At(x, y).RGBA()
			r8, g8, b8 := to8bit(r, g, b)
			// fmt.Printf("%v", color)

			col := utils.RoundFloat64(float64(r8+g8+b8)/(256*3), 3)
			if col < 0.002 {
				bpix++
			}
			//if col != 0 {
			//	fmt.Printf("(%v-%v) red %d  green %d  blue %d  alpha %d | BW COLOR: %v\n", x, y, r/256, g/256, b/256, a/256, col)
			//}

		}
	}
	blPix := bpix * 100 / totalPix
	fmt.Println("total", blPix, "%")
	//t1 := time.Now().Nanosecond()
	dur := time.Duration(t0)
	fmt.Println("duration Nano", dur)
}

func to8bit(r, g, b uint32) (r8 uint32, g8 uint32, b8 uint32) {
	r8, g8, b8 = r/256, g/256, b/256
	return r8, g8, b8
}

/*
1920 x 1080
 960 x  540
 480 x  270
 240 x  135

*/
