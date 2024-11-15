package commands

import (
	"github.com/Galdoba/ffstuff/app/grabber/commands/grabberflag"
	"github.com/Galdoba/ffstuff/pkg/logman"
	"github.com/urfave/cli/v2"
)

func Queue() *cli.Command {
	return &cli.Command{
		Name:        "queue",
		Aliases:     []string{},
		Usage:       "Direct command for transfering files",
		UsageText:   "grabber queue [command options] args...",
		Description: "Setup and put to the process queue single copy process.",
		Args:        false,
		ArgsUsage:   "Args Usage Text",
		Category:    "",
		Before: func(c *cli.Context) error {
			return commandInit(c)
		},
		Action: func(c *cli.Context) error {
			return queue(c)
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

func queue(c *cli.Context) error {
	copyProc, err := preapareProcess(c)
	if err != nil {
		return logman.Errorf("process setup failed: %v", err)
	}
	if err = copyProc.ErrorReport(); err != nil {
		return logman.Errorf("process error: %v", err)
	}

	if err = copyProc.AddToQueue(); err != nil {
		return logman.Errorf("failed to add process to queue: %v", err)
	}

	return nil
}
