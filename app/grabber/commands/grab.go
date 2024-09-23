package commands

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/Galdoba/ffstuff/app/grabber/config"
	"github.com/Galdoba/ffstuff/app/grabber/internal/actions"
	"github.com/Galdoba/ffstuff/app/grabber/internal/validation"
	"github.com/Galdoba/ffstuff/pkg/logman"
	"github.com/urfave/cli/v2"
)

var dest string

var cfg *config.Configuration

func Grab() *cli.Command {
	return &cli.Command{
		Name:        "grab",
		Aliases:     []string{},
		Usage:       "TODO: Direct command for transfering operation(s)",
		UsageText:   "",
		Description: "",
		Args:        false,
		ArgsUsage:   "",
		Category:    "",
		Before: func(*cli.Context) error {
			cfgLoaded, err := config.Load()
			if err != nil {
				return err
			}
			errs := config.Validate(cfgLoaded)
			switch len(errs) {
			default:
				for _, err := range errs {
					if err != nil {
						fmt.Println(err)
					}
				}
				return fmt.Errorf("config errors detected")
			case 0:
				cfg = cfgLoaded
				setupLogger(cfg.LOG_LEVEL, cfg.LOG)

			}

			return nil
		},
		Action: func(c *cli.Context) error {
			logman.Debug(logman.NewMessage("start 'grub'"))

			if err := validateDestination(c); err != nil {
				return err
			}
			args, err := validateArguments(c)
			if err != nil {
				return err
			}
			logman.Printf("grab to %v\n", dest)
			sorted := sortOrder(args...)
			fmt.Println("sorting...")
			for _, arg := range sorted {
				actions.CopyFile(arg, dest, true)
			}

			logman.Debug(logman.NewMessage("end   'grub'"))
			return nil
		},

		Subcommands: []*cli.Command{},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "dest",
				Usage:       "destination where files will be downloaded to\n",
				DefaultText: "from config",
				Required:    false,
				Hidden:      false,
				HasBeenSet:  false,
				Value:       "",
				Destination: new(string),
				Aliases:     []string{"d"},
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

func validateDestination(c *cli.Context) error {
	switch c.String("dest") {
	case "":
		dest = cfg.DEFAULT_DESTINATION
	default:
		dest = c.String("dest")
	}
	if err := validation.DirectoryValidation(dest); err != nil {
		return fmt.Errorf("destination directory validation failed: %v")
	}
	sep := string(filepath.Separator)
	dest = strings.TrimSuffix(dest, sep)
	dest += sep
	return nil
}

func validateArguments(c *cli.Context) ([]string, error) {
	args := c.Args().Slice()
	if len(args) == 0 {
		return args, logman.Errorf("no argumets provided: grabber grab expect source files as arguments")
	}
	return args, nil
}

func sortOrder(args ...string) []string {
	sortedByScore := make(map[int][]string)
	scoreMax := -1000000
	scoreMin := 1000000
	for _, arg := range args {
		upArg := strings.ToUpper(arg)
		scoreTotal := 0
		for key, score := range cfg.PRIORITY_MAP {
			upKey := strings.ToUpper(key)
			if strings.Contains(upArg, upKey) {
				scoreTotal += score
			}
		}
		sortedByScore[scoreTotal] = append(sortedByScore[scoreTotal], arg)
		if scoreTotal > scoreMax {
			scoreMax = scoreTotal
		}
		if scoreTotal < scoreMin {
			scoreMin = scoreTotal
		}
	}
	sorted := []string{}
	for i := scoreMax; i >= scoreMin; i-- {
		if list, ok := sortedByScore[i]; ok {
			for _, path := range list {
				sorted = append(sorted, path)
			}
		}
	}
	return sorted
}
