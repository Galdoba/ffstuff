package commands

import (
	"fmt"

	"github.com/Galdoba/ffstuff/app/aue/config"
	log "github.com/Galdoba/ffstuff/app/aue/logger"
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
			cfgLoaded, err := config.Load()
			if err != nil {
				return fmt.Errorf("config loading failed: %v", err)
			}
			cfg = cfgLoaded
			return log.Setup(
				log.LogFilepath(cfg.AssetFiles[config.Asset_File_Log]),
				log.DebugMode(cfg.DebugMode),
			)
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
	log.Warn("warn")
	log.Error("Errr")
	log.Info("info text")
	log.Fatal("Kill me")
	log.Debug("shout!!!", 3, 3.14)
}
