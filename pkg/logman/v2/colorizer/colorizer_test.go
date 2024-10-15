package colorizer

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/gookit/color"
)

func TestColorizer(t *testing.T) {
	var colTest uint8
	for i := colTest; i <= 255; i++ {
		if i%16 == 0 {
			fmt.Println("")
		}
		col := color.S256(i, i)
		col.Print("12")
		fmt.Print(" ", i, " ")
		if i < 10 {
			fmt.Print(" ")
		}
		if i < 100 {
			fmt.Print(" ")
		}

		if i == 255 {
			fmt.Println("")
			break
		}
	}

	cl := DefaultScheme()
	for n, arg := range []interface{}{
		"test string",
		6,
		3.141597,
		true,
		rune(64),
		byte(64),
		[]string{"asd", "123"},
		[]int{11, 45},
		[]bool{},
		fgKey("aaa"),
		[]colorKey{fgKey("1"), bgKey("2")},
		//map[string]int{"11": 11},
		testStr(),
	} {
		fmt.Println("--------", n)
		s := cl.Colorize(arg)
		argType := reflect.ValueOf(arg).Type().String()
		fmt.Printf("this is arg '%v' (%v %v) of type '%v'\n", s, cl.color256[fgKey(argType)], cl.color256[bgKey(argType)], argType)
	}
}
