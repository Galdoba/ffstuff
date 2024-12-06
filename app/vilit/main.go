package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

func main() {
	args := os.Args
	str := ""
	for i, arg := range args {
		switch i {
		case 0:
			continue
		case 1:
		default:
			str += " "
		}
		str += arg

	}
	literas := strings.Split(str, "")
	for _, lit := range literas {
		lower := strings.ToLower(lit)
		if strings.Contains("0123456789\\/!@#$^&*(){}[],._-=+<>~`':;?"+`"`+"%", lower) {
			fmt.Fprintf(os.Stderr, "%v", color.YellowString("%v", lit))
			continue
		}
		if strings.Contains("abcdefghijklmnopqrstuvwxyz", lower) {
			fmt.Fprintf(os.Stderr, "%v", lit)
			continue
		}
		// if strings.Contains("абвгдеёжзийклмнопрстуфхцчшщьыъэюя", lower) {
		// 	fmt.Fprintf(os.Stderr, "%v", color.HiRedString("%v", lit))
		// 	continue
		// }
		fmt.Fprintf(os.Stderr, "%v", color.HiRedString("%v", lit))
	}
}
