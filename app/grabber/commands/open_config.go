package commands

import (
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/Galdoba/ffstuff/pkg/stdpath"
	"github.com/urfave/cli/v2"
)

func Open_Config() *cli.Command {
	return &cli.Command{
		Name:        "open_config",
		Aliases:     []string{},
		Usage:       "Open config file",
		UsageText:   "",
		Description: "",
		Args:        false,
		ArgsUsage:   "",
		Category:    "",

		Action: func(c *cli.Context) error {
			args := c.Args().Slice()
			if len(args) == 0 {
				args = append(args, "default")
			}
			cfgDir := stdpath.ConfigDir()
			for _, arg := range args {
				path := cfgDir + arg + ".config"
				editor := filepath.ToSlash("C:\\Windows\\system32\\notepad.exe")
				file := filepath.ToSlash(path)
				cm := exec.Command(editor, file)
				err := cm.Run()
				if err != nil {
					fmt.Printf("error opening %v: %v\n", path, err)
				}
			}
			return nil
		},
		SkipFlagParsing: true,
	}
}
