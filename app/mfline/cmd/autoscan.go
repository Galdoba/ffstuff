package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Galdoba/ffstuff/app/mfline/config"
	"github.com/Galdoba/ffstuff/app/mfline/internal/files"
	"github.com/Galdoba/ffstuff/app/mfline/ump"
	"github.com/urfave/cli/v2"
)

func FullScan() *cli.Command {
	return &cli.Command{
		Name:        "fullscan",
		Aliases:     []string{"fs"},
		Usage:       "make all possible scans for all files in tracked directory",
		UsageText:   "mfline fullscan",
		Description: "",
		Args:        false,
		ArgsUsage:   "",
		Category:    "",
		BashComplete: func(*cli.Context) {
		},
		Before: func(c *cli.Context) error {
			return nil
		},
		After: func(c *cli.Context) error {
			return nil
		},
		Action: func(c *cli.Context) error {
			cfg, _ := config.Load(c.App.Name)
			stored := ump.MapStorage(cfg.StorageDir)
			fileList := []string{}
			for _, trcDir := range cfg.TrackDirs {
				fileList = append(fileList, files.ListDir(trcDir)...)
			}
			for _, fl := range fileList {
				fmt.Fprintf(os.Stderr, "file: %v\n", fl)
				mp := ump.NewProfile()
				dataChanged := false
				name := filepath.Base(fl)
				if instore, ok := stored[name]; ok {
					mp = instore
				}
				err := mp.ScanBasic(fl)
				switch err {
				case nil:
					dataChanged = true
				default:
					fmt.Fprintf(os.Stderr, "scan basic: %v\n", err.Error())
				}
				err = mp.ScanInterlace(fl)
				switch err {
				case nil:
					dataChanged = true
				default:
					fmt.Fprintf(os.Stderr, "scan interlace: %v\n", err.Error())
				}
				if !dataChanged {
					continue
				}

				if err := mp.SaveAs(cfg.StorageDir + filepath.Base(fl) + ".json"); err != nil {
					fmt.Println(err.Error())
				}
			}
			// //список файлов для работы
			//по списку запускаем basic

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
