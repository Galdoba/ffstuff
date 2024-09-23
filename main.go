package main

import (
	"fmt"
	"os"
	"syscall"
	"time"
)

type btnList struct {
	btns []button
}

type button struct {
	key   string
	runes []rune
}

func main() {
	f, err := os.Stat(`\\192.168.31.4\buffer\IN\_AWAIT\___Rezident_chief_of_station--FILM--Chief_of_Station_Forced_25fps.srt`)
	fmt.Println(err)
	d := f.Sys().(*syscall.Win32FileAttributeData)
	fmt.Println(d.CreationTime)
	cTime := time.Unix(0, d.CreationTime.Nanoseconds())
	fmt.Println(cTime)
	//cTime = time.Unix(0, d.CreationTime.Nanoseconds())
}
