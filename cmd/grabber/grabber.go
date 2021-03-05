package main

import (
	"fmt"
	"os"
)

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
