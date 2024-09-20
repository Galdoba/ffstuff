package commands

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Galdoba/devtools/decidion/operator"
	"github.com/Galdoba/ffstuff/app/grabber/config"
	"github.com/Galdoba/ffstuff/pkg/stdpath"
	"github.com/urfave/cli/v2"
)

func Setup() *cli.Command {
	return &cli.Command{
		Name:  "setup",
		Usage: "TODO: Setup program files",

		BashComplete: func(*cli.Context) {
		},
		Before: func(c *cli.Context) error {
			return nil
		},
		After: func(c *cli.Context) error {
			return nil
		},
		Action: func(c *cli.Context) error {
			cfgPath := stdpath.ConfigFile()
			f, err := os.OpenFile(cfgPath, os.O_RDWR, 0666)
			defer f.Close()
			if err != nil {
				fmt.Print("can't open config file: ")
				if errors.Is(err, os.ErrNotExist) {
					fmt.Print("config file does not exist\n")
					switch operator.Confirm(fmt.Sprintf("create default config?\n%v", cfgPath)) {
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
					case false:
						fmt.Println("refused create: exit")
					}
				}
			} else {
				fmt.Println("config file detected:", cfgPath)
			}
			cfg, errC := config.Load()
			if errC != nil {
				fmt.Println("loading failed:", errC.Error())
			}
			/////////
			// if err := directoryValidation(cfg.SOURCE_ROOT_PATH); err != nil {
			// 	fmt.Println("source directory validation error:", err.Error())
			// 	sourceDir, err := operator.Input("Input source root directory:", directoryValidation)
			// 	if err != nil {
			// 		fmt.Printf("failed to get input from prompt: %v", err.Error())
			// 	}
			// 	cfg.SOURCE_ROOT_PATH = sourceDir
			// 	fmt.Printf("%v was set as a source directory\n", cfg.SOURCE_ROOT_PATH)
			// }

			if err := directoryValidation(cfg.DEFAULT_DESTINATION); err != nil {
				fmt.Println("destination directory validation error:", err.Error())
				destinationDir, err := operator.Input("Input destination root directory:", directoryValidation)
				if err != nil {
					fmt.Printf("failed to get input from prompt: %v\n", err.Error())
				}
				cfg.DEFAULT_DESTINATION = destinationDir
				fmt.Printf("%v was set as a destination directory\n", cfg.DEFAULT_DESTINATION)
			}
			logPath := stdpath.LogFile()
			if cfg.LOG == "" {
				fmt.Println("log file is not set")
				switch operator.Confirm(fmt.Sprintf("create default logfile?\n%v", logPath)) {
				case true:
					dir := filepath.Dir(logPath)
					if err := os.MkdirAll(dir, 0666); err != nil {
						fmt.Printf("log directory creation failed: %v\n", err)
						os.Exit(1)
					}
					f, err := os.Create(logPath)
					defer f.Close()
					if err != nil {
						fmt.Printf("log file creation failed: %v\n", err)
						os.Exit(1)
					}
				default:
					switch operator.Confirm(fmt.Sprintf("set custom logfile?")) {
					case true:
						logPath, err = operator.Input("Input custom log filepath:", fileValidation)
						if err != nil {
							fmt.Printf("failed to get input from prompt: %v", err.Error())
						}
					default:
						fmt.Printf("log file was not set\n")
						os.Exit(1)
					}

				}
				cfg.LOG = logPath
				fmt.Printf("%v was set as log file\n", cfg.LOG)
			}

			if err := fileValidation(cfg.LOG); err != nil {
				fmt.Println("log file validation error:", err.Error())
				logPath, err := operator.Input("Input logfile path:", fileValidation)
				if err != nil {
					fmt.Printf("failed to get input from prompt: %v", err.Error())
				}
				cfg.LOG = logPath
				fmt.Printf("%v was set as a destination directory", cfg.DEFAULT_DESTINATION)
			}
			cfg.Version = c.App.Version
			fmt.Println("config:")
			fmt.Println("version    :", cfg.Version)
			//fmt.Println("source     :", cfg.SOURCE_ROOT_PATH)
			fmt.Println("destination:", cfg.DEFAULT_DESTINATION)
			fmt.Println("log        :", cfg.LOG)
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

func directoryValidation(path string) error {
	if path == "" {
		return fmt.Errorf("directory is not set")
	}
	switch filepath.IsAbs(path) {
	case true:
		fi, err := os.Stat(path)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return fmt.Errorf("path not exists")
			}
			return fmt.Errorf("path exists, but '%v'", err)
		}
		if !fi.IsDir() {
			return fmt.Errorf("path exists, but is not directory")
		}
		return nil
	default:
		return fmt.Errorf("absolute path expected")
	}
}

func fileValidation(path string) error {
	if path == "" {
		return fmt.Errorf("directory is not set")
	}
	switch filepath.IsAbs(path) {
	case true:
		fi, err := os.Stat(path)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {

				return fmt.Errorf("path not exists")
			}
			return fmt.Errorf("path exists, but '%v'", err)
		}
		if fi.IsDir() {
			return fmt.Errorf("path exists, but is directory")
		}
		f, err := os.OpenFile(path, os.O_RDWR, 0666)
		defer f.Close()
		if err != nil {
			return fmt.Errorf("file can't be opened: %v", err)
		}

		return nil
	default:
		return fmt.Errorf("absolute path expected")
	}
}
