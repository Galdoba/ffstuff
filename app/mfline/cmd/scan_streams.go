package cmd

import (
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"
)

func ScanStreams() *cli.Command {
	cm := &cli.Command{
		Name:      "scan",
		Usage:     "scan file ('mfline scan --help' for discription)",
		UsageText: "mfline scan subcommand --source PATH1 --destination PATH2 [flags]...\nno arguments expected",
		Description: strings.Join([]string{
			"scan file with designated target and write all data about 'source' file is stored written 'destination' file",
			"each stage other than basic REQUIRE basic stage to be completed",
			"scan stages are:",
			"  basic      - get common data from file header",
			"  interlace  - scan video streams with ffmpeg's 'idet' filter",
			"  silence    - scan audio streams for segments with loudness less than X Lufs (TODO)",
		}, "\n"),
		ArgsUsage: "args used this way",
		//Category:    "Scan",
		BashComplete: func(*cli.Context) {
		},
		Before: func(*cli.Context) error {
			return nil
		},
		After: func(*cli.Context) error {
			return nil
		},
		Action: func(c *cli.Context) error {
			return fmt.Errorf("bad usage\n  command 'scan' must be used in conjunction with subcommands:\n  enter 'mfline scan --help' for description")
		},
		OnUsageError: func(cCtx *cli.Context, err error, isSubcommand bool) error {
			return err
		},
		Subcommands: []*cli.Command{&cli.Command{
			Name:        "basic",
			Aliases:     []string{},
			Usage:       "get common data from file header",
			UsageText:   "mfline scan basic --source FILE [--destination DIR]\n\ndestination priority:\n -flag        (special case: 'local' = runtime dir)\n -config value\n -runtime dir\n\nno arguments expected\n\nmfline scan basic -? - for detailed args usage",
			Description: "DESCR",
			ArgsUsage:   "ARGS",
			Category:    "",
			Action: func(c *cli.Context) error {
				fmt.Println("sub1 command used")
				fmt.Println("global fl", c.String("test"))
				fmt.Println("sub fl", c.String("fl1"))
				fmt.Println(c.Command.HelpName)
				fmt.Println("sub fl", c.String("fl2"))
				if c.String("?") != "" {
					fmt.Println("Args usage:", c.Command.ArgsUsage)
				}
				return nil
			},
			OnUsageError: func(cCtx *cli.Context, err error, isSubcommand bool) error {
				fmt.Println("usage error")
				fmt.Println(err.Error())
				fmt.Println("sub =", isSubcommand)
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

				&cli.StringFlag{
					Name: "destination",
					//Category:  "path",
					Usage:     "path to file data will be written to           \n\nother:",
					Required:  false,
					TakesFile: true,
				},
				&cli.BoolFlag{
					Name:               "overwrite",
					Usage:              "rewrite destination file",
					DisableDefaultText: true,
					Aliases:            []string{"o"},
				},
				&cli.BoolFlag{
					Name:               "args",
					Usage:              "show detailed args usage",
					Aliases:            []string{"?"},
					DisableDefaultText: true,
				},
			},
			SkipFlagParsing:        false,
			HideHelp:               false,
			HideHelpCommand:        false,
			Hidden:                 false,
			UseShortOptionHandling: false,
			HelpName:               "",
		}},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "fl1",
				Category:    "",
				DefaultText: "",
				FilePath:    "",
				Usage:       "",
				Required:    false,
				Hidden:      false,
				HasBeenSet:  false,
				Value:       "",
				Destination: new(string),
				Aliases:     []string{},
				EnvVars:     []string{},
				TakesFile:   false,
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
	return cm
}
