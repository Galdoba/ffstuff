package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/Galdoba/ffstuff/pkg/ump"
	"github.com/urfave/cli/v2"
)

func AudioStreamsData() *cli.Command {
	return &cli.Command{
		Name:        "audio_streams_data",
		Aliases:     []string{"asd"},
		Usage:       "print audio channel layout and number of channels",
		UsageText:   "USAGE TEXT HERE",
		Description: "DESCRIPTION HERE",
		Category:    "INFO",

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
			data := []string{}
			for _, audio := range mp.Streams {
				if audio.Codec_type != "audio" {
					continue
				}
				data = append(data, fmt.Sprintf("%v", audio.Channels))
				data = append(data, fmt.Sprintf("%v", audio.Channel_layout))
			}

			fmt.Fprint(os.Stdout, strings.Join(data, "\n"))
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
