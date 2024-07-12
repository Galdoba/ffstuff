package cmd

import (
	"errors"
	"fmt"
	"sort"

	"github.com/Galdoba/ffstuff/app/mfline/internal/scan"
	"github.com/urfave/cli/v2"
)

func RWCheck() *cli.Command {
	return &cli.Command{
		Name:        "rwcheck",
		Aliases:     []string{"rw"},
		Usage:       "TODO USAGE",
		UsageText:   "TODO Usage text",
		Description: "",
		Args:        false,
		ArgsUsage:   "",
		Category:    "",

		BashComplete: func(*cli.Context) {
		},
		Before: func(c *cli.Context) error {

			return nil
		},
		After: func(c *cli.Context) error {
			return nil
		},
		Action: func(c *cli.Context) error {
			// db, err := db.New(`c:\Users\pemaltynov\.ffstuff\data\mfline\`)
			// if err != nil {
			// 	return err
			// }
			args := c.Args().Slice()
			sort.Strings(args)
			for _, arg := range args {
				err := scan.ReadWrite(arg)
				if err != nil && errors.Is(err, scan.ErrRWCheck) {
					fmt.Println("error:", err.Error())
					fmt.Println("Write to Db")
				}
			}
			return nil
		},
		OnUsageError: func(cCtx *cli.Context, err error, isSubcommand bool) error {
			return nil
		},
		Subcommands:            []*cli.Command{},
		Flags:                  []cli.Flag{},
		SkipFlagParsing:        false,
		HideHelp:               false,
		HideHelpCommand:        false,
		Hidden:                 false,
		UseShortOptionHandling: false,
		HelpName:               "",
		CustomHelpTemplate:     "",
	}

}
