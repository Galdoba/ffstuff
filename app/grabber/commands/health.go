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
			it := newIssueTracker()
			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("no config detected\nsolution: run 'grabber setup'")
			}
			// if cfg != nil {
			// 	if cfg.Version != c.App.Version {
			// 		it.addIssue(
			// 			newIssue(fmt.Sprintf("Config: config version (%v) does not match with app version (%v)", cfg.Version, c.App.Version),
			// 				issueSolution(
			// 					fmt.Sprintf("check config file (%v) and set version to '%v' if all is valid", stdpath.ConfigFile(), c.App.Version)),
			// 			),
			// 		)
			// 	}
			// }
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
	issues []issue
}

func newIssueTracker() *issueTracker {
	return &issueTracker{}
}

func (it *issueTracker) addIssue(issue issue) {
	it.issues = append(it.issues, issue)
}

func (it *issueTracker) report() {
	switch len(it.issues) {
	case 0:
		fmt.Fprintf(os.Stdout, "Health: ok")
		return
	case 1:
		fmt.Fprintf(os.Stderr, "Health: 1 issue\n")
	default:
		fmt.Fprintf(os.Stderr, "Health: %v issues\n", len(it.issues))
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
