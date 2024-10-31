package commands

import (
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

func Search() *cli.Command {
	return &cli.Command{
		Name:  "search",
		Usage: color.HiBlackString("TODO: Search and return marker files"),
	}
}
