package commands

import (
	"fmt"
	"os"

	"github.com/Galdoba/ffstuff/app/demux/internal/filemarkers"
	"github.com/Galdoba/ffstuff/app/demux/internal/prompt"
	"github.com/Galdoba/ffstuff/app/demux/internal/shell"
	"github.com/urfave/cli/v2"
)

func Prompt() *cli.Command {
	return &cli.Command{
		Name: "prompt",
		//Aliases:     []string{"fs"},
		Usage:     "use prompt mode",
		UsageText: "demux prompt [args]",

		BashComplete: func(*cli.Context) {
		},
		Before: func(c *cli.Context) error {
			return nil
		},
		After: func(c *cli.Context) error {
			return nil
		},
		Action: func(c *cli.Context) error {
			fmt.Println("multiselect files")
			selectedFiles := selectFilesFromIN()
			if selectedFiles.Error != nil {
				return fmt.Errorf("file selection failed: %v", selectedFiles.Error)
			}
			fmt.Println(selectedFiles)
			for _, file := range selectedFiles.Selections {
				fmt.Println("check edit name")
				mark := filemarkers.New(file)
				fmt.Println(mark)
			}

			sh := shell.NewShell()
			fmt.Println(sh.Text())

			return nil
		},
		OnUsageError: func(cCtx *cli.Context, err error, isSubcommand bool) error {
			return nil
		},
		Subcommands:            []*cli.Command{},
		Flags:                  []cli.Flag{},
		SkipFlagParsing:        false,
		HideHelp:               false,
		HideHelpCommand:        false,
		Hidden:                 false,
		UseShortOptionHandling: false,
		HelpName:               "",
		CustomHelpTemplate:     "",
	}

}

type SelectionResult struct {
	Pick       string
	Selections []string
	Error      error
}

func selectFilesFromIN() SelectionResult {
	selRes := SelectionResult{}
	fi, err := os.ReadDir(`\\192.168.31.4\buffer\IN\`)
	if err != nil {
		return SelectionResult{Error: fmt.Errorf("can't read INPUT DIR")}
	}
	files := []string{}
	for _, f := range fi {
		if f.IsDir() {
			continue
		}
		files = append(files, f.Name())
	}
	selRes.Selections, selRes.Error = prompt.MultiSelect("Which files?", files...)
	return selRes
}
