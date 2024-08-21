package commands

import (
	"fmt"
	"os"

	"github.com/Galdoba/devtools/decidion/operator"
	"github.com/urfave/cli/v2"
)

var in_dir = ""

func Menu() *cli.Command {

	cmnd := &cli.Command{
		Name: "menu",
		//Aliases:     []string{"fs"},
		Usage:     "use prompt mode",
		UsageText: "aue menu",

		BashComplete: func(*cli.Context) {
		},
		Before: func(c *cli.Context) error {
			return nil
		},
		After: func(c *cli.Context) error {
			return nil
		},
		Action: func(c *cli.Context) error {
			fmt.Println("start menu")
			result, err := operator.Select("what is nex action?", "action 1", "action 2", "action 3", "shout", "exit")
			if err != nil {
				return fmt.Errorf("selection error: %v", err)
			}
			switch result {
			case "exit":
				fmt.Println("exit now")
				return nil
			case "shout":

				return Shout().Action(c)
			default:
				fmt.Printf("%v was selected\n", result)
			}
			fmt.Println("end menu")

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

	return cmnd

}

func isDir(dir string) bool {
	f, _ := os.Stat(dir)
	return f.IsDir()
}
