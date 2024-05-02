package cmd

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func NewCommand() *cli.Command {
	cmnd := &cli.Command{
		Name:        "new_ommand",
		Aliases:     []string{},
		Usage:       "new_command usage",
		UsageText:   "short text",
		Description: "long text",
		// BashComplete: func(*cli.Context) {
		// },
		// Before: func(*cli.Context) error {
		// },
		// After: func(*cli.Context) error {
		// },
		Action: func(c *cli.Context) error {
			fmt.Println("new_command called")
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
