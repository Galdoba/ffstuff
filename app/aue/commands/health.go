package commands

import (
	"fmt"

	"github.com/Galdoba/ffstuff/app/aue/config"
	"github.com/urfave/cli/v2"
)

func Health() *cli.Command {
	return &cli.Command{
		Name: "health",
		//Aliases:     []string{"fs"},
		Usage: "check program files",

		BashComplete: func(*cli.Context) {
		},
		Before: func(c *cli.Context) error {
			return nil
		},
		After: func(c *cli.Context) error {
			return nil
		},
		Action: func(c *cli.Context) error {
			it := newIssueTracker()
			_, err := config.Load()
			if err != nil {
				it.addIssue(fmt.Sprintf("Config: %v", err))
			}
			it.report()
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

type issueTracker struct {
	issues []error
}

func newIssueTracker() *issueTracker {
	return &issueTracker{}
}

func (it *issueTracker) addIssue(message string) {
	issNum := len(it.issues) + 1
	it.issues = append(it.issues, fmt.Errorf("issue %v: %v", issNum, message))
}

func (it *issueTracker) report() {
	if len(it.issues) == 0 {
		fmt.Println("No issues detected. aue is ready to go.")
		return
	}
	text := "Issues detected:"
	for _, issue := range it.issues {
		text += "\n" + issue.Error()
	}
	fmt.Println(text)
}
