package format

import (
	"fmt"
	"testing"
)

func inputFormat() []string {
	return []string{
		Trailer4K,
		TrailerHD,
		TrailerSD,
		Film4K,
		FilmHD,
		FilmSD,
		PureSound,
		"",
		"[Invalid]",
	}
}

type dimentionInput struct {
	width  int
	height int
}

func inputDimentions() []dimentionInput {
	di := []dimentionInput{}
	for i := -1; i < 11; i++ {
		di = append(di, dimentionInput{i * 400, i * 300})
	}
	di = append(di, dimentionInput{720, 576})
	di = append(di, dimentionInput{1920, 1080})
	di = append(di, dimentionInput{3840, 2160})
	return di
}

func TestDimentions(t *testing.T) {
	for _, dim := range inputDimentions() {
		di := NewDimention(dim.width, dim.height)

		fmt.Println(di)
	}

}

func TestFormat(t *testing.T) {

}
