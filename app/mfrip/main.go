package main

import (
	"fmt"
	"os"

	"github.com/Galdoba/ffstuff/app/mfrip/cmd"
	"github.com/urfave/cli/v2"
)

const (
	CONFIG = "cfg"
)

/*




ripper -streams a:?:? -index -acodec alac [FILE]



*/

func main() {
	app := cli.NewApp()

	app.Version = "v 0.0.1"
	app.Usage = "rip streams and channels from mediafile"
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:      "use_config",
			Usage:     "use alternative config",
			TakesFile: false,
			Action: func(*cli.Context, string) error {
				return nil
			},
		},
	}

	//ДО НАЧАЛА ДЕЙСТВИЯ
	app.Before = func(c *cli.Context) error {
		return nil
	}

	app.Commands = []*cli.Command{
		cmd.RipStreams(),
		cmd.RipChannels(),
	}

	// defComm := func(c *cli.Context, s string) {
	// 	fmt.Println("defcom")
	// 	fmt.Println("action start")
	// 	fmt.Println(c.String("use_config"))
	// 	args := c.Args().Slice()
	// 	fmt.Println("args:", args)
	// }

	// app.CommandNotFound = defComm

	//ПО ОКОНЧАНИЮ ДЕЙСТВИЯ

	args := os.Args
	if err := app.Run(args); err != nil {
		errOut := fmt.Sprintf("%v error: %v", app.Name, err.Error())
		println(errOut)
		os.Exit(1)
	}

}
