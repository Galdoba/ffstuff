package commands

import (
	"fmt"
	"os"
	"time"

	"github.com/Galdoba/ffstuff/app/aue/config"
	"github.com/Galdoba/ffstuff/app/aue/internal/actions"
	"github.com/Galdoba/ffstuff/app/aue/internal/job"
	log "github.com/Galdoba/ffstuff/pkg/logman"
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
			err = log.Setup(
				log.WithAppLogLevelImportance(log.LvlALL),
			)
			log.SetOutput(cfg.AssetFiles[config.Asset_File_Log], log.ALL)
			if err != nil {
				return fmt.Errorf("logger setup failed: %v", err)
			}
			fmt.Println("logger started")
			defer fmt.Println("initiation completed")
			return nil
		},
		After: func(c *cli.Context) error {
			return nil
		},
		Action: func(c *cli.Context) error {
			log.Info("start: aue run")
			in_dir := cfg.IN_DIR
			for {
				//return fmt.Errorf("testining config exit")
				log.Debug(fmt.Sprintf("read directory:", in_dir), time.Now())
				fi, err := os.ReadDir(in_dir)
				if err != nil {
					log.Fatal(fmt.Sprintf("directory '%v' reading failed: %v", in_dir, err.Error()))
					return err
				}
				projects := []string{}
				for _, f := range fi {
					if f.IsDir() {
						projects = append(projects, fmt.Sprintf("%v%v", cfg.IN_DIR, f.Name()))
					}
				}
				if len(projects) == 0 {
					log.Info("no projects detected")
				}

				for _, project := range projects {
					log.Info("start project: %v", project)
					sources, err := actions.SetupSources(project, cfg.BUFFER_DIR, cfg.AssetFiles[config.Asset_File_Serial_data])
					if len(sources) == 0 {

						log.Warn("project %v: no sources created", project)
						continue
					}
					if err != nil {
						log.Error(fmt.Errorf("project %v: source setup: %v", project, err))
						continue
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
						log.Error(err)
						continue
					}

					log.Info("job creation complete: %v", project)
					if err := ja.DecideType(); err != nil {
						log.Error(err)
						continue
						//return err
					}
					log.Info("job type decided: %v", ja.TypeDecided())
					if err := ja.CompileTasks(); err != nil {
						log.Error(err)
						continue
						return err
					}
					if err := ja.Execute(); err != nil {

						log.Error(err)
						continue
						//return err
					}
					log.Info("job execution completed")

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
