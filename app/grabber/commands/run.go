package commands

import (
	"fmt"
	"os"

	"github.com/Galdoba/ffstuff/app/grabber/commands/grabberflag"
	"github.com/Galdoba/ffstuff/app/grabber/internal/copyprocess"
	"github.com/Galdoba/ffstuff/pkg/logman"
	"github.com/Galdoba/ffstuff/pkg/stdpath"
	"github.com/urfave/cli/v2"
)

func Run() *cli.Command {
	return &cli.Command{
		Name:        "run",
		Aliases:     []string{},
		Usage:       "Transfering files from queue",
		UsageText:   "grabber run",
		Description: "Pull and setup processes from queue.",
		Args:        false,
		ArgsUsage:   "Args Usage Text",
		Category:    "",
		Before: func(c *cli.Context) error {
			return commandInit(c)
		},
		Action: func(c *cli.Context) error {
			return pullAndExecute(c)
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

func pullAndExecute(c *cli.Context) error {
	storage := stdpath.ProgramDir("queue")
	di, err := os.ReadDir(storage)
	copyProcs := []copyprocess.CopyProcess{}
	errors := []error{}
	if err != nil {
		return logman.Errorf("failed to read queue directory: %v", err)
	}
	for _, fi := range di {
		if fi.IsDir() {
			logman.Debug(logman.NewMessage("%v is directory: skip", fi.Name()))
		}
		cp, err := copyprocess.Reconstruct(storage + fi.Name())
		if err != nil {
			logman.Warn("failed to reconstruct process from %v: %v", fi.Name(), err)
			errors = append(errors, fmt.Errorf("failed to reconstruct process from %v: %v", fi.Name(), err))
			continue
		}

		copyProcs = append(copyProcs, cp)
	}

	for _, cp := range copyProcs {
		if err := cp.Start(); err != nil {
			logman.Warn("failed to complete process %v: %v", cp, err)
			errors = append(errors, fmt.Errorf("failed to complete process %v: %v", cp, err))
		}
	}
	if len(errors) != 0 {
		errFormat := "errors detected:\n"
		for i := range errors {
			errFormat += fmt.Sprintf("  error %0d: ", i+1) + `%v` + "\n"
		}
		return logman.Errorf(errFormat, errors)
	}

	return nil
}
