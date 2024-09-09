package main

import (
	"fmt"
	"os"

	"github.com/Galdoba/ffstuff/app/reporter/commands"
	"github.com/urfave/cli/v2"
)

func main() {

	app := cli.NewApp()

	app.Version = "0.1.0"
	app.Usage = "handle report files"
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:     "file",
			Usage:    "set report file path",
			Required: true,
			Aliases:  []string{"f"},
		},
	}
	app.Commands = []*cli.Command{
		commands.New(),
		commands.Add(),
		commands.Find(),
		commands.FindAll(),
	}
	args := os.Args
	if err := app.Run(args); err != nil {
		errOut := fmt.Sprintf("%v error: %v", app.Name, err.Error())
		println(errOut)
		os.Exit(1)
	}

}

/*

reporter -f path/to/file new -fields "key1=value1;key2=value2"
reporter -f path/to/file add -fields "key1=value1;key2=value2"
reporter -f path/to/file find [keys]...
reporter -f path/to/file findall -keyPart


*/
