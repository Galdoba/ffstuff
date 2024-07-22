package cmd

import (
	"fmt"
	"sort"

	"github.com/Galdoba/ffstuff/app/mfline/internal/database"
	"github.com/Galdoba/ffstuff/app/mfline/internal/files"
	"github.com/Galdoba/ffstuff/app/mfline/internal/scan"
	"github.com/urfave/cli/v2"
)

func InterlaceCheck() *cli.Command {
	return &cli.Command{
		Name:        "interlacecheck",
		Aliases:     []string{"intr"},
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
				err = scanInterlace(db, entry, arg)
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

func scanInterlace(db *database.DBjson, entry *database.Entry, arg string) error {
	videoStreams := 0
	for _, stream := range entry.Profile.Streams {
		switch stream.Codec_type {
		case "video":
			videoStreams++
		}
	}
	switch videoStreams {
	case 0:
		return scan.ErrNoScanNeeded
	case 1:
	default:
		return fmt.Errorf("multiple video streams not implemented")
	}
	switch entry.Profile.Interlace_Status {
	case "":
		fmt.Printf("interlace check: %v\n", arg)
		markedArg, err := files.MarkScan(arg, files.Interlace_Marker)
		if err != nil {
			return fmt.Errorf("mark file: %v", err.Error())
		}
		report, procErr := scan.Interlace(markedArg)
		switch procErr {
		case nil:
			entry.Profile.Interlace_Status = report
		default:
			fmt.Printf("error: '%v'\n", err.Error())
			panic(1)
		}
		if err := db.Update(entry); err != nil {
			return fmt.Errorf("db.Update: %v\n")
		}
		if err := files.ClearMarkers(markedArg); err != nil {
			fmt.Println("unmark file:", err.Error())
		}
		return nil
	default:
		return scan.ErrNoScanNeeded
	}
	return nil
}

/*
<(\ |А|а|Б|б|В|в|Г|г|Д|д|Е|е|Ё|ё|Ж|ж|З|з|И|и|Й|й|К|к|Л|л|М|м|Н|н|О|о|П|п|Р|р|С|с|Т|т|У|у|Ф|ф|Х|х|Ц|ц|Ч|ч|Ш|ш|Щ|щ|Ъ|ъ|Ы|ы|Ь|ь|Э|э|Ю|ю|Я|я)

<(П|п|Р|р|С|с|Т|т|У|у|Ф|ф|Х|х|Ц|ц|Ч|ч|Ш|ш|Щ|щ|Ъ|ъ|Ы|ы|Ь|ь|Э|э|Ю|ю|Я|я)

*/
