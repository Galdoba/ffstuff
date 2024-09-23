package commands

import (
	"github.com/urfave/cli/v2"
)

func Search() *cli.Command {
	return &cli.Command{
		Name:  "search",
		Usage: "TODO: Search and return marker files",
	}
}

func Run() *cli.Command {
	return &cli.Command{
		Name:  "run",
		Usage: "TODO: Endless loop to handle processing operations",
	}
}

func Queue() *cli.Command {
	return &cli.Command{
		Name:  "queue",
		Usage: "TODO: Add delayed operation(s)",
	}
}
