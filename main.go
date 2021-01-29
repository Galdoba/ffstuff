package main

import (
	"fmt"

	"github.com/Galdoba/ffstuff/clipmaker"
	"github.com/Galdoba/ffstuff/ediread"
)

func main() {
	edi, err := ediread.NewEdlData("f:\\Work\\petr_proj\\___IN\\IN_2021-01_29\\Vozvrashenie_v_golubuyu_lagunu\\Test.edl")
	fmt.Println(err)
	for _, clipData := range edi.Entry() {
		fmt.Println(clipData)
		fmt.Println(edi.Folder())
		cl, err := clipmaker.NewClip(clipData, edi.Folder())
		fmt.Println(err, cl)
		clipmaker.Create(cl)
	}
}
