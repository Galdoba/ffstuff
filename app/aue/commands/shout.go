package commands

import (
	"errors"
	"fmt"

	"github.com/Galdoba/ffstuff/app/aue/config"
	"github.com/Galdoba/ffstuff/pkg/logman"
	"github.com/urfave/cli/v2"
)

func Shout() *cli.Command {
	return &cli.Command{
		Name: "shout",
		//Aliases:     []string{"fs"},
		Usage: "test commands",

		BashComplete: func(*cli.Context) {
		},
		Before: func(c *cli.Context) error {
			cfgLoaded := config.NewConfig(c.App.Version)

			cfg = cfgLoaded
			fmt.Println(cfg.CONSOLE_LOG_LEVEL)
			// log.Setup(
			// 	log.WithAppLogLevelImportance(log.ImportanceALL),
			// )
			// log.SetOutput(cfg.AssetFiles[config.Asset_File_Log], log.ALL)

			return setupLogger(cfg)
		},
		After: func(c *cli.Context) error {
			return nil
		},
		Action: func(c *cli.Context) error {
			fmt.Println("Shout !!!")
			internalF()
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

func internalF() {
	logman.Warn("warn")
	logman.Error(errors.New("error"))
	logman.Info("info text")
	logman.Fatalf("Kill me")
	logman.Debug(logman.NewMessage("shout!!!").WithArgs(3, 3.14))
}
