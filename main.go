package main

import (
	"fmt"
	"image"
	"image/png"
	"os"
)

func main() {
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
