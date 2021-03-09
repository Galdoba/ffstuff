package main

import (
	"fmt"
	"os"
	"strings"
)

/*
TZ:
>> grab filename.ready			-- забрать все связанное с ready файлом
>> grab -h						-- вывести на экран помогалку 							--help
>> grab -n						-- забрать все новое (предварительное сканирование)		--new
>> grab -v						-- забрать только если одобряет инчекер					--valid
>> grab -p						-- забрать только звук и прокси							--proxy
>> grab -fc						-- забрать только если одобряет fflite @check0			--fflitecheck0

пред проверки:
-папка куда копировать
-отсуствие файла с таким же именем и размером
-наличие свободного места для копии

пост проверки:
-копия равна по имени и размеру с источником

*/

func main() {
	// logger := logfile.New(fldr.MuxPath()+"logfile.txt", logfile.LogLevelINFO)
	// args := os.Args

}

func argsReceived() []string {
	outArgs := []string{}
	for i, val := range os.Args {
		if len(os.Args) == 1 {
			fmt.Println("No аrguments received")
		}
		if i == 0 {
			continue
		}
		outArgs = append(outArgs, val)
	}
	return outArgs
}

func drawProgress(c, max int) {
	bar := []string{}
	for i := 0; i < 50; i++ {
		bar = append(bar, "-")
	}
	lim := max / c
	for i := 0; i < lim; i++ {
		bar[i] = "+"
	}
	fmt.Print(strings.Join(bar, ""), "\r")
}
