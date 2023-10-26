package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/Galdoba/devtools/cli/command"
)

func main() {
	//
	args := os.Args

	if len(args) > 3 {
		fmt.Println("changestatus failed: to many arguments")
		return
	}
	if len(args) < 3 {
		fmt.Println("changestatus failed: not enough arguments")
		return
	}
	if badStatusVal(args[2]) {
		fmt.Printf("changestatus failed: bad value '%v' for new status (expecting hex value 0-F)", args[2])
		return
	}
	if err := readStatus(args[1]); err != nil {
		fmt.Println(err.Error())
		return
	}
	err := writeStatus(args[1], args[2])
	if err != nil {
		fmt.Println("changestatus failed:", err.Error())
	}
}

func readStatus(key string) error {
	cm, err := command.New(
		command.CommandLineArguments("kval", fmt.Sprintf("read -from fftasks_status -k %v", key)),
		command.Set(command.BUFFER_ON),
	)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	if err = cm.Run(); err != nil {
		fmt.Println(err.Error())
		return err
	}
	stat := strings.TrimSuffix(cm.StdOut(), "\n")
	if badStatusVal(stat) {
		return fmt.Errorf(stat)
	}

	return nil
}

func writeStatus(key, status string) error {
	cm, err := command.New(
		command.CommandLineArguments("kval", fmt.Sprintf("write -to fftasks_status -k %v %v", key, status)),
		command.Set(command.BUFFER_ON),
	)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	if err = cm.Run(); err != nil {
		fmt.Println(err.Error())
		return err
	}
	out := cm.StdOut()
	if strings.TrimSuffix(out, "\n") != "" {
		fmt.Println(out)
		return fmt.Errorf("bad write output")
	}
	return nil
}

func badStatusVal(s string) bool {
	switch s {
	default:
		return true
	case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F":
		return false
	}
}
