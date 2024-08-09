package commands

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

var in_dir = `\\192.168.31.4\buffer\IN\`

func Process() *cli.Command {
	return &cli.Command{
		Name: "process",
		//Aliases:     []string{"fs"},
		Usage:     "use prompt mode",
		UsageText: "demux prompt [args]",

		BashComplete: func(*cli.Context) {
		},
		Before: func(c *cli.Context) error {
			return nil
		},
		After: func(c *cli.Context) error {
			return nil
		},
		Action: func(c *cli.Context) error {
			fmt.Println("rename")

			//rename
			if err := renameAmediaFiles(in_dir); err != nil {
				return fmt.Errorf("rename stage: %v", err)
			}
			groups, err := sortAmediaFilesByEpisodes()

			if err != nil {
				return fmt.Errorf("sort stage: %v", err)
			}
			for k, v := range groups {
				fmt.Println(k)
				sc, err := NewCompiler(v)
				fmt.Println(err)
				fmt.Println(sc)

			}

			//analize
			//make shell

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
