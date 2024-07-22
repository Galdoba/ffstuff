package cmd

import (
	"fmt"
	"time"

	"github.com/Galdoba/devtools/helpers"
	"github.com/Galdoba/ffstuff/app/mfline/config"
	"github.com/Galdoba/ffstuff/app/mfline/internal/database"
	"github.com/Galdoba/ffstuff/app/mfline/internal/files"
	"github.com/Galdoba/ffstuff/app/mfline/internal/scan"
	"github.com/urfave/cli/v2"
)

func AutoScan() *cli.Command {
	return &cli.Command{
		Name: "autoscan",
		//Aliases:     []string{"fs"},
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
			db, err := database.New(cfg.StorageDir)
			if err != nil {
				return err
			}
			var loop time.Duration = time.Second * 15
			dir := cfg.TrackDirs[0]
			for {
				fi := files.ListDir(dir)
				helpers.ClearTerminal()
				fmt.Printf("Tracking %v files:\n", len(fi))
				for _, file := range fi {
					path := file
					if files.BadName(path) {
						err := files.MarkAsBad(path)
						if err != nil {
							fmt.Println(err.Error())
						}
						continue
					}
					// f, err := os.Open(path)
					// if err != nil {
					// 	fmt.Printf("os.Open: %v\n  %v\n", path, err.Error())
					// 	continue
					// }
					// f.Close()
					entry, err := getEntry(db, path)
					if err != nil {
						switch err {
						case database.ErrNotFound:
							fmt.Println("new entry", path)
						default:
							fmt.Println(err.Error())
							continue
						}

					}
					if err = scanBasic(db, entry, path); err != nil {
						switch err {
						case scan.ErrNoScanNeeded:
						default:
							fmt.Println(err.Error())
							continue
						}
					}
					entry, err = db.Read(path)
					if err != nil {
						fmt.Println(err.Error())
						continue
					}
					if err = scanRW(db, entry, path); err != nil {
						switch err {
						case scan.ErrNoScanNeeded:
						default:
							fmt.Println(err.Error())
							continue
						}
					}
					entry, err = db.Read(path)
					if err != nil {
						fmt.Println(err.Error())
						continue
					}
					if err = scanInterlace(db, entry, path); err != nil {
						switch err {
						case scan.ErrNoScanNeeded:
						default:
							fmt.Println(err.Error())
							continue
						}
					}

					// f.Close()
				}

				time.Sleep(loop)

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
