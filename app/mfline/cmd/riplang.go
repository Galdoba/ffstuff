package cmd

import (
	"fmt"
	"os"

	"github.com/Galdoba/ffstuff/pkg/ump"
	"github.com/urfave/cli/v2"
)

func RipLang() *cli.Command {
	return &cli.Command{
		Name:        "riplang",
		Aliases:     []string{"rl"},
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
			//CHECK SOURCE
			sourceFile := c.String("source")
			if err := checkSource(c.String("source")); err != nil {
				return err
			}
			//CHECK DESTINATION

			//COMENCE BASIC SCAN
			mp := ump.NewProfile()
			err := mp.ScanBasic(sourceFile)
			if err != nil {
				return err
			}
			//SAVE TO TARGET FILE
			out := ""

			/*
							-map 0:v ^
				    -map 1:a -metadata:s:a:0 language=rus ^
				    -map 2:a -metadata:s:a:1 language=eng ^
				    -map 3:s -metadata:s:s:0 language=rus ^
			*/
			for i, stream := range mp.Streams {
				out += fmt.Sprintf(" -map 0:%v -metadata:s:%v language=%v", i, i, stream.Tags["language"])

				//fmt.Println(i, stream.Tags["language"])
			}
			fmt.Fprint(os.Stdout, out)
			return nil
		},
		OnUsageError: func(cCtx *cli.Context, err error, isSubcommand bool) error {
			return nil
		},
		Subcommands: []*cli.Command{},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name: "source",
				//Category:  "path",
				Usage:     "path to file which will be scanned   (required)",
				Required:  true,
				Aliases:   []string{},
				TakesFile: true,
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
