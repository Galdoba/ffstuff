package commands

import (
	"fmt"
	"os"

	"github.com/Galdoba/ffstuff/app/grabber/config"
	"github.com/urfave/cli/v2"
)

func Health() *cli.Command {
	return &cli.Command{
		Name: "health",
		//Aliases:     []string{"fs"},
		Usage: "Check program files",

		BashComplete: func(*cli.Context) {
		},
		Before: func(c *cli.Context) error {
			return nil
		},
		After: func(c *cli.Context) error {
			return nil
		},
		Action: func(c *cli.Context) error {
			//stdpath.SetAppName(c.App.Name)
			it := newIssueTracker(c.App.Name)
			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("no config detected\nsolution: run 'grabber setup'")
			}
			for i, err := range config.Validate(cfg) {
				if err != nil {
					it.addIssue(newIssue(fmt.Sprintf("config issue %v", i+1), issueErr(err)))
				}
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
	issues  []issue
	appName string
}

func newIssueTracker(name string) *issueTracker {
	return &issueTracker{appName: name}
}

func (it *issueTracker) addIssue(issue issue) {
	it.issues = append(it.issues, issue)
}

func (it *issueTracker) report() {
	switch len(it.issues) {
	case 0:
		fmt.Fprintf(os.Stdout, "%v health: ok", it.appName)
		return
	case 1:
		fmt.Fprintf(os.Stderr, "%v health: 1 issue\n", it.appName)
	default:
		fmt.Fprintf(os.Stderr, "%v health: %v issues\n", it.appName, len(it.issues))
	}
	for i, issue := range it.issues {
		text := fmt.Sprintf("  issue %v: %v\n", i+1, issue.message)
		if issue.err != nil {
			text += fmt.Sprintf("    error: %v\n", issue.err.Error())
		}
		if issue.solution != "" {
			text += fmt.Sprintf("    solution: %v\n", issue.solution)
		}
		fmt.Printf("%v", text)
	}
}

type issue struct {
	message  string
	err      error
	solution string
}

type issueField func(*iFld)

type iFld struct {
	err      error
	solution string
}

func newIssue(message string, flds ...issueField) issue {
	is := issue{}
	is.message = message
	known := iFld{nil, ""}
	for _, set := range flds {
		set(&known)
	}
	is.err = known.err
	is.solution = known.solution
	return is
}

func issueErr(err error) issueField {
	return func(isf *iFld) {
		isf.err = err
	}
}

func issueSolution(sol string) issueField {
	return func(isf *iFld) {
		isf.solution = sol
	}
}
