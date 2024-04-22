package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/Galdoba/utils"
)

func main() {
	args := os.Args[1:]
	for _, arg := range args {
		for i, line := range utils.LinesFromTXT(arg) {
			line = strings.TrimSpace(line)
			fmt.Printf("'%v'\n", line)
			utils.EditLineInFile(arg, i, line)
		}
	}
}
