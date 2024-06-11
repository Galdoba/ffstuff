package cmd

import (
	"fmt"

	"github.com/Galdoba/devtools/printer"
	"github.com/Galdoba/ffstuff/app/autogen/config"
	"github.com/urfave/cli/v2"
)

var TiketFileStorage string
var TableFile string
var pm printer.Printer

func Serial(cfg config.Config) *cli.Command {
	cmnd := &cli.Command{
		Name:      "ser",
		Aliases:   []string{},
		Usage:     "generate script for serial",
		UsageText: "bashgen srt [parameters]",
		// Description: "Track files in InputDir and manage job tickets. TODO: steps descr",
		// BashComplete: func(*cli.Context) {
		// },
		Before: func(*cli.Context) error {
			fmt.Println("Validate parameters")
			fmt.Println("DONE")
			return nil
		},
		// After: func(*cli.Context) error {
		// },
		Action: func(c *cli.Context) error {
			fmt.Println("load config")
			fmt.Println("setup logger")
			fmt.Println("generate bash script")
			return nil
		},
		// OnUsageError: func(cCtx *cli.Context, err error, isSubcommand bool) error {
		// },
		Subcommands: []*cli.Command{},
		Flags:       []cli.Flag{
			// &cli.StringFlag{},
		},
	}
	return cmnd
}
