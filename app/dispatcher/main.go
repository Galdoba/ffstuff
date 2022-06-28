package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Version = "v 0.0.0"
	app.Name = "director"
	app.Usage = "Scans media streams to tell the story..."
	app.Flags = []cli.Flag{}
	app.Commands = []cli.Command{}
	args := []string{""}
	if len(args) == 0 {
		fmt.Println("No arguments provided")
		os.Exit(0)
	}
	if err := app.Run(args); err != nil {
		fmt.Println("Here is my error:")
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Println("Here is my story: ...")

}
