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
			for _, f := range fi {
				if f.IsDir() {
					projects = append(projects, fmt.Sprintf("%v%v", cfg.IN_DIR, f.Name()))

				}
			}

			for _, project := range projects {
				fmt.Println("\n--------\nStart Project:", project)
				sources, err := actions.SetupSources(project, cfg.BUFFER_DIR)
				if len(sources) == 0 {
					fmt.Println("LOG ERROR:", "no sources created")
					continue
				}
				if err != nil {
					fmt.Println("||||||", err)
				}

				ja, err := job.New(sources, nil,
					job.WithInputDir(cfg.BUFFER_DIR),
					job.WithProcessingDir(cfg.IN_PROGRESS_DIR),
					job.WithDoneDir(cfg.DONE_DIR),
					job.WithOutDir(cfg.OUT_DIR),
				)
				if err != nil {
					fmt.Println("LOG ERROR: job creation", err.Error())
				}
				if err := ja.DecideType(); err != nil {
					fmt.Println("LOG ERROR: job decide type", err.Error())
					return err
				}
				if err := ja.CompileTasks(); err != nil {
					fmt.Println("LOG ERROR: job compile tasks", err.Error())
					return err
				}
				if err := ja.Execute(); err != nil {
					fmt.Println("LOG ERROR: job execute", err.Error())
					//return err
				}

			}
			fmt.Println("GOOD ENDING!!!")
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
