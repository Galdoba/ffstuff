package main

import (
	"fmt"
	"os"

	"github.com/Galdoba/ffstuff/cli"
)

type InChecker interface {
	CheckValidity(string) error
}

func main() {
	for i, val := range os.Args {
		fmt.Print(i, " = '", val, "'\n")
		fmt.Print(" ===============\n")

		stdout, stderr, err := cli.RunConsole("ffinfo", val)
		fmt.Println(stdout)
		fmt.Println(stderr)
		fmt.Println(err)
		fmt.Print(" ===============\n")
	}
}

//go run inchecker.go f:\Work\petr_proj\___IN\IN_2021-02-05\Zloy_duh_AUDIORUS51.m4a f:\Work\petr_proj\___IN\IN_2021-02-05\Zvonok_poslednyaya_glava_AUDIORUS51.m4a f:\Work\petr_proj\___IN\IN_2021-02-05\Daniel_Isnt_Real_HD.mp4 f:\Work\petr_proj\___IN\IN_2021-02-05\Daniel_Isnt_Real_SD.mp4
