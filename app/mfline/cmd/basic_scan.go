package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/Galdoba/ffstuff/app/mfline/config"
	"github.com/Galdoba/ffstuff/app/mfline/internal/database"
	"github.com/Galdoba/ffstuff/app/mfline/internal/scan"
	"github.com/Galdoba/ffstuff/pkg/ump"
	"github.com/urfave/cli/v2"
)

func BasicScan() *cli.Command {
	return &cli.Command{
		Name: "basic",
		//Aliases:     []string{"fs"},
		Usage:       "make initial scan",
		UsageText:   "mfline basic",
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
			if len(c.Args().Slice()) == 0 {
				return fmt.Errorf("expect args")
			}

			args := c.Args().Slice()
			sort.Strings(args)

			for _, arg := range args {
				f, err := os.Open(arg)
				if err != nil {
					fmt.Errorf("os.Open: %v\n  %v\n", arg, err.Error())
					continue
				}
				defer f.Close()
				entry, _ := getEntry(db, arg)
				err = scanBasic(db, entry, arg)
				if err != nil {
					fmt.Println(err.Error())
				}

			}

			// data, err := db.Read(args[0])
			// fmt.Println(data)
			// fmt.Println(data.File)
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

func scanBasic(db *database.DBjson, entry *database.Entry, arg string) error {
	switch entry {
	case nil:
		if err := db.Create(arg); err != nil {
			return fmt.Errorf("db.Create: %v\n  %v\n", arg, err.Error())
		}
		entry, err := db.Read(arg)
		if err != nil {
			return fmt.Errorf("db.Read: %v\n  %v\n", arg, err.Error())
		}
		scanResult := scan.Basic(arg)

		entry.Profile = scanResult.Profile
		entry.Profile.BasicStatus = ump.Status_OK
		if err := db.Update(entry); err != nil {
			return fmt.Errorf("can't update %v", entry.File)
		}
	default:
		switch entry.Profile.BasicStatus {
		case ump.Status_OK:
			return scan.ErrNoScanNeeded
		default:
			scanResult := scan.Basic(arg)
			entry.Profile = scanResult.Profile
			entry.Profile.BasicStatus = ump.Status_OK
		}
		if err := db.Update(entry); err != nil {
			return fmt.Errorf("can't update %v", entry.File)
		}
	}
	return nil
}

func nameScan(db *database.DBjson, entry *database.Entry, arg string) error {
	if hasAnySuffix(arg, nonMediaSuffixes()...) {
		return scan.ErrNoScanNeeded
	}
	return nil
}

func nonMediaSuffixes() []string {
	return []string{
		".sh", ".bat", ".txt", ".jpeg",
	}
}

func hasAnySuffix(str string, suffixes ...string) bool {
	for _, suff := range suffixes {
		if strings.HasSuffix(str, suff) {
			return true
		}
	}
	return false
}
