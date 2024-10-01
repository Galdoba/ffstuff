package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Galdoba/ffstuff/app/grabber/commands/grabberargs"
	"github.com/Galdoba/ffstuff/app/grabber/commands/grabberflag"
	"github.com/Galdoba/ffstuff/app/grabber/config"
	"github.com/Galdoba/ffstuff/app/grabber/internal/actions"
	"github.com/Galdoba/ffstuff/app/grabber/internal/origin"
	"github.com/Galdoba/ffstuff/app/grabber/internal/process"
	"github.com/Galdoba/ffstuff/app/grabber/internal/sourcesort"
	"github.com/Galdoba/ffstuff/pkg/logman"
	"github.com/urfave/cli/v2"
)

var cfg *config.Configuration

func Grab() *cli.Command {
	return &cli.Command{
		Name:      "grab",
		Aliases:   []string{},
		Usage:     "direct command for transfering files",
		UsageText: "grabber grab [command options] args...",
		Description: "Setup single process to copy source files to destination directory. Command receives filepaths as arguments.\n" +
			"If argument is 'marker' file grabber will search all related files in the same directory and add them for transfering.",
		Args:      false,
		ArgsUsage: "Args Usage Text",
		Category:  "",
		Before: func(*cli.Context) error {
			return commandInit()
		},
		Action: func(c *cli.Context) error {

			logman.Debug(logman.NewMessage("start 'grub'"))
			//Setup process
			logman.Debug(logman.NewMessage("check flags"))
			if err := grabberflag.ValidateGrabFlags(c); err != nil {
				return logman.Errorf("flag validation failed: %v", err)
			}

			logman.Debug(logman.NewMessage("check arguments"))
			if err := grabberargs.ValidateGrabArguments(c.Args().Slice()...); err != nil {
				return logman.Errorf("argument validation failed: %v", err)
			}

			logman.Debug(logman.NewMessage("set process options"))
			options := process.DefineGrabOptions(c, cfg)
			process, err := process.New(options...)
			if err != nil {
				return logman.Errorf("process creation failed")
			}
			//Setup sources
			if err := origin.ConstructorSetup(
				origin.WithFilePriority(cfg.FILE_PRIORITY_WEIGHTS),
				origin.WithDirectoryPriority(cfg.DIRECTORY_PRIORITY_WEIGHTS),
				origin.KillSignal(process.DeleteDecidion),
				origin.WithMarkerExt(cfg.MARKER_FILE_EXTENTION),
			); err != nil {
				return logman.Errorf("source constructor setup failed: %v", err)
			}
			sources := []origin.Origin{}
			for grNum, arg := range c.Args().Slice() {
				gr := fmt.Sprintf("group_%02d", grNum)
				sources = append(sources, origin.New(arg, gr))
				related, err := actions.DiscoverRelatedFiles(arg)
				if err != nil {
					logman.Warn("failed to discover related files: %v", err)
				}
				for _, found := range related {
					sources = append(sources, origin.New(found, gr))
				}
			}
			logman.Printf("%v sorce files received", len(sources))

			//Sort
			switch process.SortDecidion {
			case grabberflag.VALUE_SORT_PRIORITY:
				sources = sourcesort.SortByPriority(process.KeepMarkerGroups, sources...)
			case grabberflag.VALUE_SORT_SIZE:
				sources = sourcesort.SortBySize(process.KeepMarkerGroups, sources...)
			case grabberflag.VALUE_SORT_NONE:
			}
			//grab
			dest := process.DestinationDir
			for i, src := range sources {
				name := filepath.Base(src.Path())
				logman.Debug(logman.NewMessage("emulating %v: %v to %v", i, name, dest))
				oldCopy, err := exists(dest + name)
				if err != nil {
					logman.Errorf("old copy check received unexpected error: %v", err)
				}
				if oldCopy {
					logman.Warn("old copy exist: renaming %v to %v", name)

				}
			}

			fmt.Println(process)
			process.ShowOrder()
			for i, src := range sources {
				fmt.Println(i, src)
			}

			logman.Debug(logman.NewMessage("end   'grub'"))
			return nil
		},
		Subcommands: []*cli.Command{},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        grabberflag.DESTINATION,
				Usage:       "sets where files will be downloaded to",
				DefaultText: "from config file",
				Required:    false,
				Hidden:      false,
				HasBeenSet:  false,
				Value:       "",
				Destination: new(string),
				Aliases:     []string{"dest"},
			},

			&cli.StringFlag{
				Name:    grabberflag.COPY,
				Usage:   "set decidion if target file with same name exists\n	  valid values: skip, remane or overwrite",
				Value:   grabberflag.VALUE_COPY_SKIP,
				Aliases: []string{"c"},
			},
			&cli.StringFlag{
				Name:    grabberflag.DELETE,
				Usage:   "set which original files will be deleted after grabbing\n	  valid values: none, marker or all",
				Value:   grabberflag.VALUE_DELETE_MARKER,
				Aliases: []string{"d"},
			},
			&cli.StringFlag{
				Name:    grabberflag.SORT,
				Usage:   "set method to decide grabbing order\n	  valid values: priority, size or none",
				Value:   grabberflag.VALUE_SORT_PRIORITY,
				Aliases: []string{"so"},
			},
			&cli.BoolFlag{
				Name:    grabberflag.KEEP_GROUPS,
				Usage:   "keep files from one group close to each other in grabbing order",
				Aliases: []string{"kg"},
			},
		},
	}
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
