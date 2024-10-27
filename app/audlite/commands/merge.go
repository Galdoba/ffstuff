package commands

import (
	"os"

	"github.com/Galdoba/ffstuff/pkg/logman"
	"github.com/urfave/cli/v2"
)

func Merge() *cli.Command {
	return &cli.Command{
		Name:        "merge",
		Aliases:     []string{},
		Usage:       "Direct command for merging audio files",
		UsageText:   "audman merge [command options] args...",
		Description: "Setup task for ffmpeg to merge audio files",
		Args:        false,
		ArgsUsage:   "Args Usage Text",
		Category:    "",
		Before: func(c *cli.Context) error {
			return commandInit(c)
		},
		Action: func(c *cli.Context) error {
			logman.Debug(logman.NewMessage("start 'merge'"))
			args := c.Args().Slice()
			switch len(args) {
			case 0:
				logman.Debug(logman.NewMessage("no args found"), "test read args")
			default:
				for i, arg := range args {
					logman.Debug(logman.NewMessage("found arg %v: %v", i, arg), "test read args")
				}
			}

			logman.Debug(logman.NewMessage("end   'merge'"))
			return nil
		},
		Subcommands: []*cli.Command{},
		Flags:       []cli.Flag{},
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
