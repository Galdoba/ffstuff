package main

import (
	"fmt"
	"io"
	"os"

	xpp "github.com/mmcdole/goxpp"
	"golang.org/x/net/html/charset"
)

func main() {
	file, err := os.Open("f:\\Work\\petr_proj\\___IN\\IN_2021-02-15\\Pretty_Little_Liars_s01_[C]__hd_q0w0_rus20.xml")
	fmt.Println(err)
	fmt.Println(file)
	parser := xpp.NewXMLPullParser(file, false, charset.NewReader(io.Reader{}, "11"))
}
