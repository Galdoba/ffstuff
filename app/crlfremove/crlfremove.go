package main

import (
	"os"

	"github.com/Galdoba/utils"
)

func main() {
	args := os.Args[1:]
	for _, arg := range args {
		for i, line := range utils.LinesFromTXT(arg) {
			utils.EditLineInFile(arg, i, line)
		}
	}
}
