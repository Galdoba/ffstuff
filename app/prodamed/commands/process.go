package commands

import (
	"fmt"
	"os"

	"github.com/Galdoba/ffstuff/app/prodamed/config"
	"github.com/urfave/cli/v2"
)

var in_dir = ""

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
			cfg, err := config.Load(c.String("config"))
			in_dir = cfg.Option.PATH[config.IN]
			if !isDir(in_dir) {
				return fmt.Errorf("config value '%v' is not a dir")
			}
			//rename
			if err := renameAmediaFiles(in_dir); err != nil {
				return fmt.Errorf("rename stage: %v", err)
			}
			//analize
			groups, err := sortAmediaFilesByEpisodes()
			if err != nil {
				return fmt.Errorf("sort stage: %v", err)
			}
			for k, v := range groups {
				fmt.Println(k)
				//make shell
				sc, err := NewCompiler(v)
				sc.injectLogistics(cfg)
				fmt.Println(err)
				fmt.Println(sc)
				if err := sc.GenerateShellFile(); err != nil {
					fmt.Println(err.Error())
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

func isDir(dir string) bool {
	f, _ := os.Stat(dir)
	return f.IsDir()
}
