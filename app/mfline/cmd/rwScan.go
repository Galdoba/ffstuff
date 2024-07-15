package cmd

import (
	"fmt"
	"sort"

	"github.com/Galdoba/ffstuff/app/mfline/internal/database"
	"github.com/Galdoba/ffstuff/app/mfline/internal/scan"
	"github.com/Galdoba/ffstuff/app/mfline/ump"
	"github.com/urfave/cli/v2"
)

func RWCheck() *cli.Command {
	return &cli.Command{
		Name:        "rwcheck",
		Aliases:     []string{"rw"},
		Usage:       "TODO USAGE",
		UsageText:   "TODO Usage text",
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
			db, err := database.New(`c:\Users\pemaltynov\.ffstuff\data\mfline\`)
			if err != nil {
				return err
			}

			args := c.Args().Slice()
			sort.Strings(args)
			for _, arg := range args {
				entry, err := getEntry(db, arg)
				if err != nil {
					fmt.Printf("entry %v:\n  %v\n", arg, err.Error())
					continue
				}
				err = scanRW(db, entry, arg)
				if err != nil {
					fmt.Println(err.Error())
					continue
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

func scanRW(db *database.DBjson, entry *database.Entry, arg string) error {
	switch entry.Profile.RDWR_Status {
	case ump.Status_OK:
		return scan.ErrNoScanNeeded
	case ump.Status_Error:
		fmt.Printf("rwcheck %v: %v\n", entry.Profile.RDWR_Status, arg)
		return scan.ErrNoScanNeeded
	default:
		fmt.Printf("rwcheck: %v\n", arg)
		err := scan.ReadWrite(arg)
		switch err {
		case nil:
			entry.Profile.RDWR_Status = ump.Status_OK
		case scan.ErrRWCheck:
			entry.Profile.RDWR_Status = ump.Status_Error
		}
		if err := db.Update(entry); err != nil {
			return fmt.Errorf("db.Update: %v\n")
		}
	}
	return nil
}

/*
<(\ |А|а|Б|б|В|в|Г|г|Д|д|Е|е|Ё|ё|Ж|ж|З|з|И|и|Й|й|К|к|Л|л|М|м|Н|н|О|о|П|п|Р|р|С|с|Т|т|У|у|Ф|ф|Х|х|Ц|ц|Ч|ч|Ш|ш|Щ|щ|Ъ|ъ|Ы|ы|Ь|ь|Э|э|Ю|ю|Я|я)

<(П|п|Р|р|С|с|Т|т|У|у|Ф|ф|Х|х|Ц|ц|Ч|ч|Ш|ш|Щ|щ|Ъ|ъ|Ы|ы|Ь|ь|Э|э|Ю|ю|Я|я)

*/
