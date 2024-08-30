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
			fmt.Println("initiation:")
			cfgLoaded, err := config.Load()
			if err != nil {
				return fmt.Errorf("config loading failed: %v", err)
			}
			cfg = cfgLoaded
			fmt.Println("config loaded")
			defer fmt.Println("initiation completed")
			return nil
		},
		After: func(c *cli.Context) error {
			return nil
		},
		Action: func(c *cli.Context) error {
			fmt.Println("start: aue run")
			for {
				//return fmt.Errorf("testining config exit")
				fmt.Println("read dir:", cfg.IN_DIR)
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
				if len(projects) == 0 {
					fmt.Println("no projects detected")
				}

				for _, project := range projects {

					fmt.Println("\n--------\nStart Project:", project)
					sources, err := actions.SetupSources(project, cfg.BUFFER_DIR, cfg.AssetFiles[config.Asset_File_Serial_data])
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
						job.WithNotificationDir(cfg.NotificationDir),
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
				for left := cfg.SleepSeconds; left > 0; left = left - 1 {
					fmt.Printf("                                          \rwake up in %v seconds          \r", left)
					time.Sleep(time.Second)
				}
				fmt.Println("leave dormant mode             ")

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
				Name:        "no_bash",
				Category:    "Processing Mode",
				DefaultText: "",
				FilePath:    "",
				Value:       false,
				Usage:       "Skip bash script generation.",
				Aliases:     []string{"nb"},
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
	if c.Bool("nb") {
		return false
	}
	return cfg.BashGeneration
}
