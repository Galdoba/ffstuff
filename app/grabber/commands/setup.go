package commands

import (
	"fmt"
	"os"

	"github.com/Galdoba/devtools/decidion/operator"
	"github.com/Galdoba/ffstuff/app/grabber/config"
	"github.com/Galdoba/ffstuff/app/grabber/internal/validation"
	"github.com/Galdoba/ffstuff/pkg/stdpath"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

func Setup() *cli.Command {
	return &cli.Command{
		Name:  "setup",
		Usage: "TODO: Setup program files",

		Action: func(c *cli.Context) error {
			cfgPath := stdpath.ConfigFile()
			fmt.Println("config path:", cfgPath)
			switch operator.Confirm(fmt.Sprintf("Create/revert to default values?")) {
			case false:
				return nil
			case true:
				if err := os.MkdirAll(stdpath.ConfigDir(), 0666); err != nil {
					fmt.Printf("config directory creation failed: %v\n", err)
					os.Exit(1)
				}
				f, err := os.Create(cfgPath)
				defer f.Close()
				if err != nil {
					fmt.Printf("config file creation failed: %v\n", err)
					os.Exit(1)
				}
				cfg := config.NewConfig(c.App.Version)
				input, errInput := operator.Input("input default destination directory:", validation.DirectoryValidation)
				if errInput != nil {
					return fmt.Errorf("input failed: %v", errInput)
				}
				cfg.DEFAULT_DESTINATION = input
				cfg.LOG = stdpath.LogFile()
				fmt.Println("log file path:", cfg.LOG)
				cfg.LOG_LEVEL = "DEBUG"
				if err := config.Save(cfg); err != nil {
					return fmt.Errorf("config saving failed: %v", err)
				}
				bt, _ := yaml.Marshal(cfg)
				fmt.Println("===================")
				fmt.Println(string(bt))
				fmt.Println("===================")
			}
			fmt.Println("Setup successful. grabber is ready to go!")
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
