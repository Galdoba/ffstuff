package commands

import (
	"os"

	"github.com/Galdoba/ffstuff/app/grabber/commands/grabberflag"
	"github.com/Galdoba/ffstuff/app/grabber/config"
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
		Description: "Setup and execute single process to copy source files to destination directory. Command receives filepaths as arguments.\n" +
			"If argument is 'marker' file grabber will search all related files in the same directory and add them for transfering.",
		Args:      false,
		ArgsUsage: "Args Usage Text",
		Category:  "",
		Before: func(c *cli.Context) error {
			return commandInit(c)
		},
		Action: func(c *cli.Context) error {
			logman.Info("begin grabbing")
			copyProc, err := preapareProcess(c)
			if err != nil {
				return logman.Errorf("copy process creation failed: %v", err)
			}
			if err := copyProc.Start(); err != nil {
				return err

			}
			logman.Info("grabbing complete")
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
