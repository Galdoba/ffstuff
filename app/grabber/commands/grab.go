package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Galdoba/ffstuff/app/grabber/config"
	"github.com/Galdoba/ffstuff/app/grabber/internal/actions"
	"github.com/Galdoba/ffstuff/app/grabber/internal/copyprocess"
	"github.com/Galdoba/ffstuff/app/grabber/internal/validation"
	"github.com/Galdoba/ffstuff/pkg/logman"
	"github.com/urfave/cli/v2"
)

var dest string
var sortMethod string

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
			sources, err := validateArguments(c)
			if err != nil {
				return err
			}

			//create new action
			copyAction := copyprocess.NewCopyAction(
				copyprocess.WithSourcePaths(sources...),
				copyprocess.WithDestination(dest),
			)

			//start action
			copyAction.Start()
			// logman.Printf("grab to %v\n", dest)
			// sorted := sortOrder(args...)
			// fmt.Println("sorting...")
			// for _, arg := range sorted {
			// 	err := actions.CopyFile(arg, dest)
			// 	if err != nil {
			// 		logman.Error(err)
			// 	}
			// 	handleOrigins(arg)
			// }

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
			&cli.BoolFlag{
				Name:    "size_sort",
				Usage:   "use sort by size method (overwrite config)",
				Aliases: []string{"ss"},
			},
			&cli.BoolFlag{
				Name:    "priority_sort",
				Usage:   "use sort by priority method (overwrite config)",
				Aliases: []string{"ps"},
			},
			&cli.BoolFlag{
				Name:    "no_sort",
				Usage:   "use no sort method (overwrite config)",
				Aliases: []string{"ns"},
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

func validateSortMethod(c *cli.Context) error {
	sortMethod = cfg.SORT_METHOD
	sum := 0
	for _, val := range []bool{c.Bool("ss"), c.Bool("ps"), c.Bool("ns")} {
		if val {
			sum++
		}
	}
	if sum > 1 {
		fmt.Errorf("flags -ss (%v), -ps (%v) and -ns (%v) are mutualy exclusive", c.Bool("ss"), c.Bool("ps"), c.Bool("ns"))
	}
	if c.Bool("ss") {
		sortMethod = config.SORT_BY_SIZE
	}
	if c.Bool("ps") {
		sortMethod = config.SORT_BY_PRIORITY
	}
	if c.Bool("ns") {
		sortMethod = config.SORT_BY_NONE
	}
	return nil
}

func validateArguments(c *cli.Context) ([]string, error) {
	if err := validateDestination(c); err != nil {
		return nil, logman.Errorf("arguments validation: %v", err)
	}

	args := c.Args().Slice()
	if len(args) == 0 {
		logman.Errorf("no arguments provided: 'grab' command expects source files as arguments")
		return nil, ErrBadArguments
	}
	sourcePaths := []string{}
	for _, arg := range args {
		if strings.HasSuffix(arg, cfg.MARKER_FILE_EXTENTION) {
			related, err := actions.DiscoverRelatedFiles(arg)
			if err != nil {
				logman.Error(err)
			}
			for _, file := range related {
				sourcePaths = append(sourcePaths, file)
			}
			sourcePaths = append(sourcePaths, arg)
		}
	}

	return sourcePaths, nil
}

// func sortOrder(c *cli.Context, sources ...string) []string {
// 	sortMethod := "NONE"
// 	if c.Bool("sort_by_size") {
// 		sortMethod = "SIZE"
// 	}

// 	return sorted
// }

func fileType(path string) string {
	if strings.HasSuffix(path, cfg.MARKER_FILE_EXTENTION) {
		return "markerFile"
	}
	return "sourceFile"
}

func handleOrigins(arg string) {
	fType := fileType(arg)
	switch fType {
	case "markerFile":
		switch cfg.DELETE_ORIGINAL_MARKER {
		case true:
			logman.Debug(nil, fmt.Sprintf("delete marker file: %v", arg))
			if errRM := os.Remove(arg); errRM != nil {
				logman.Warn("failed to remove marker file: %v", errRM.Error())
			}
		case false:
			logman.Debug(logman.NewMessage("keeping marker file: %v", arg))
		}
	case "sourceFile":
		switch cfg.DELETE_ORIGINAL_SOURCE {
		case true:
			logman.Debug(nil, fmt.Sprintf("delete source file: %v", arg))
			if errRM := os.Remove(arg); errRM != nil {
				logman.Warn("failed to remove source file: %v", errRM.Error())
			}
		case false:
			logman.Debug(logman.NewMessage("keeping source file: %v", arg))
		}
	}
}
