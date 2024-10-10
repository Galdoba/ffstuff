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

func Run() *cli.Command {
	return &cli.Command{
		Name:  "run",
		Usage: color.HiBlackString("TODO: Endless loop to handle processing operations"),
	}
}

func Queue() *cli.Command {
	return &cli.Command{
		Name:  "queue",
		Usage: color.HiBlackString("TODO: Add delayed operation(s)"),
	}
}
