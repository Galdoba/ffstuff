package commands

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
		Usage: "Setup program files",

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
				input := stdpath.ProgramDir()
				input = strings.TrimSpace(input)
				sep := string(filepath.Separator)
				input = strings.TrimSuffix(input, sep) + sep
				fmt.Println("default destination directory is", stdpath.ProgramDir())
				switch operator.Confirm(fmt.Sprintf("set custom destination directory?")) {
				case true:
					input, err = operator.Input("input default destination directory:", validation.DirectoryValidation)
					if err != nil {
						return fmt.Errorf("input failed: %v", err)
					}
				case false:
					if err := os.MkdirAll(stdpath.ProgramDir(), 0666); err != nil {
						return fmt.Errorf("program directory creation failed: %v", err)
					}
				}
				input = strings.TrimSpace(input)
				input = strings.TrimSuffix(input, sep) + sep
				cfg.DEFAULT_DESTINATION = input
				cfg.LOG = stdpath.LogFile()
				f, err = os.OpenFile(cfg.LOG, os.O_WRONLY|os.O_APPEND, 0666)
				if err != nil {
					if errors.Is(err, os.ErrNotExist) {
						fmt.Println("log file does not exist")
						label := fmt.Sprintf("Log file: %v\nCreate?", cfg.LOG)
						switch operator.Confirm(label) {
						case false:
						case true:
							if err := os.MkdirAll(stdpath.LogDir(), 0666); err != nil {
								fmt.Println(err.Error())
							}
							f, err := os.Create(stdpath.LogFile())
							if err != nil {
								fmt.Println(err.Error())
							}
							f.Close()
						}
					}
				}
				f.Close()
				cfg.CONSOLE_LOG_LEVEL = "INFO"
				cfg.FILE_LOG_LEVEL = "INFO"
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
