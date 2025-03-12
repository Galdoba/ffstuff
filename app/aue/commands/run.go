package commands

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Galdoba/ffstuff/app/aue/config"
	"github.com/Galdoba/ffstuff/app/aue/internal/actions"
	"github.com/Galdoba/ffstuff/app/aue/internal/files/sourcefile"
	"github.com/Galdoba/ffstuff/app/aue/internal/job"
	"github.com/Galdoba/ffstuff/pkg/logman"
	log "github.com/Galdoba/ffstuff/pkg/logman"
	"github.com/urfave/cli/v2"
)

var cfg *config.Configuration
var sessionID = timestamplogFormat(time.Now())

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
				log.WithAppLogLevelImportance(log.ImportanceTRACE),
			)
			// log.ClearOutput(log.DEBUG)
			// logPathDefault := cfg.AssetFiles[config.Asset_File_Log]
			// logPathSession := strings.ReplaceAll(logPathDefault, "default.log", "session_"+sessionID+".log")
			// f, err := os.Create(logPathSession)
			// if err != nil {
			// 	fmt.Println(err)
			// }
			// f.Close()
			// log.SetOutput(logPathDefault, log.ALL)
			// log.SetOutput(logPathSession, log.ALL)
			// if err != nil {
			// 	return fmt.Errorf("logger setup failed: %v", err)
			// }
			if err := setupLogger(cfg); err != nil {
				return err
			}
			fmt.Println("logger started")
			defer fmt.Println("initiation completed")
			return nil
		},
		After: func(c *cli.Context) error {
			return nil
		},
		Action: func(c *cli.Context) error {
			log.Info("aue start session: %v", sessionID)
			in_dir := cfg.IN_DIR
			failCounter := 0
			for {
				haltIfLocked(in_dir)
				//return fmt.Errorf("testining config exit")
				log.Debug(log.NewMessage(fmt.Sprintf("scan root directory")))
				rootFiles, projects, err := actions.ScanDir(in_dir)
				if err != nil {
					log.Errorf(fmt.Sprintf("directory '%v' reading failed: %v", in_dir, err.Error()))
					failCounter++
					if failCounter%5 == 0 {
						fmt.Println("TODO: Alert!!")
						return err
					}
					continue

				}
				if len(rootFiles) == 0 {
					log.Errorf("no metadata files detected")
				}

				if len(projects) == 0 {
					log.Debug(log.NewMessage("no projects detected"))
				}
				in_dir = strings.TrimSuffix(in_dir, "/") + "/"
				for _, project := range projects {

					proj_dir := in_dir + project + "/"

					log.Debug(log.NewMessage("start project: %v", project))
					sources := []*sourcefile.SourceFile{}
					//Project Setup
					selectedAction, err := actions.SelectSourceAction(proj_dir)
					if err != nil {
						log.Error(err)
						continue
					}
					switch selectedAction {
					default:
						log.Error(fmt.Errorf("unimplemented action selected: '%v'", selectedAction))
						continue
					case actions.ActionRemoveProject:
						if err := actions.RemoveProject(proj_dir); err != nil {
							log.Error(err)
						}
						log.Info("end project: %v (clean)", project)
						continue
					case actions.ActionSkip:
						log.Debug(log.NewMessage("end project: %v (skip)", project))
						continue
					case actions.ActionSetup:
						sourcesCreated, err := actions.SetupSources(proj_dir, cfg.BUFFER_DIR, cfg.AssetFiles[config.Asset_File_Serial_data])
						if err != nil {
							log.Warn("end project: %v (failed)", project)
							continue
						}
						if len(sourcesCreated) == 0 {
							log.Debug(log.NewMessage("project %v: no sources created", project))
							log.Warn("end project: %v (skip)", project)
							continue
						}
						sources = sourcesCreated
						f, err := os.Create(proj_dir + "lock")
						if err != nil {
							log.Warn("failed to lock %v", project)
						}
						f.Close()
						log.Debug(log.NewMessage("%v locked", project))
					}
					//Project Execution
					ja, err := job.New(sources, nil,
						job.WithDirectProcessing(processingValue(c, cfg)),
						job.WithBashGeneration(bashGenValue(c, cfg)),
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

					log.Info("job creation complete: %v", ja.ProjectName())
					if err := ja.DecideType(); err != nil {
						log.Warn("job type decidion failed")
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
					log.Info("end project: %v (success)", project)

				}

				for left := cfg.SleepSeconds; left > 0; left = left - 1 {
					fmt.Fprintf(os.Stderr, "dormant mode for %v    \r", timer(left))
					time.Sleep(time.Second)
				}

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

func timer(seconds int) string {
	h, m, s := seconds/3600, seconds/60, seconds%60
	return fmt.Sprintf("%v:%v:%v", numToStr(h), numToStr(m), numToStr(s))
}

func numToStr(n int) string {
	s := fmt.Sprintf("%v", n)
	for len(s) < 2 {
		s = "0" + s
	}
	return s
}

func timestamplogFormat(tm time.Time) string {
	format := "060102150405"
	formatLen := len(format)
	stamp := tm.Format(format)
	for len(stamp) < formatLen {
		stamp += "0"
	}
	return stamp
}

func haltIfLocked(in_dir string) {
mainLoop:
	for {
		fi, err := os.ReadDir(in_dir)
		if err != nil {
			logman.Warn("failed to check lock: %v")
			time.Sleep(time.Second * 5)
			continue
		}
		for _, f := range fi {
			if f.Name() == "lock" {
				fmt.Printf("aue is locked                               \r")
				time.Sleep(time.Second)
				continue mainLoop
			}
		}
		return
	}
}
