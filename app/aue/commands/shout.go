package commands

import (
	"errors"
	"fmt"

	"github.com/Galdoba/ffstuff/app/aue/config"
	log "github.com/Galdoba/ffstuff/pkg/logman"
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
			log.Setup(
				log.WithAppLogLevelImportance(log.ImportanceALL),
			)
			log.SetOutput(cfg.AssetFiles[config.Asset_File_Log], log.ALL)
			return nil
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
	log.Error(errors.New("error"))
	log.Info("info text")
	log.Fatalf("Kill me")
	log.Debug(log.NewMessage("shout!!!").WithArgs(3, 3.14), nil)
}
