package commands

import (
	"fmt"

	"github.com/Galdoba/ffstuff/app/grabber/commands/grabberargs"
	"github.com/Galdoba/ffstuff/app/grabber/commands/grabberflag"
	"github.com/Galdoba/ffstuff/app/grabber/config"
	"github.com/Galdoba/ffstuff/app/grabber/internal/actions"
	"github.com/Galdoba/ffstuff/app/grabber/internal/origin"
	"github.com/Galdoba/ffstuff/app/grabber/internal/process"
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

			fmt.Println(process)
			process.ShowOrder()
			for i, src := range sources {
				fmt.Println(i, src)
			}

			// grabProcessSettings := grab.NewProcessControl(cfg)
			// logman.Debug(logman.NewMessage("operation processing settings created"))
			// if err := grabProcessSettings.Modify(c); err != nil {
			// 	return fmt.Errorf("operation process setup failed: %v", err)
			// }
			// if err := grabProcessSettings.Assert(); err != nil {
			// 	return fmt.Errorf("operation setup invalid")
			// }
			// logman.Debug(logman.NewMessage("operation processing settings asserted"))

			// //Setup Sources
			// sourcePaths, validationErrors := grab.ValidateArgs(c.Args().Slice()...)
			// if len(validationErrors) != 0 {
			// 	logman.Warn("%v validation errors", len(validationErrors))
			// 	if len(sourcePaths) == 0 {
			// 		return logman.Errorf("no valid arguments received")
			// 	}
			// }
			// sourceList := grab.SetupSourceList(grabProcessSettings, sourcePaths...)
			// logman.Printf("total %v source files set", len(sourceList))
			// logman.Debug(logman.NewMessage("emulatng sorting"))
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
		},
	}
}
