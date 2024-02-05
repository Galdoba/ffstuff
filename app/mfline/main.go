package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/Galdoba/ffstuff/app/mfline/cmd"
	"github.com/Galdoba/ffstuff/app/mfline/config"
	"github.com/urfave/cli/v2"
)

const (
	CONFIG = "cfg"
)

func main() {
	app := cli.NewApp()

	app.Version = "v 0.1.2"
	app.Usage = "Parse media stream data from file\nRequires ffprobe to work"
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:      "use_config",
			Usage:     "use alternative config",
			TakesFile: false,
			Action: func(*cli.Context, string) error {
				return nil
			},
		},
	}

	//ДО НАЧАЛА ДЕЙСТВИЯ
	app.Before = func(c *cli.Context) error {
		cfg, err := config.Load(app.Name)
		if err != nil {
			switch {
			case strings.Contains(err.Error(), "The system cannot find "):
				cfg, err = config.NewConfig(c.App.Name)
				if err = cfg.Save(); err != nil {
					return fmt.Errorf("can't setup config: %v", err.Error())
				}
				fmt.Printf("default config created at %v: restart %v\n", cfg.Location, cfg.AppName)
				os.Exit(0)
			default:
				return err
			}
		}
		if _, err := os.ReadDir(cfg.StorageDir); err != nil {
			return fmt.Errorf("can't read storage dir: %v", err.Error())
		}

		if cfg.WriteLogs {
			f, err := os.OpenFile(cfg.LogFile, os.O_CREATE|os.O_WRONLY, 0777)
			if err != nil {
				return fmt.Errorf("can't write to log file: %v", err.Error())
			}
			defer f.Close()
		}
		for _, dir := range cfg.TrackDirs {
			if dir == "[TEMPLATE]" {
				continue
			}
			if _, err := os.ReadDir(dir); err != nil {
				return fmt.Errorf("can't read tracked dir: %v", err.Error())
			}
		}
		return nil
	}
	app.Commands = []*cli.Command{
		cmd.Sync(),
		cmd.Config(),
		cmd.Show(),
		cmd.ScanStreams(),
	}

	//ПО ОКОНЧАНИЮ ДЕЙСТВИЯ
	app.After = func(c *cli.Context) error {
		return nil
	}
	args := os.Args
	if err := app.Run(args); err != nil {
		errOut := fmt.Sprintf("%v error: %v", app.Name, err.Error())
		println(errOut)
		os.Exit(1)
	}

}
