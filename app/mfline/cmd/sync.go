package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Galdoba/ffstuff/app/mfline/config"
	"github.com/Galdoba/ffstuff/app/mfline/ump"
	"github.com/urfave/cli/v2"
)

func Sync() *cli.Command {

	return &cli.Command{
		Name:        "sync",
		Aliases:     []string{},
		Usage:       "sync scan data files is original files get prefix (auto feature)",
		UsageText:   "",
		Description: "",
		Args:        false,
		ArgsUsage:   "",
		Category:    "",
		BashComplete: func(*cli.Context) {
		},
		Before: func(*cli.Context) error {
			return nil
		},
		After: func(*cli.Context) error {
			return nil
		},
		Action: func(c *cli.Context) error {
			cfg, _ := config.Load(c.App.Name)
			stored := ump.MapStorage(cfg.StorageDir)
			// SYNC NAMES WITH NEW PREFIXES
			for _, trackedDir := range cfg.TrackDirs {
				fls, err := os.ReadDir(trackedDir)
				if err != nil {
					return nil
				}
				for _, fl := range fls {
					if fl.IsDir() {
						continue
					}
					for k, v := range stored {
						if fl.Name() != k && strings.Contains(fl.Name(), k) {
							v.Format.Filename = fl.Name()
							v.SaveAs(cfg.StorageDir+fl.Name()+".json", false)
							os.Remove(cfg.StorageDir + k + ".json")
						}
					}
				}
			}
			//DELETE OLD FILES
			fls, err := os.ReadDir(cfg.StorageDir)
			if err != nil {
				return nil
			}
			for _, fl := range fls {
				if fl.IsDir() {
					continue
				}
				fi, err := fl.Info()
				if err != nil {
					return err
				}
				lastMod := fi.ModTime()
				age := time.Since(lastMod).Minutes()
				if age > cfg.OldScan {
					fmt.Fprintf(os.Stderr, "scandata %v is %v hours old: deleted\n", fl.Name(), int(age))
					os.Remove(cfg.StorageDir + fl.Name())
				}
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
