package commands

import (
	"fmt"
	"os"

	"github.com/Galdoba/ffstuff/app/grabber/commands/grabberflag"
	"github.com/Galdoba/ffstuff/app/grabber/config"
	"github.com/Galdoba/ffstuff/app/grabber/internal/actions"
	"github.com/Galdoba/ffstuff/app/grabber/internal/copyprocess"
	"github.com/Galdoba/ffstuff/app/grabber/internal/origin"
	"github.com/Galdoba/ffstuff/app/grabber/internal/process"
	"github.com/Galdoba/ffstuff/app/grabber/internal/sourcesort"
	"github.com/Galdoba/ffstuff/app/grabber/internal/target"
	logman "github.com/Galdoba/ffstuff/pkg/logman"
	"github.com/urfave/cli/v2"
)

var cfg *config.Configuration

func Grab() *cli.Command {
	return &cli.Command{
		Name:      "grab",
		Aliases:   []string{},
		Usage:     "Direct command for transfering files",
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

			logman.Debug(logman.NewMessage("start 'grab'"))
			//Setup process
			logman.Debug(logman.NewMessage("check flags"))
			if err := grabberflag.ValidateGrabFlags(c); err != nil {
				return logman.Errorf("flag validation failed: %v", err)
			}

			logman.Debug(logman.NewMessage("check arguments"))
			// if err := grabberargs.ValidateGrabArguments(c.Args().Slice()...); err != nil {
			// 	return logman.Errorf("argument validation failed: %v", err)
			// }

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
			//Setup target manager
			tm, err := target.NewTargetManager(cfg,
				target.WithDestination(process.DestinationDir),
				target.WithCopyHandling(process.CopyDecidion))
			if err != nil {
				return logman.Errorf("target manager setup failed: %v", err)
			}
			//compile source list
			sources := []origin.Origin{}
			for grNum, arg := range c.Args().Slice() {
				gr := fmt.Sprintf("group_%02d", grNum)
				src := origin.New(arg, gr)
				sources = append(sources, src)
				related, err := actions.DiscoverRelatedFiles(src)
				if err != nil {
					logman.Warn("failed to discover related files: %v", err)
				}
				for _, found := range related {
					sources = append(sources, origin.New(found, gr))
				}
			}
			logman.Printf("%v source files received", len(sources))

			//Sort
			sources, err = sourcesort.Sort(process, sources...)
			if err != nil {
				return logman.Errorf("sort error: %v", err)
			}
			//targeting
			filteredSources := []origin.Origin{}
			for _, src := range sources {
				tgtName, err := tm.NewTarget(src)
				if err != nil {
					logman.Errorf("failed to create target for source '%v': %v", src.Name())
					continue
				}
				if tgtName == "" {
					switch process.CopyDecidion {
					case grabberflag.VALUE_COPY_SKIP:
						logman.Warn("skip %v", src.Name())
						continue
					default:
						logman.Errorf("failed to compile target for source '%v'", src.Name())
						continue
					}
				}
				process.SourceTargetMap[src] = tgtName
				filteredSources = append(filteredSources, src)
			}
			copyProc := copyprocess.NewCopyAction(process.SourceTargetMap,
				copyprocess.WithDestination(process.DestinationDir),
				copyprocess.WithMarkerExt(cfg.MARKER_FILE_EXTENTION),
				copyprocess.WithSourcePaths(filteredSources...),
			)
			if err := copyProc.Start(); err != nil {
				return err

			}

			logman.Debug(logman.NewMessage("end   'grab'"))
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
				Usage:   "set decidion if target file with same name exists\n	  valid values: skip, rename or overwrite",
				Aliases: []string{"c"},
			},
			&cli.StringFlag{
				Name:    grabberflag.DELETE,
				Usage:   "set which original files will be deleted after grabbing\n	  valid values: none, marker or all",
				Aliases: []string{"d"},
			},
			&cli.StringFlag{
				Name:    grabberflag.SORT,
				Usage:   "set method to decide grabbing order\n	  valid values: priority, size or none",
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
