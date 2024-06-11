package cmd

import (
	"fmt"

	"github.com/Galdoba/ffstuff/app/autogen/config"
	"github.com/urfave/cli/v2"
)

func Config(cfg config.Config) *cli.Command {
	cmnd := &cli.Command{
		Name:      "config",
		Aliases:   []string{},
		Usage:     "print config file content",
		UsageText: "autogen config",
		// BashComplete: func(*cli.Context) {
		// },
		// Before: func(*cli.Context) error {
		// },
		// After: func(*cli.Context) error {
		// },
		Action: func(c *cli.Context) error {

			fmt.Println("config command called")
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
