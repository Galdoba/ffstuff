package commands

import (
	"fmt"
	"os"
	"time"

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
			fmt.Println("start before run")
			cfgDef := config.Default()
			cfgDef.Save()
			cfgLoaded, err := config.Load()
			if err != nil {
				return fmt.Errorf("config loading failed: %v", err)
			}
			cfg = cfgLoaded
			defer fmt.Println("end before run")
			return nil
		},
		After: func(c *cli.Context) error {
			return nil
		},
		Action: func(c *cli.Context) error {
			fmt.Println("start run")
			for {
				//return fmt.Errorf("testining config exit")
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
					processingMode := processingValue(c, cfg)
					bashGen := bashGenValue(c, cfg)

					ja, err := job.New(sources, nil,
						job.WithDirectProcessing(processingMode),
						job.WithBashGeneration(bashGen),
						job.WithBashDestination(cfg.BUFFER_DIR),
						job.WithBashTranslationMap(cfg.BashPathTranslation),
						job.WithInputDir(cfg.BUFFER_DIR),
						job.WithProcessingDir(cfg.IN_PROGRESS_DIR),
						job.WithDoneDir(cfg.DONE_DIR),
						job.WithOutDir(cfg.OUT_DIR),
					)
					if err != nil {
						fmt.Println("LOG ERROR: job creation", err.Error())
						continue
					}
					if err := ja.DecideType(); err != nil {
						fmt.Println("LOG ERROR: job decide type", err.Error())
						continue
						return err
					}
					if err := ja.CompileTasks(); err != nil {
						fmt.Println("LOG ERROR: job compile tasks", err.Error())
						continue
						return err
					}
					if err := ja.Execute(); err != nil {
						fmt.Println("LOG ERROR: job execute", err.Error())
						//return err
					}

				}
				fmt.Println("entering dormant mode:")
				for left := 15; left > 0; left = left - 1 {
					fmt.Printf("                                          \rwake up in %v seconds          \r", left)
					time.Sleep(time.Second)
				}
				fmt.Println("leave dormant mode:")
				fmt.Println("GOOD ENDING!!!")
			}

			return nil
		},
		OnUsageError: func(cCtx *cli.Context, err error, isSubcommand bool) error {
			return nil
		},
		Subcommands: []*cli.Command{},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:               "direct_processing",
				Category:           "Processing Mode",
				FilePath:           "",
				Usage:              "Process files directly with ffmpeg. ",
				Aliases:            []string{"dp"},
				DisableDefaultText: false,
			},
			&cli.BoolFlag{
				Name:        "bash_generation",
				Category:    "Processing Mode",
				DefaultText: "",
				FilePath:    "",
				Value:       true,
				Usage:       "Generate bash script for processing.",
				Aliases:     []string{"bg"},
			},
		},
		SkipFlagParsing:        false,
		HideHelp:               false,
		HideHelpCommand:        false,
		Hidden:                 false,
		UseShortOptionHandling: false,
		HelpName:               "",
		CustomHelpTemplate:     "",
	}

}

func processingValue(c *cli.Context, cfg *config.Configuration) bool {
	if c.Bool("dp") {
		return true
	}
	return cfg.DirectProcessing
}

func bashGenValue(c *cli.Context, cfg *config.Configuration) bool {
	if c.Bool("bg") {
		return true
	}
	return cfg.BashGeneration
}
