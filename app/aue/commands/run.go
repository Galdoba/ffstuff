package commands

import (
	"fmt"
	"os"

	"github.com/Galdoba/ffstuff/app/aue/config"
	"github.com/Galdoba/ffstuff/app/aue/internal/actions"
	"github.com/Galdoba/ffstuff/app/aue/internal/job"
	"github.com/urfave/cli/v2"
)

var cfg *config.Configuration

func Run() *cli.Command {
	return &cli.Command{
		Name: "run",
		//Aliases:     []string{"fs"},
		Usage: "main command",

		BashComplete: func(*cli.Context) {
		},
		Before: func(c *cli.Context) error {
			cfgDef := config.Default()
			cfgDef.Save()
			cfgLoaded, err := config.Load()
			if err != nil {
				return fmt.Errorf("config loading failed: %v", err)
			}
			cfg = cfgLoaded
			return nil
		},
		After: func(c *cli.Context) error {
			return nil
		},
		Action: func(c *cli.Context) error {
			fmt.Println("start run")
			fmt.Println(cfg.IN_DIR)
			fi, err := os.ReadDir(cfg.IN_DIR)
			if err != nil {
				return err
			}
			projects := []string{}
			for i, f := range fi {
				fmt.Println(i, f.Name(), "is dir:", f.IsDir())
				if f.IsDir() {
					projects = append(projects, fmt.Sprintf("%v%v", cfg.IN_DIR, f.Name()))
				}
			}
			fmt.Println(projects)
			for _, project := range projects {
				sources, err := actions.SetupSources(project, cfg.BUFFER_DIR)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println(sources)
				ja, err := job.New(sources, nil)
				fmt.Println(err)
				fmt.Println(&ja)
				fmt.Println(ja.DecideType())
				fmt.Println(&ja)

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
