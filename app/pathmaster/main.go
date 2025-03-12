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
	app.Name = "mfrip"
	app.Version = "v 0.0.1"
	app.Usage = "rip streams and channels from mediafile"
	app.Flags = []cli.Flag{}

	app.Commands = []*cli.Command{
		cmd.RipStreams(),
		cmd.RipChannels(),
	}

	args := os.Args
	if err := app.Run(args); err != nil {
		errOut := fmt.Sprintf("%v error: %v", app.Name, err.Error())
		println(errOut)
		os.Exit(1)
	}

}
