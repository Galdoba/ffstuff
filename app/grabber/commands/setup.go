package commands

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Galdoba/devtools/decidion/operator"
	"github.com/Galdoba/ffstuff/app/grabber/config"
	"github.com/Galdoba/ffstuff/pkg/stdpath"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

func Setup() *cli.Command {
	return &cli.Command{
		Name:  "setup",
		Usage: "Setup program files",

		Action: func(c *cli.Context) error {
			cfgPath := stdpath.ConfigFile()
			fmt.Printf("%v: %v\n", color.HiCyanString("Config Path"), cfgPath)
			wantCreateFile := false
			f, err := os.OpenFile(cfgPath, os.O_RDWR, 0666)
			if err != nil {
				switch errors.Is(err, os.ErrNotExist) {
				default:
					fmt.Printf("failed to open config file: %v\n", err)
				case true:
					wantCreateFile = true
				}
			}
			f.Close()

			if !wantCreateFile {
				fmt.Println("Config file found")
				switch operator.Confirm("Revert to default values?") {
				case false:
					fmt.Println("Setup done")
					os.Exit(2)
				case true:
					fmt.Println("Reverting config to default values")
					wantCreateFile = true
				}
			}
			if wantCreateFile {
				if err := os.MkdirAll(filepath.Dir(cfgPath), 0666); err != nil {
					fmt.Printf("Config directory creation failed: %v\n", err)
					os.Exit(1)
				}
				_, err := os.Create(cfgPath)
				if err != nil {
					fmt.Printf("Config file creation failed: %v\n", err)
					os.Exit(1)
				}
			}
			cfg := config.NewConfig(c.App.Version)
			cfg.LOG = stdpath.LogFile()
			cfg.DEFAULT_DESTINATION = stdpath.ProgramDir()
			config.Save(cfg)
			os.MkdirAll(cfg.DEFAULT_DESTINATION, 0666)
			os.MkdirAll(filepath.Dir(cfg.LOG), 0666)
			os.Create(cfg.LOG)
			errs := config.Validate(cfg)
			if len(errs) != 0 {
				fmt.Println("Config contains errors:")
				for _, err := range errs {
					fmt.Println(" ", err)
				}
			}

			bt, _ := yaml.Marshal(cfg)
			fmt.Println("===================")
			fmt.Println(string(bt))
			fmt.Println("===================")

			fmt.Printf("Setup successful. %v is ready to go!", c.App.Name)
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
