package commands

import (
	"fmt"

	"github.com/Galdoba/ffstuff/app/grabber/commands/grab"
	"github.com/Galdoba/ffstuff/app/grabber/config"
	"github.com/Galdoba/ffstuff/app/grabber/internal/origin"
	"github.com/Galdoba/ffstuff/pkg/logman"
	"github.com/urfave/cli/v2"
)

var cfg *config.Configuration

func Grab() *cli.Command {
	return &cli.Command{
		Name:        "grab",
		Aliases:     []string{},
		Usage:       "TODO: Direct command for transfering operation(s)",
		UsageText:   "grabber grab [command options] args...",
		Description: "Desc",
		Args:        false,
		ArgsUsage:   "Args Usage Text",
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
			//Setup process
			grabProcessSettings := grab.NewProcessControl(cfg)
			if err := grabProcessSettings.Modify(c); err != nil {
				return fmt.Errorf("operation process setup failed: %v", err)
			}
			if err := grabProcessSettings.Assert(); err != nil {
				return fmt.Errorf("operation setup invalid")
			}
			grabProcessSettings.Status()
			//Setup Sources
			sourcePaths, validationErrors := grab.ValidateArgs(c.Args().Slice()...)
			origin.ConstructorSetup(
				origin.WithFilePriority(grabProcessSettings.FileWeights),
				origin.WithDirectoryPriority(grabProcessSettings.DirWeights),
				origin.KillAll(grabProcessSettings.KillAll),
				origin.KillMarkers(grabProcessSettings.KillMarker),
				origin.WithMarkerExt(grabProcessSettings.MarkerExt),
			)
			fmt.Println(sourcePaths, validationErrors)

			// if err := validateDestination(c); err != nil {
			// 	return err
			// }
			// sources, err := validateArguments(c)
			// if err != nil {
			// 	return err
			// }

			//create new action
			// copyAction := copyprocess.NewCopyAction(
			// 	copyprocess.WithSourcePaths(sources...),
			// 	copyprocess.WithDestination(dest),
			// )

			// //start action
			// copyAction.Start()
			// // logman.Printf("grab to %v\n", dest)
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
				Category:    "PATH:",
				Usage:       "destination where files will be downloaded to\n",
				DefaultText: "",
				Required:    false,
				Hidden:      false,
				HasBeenSet:  false,
				Value:       "",
				Destination: new(string),
				Aliases:     []string{"d"},
			},
			&cli.BoolFlag{
				Name:               "size_sort",
				Category:           "SORTING:",
				Usage:              "sort order by size (supress config settings)",
				DisableDefaultText: true,
				Aliases:            []string{"ss"},
			},
			&cli.BoolFlag{
				Name:               "priority_sort",
				Category:           "SORTING:",
				Usage:              "sort by calculating priority scores (supress config settings)",
				DisableDefaultText: true,
				Aliases:            []string{"ps"},
			},
			&cli.BoolFlag{
				Name:               "no_sort",
				Category:           "SORTING:",
				Usage:              "do not sort (supress config settings)",
				DisableDefaultText: true,
				Aliases:            []string{"ns"},
			},

			&cli.BoolFlag{
				Name:     "copy_skip",
				Category: "COPY HANDLING",
				Usage:    "if copy exist file will not be grabbed",
				Aliases:  []string{"cs"},
			},
			&cli.BoolFlag{
				Name:     "copy_rename",
				Category: "COPY HANDLING",
				Usage:    "if copy exist file will be renamed",
				Aliases:  []string{"cr"},
			},
			&cli.BoolFlag{
				Name:     "copy_overwrite",
				Category: "COPY HANDLING",
				Usage:    "if copy exist file will be overwritten",
				Aliases:  []string{"co"},
			},

			&cli.BoolFlag{
				Name:     "delete_marker",
				Category: "DELETE ORIGINAL FILES",
				Usage:    "delete marker files after grabbing",
				Aliases:  []string{"dm"},
			},
			&cli.BoolFlag{
				Name:     "delete_all",
				Category: "DELETE ORIGINAL FILES",
				Usage:    "delete all original files after grabbing",
				Aliases:  []string{"da"},
			},
		},
	}
}

func setupSourceConstructor() error {
	if err := origin.ConstructorSetup(
		origin.WithFilePriority(cfg.FILE_PRIORITY_WEIGHTS),
		origin.WithDirectoryPriority(cfg.DIRECTORY_PRIORITY_WEIGHTS),
		origin.KillAll(cfg.DELETE_ORIGINAL_MARKER),
		origin.KillAll(cfg.DELETE_ORIGINAL_SOURCE),
		origin.WithMarkerExt(cfg.MARKER_FILE_EXTENTION),
	); err != nil {
		return logman.Errorf("source constructor setup failed: %v", err)
	}
	return nil
}
