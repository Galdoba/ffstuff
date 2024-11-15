package commands

import (
	"github.com/Galdoba/devtools/decidion/operator"
	"github.com/Galdoba/ffstuff/app/aue/config"
	"github.com/urfave/cli/v2"
)

func Setup() *cli.Command {
	return &cli.Command{
		Name: "setup",
		//Aliases:     []string{"fs"},
		Usage: "test commands",

		BashComplete: func(*cli.Context) {
		},
		Before: func(c *cli.Context) error {

			return nil
		},
		After: func(c *cli.Context) error {
			return nil
		},
		Action: func(c *cli.Context) error {
			if !operator.Confirm("Create config?") {
				return nil
			}
			cfg = config.NewConfig(c.App.Version)

			// log.Setup(
			// 	log.WithAppLogLevelImportance(log.ImportanceALL),
			// )
			// log.SetOutput(cfg.AssetFiles[config.Asset_File_Log], log.ALL)

			return config.Save(cfg)
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
